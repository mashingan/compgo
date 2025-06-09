package comp

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store  map[string]Symbol
	numdef int
	scoped *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: map[string]Symbol{}}
}

func (s *SymbolTable) Define(sym string) Symbol {
	syms := Symbol{sym, GlobalScope, s.numdef}
	if s.scoped != nil {
		syms.Scope = LocalScope
	}
	s.numdef++
	s.store[sym] = syms
	return syms
}

func (s *SymbolTable) Resolve(sym string) (Symbol, bool) {
	st := s
	ss, ok := st.store[sym]
	for !ok && st != nil {
		ss, ok = st.store[sym]
		if !ok {
			st = st.scoped
		}
	}
	return ss, ok
}

func NewFrameSymbolTable(outer *SymbolTable) *SymbolTable {
	st := NewSymbolTable()
	st.scoped = outer
	return st
}
