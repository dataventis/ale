package env

import (
	"fmt"
	"regexp"
	"sort"
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.Name
		Declared() []data.Name
		Declare(data.Name) Entry
		Resolve(data.Name) (Entry, bool)
	}

	// Entry represents a namespace entry
	Entry interface {
		Owner() Namespace
		Name() data.Name
		Value() data.Value
		IsBound() bool
		Bind(data.Value)
	}

	namespace struct {
		environment *Environment
		domain      data.Name
		entries     entries
		mutex       sync.RWMutex
	}

	anonymous struct {
		Namespace
	}

	entry struct {
		owner Namespace
		name  data.Name
		value data.Value
		bound bool
		mutex sync.RWMutex
	}

	entries map[data.Name]Entry
)

// Error messages
const (
	ErrNameAlreadyBound = "name is already bound in namespace: %s"
	ErrNameNotBound     = "name is not bound in namespace: %s"
)

var privateSymbol = regexp.MustCompile(`^\^.+$`)

func (ns *namespace) Environment() *Environment {
	return ns.environment
}

func (ns *namespace) Domain() data.Name {
	return ns.domain
}

func (ns *namespace) Declared() []data.Name {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()
	e := ns.entries
	res := make([]data.Name, 0, len(e))
	for k := range e {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool {
		return string(res[i]) < string(res[j])
	})
	return res
}

func (ns *namespace) Declare(n data.Name) Entry {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	if res, ok := ns.entries[n]; ok {
		return res
	}
	e := &entry{
		owner: ns,
		name:  n,
		value: data.Nil,
		bound: false,
	}
	ns.entries[n] = e
	return e
}

func (ns *namespace) Resolve(n data.Name) (Entry, bool) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()
	if res, ok := ns.entries[n]; ok {
		return res, true
	}
	return nil, false
}

func (e *entry) Owner() Namespace {
	return e.owner
}

func (e *entry) Name() data.Name {
	return e.name
}

func (e *entry) Value() data.Value {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	if e.bound {
		return e.value
	}
	panic(fmt.Errorf(ErrNameNotBound, e.name))
}

func (e *entry) IsBound() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bound
}

func (e *entry) Bind(v data.Value) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.bound {
		panic(fmt.Errorf(ErrNameAlreadyBound, e.name))
	}
	e.value = v
	e.bound = true
}

func resolvePublic(from, in Namespace, n data.Name) (Entry, bool) {
	if isPrivateSymbol(n) && from != in {
		return nil, false
	}
	return in.Resolve(n)
}

func isPrivateSymbol(n data.Name) bool {
	return privateSymbol.MatchString(string(n))
}
