package env

import "github.com/kode4food/ale/data"

type chainedNamespace struct {
	child  Namespace
	parent Namespace
}

func newChild(parent Namespace, n data.LocalSymbol) Namespace {
	e := parent.Environment()
	return chain(parent, e.New(n))
}

func chain(parent Namespace, child Namespace) *chainedNamespace {
	return &chainedNamespace{
		parent: parent,
		child:  child,
	}
}

func (ns *chainedNamespace) Snapshot(e *Environment) (Namespace, error) {
	p, err := ns.parent.Snapshot(e)
	if err != nil {
		return nil, err
	}
	c, err := ns.child.Snapshot(e)
	if err != nil {
		return nil, err
	}
	return &chainedNamespace{
		parent: p,
		child:  c,
	}, nil
}

func (ns *chainedNamespace) Environment() *Environment {
	return ns.child.Environment()
}

func (ns *chainedNamespace) Domain() data.LocalSymbol {
	return ns.child.Domain()
}

func (ns *chainedNamespace) Declared() data.LocalSymbols {
	return ns.child.Declared()
}

func (ns *chainedNamespace) Declare(n data.LocalSymbol) Entry {
	return ns.child.Declare(n)
}

func (ns *chainedNamespace) Private(n data.LocalSymbol) Entry {
	return ns.child.Private(n)
}

func (ns *chainedNamespace) Resolve(n data.LocalSymbol) (Entry, bool) {
	if e, ok := ns.child.Resolve(n); ok {
		return e, true
	}
	return resolvePublic(ns, ns.parent, n)
}
