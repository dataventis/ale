package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"gitlab.com/kode4food/ale"
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/bootstrap"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/docstring"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/read"
)

// UserDomain is the name of the namespace that the REPL starts in
const UserDomain = api.Name("user")

const (
	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> " + code
	cont   = domain + "[%d]" + dgray + nlMarker + "   " + code

	output = bold + "%s" + reset
	good   = domain + result + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
)

type (
	any      interface{}
	sentinel struct{}

	// REPL manages a FromScanner-Eval-Print Loop
	REPL struct {
		buf bytes.Buffer
		rl  *readline.Instance
		idx int
	}
)

var (
	anyChar = regexp.MustCompile(".")

	nothing = &sentinel{}

	openers = map[rune]rune{')': '(', ']': '[', '}': '{'}
	closers = map[rune]rune{'(': ')', '[': ']', '{': '}'}

	ns = bootstrap.TopLevelManager().GetQualified(UserDomain)
)

// NewREPL instantiates a new REPL instance
func NewREPL() *REPL {
	repl := &REPL{}

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: getHistoryFile(),
		Painter:     repl,
	})

	if err != nil {
		panic(err)
	}

	repl.rl = rl
	repl.idx = 1

	return repl
}

// Run will perform the Eval-Print-Loop
func (r *REPL) Run() {
	defer r.rl.Close()

	fmt.Println(ale.AppName, ale.Version)
	help()
	r.setInitialPrompt()

	for {
		line, err := r.rl.Readline()
		r.buf.WriteString(line + "\n")
		fmt.Print(reset)

		if err != nil {
			emptyBuffer := isEmptyString(r.buf.String())
			if err == readline.ErrInterrupt && !emptyBuffer {
				r.reset()
				continue
			}
			break
		}

		if isEmptyString(line) {
			continue
		}

		if !r.evalBuffer() {
			r.setContinuePrompt()
			continue
		}

		r.reset()
	}
	shutdown()
}

func (r *REPL) reset() {
	r.buf.Reset()
	r.idx++
	r.setInitialPrompt()
}

func (r *REPL) setInitialPrompt() {
	name := ns.Domain()
	r.setPrompt(fmt.Sprintf(prompt, name, r.idx))
}

func (r *REPL) setContinuePrompt() {
	r.setPrompt(fmt.Sprintf(cont, r.nsSpace(), r.idx))
}

func (r *REPL) setPrompt(s string) {
	r.rl.SetPrompt(s)
}

func (r *REPL) nsSpace() string {
	ns := string(ns.Domain())
	return anyChar.ReplaceAllString(ns, " ")
}

func (r *REPL) evalBuffer() (completed bool) {
	defer func() {
		if err := toError(recover()); err != nil {
			if isRecoverable(err) {
				completed = false
				return
			}
			r.outputError(err)
			completed = true
		}
	}()

	res := eval.String(ns, api.String(r.buf.String()))
	r.outputResult(res)
	return true
}

func (r *REPL) outputResult(v any) {
	if v == nothing {
		return
	}
	var sv any
	if s, ok := v.(api.Value); ok {
		sv = api.MaybeQuoteString(s)
	} else {
		sv = v
	}
	res := fmt.Sprintf(good, r.nsSpace(), r.idx, sv)
	fmt.Println(res)
}

func (r *REPL) outputError(err error) {
	msg := err.Error()
	res := fmt.Sprintf(bad, r.nsSpace(), r.idx, msg)
	fmt.Println(res)
}

func (s *sentinel) String() string {
	return ""
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func toError(i interface{}) error {
	if i == nil {
		return nil
	}
	switch typed := i.(type) {
	case error:
		return typed
	case api.Value:
		return fmt.Errorf(typed.String())
	default:
		panic(fmt.Errorf("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	msg := err.Error()
	return msg == read.ListNotClosed ||
		msg == read.VectorNotClosed ||
		msg == read.MapNotClosed ||
		msg == read.StringNotTerminated
}

func use(args ...api.Value) api.Value {
	arity.AssertFixed(1, len(args))
	n := args[0].(api.LocalSymbol).Name()
	old := ns
	ns = ns.Manager().GetQualified(n)
	if old != ns {
		fmt.Println()
	}
	return nothing
}

func shutdown(args ...api.Value) api.Value {
	arity.AssertFixed(0, len(args))
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	rg := rand.New(rs)
	idx := rg.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func debugInfo(args ...api.Value) api.Value {
	arity.AssertFixed(0, len(args))
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
	return nothing
}

func cls(args ...api.Value) api.Value {
	arity.AssertFixed(0, len(args))
	fmt.Println(clear)
	return nothing
}

func formatForREPL(s string) string {
	md := formatMarkdown(s)
	lines := strings.Split(md, "\n")
	var out []string
	out = append(out, "")
	for _, l := range lines {
		if isEmptyString(l) {
			out = append(out, l)
		} else {
			out = append(out, "  "+l)
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n")
}

func help(args ...api.Value) api.Value {
	arity.AssertFixed(0, len(args))
	md := string(docstring.Get("help"))
	fmt.Println(formatForREPL(md))
	return nothing
}

func doc(args ...api.Value) api.Value {
	arity.AssertFixed(1, len(args))
	sym := args[0].(api.LocalSymbol)
	name := string(sym.Name())
	if docstring.Exists(name) {
		docStr := docstring.Get(name)
		f := formatForREPL(string(docStr))
		fmt.Println(f)
		return nothing
	}
	panic(fmt.Errorf("symbol is not documented: %s", sym))
}

func getBuiltInsNamespace() namespace.Type {
	return ns.Manager().GetRoot()
}

func registerBuiltIn(n api.Name, v api.Value) {
	ns := getBuiltInsNamespace()
	ns.Bind(n, v)
}

// GetNS allows the tests to get at the namespace
func GetNS() namespace.Type {
	return ns
}

func registerREPLBuiltIns() {
	registerBuiltIn("use", api.NormalFunction(use))
	registerBuiltIn("quit", api.ApplicativeFunction(shutdown))
	registerBuiltIn("debug", api.ApplicativeFunction(debugInfo))
	registerBuiltIn("cls", api.ApplicativeFunction(cls))
	registerBuiltIn("help", api.ApplicativeFunction(help))
	registerBuiltIn("doc", api.NormalFunction(doc))
}

func getScreenWidth() int {
	return readline.GetScreenWidth()
}

func getHistoryFile() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return path.Join(usr.HomeDir, ".ale-history")
}

func init() {
	bootstrap.Into(ns.Manager())
	registerREPLBuiltIns()
}
