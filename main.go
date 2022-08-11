package main

import (
	"io/ioutil"
	"sem4/figure"
	"sem4/parser"
)

func main() {
	ar, _ := ioutil.ReadFile("code")
	r, _ := parser.Parse(string(ar))
	figure.DrawBlock(r, "out.png")
}
