package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNull,
		isa.True:  isa.RetTrue,
	}

	literalKeys = basics.MapKeys(literalReturnMap)

	literalReturnPatterns = visitor.Pattern{
		literalKeys,
		{isa.Return},
	}
)

func literalReturns(root visitor.Node) visitor.Node {
	visitor.Replace(root, literalReturnPatterns, literalReturnMapper)
	return root
}

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc := i[0].Opcode()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
