package backtest

//experimentally not finished

type FilterOperation struct {
	Name string
	Op   Filter
}

type Tree struct {
	P    *Tree
	C    []*Tree
	Op   FilterOperation
	Tr   []*Trade
	Info string
}

func NewTree(tr []*Trade, op []FilterOperation) *Tree {
	return &Tree{
		P:    nil,
		C:    nil,
		Op:   FilterOperation{},
		Tr:   nil,
		Info: "",
	}
}
