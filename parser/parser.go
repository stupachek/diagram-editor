package parser

import (
	"bufio"
	"sem4/figure"
	"strings"
)

type parser struct {
	scaner   *bufio.Scanner
	level    int
	line     string
	lineBuf  *string
	levelBuf int
}

func (p *parser) Scan() bool {
	if p.lineBuf != nil {
		p.line = *p.lineBuf
		p.level = p.levelBuf
		p.lineBuf = nil
		return true
	}
	if p.scaner.Scan() {
		p.line = strings.TrimLeft(p.scaner.Text(), "\t ")
		p.level = len(p.scaner.Text()) - len(p.line)
		return true
	}
	return false
}

func (p *parser) Return() {
	p.lineBuf = &p.line
	p.levelBuf = p.level
}

func (p *parser) ParseBox() figure.AstBox {
	return figure.AstBox{
		Text: p.line,
	}

}

func (p *parser) ParseInput() figure.AstInput {
	text := strings.TrimPrefix(p.line, "input ")
	return figure.AstInput{
		Text: text,
	}

}

func (p *parser) ParseElement() figure.AstElement {
	if strings.HasPrefix(p.line, "if ") {
		ifStmt := p.ParseIf()
		return &ifStmt
	} else if strings.HasPrefix(p.line, "input ") {
		input := p.ParseInput()
		return &input
	} else {
		box := p.ParseBox()
		return &box
	}

}

func (p *parser) ParseBlock(level int) figure.AstBlock {
	children := make([]figure.AstElement, 0)
	if level == 0 {
		children = append(children, &figure.AstStartStop{
			Text: "Begin",
		})
	}

	for {
		children = append(children, p.ParseElement())
		if !p.Scan() {
			break
		}
		if p.level != level {
			p.Return()
			break
		}
	}

	if level == 0 {
		children = append(children, &figure.AstStartStop{
			Text: "Stop",
		})
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
	p.Scan()
	return p.ParseBlock(0), nil
}

func (p *parser) ParseIf() figure.AstIf {
	text := strings.TrimPrefix(p.line, "if ")
	if !p.Scan() {
		return figure.AstIf{
			Left:  figure.AstBlock{},
			Right: figure.AstBlock{},
			Text:  text,
		}
	}
	left := p.ParseBlock(p.level)
	right := figure.AstBlock{}
	if strings.HasPrefix(p.line, "else ") && p.Scan() && p.Scan() {
		right = p.ParseBlock(p.level)
	}
	return figure.AstIf{
		Left:  left,
		Right: right,
		Text:  text,
	}

}
