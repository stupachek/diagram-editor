package figure

import (
	"bytes"

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

type StartStop struct {
	x, y          int
	width, height int
	text          string
	radius        float64
}

type Rhombus struct {
	x, y, width, height int
	text                string
}

type Input struct {
	x, y, width, height int
	text                string
}

func (i *Input) draw(canvas *gg.Context) {
	x := i.x - i.width/2 + horizontalMargins/2
	y := i.height + i.y
	canvas.MoveTo(float64(x), float64(i.y))
	x2 := i.x + i.width/2
	canvas.LineTo(float64(x2), float64(i.y))
	x3 := x2 - horizontalMargins/2
	canvas.LineTo(float64(x3), float64(y))
	x4 := x - horizontalMargins/2
	canvas.LineTo(float64(x4), float64(y))
	canvas.LineTo(float64(x), float64(i.y))
	canvas.Stroke()
	middleY := i.height/2 + i.y
	canvas.DrawStringAnchored(i.text, float64(i.x), float64(middleY), 0.5, 0.35)
}

func (b *Box) draw(canvas *gg.Context) {
	x := b.x - b.width/2
	canvas.DrawRectangle(float64(x), float64(b.y), float64(b.width), float64(b.height))
	canvas.Stroke()
	middleY := b.height/2 + b.y
	canvas.DrawStringAnchored(b.text, float64(b.x), float64(middleY), 0.5, 0.35)
}

func (i *Input) position(x, y int) {
	i.x += x
	i.y += y
}

func (i *Input) connectTo(x, y int, canvas *gg.Context) {
	bottom := i.y + i.height
	canvas.DrawLine(float64(i.x), float64(bottom), float64(x), float64(y))
	canvas.Stroke()
}

func (i *Input) drawLines(canvas *gg.Context) {}

func (b *Box) position(x, y int) {
	b.x += x
	b.y += y
}

func (b *Box) connectTo(x, y int, canvas *gg.Context) {
	bottom := b.y + b.height
	canvas.DrawLine(float64(b.x), float64(bottom), float64(x), float64(y))
	canvas.Stroke()
}
func (b *Box) drawLines(canvas *gg.Context) {}

func (s *StartStop) connectTo(x, y int, canvas *gg.Context) {
	bottom := s.y + s.height
	canvas.DrawLine(float64(s.x), float64(bottom), float64(x), float64(y))
	canvas.Stroke()
}

func (s *StartStop) drawLines(canvas *gg.Context) {}

func (s *StartStop) position(x, y int) {
	s.x += x
	s.y += y
}

func (r *Rhombus) draw(canvas *gg.Context) {
	bottomY := r.y + r.height
	leftX, middleY := r.left()
	rightX, _ := r.right()
	canvas.MoveTo(float64(r.x), float64(r.y))
	canvas.LineTo(float64(rightX), float64(middleY))
	canvas.LineTo(float64(r.x), float64(bottomY))
	canvas.LineTo(float64(leftX), float64(middleY))
	canvas.LineTo(float64(r.x), float64(r.y))
	canvas.Stroke()
	middleY_2 := r.height/2 + r.y
	canvas.DrawStringAnchored(r.text, float64(r.x), float64(middleY_2), 0.5, 0.35)
}

func (s *StartStop) draw(canvas *gg.Context) {
	x := s.x - s.width/2
	canvas.DrawRoundedRectangle(float64(x), float64(s.y), float64(s.width), float64(s.height), s.radius)
	canvas.Stroke()
	middleY := s.height/2 + s.y
	canvas.DrawStringAnchored(s.text, float64(s.x), float64(middleY), 0.5, 0.35)
}

func (r *Rhombus) left() (x, y int) {
	return r.x - r.width/2, r.y + r.height/2
}
func (r *Rhombus) right() (x, y int) {
	return r.x + r.width/2, r.y + r.height/2
}

func (r *Rhombus) position(x, y int) {
	r.x += x
	r.y += y
}
func (r *Rhombus) top() (int, int) {
	return r.x, r.y
}

func (r *Rhombus) size() (int, int) {
	return r.width, r.height
}

func (b *Block) position(x, y int) {
	for _, child := range b.children {
		child.position(x, y)
	}
	b.x += x
	b.y += y
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

func (b *Block) top() (int, int) {
	return b.x, b.y
}

func (b *Block) size() (int, int) {
	width, height := 0, 0
	for _, child := range b.children {
		w, h := child.size()
		height += h + blockSpacing
		if w > width {
			width = w
		}
	}
	if len(b.children) != 0 {
		height -= blockSpacing
	}
	return width, height
}
func (b *Block) draw(canvas *gg.Context) {
	for _, child := range b.children {
		child.draw(canvas)
	}
}

func (b *Block) bottom() (int, int) {
	topX, topY := b.top()
	_, h := b.size()
	return topX, topY + h
}

func (b *Block) drawLines(canvas *gg.Context) {
	for _, child := range b.children {
		child.drawLines(canvas)
	}
	for i := 0; i < len(b.children)-1; i++ {
		x, y := b.children[i+1].top()
		b.children[i].connectTo(x, y, canvas)
		Arrow(x, y, canvas)
	}
}

func (i *Input) top() (int, int) {
	return i.x, i.y
}

func (i *Input) size() (int, int) {
	return i.width, i.height
}
func (b *Box) top() (int, int) {
	return b.x, b.y
}

func (b *Box) size() (int, int) {
	return b.width, b.height
}

func (s *StartStop) top() (int, int) {
	return s.x, s.y
}

func (s *StartStop) size() (int, int) {
	return s.width, s.height
}

type If struct {
	cond        Rhombus
	left, right Block
}

func (i *If) draw(canvas *gg.Context) {
	i.left.draw(canvas)
	i.right.draw(canvas)
	i.cond.draw(canvas)

}
func (i *If) top() (int, int) {
	return i.cond.top()
}

func (i *If) size() (int, int) {
	_, rhombusHeigth := i.cond.size()
	leftWidth, leftHeigth := i.left.size()
	rightWidth, rightHeigth := i.right.size()
	height := rhombusHeigth + blockSpacing
	if leftHeigth > rightHeigth {
		height += leftHeigth
	} else {
		height += rightHeigth
	}
	leftX, _ := i.left.top()
	rightX, _ := i.right.top()
	width := Abs(leftX-rightX) + leftWidth/2 + rightWidth/2
	return width, height

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (i *If) drawLines(canvas *gg.Context) {
	ifRX, ifRY := i.cond.right()
	ifLX, ifLY := i.cond.left()
	leftX, leftY := i.left.top()
	rightX, rightY := i.right.top()
	canvas.DrawLine(float64(ifRX), float64(ifRY), float64(rightX), float64(ifRY))
	canvas.Stroke()
	canvas.DrawLine(float64(ifLX), float64(ifLY), float64(leftX), float64(ifLY))
	canvas.Stroke()
	canvas.DrawLine(float64(rightX), float64(ifRY), float64(rightX), float64(rightY))
	canvas.Stroke()
	canvas.DrawLine(float64(leftX), float64(ifLY), float64(leftX), float64(leftY))
	canvas.Stroke()
	i.left.drawLines(canvas)
	if !i.left.IsEmpty() {
		Arrow(leftX, leftY, canvas)
	}
	if !i.right.IsEmpty() {
		Arrow(rightX, rightY, canvas)
	}
	i.right.drawLines(canvas)
	canvas.DrawStringAnchored("1", float64(ifLX-horizontalMargins/2), float64(ifLY-verticalMargins/2), 1, 0)
	canvas.DrawStringAnchored("0", float64(ifRX+horizontalMargins/2), float64(ifRY-verticalMargins/2), 0, 0)

}

func (b *Block) IsEmpty() bool {
	return len(b.children) == 0
}

func (i *If) position(x int, y int) {
	i.cond.position(x, y)
	i.left.position(x, y)
	i.right.position(x, y)
}

func (i *If) connectTo(x, y int, canvas *gg.Context) {
	leftX, leftY := i.left.bottom()
	rightX, rightY := i.right.bottom()
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
	ifX, _ := i.top()
	canvas.DrawLine(float64(ifX), float64(middleY), float64(x), float64(y))
	canvas.Stroke()
}

func newInput(text string, x, y int) Input {
	w, h := TextSize(text)
	return Input{
		x:      x,
		y:      y,
		width:  w + horizontalMargins,
		height: h + verticalMargins,
		text:   text,
	}
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

func newStartStop(text string, x, y int) StartStop {
	w, h := TextSize(text)
	r := float64((h + horizontalMargins) / 2)
	return StartStop{
		x:      x,
		y:      y,
		width:  w + horizontalMargins,
		height: h + verticalMargins,
		text:   text,
		radius: r,
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

func init() {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	face = truetype.NewFace(font, &truetype.Options{Size: textHeight})
}

func DrawBlock(block AstBlock) []byte {
	blockFigure := block.toFigure(0, 0)
	w, h := blockFigure.size()
	canvas := gg.NewContext(2*w, h)
	canvas.SetFontFace(face)
	blockFigure.position(w, 0)
	canvas.SetRGB(0, 0, 0)
	blockFigure.draw(canvas)
	blockFigure.drawLines(canvas)
	buf := bytes.NewBuffer(nil)
	canvas.EncodePNG(buf)
	return buf.Bytes()
}
