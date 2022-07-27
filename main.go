package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

const unit int = 100

type Box struct {
	x, y          int
	width, height int
}

func (box *Box) draw(canvas *svg.SVG) {
	x := box.x - box.width/2
	canvas.Rect(x, box.y, box.width, box.height)
}

func (box *Box) position(x, y int) {
	box.x += x
	box.y += y
}

func (block *Block) position(x, y int) {
	for _, child := range block.children {
		child.position(x, y)
	}
}

type Figure interface {
	draw(*svg.SVG)
	top() (int, int)
	size() (int, int)
	position(x, y int)
}

type Block struct {
	children []Figure
}

func (block *Block) top() (int, int) {
	//return block.x, block.y
	panic("TODO")
}

func (block *Block) size() (int, int) {
	//return block.width, block.height
	panic("TODO")
}
func (block *Block) draw(canvas *svg.SVG) {
	for _, child := range block.children {
		child.draw(canvas)
	}
}

func (box *Box) top() (int, int) {
	return box.x, box.y
}

func (box *Box) size() (int, int) {
	return box.width, box.height
}

func main() {
	width := 1000
	height := 1000
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	//canvas.Circle(width/2, height/2, 100)
	//canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	box := Box{
		x:      0,
		y:      0,
		width:  200,
		height: 100,
	}
	box2 := Box{
		x:      0,
		y:      300,
		width:  300,
		height: 209,
	}
	block := Block{
		children: []Figure{&box, &box2},
	}
	block.position(250, 250)
	block.draw(canvas)
	canvas.End()
}
