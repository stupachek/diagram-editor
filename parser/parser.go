package parser

import (
	"bufio"
	"sem4/figure"
	"strings"
)

type parser struct {
	scaner *bufio.Scanner
	level  int
	line   string
}

func (p *parser) Scan() bool {
	if p.scaner.Scan() {
		p.line = strings.TrimLeft(p.scaner.Text(), "\t ")
		p.level = len(p.scaner.Text()) - len(p.line)
		return true
	}
	return false
}

func (p *parser) ParseBox() figure.AstBox {
	return figure.AstBox{
		Text: p.line,
	}

}

func (p *parser) ParseElement() figure.AstElement {
	box := p.ParseBox()
	return &box
}

func (p *parser) ParseBlock(level int) figure.AstBlock {
	children := make([]figure.AstElement, 0)
	for p.Scan() {
		children = append(children, p.ParseElement())
	}
	return figure.AstBlock{
		Children: children,
	}
}

func Parse(text string) (figure.AstBlock, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)
	p := parser{
		scaner: scanner,
		level:  0,
		line:   "",
	}
	return p.ParseBlock(0), nil
}
