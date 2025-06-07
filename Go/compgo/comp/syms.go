package comp

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store  map[string]Symbol
	numdef int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: map[string]Symbol{}}
}

func (s *SymbolTable) Define(sym string) Symbol {
	syms := Symbol{sym, GlobalScope, s.numdef}
	s.numdef++
	s.store[sym] = syms
	return syms
}

func (s *SymbolTable) Resolve(sym string) (Symbol, bool) {
	ss, ok := s.store[sym]
	return ss, ok
}
