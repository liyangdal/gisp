package gisp

import (
	//"fmt"
	"reflect"
	"testing"
	tm "time"

	p "github.com/Dwarfartisan/goparsec"
)

func TestDotTime(t *testing.T) {
	now := tm.Now()
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	slot := VarSlot(TIMEOPTION)
	slot.Set(now)
	g.Defvar("now", slot)
	year := Int(now.Year())
	y, err := g.Parse("(now.Year)") //g.Eval(List{AA("now.Year")})
	if err != nil {
		t.Fatalf("except (now.Year) equal to now.Year() as %v but got error %v", year, err)
	}
	if !reflect.DeepEqual(year, Int(y.(int))) {
		t.Fatalf("except (now.Year) equal to now.Year() but got %v and %v", year, y)
	}
}

func TestDotParser(t *testing.T) {
	data := "now.Year"
	st := p.MemoryParseState(data)
	re, err := p.Bind(AtomParser, DotSuffixParser)(st)
	if err != nil {
		t.Fatalf("except a Dot but error %v", err)
	}
	t.Log(re)
}

type Box map[string]interface{}

func (b Box) Get(name string) interface{} {
	return b[name]
}

func TestDotMap(t *testing.T) {
	box := Box{
		"a": Quote{AA("a")},
		"b": Quote{AA("bb")},
		"c": Quote{AA("ccc")},
	}
	bv := reflect.ValueOf(box)
	get := bv.MethodByName("Get")
	res := get.Call([]reflect.Value{reflect.ValueOf("b")})
	if !reflect.DeepEqual(res[0].Interface(), box["b"]) {
		t.Fatalf("except %v but got %v", box["b"], res[0].Interface())
	}
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.DefAs("box", box)
	c, err := g.Parse(`(box.Get "c")`)
	if err != nil {
		t.Fatalf("excpet got b but error %v", err)
	}
	if !reflect.DeepEqual(c, box["c"]) {
		t.Fatalf("except %v but got %v", box["c"], c)
	}
}
