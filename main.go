package main

import (
	"github.com/fogleman/gg"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const unit int = 100
const blockSpacing int = unit
const blockSpacingWidth = unit
const textHeight = 30
const verticalMargins = 30
const horizontalMargins = 30

var face truetype.IndexableFace

type Box struct {
	x, y          int
	width, height int
	text          string
}

type Rhombus struct {
	x, y, width, height int
	text                string
}

func (box *Box) draw(canvas *gg.Context) {
	x := box.x - box.width/2
	canvas.DrawRectangle(float64(x), float64(box.y), float64(box.width), float64(box.height))
	canvas.Stroke()
	middleY := box.height/2 + box.y
	canvas.DrawStringAnchored(box.text, float64(box.x), float64(middleY), 0.5, 0.35)
}

func (box *Box) position(x, y int) {
	box.x += x
	box.y += y
}
func (box *Box) connectTo(x, y int, canvas *gg.Context) {
	bottom := box.y + box.height
	canvas.DrawLine(float64(box.x), float64(bottom), float64(x), float64(y))
	canvas.Stroke()
}
func (box *Box) drawLines(canvas *gg.Context) {}

func (rhombus *Rhombus) draw(canvas *gg.Context) {
	bottomY := rhombus.y + rhombus.height
	leftX, middleY := rhombus.left()
	rightX, _ := rhombus.right()
	canvas.MoveTo(float64(rhombus.x), float64(rhombus.y))
	canvas.LineTo(float64(rightX), float64(middleY))
	canvas.LineTo(float64(rhombus.x), float64(bottomY))
	canvas.LineTo(float64(leftX), float64(middleY))
	canvas.LineTo(float64(rhombus.x), float64(rhombus.y))
	canvas.Stroke()
	middleY_2 := rhombus.height/2 + rhombus.y
	canvas.DrawStringAnchored(rhombus.text, float64(rhombus.x), float64(middleY_2), 0.5, 0.35)
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
	draw(*gg.Context)
	top() (int, int)
	size() (int, int)
	position(x, y int)
	connectTo(x, y int, canvas *gg.Context)
	drawLines(canvas *gg.Context)
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
func (block *Block) draw(canvas *gg.Context) {
	for _, child := range block.children {
		child.draw(canvas)
	}
}

func (block *Block) bottom() (int, int) {
	topX, topY := block.top()
	_, h := block.size()
	return topX, topY + h
}

func (block *Block) drawLines(canvas *gg.Context) {
	for _, child := range block.children {
		child.drawLines(canvas)
	}
	for i := 0; i < len(block.children)-1; i++ {
		x, y := block.children[i+1].top()
		block.children[i].connectTo(x, y, canvas)
		Arrow(x, y, canvas)
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

func (ifStmt *If) draw(canvas *gg.Context) {
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

func (ifStmt *If) drawLines(canvas *gg.Context) {
	ifRX, ifRY := ifStmt.cond.right()
	ifLX, ifLY := ifStmt.cond.left()
	leftX, leftY := ifStmt.left.top()
	rightX, rightY := ifStmt.right.top()
	canvas.DrawLine(float64(ifRX), float64(ifRY), float64(rightX), float64(ifRY))
	canvas.Stroke()
	canvas.DrawLine(float64(ifLX), float64(ifLY), float64(leftX), float64(ifLY))
	canvas.Stroke()
	canvas.DrawLine(float64(rightX), float64(ifRY), float64(rightX), float64(rightY))
	canvas.Stroke()
	canvas.DrawLine(float64(leftX), float64(ifLY), float64(leftX), float64(leftY))
	canvas.Stroke()
	ifStmt.left.drawLines(canvas)
	if !ifStmt.left.IsEmpty() {
		Arrow(leftX, leftY, canvas)
	}
	if !ifStmt.right.IsEmpty() {
		Arrow(rightX, rightY, canvas)
	}
	ifStmt.right.drawLines(canvas)
}

func (block *Block) IsEmpty() bool {
	return len(block.children) == 0
}

func (ifStmt *If) position(x int, y int) {
	ifStmt.cond.position(x, y)
	ifStmt.left.position(x, y)
	ifStmt.right.position(x, y)
}

func (ifStmt *If) connectTo(x, y int, canvas *gg.Context) {
	leftX, leftY := ifStmt.left.bottom()
	rightX, rightY := ifStmt.right.bottom()
	middleY := 0
	if leftY > rightY {
		middleY = (leftY + y) / 2
	} else {
		middleY = (rightY + y) / 2
	}
	canvas.DrawLine(float64(leftX), float64(leftY), float64(leftX), float64(middleY))
	canvas.Stroke()
	canvas.DrawLine(float64(rightX), float64(rightY), float64(rightX), float64(middleY))
	canvas.Stroke()
	canvas.DrawLine(float64(leftX), float64(middleY), float64(rightX), float64(middleY))
	canvas.Stroke()
	ifX, _ := ifStmt.top()
	canvas.DrawLine(float64(ifX), float64(middleY), float64(x), float64(y))
	canvas.Stroke()
}

func newBox(text string, x, y int) Box {
	w, h := TextSize(text)
	return Box{
		x:      x,
		y:      y,
		width:  w + horizontalMargins,
		height: h + verticalMargins,
		text:   text,
	}
}

func TextSize(text string) (width, height int) {
	d := &font.Drawer{
		Face: face,
	}
	w := d.MeasureString(text)
	return int(w >> 6), textHeight
}

func Arrow(x, y int, canvas *gg.Context) {
	h := 10
	w := 15
	canvas.MoveTo(float64(x-w/2), float64(y-h))
	canvas.LineTo(float64(x+w/2), float64(y-h))
	canvas.LineTo(float64(x), float64(y))
	canvas.LineTo(float64(x-w/2), float64(y-h))
	canvas.Fill()
}

func newRhombus(text string, x, y int) Rhombus {

	w, h := TextSize(text)
	h_rhombus := h * 3
	w_hrombus := Max((w * h_rhombus / (h_rhombus - h)), h_rhombus)
	return Rhombus{
		x:      x,
		y:      y,
		width:  w_hrombus,
		height: h_rhombus,
		text:   text,
	}
}

func main() {
	width := 2000
	height := 2000
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	face = truetype.NewFace(font, &truetype.Options{Size: textHeight})
	canvas := gg.NewContext(width, height)
	canvas.SetFontFace(face)
	astBlock := AstBlock{
		children: []AstElement{
			&AstBox{"123"},
			&AstIf{text: "xuy_zalupa\nqweqe\nqweqe",
				left: AstBlock{
					children: []AstElement{
						&AstBox{"*"},
						&AstBox{"KISS U"},
					},
				},
				right: AstBlock{
					children: []AstElement{
						&AstBox{"12312312312321"},
						&AstBox{"12321321321312312313123"},
						&AstBox{"QWEQWEQWEQ"},
						&AstIf{text: "<3",
							left: AstBlock{
								children: []AstElement{
									&AstBox{"dfghjkl"},
								},
							},
							right: AstBlock{},
						},
						&AstBox{"iop"},
					},
				},
			},

			&AstBox{"uio"},
			&AstBox{"i"},
		},
	}
	b := astBlock.toFigure(500, 250)
	canvas.SetRGB(0, 0, 0)
	b.draw(canvas)
	b.drawLines(canvas)
	canvas.SavePNG("out.png")
}
