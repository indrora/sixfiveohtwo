package assembler

import "fmt"

type Symbol struct {
	Name    string
	Address uint16
	Defined bool
}

type SymbolTable struct {
	symbols map[string]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
	}
}

func (st *SymbolTable) Define(name string, address uint16) {
	symbol, exists := st.symbols[name]
	if exists {
		symbol.Address = address
		symbol.Defined = true
	} else {
		st.symbols[name] = &Symbol{
			Name:    name,
			Address: address,
			Defined: true,
		}
	}
}

func (st *SymbolTable) Reference(name string) *Symbol {
	symbol, exists := st.symbols[name]
	if !exists {
		symbol = &Symbol{
			Name:    name,
			Address: 0,
			Defined: false,
		}
		st.symbols[name] = symbol
	}
	return symbol
}

func (st *SymbolTable) Resolve(name string) (uint16, error) {
	symbol, exists := st.symbols[name]
	if !exists || !symbol.Defined {
		return 0, fmt.Errorf("undefined symbol: %s", name)
	}
	return symbol.Address, nil
}

func (st *SymbolTable) IsResolved() bool {
	for _, symbol := range st.symbols {
		if !symbol.Defined {
			return false
		}
	}
	return true
}

func (st *SymbolTable) GetUndefined() []string {
	var undefined []string
	for name, symbol := range st.symbols {
		if !symbol.Defined {
			undefined = append(undefined, name)
		}
	}
	return undefined
}