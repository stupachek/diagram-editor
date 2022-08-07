package main

type AstBox struct {
	text string
}

type AstIf struct {
	left, right AstBlock
	text        string
}

type AstBlock struct {
	children []AstElement
}

type AstElement interface {
	toFigure(x, y int) Figure
}

func (astBox *AstBox) toFigure(x, y int) Figure {
	newB := newBox(astBox.text, x, y)
	return &newB
}

func (astIf *AstIf) toFigure(x, y int) Figure {
	rhombus := newRhombus(astIf.text, x, y)
	left := astIf.left.toFigure(0, 0)
	right := astIf.right.toFigure(0, 0)
	blockY := y + rhombus.height + blockSpacing
	widthLeft, _ := left.size()
	widthRigth, _ := right.size()
	leftX := x - widthLeft/2 - blockSpacingWidth
	rightX := x + widthRigth/2 + blockSpacingWidth
	left.position(leftX, blockY)
	right.position(rightX, blockY)
	return &If{
		cond:  rhombus,
		left:  left,
		right: right,
	}
}

func (astBlock *AstBlock) toFigure(x, y int) Block {
	children := []Figure{}
	yStart := y
	for _, child := range astBlock.children {
		childFigure := child.toFigure(x, y)
		children = append(children, childFigure)
		_, h := childFigure.size()
		y += h + blockSpacing
	}
	return Block{
		children: children, x: x, y: yStart,
	}
}
