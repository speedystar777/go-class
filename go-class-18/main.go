package main

import (
	"fmt"
	"image/color"
	"math"
)

type Point struct {
	X, Y float64
}

type Line struct {
	Begin, End Point
}

type Path []Point

func (l Line) Distance() float64 {
	return math.Hypot(l.End.X-l.Begin.X, l.End.Y-l.Begin.Y)
}

func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

type ColoredPoint struct { // composition
	Point
	Color color.RGBA
}

// receiver must be a pointer if we don't have return value, so we can directly modify receiver
func (l *Line) ScaleBy(f float64) {
	l.End.X += (f - 1) * (l.End.X - l.Begin.X)
	l.End.Y += (f - 1) * (l.End.Y - l.Begin.Y)
}

func (p Path) Distance() (sum float64) {
	for i := 1; i < len(p); i++ {
		sum += Line{p[i-1], p[i]}.Distance()
	}
	return sum // could also just to naked return with `return` since we have named return value
}

// both Line and Path satisfy Distancer interface
type Distancer interface {
	Distance() float64
}

func PrintDistance(d Distancer) {
	fmt.Println(d.Distance())
}

func main() {
	side := Line{Point{1, 2}, Point{4, 6}}
	perimeter := Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}} // don't need to include Point here since Path knows it is slice of Points

	side.ScaleBy(2)

	PrintDistance(side)
	PrintDistance(perimeter)
	// Code below will not work for 2 reasons
	// 1. ScaleBy can only be called on by Line pointer
	// 2. ScaleBy does not return anything that Distance can be called on
	// fmt.Println(Line{Point{1, 2}, Point{4, 6}}.ScaleBy(2).Distance())

	p, q := Point{1, 1}, ColoredPoint{Point{5, 4}, color.RGBA{255, 0, 0, 255}}
	l1 := q.Distance(p) // automatically uses q.Point
	l2 := p.Distance(q.Point)
	// l2 := p.Distance(q) -> this is not allowed since p.Distance only accepts Point and not ColoredPoint
	fmt.Println(l1, l2)
}
