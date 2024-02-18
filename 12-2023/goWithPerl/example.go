package main

import (
	"fmt"

	"github.com/bradfitz/campher/perl"
)

func main() {
	p := perl.NewInterpreter()
	p.Eval(`$foo = "bar";`)
	fmt.Println("foo is:", p.EvalString("$foo"))

	cv := p.Eval(`sub { my ($a1, $a2, $func) = @_; $func->($a1, $a2); }`).CV()
	var ret *SV
	p := perl.NewInterpreter()
	foo := p.Eval(`$foo = sub {
 my ($op, $v1, $v2) = @_;
 return "Perl says: " . $op->($v1, $v2);
};`)
	concat := func(a, b string) string { return a + b }
	fmt.Println("concat:", foo.CV().Call(concat, "foo", "bar"))

	base := 0
	add := func(a, b int) int {
		base++
		return a + b + base
	}
	fmt.Println("add:", foo.CV().Call(add, 1, "40"))
	fmt.Println("add:", foo.CV().Call(add, 1, "40"))
}
