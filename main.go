package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

const unit int = 100
const blockSpacing int = unit
const blockSpacingWidth = unit
const lineStyle = `stroke="black"`

type Box struct {
	x, y          int
	width, height int
}

type Rhombus struct {
	x, y, width, height int
}

func (box *Box) draw(canvas *svg.SVG) {
	x := box.x - box.width/2
	canvas.Rect(x, box.y, box.width, box.height)
}

func (box *Box) position(x, y int) {
	box.x += x
	box.y += y
}
func (box *Box) connectTo(x, y int, canvas *svg.SVG) {
	bottom := box.y + box.height
	canvas.Line(box.x, bottom, x, y, lineStyle)
}
func (box *Box) drawLines(canvas *svg.SVG) {}

func (rhombus *Rhombus) draw(canvas *svg.SVG) {
	bottomY := rhombus.y + rhombus.height
	leftX, middleY := rhombus.left()
	rightX, _ := rhombus.right()
	canvas.Polygon([]int{rhombus.x, rightX, rhombus.x, leftX}, []int{rhombus.y, middleY, bottomY, middleY})
}
func (rhombus *Rhombus) left() (x, y int) {
	return rhombus.x - rhombus.width/2, rhombus.y + rhombus.height/2
}
func (rhombus *Rhombus) right() (x, y int) {
	return rhombus.x + rhombus.width/2, rhombus.y + rhombus.height/2
}

func (rhombus *Rhombus) position(x, y int) {
	rhombus.x += x
	rhombus.y += y
}
func (rhoumbus *Rhombus) top() (int, int) {
	return rhoumbus.x, rhoumbus.y
}

func (roumbus *Rhombus) size() (int, int) {
	return roumbus.width, roumbus.height
}

func (block *Block) position(x, y int) {
	for _, child := range block.children {
		child.position(x, y)
	}
	block.x += x
	block.y += y
}

type Figure interface {
	draw(*svg.SVG)
	top() (int, int)
	size() (int, int)
	position(x, y int)
	connectTo(x, y int, canvas *svg.SVG)
	drawLines(canvas *svg.SVG)
}

type Block struct {
	children []Figure
	x, y     int
}

func (block *Block) top() (int, int) {
	return block.x, block.y
}

func (block *Block) size() (int, int) {
	width, height := 0, 0
	for _, child := range block.children {
		w, h := child.size()
		height += h + blockSpacing
		if w > width {
			width = w
		}
	}
	if len(block.children) != 0 {
		height -= blockSpacing
	}
	return width, height
}
func (block *Block) draw(canvas *svg.SVG) {
	for _, child := range block.children {
		child.draw(canvas)
	}
}

func (block *Block) bottom() (int, int) {
	topX, topY := block.top()
	_, h := block.size()
	return topX, topY + h
}

func (block *Block) drawLines(canvas *svg.SVG) {
	for _, child := range block.children {
		child.drawLines(canvas)
	}
	for i := 0; i < len(block.children)-1; i++ {
		x, y := block.children[i+1].top()
		block.children[i].connectTo(x, y, canvas)
	}
}
func (box *Box) top() (int, int) {
	return box.x, box.y
}

func (box *Box) size() (int, int) {
	return box.width, box.height
}

type If struct {
	cond        Rhombus
	left, right Block
}

func (ifStmt *If) draw(canvas *svg.SVG) {
	ifStmt.left.draw(canvas)
	ifStmt.right.draw(canvas)
	ifStmt.cond.draw(canvas)

}
func (ifStmt *If) top() (int, int) {
	return ifStmt.cond.top()
}

func (ifStmt *If) size() (int, int) {
	_, rhombusHeigth := ifStmt.cond.size()
	leftWidth, leftHeigth := ifStmt.left.size()
	rightWidth, rightHeigth := ifStmt.right.size()
	height := rhombusHeigth + blockSpacing
	if leftHeigth > rightHeigth {
		height += leftHeigth
	} else {
		height += rightHeigth
	}
	width := leftWidth + rightWidth + 2*blockSpacingWidth
	return width, height

}

func (ifStmt *If) drawLines(canvas *svg.SVG) {
	ifRX, ifRY := ifStmt.cond.right()
	ifLX, ifLY := ifStmt.cond.left()
	leftX, leftY := ifStmt.left.top()
	rightX, rightY := ifStmt.right.top()
	canvas.Line(ifRX, ifRY, rightX, ifRY, lineStyle)
	canvas.Line(ifLX, ifLY, leftX, ifLY, lineStyle)
	canvas.Line(rightX, ifRY, rightX, rightY, lineStyle)
	canvas.Line(leftX, ifLY, leftX, leftY, lineStyle)
	ifStmt.left.drawLines(canvas)
	ifStmt.right.drawLines(canvas)
}

func (ifStmt *If) position(x int, y int) {
	ifStmt.cond.position(x, y)
	ifStmt.left.position(x, y)
	ifStmt.right.position(x, y)
}

func (ifStmt *If) connectTo(x, y int, canvas *svg.SVG) {
	leftX, leftY := ifStmt.left.bottom()
	rightX, rightY := ifStmt.right.bottom()
	middleY := 0
	if leftY > rightY {
		middleY = (leftY + y) / 2
	} else {
		middleY = (rightY + y) / 2
	}
	canvas.Line(leftX, leftY, leftX, middleY, lineStyle)
	canvas.Line(rightX, rightY, rightX, middleY, lineStyle)
	canvas.Line(leftX, middleY, rightX, middleY, lineStyle)
	ifX, _ := ifStmt.top()
	canvas.Line(ifX, middleY, x, y, lineStyle)
}

func main() {
	width := 2000
	height := 2000
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	//canvas.Circle(width/2, height/2, 100)
	//canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	astBlock := AstBlock{
		children: []AstElement{
			&AstBox{},
			&AstIf{
				left: AstBlock{
					children: []AstElement{
						&AstBox{},
						&AstBox{},
					},
				},
				right: AstBlock{
					children: []AstElement{
						&AstBox{},
						&AstBox{},
						&AstBox{},
						&AstIf{
							left: AstBlock{
								children: []AstElement{
									&AstBox{},
								},
							},
							right: AstBlock{},
						},
						&AstBox{},
					},
				},
			},

			&AstBox{},
			&AstBox{},
		},
	}
	b := astBlock.toFigure(500, 250)

	b.draw(canvas)
	b.drawLines(canvas)
	canvas.End()
}
