package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/docstring"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/console"
	"github.com/kode4food/ale/internal/markdown"
	"github.com/kode4food/ale/read"
)

type (
	sentinel struct{}

	// REPL manages a FromScanner-Eval-Print Loop
	REPL struct {
		buf bytes.Buffer
		rl  *readline.Instance
		idx int
	}
)

// Error messages
const (
	ErrSymbolNotDocumented = "symbol not documented: %s"
)

const (
	// UserDomain is the name of the namespace that the REPL starts in
	UserDomain = data.Name("user")

	domain = console.Domain + "%s" + console.Reset + " "
	prompt = domain + "[%d]> " + console.Code
	cont   = domain + "[%d]" + console.NewLine + "   " + console.Code

	output = console.Bold + "%s" + console.Reset
	good   = domain + console.Result + "[%d]= " + output
	bad    = domain + console.Error + "[%d]! " + output
)

var (
	anyChar   = regexp.MustCompile(".")
	notPaired = fmt.Sprintf(read.ErrPrefixedNotPaired, "")
	nothing   = new(sentinel)

	ns = bootstrap.TopLevelEnvironment().GetQualified(UserDomain)
)

// NewREPL instantiates a new REPL instance
func NewREPL() *REPL {
	repl := new(REPL)

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: getHistoryFile(),
		Painter:     console.Painter(),
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
	defer func() {
		_ = r.rl.Close()
	}()

	fmt.Println(ale.AppName, ale.Version)
	help()
	r.setInitialPrompt()

	for {
		line, err := r.rl.Readline()
		r.buf.WriteString(line + "\n")
		fmt.Print(console.Reset)

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

	res := eval.String(ns, data.String(r.buf.String()))
	r.outputResult(res)
	return true
}

func (r *REPL) outputResult(v any) {
	if v == nothing {
		return
	}
	var sv any
	if s, ok := v.(data.Value); ok {
		sv = data.MaybeQuoteString(s)
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

func (*sentinel) Equal(data.Value) bool {
	return false
}

func (*sentinel) String() string {
	return ""
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func toError(i any) error {
	if i == nil {
		return nil
	}
	switch i := i.(type) {
	case error:
		return i
	case data.Value:
		return errors.New(i.String())
	default:
		// Programmer error
		panic(fmt.Sprintf("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
	}
	msg := err.Error()
	return msg == read.ErrListNotClosed ||
		msg == read.ErrVectorNotClosed ||
		msg == read.ErrMapNotClosed ||
		msg == read.ErrStringNotTerminated ||
		strings.HasPrefix(msg, notPaired)
}

func use(args ...data.Value) data.Value {
	data.AssertFixed(1, len(args))
	n := args[0].(data.LocalSymbol).Name()
	old := ns
	ns = ns.Environment().GetQualified(n)
	if old != ns {
		fmt.Println()
	}
	return nothing
}

func shutdown(...data.Value) data.Value {
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	rg := rand.New(rs)
	idx := rg.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func debugInfo(...data.Value) data.Value {
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
	return nothing
}

func cls(...data.Value) data.Value {
	fmt.Println(console.Clear)
	return nothing
}

func formatForREPL(s string) string {
	md := markdown.FormatMarkdown(s)
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

func help(...data.Value) data.Value {
	md, err := docstring.Get("help")
	if err != nil {
		// Programmer error, help is missing
		panic(err)
	}
	fmt.Println(formatForREPL(md))
	return nothing
}

var docTemplate = strings.Join([]string{
	"---",
	`description: "display documentation of a form"`,
	`usage: "(doc form)"`,
	"---",
	"Where `form` is any of the core language symbols:",
	"",
	"%s",
}, "\n")

func doc(args ...data.Value) data.Value {
	if len(args) != 0 {
		docSymbol(args[0].(data.LocalSymbol))
		return nothing
	}
	docSymbolList()
	return nothing
}

func docSymbol(sym data.Symbol) {
	name := string(sym.Name())
	if docStr, err := docstring.Get(name); err == nil {
		f := formatForREPL(docStr)
		fmt.Println(f)
		return
	}
	if name == "doc" {
		docSymbolList()
		return
	}
	panic(fmt.Errorf(ErrSymbolNotDocumented, sym))
}

func docSymbolList() {
	names := docstring.Names()
	slices.Sort(names)
	escapeNames(names)
	joined := strings.Join(names, ", ")
	f := formatForREPL(fmt.Sprintf(docTemplate, joined))
	fmt.Println(f)
}

func escapeNames(names []string) {
	for i, n := range names {
		if strings.Contains("`*_", n[0:1]) {
			names[i] = "\\" + n
		}
	}
}

func getBuiltInsNamespace() env.Namespace {
	return ns.Environment().GetRoot()
}

func registerBuiltIn(n data.Name, v data.Value) {
	ns := getBuiltInsNamespace()
	ns.Declare(n).Bind(v)
}

// GetNS allows the tests to get at the namespace
func GetNS() env.Namespace {
	return ns
}

func registerREPLBuiltIns() {
	registerBuiltIn("quit", data.Applicative(shutdown, 0))
	registerBuiltIn("debug", data.Applicative(debugInfo, 0))
	registerBuiltIn("cls", data.Applicative(cls, 0))
	registerBuiltIn("help", data.Applicative(help, 0))
	registerBuiltIn("use", data.Normal(use, 1))
	registerBuiltIn("doc", data.Normal(doc, 0, 1))
}

func getHistoryFile() string {
	if usr, err := user.Current(); err == nil {
		return path.Join(usr.HomeDir, ".ale-history")
	}
	return ""
}

func init() {
	registerREPLBuiltIns()
}
