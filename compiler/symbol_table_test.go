// compiler/symbol_table_test.go

package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Index: 0, Scope: GlobalScope},
		"b": {Name: "b", Index: 1, Scope: GlobalScope},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected=%+v, got=%+v", expected["b"], b)
	}

}

func TestResolveGlobal(t *testing.T) {
	s := NewSymbolTable()

	s.Define("a")
	s.Define("b")

	expected := []Symbol{
		{Name: "a", Index: 0, Scope: GlobalScope},
		{Name: "b", Index: 1, Scope: GlobalScope},
	}

	for i, sym := range expected {
		sym, ok := s.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %q not resolvable", sym.Name)
		}
		if sym != expected[i] {
			t.Errorf("expected=%+v, got=%+v", expected[i], sym)
		}
	}
}
