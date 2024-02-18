package main

import "fmt"

var concat = func(a, b string) string { return a + b }
var myPerlFunc = perlsub{
	my ( $op, $v1, $v2)=@_;
	return "Perl says: " . $op->($v1, $v2);
}

func main() {
	fmt.Println(myPerlFunc(concat, "foo", "bar"))
}
