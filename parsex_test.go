package gisp

import (
	//"fmt"
	"reflect"
	"testing"

	px "github.com/Dwarfartisan/goparsec/parsex"
)

func TestParsexBasic(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "px": Parsex})

	digit := px.Bind(px.Many1(px.Digit), px.ReturnString)
	data := "344932454094325"
	state := NewStringState(data)
	pxre, err := digit(state)
	if err != nil {
		t.Fatalf("except \"%v\" pass test many1 digit but error:%v", data, err)
	}

	src := "(let ((st (px.state \"" + data + `")))
	(var data ((px.many1 px.digit) st))
	(px.s2str data))
	`
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("except \"%v\" pass gisp many1 digit but error:%v", src, err)
	}
	t.Logf("from gisp: %v", gre)
	t.Logf("from parsex: %v", pxre)
	if !reflect.DeepEqual(pxre, gre) {
		t.Fatalf("except got \"%v\" from gisp equal \"%v\" from parsex", gre, pxre)
	}
}

func TestParsexRune(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "px": Parsex})
	//data := "Here is a Rune : 'a' and a is't a rune. It is a word in sentence."
	data := "'a' and a is't a rune. It is a word in sentence."
	state := NewStringState(data)
	pre, err := px.Between(px.Rune('\''), px.Rune('\''), px.AnyRune)(state)
	if err != nil {
		t.Fatalf("except found rune expr from \"%v\" but error:%v", data, err)
	}
	src := `
	(let ((st (px.state "` + data + `")))
		((px.between (px.rune '\'') (px.rune '\'') px.anyone) st))
	`

	//fmt.Println(src)
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("except \"%v\" pass gisp '<rune>' but error:%v", src, err)
	}
	t.Logf("from gisp: %v", gre)
	t.Logf("from parsec: %v", pre)
	if !reflect.DeepEqual(pre, gre) {
		t.Fatalf("except got \"%v\" from gisp equal \"%v\" from parsec", gre, pre)
	}
}
