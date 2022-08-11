package figure

type AstBox struct {
	Text string
}

type AstIf struct {
	Left, Right AstBlock
	Text        string
}

type AstBlock struct {
	Children []AstElement
}

type AstElement interface {
	toFigure(x, y int) Figure
}

func (astBox *AstBox) toFigure(x, y int) Figure {
	newB := newBox(astBox.Text, x, y)
	return &newB
}

func (astIf *AstIf) toFigure(x, y int) Figure {
	rhombus := newRhombus(astIf.Text, x, y)
	left := astIf.Left.toFigure(0, 0)
	right := astIf.Right.toFigure(0, 0)
	blockY := y + rhombus.height + blockSpacing
	widthLeft, _ := left.size()
	widthRigth, _ := right.size()
	leftX := x - Max(widthLeft, rhombus.width)/2 - blockSpacingWidth
	rightX := x + Max(widthRigth, rhombus.width)/2 + blockSpacingWidth
	left.position(leftX, blockY)
	right.position(rightX, blockY)
	return &If{
		cond:  rhombus,
		left:  left,
		right: right,
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func (astBlock *AstBlock) toFigure(x, y int) Block {
	children := []Figure{}
	yStart := y
	for _, child := range astBlock.Children {
		childFigure := child.toFigure(x, y)
		children = append(children, childFigure)
		_, h := childFigure.size()
		y += h + blockSpacing
	}
	return Block{
		children: children, x: x, y: yStart,
	}
}
