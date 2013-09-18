package geometry

import (
	"math"
)

type Shape interface {
	Contains(Point) bool
	Area() float64
	Perimeter() float64
}

func init() {

	// require that some structs are Shapes
	var c Circle
	var b Box

	_ = Shape(&c)
	_ = Shape(&b)
}

// ----------

type Point struct {
	x, y float64
}

var ORIGIN = NewPoint(0, 0)

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func distanceBetween(p1, p2 Point) float64 {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Point) DistanceTo(other Point) float64 {
	return distanceBetween(*p, other)
}

func (p *Point) VectorTo(other Point) Vector {
	dx := other.x - p.x
	dy := other.y - p.y
	return Vector{x: dx, y: dy}
}

func (p *Point) Translate(v Vector) Point {
	return Point{x: p.x + v.x, y: p.y + v.y}
}

func (p *Point) SegmentOf(v Vector) Segment {
	return NewSegment(*p, Point{x: p.x + v.x, y: p.y + v.y})
}

// ----------

// same structure as point (and stored as a Point in the db), but used in different situations
type Vector Point

var V_ZERO = NewVector(0, 0)
var V_BASIS_X = NewVector(1, 0)
var V_BASIS_Y = NewVector(0, 1)

func NewVector(x, y float64) Vector {
	return Vector{x: x, y: y}
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Unit() Vector {
	m := v.Magnitude()
	return Vector{x: v.x / m, y: v.y / m}
}

func (v *Vector) Scale(n float64) Vector {
	return Vector{v.x * n, v.y * n}
}

func (v *Vector) Angle() float64 {
	return math.Atan2(v.y, v.x)
}

func (v *Vector) Plus(vs ...Vector) Vector {
	dx := v.x
	dy := v.y

	for _, v2 := range vs {
		dx += v2.x
		dy += v2.y
	}

	return Vector{x: dx, y: dy}
}

// ----------

type Circle struct {
	center Point
	radius float64
}

func NewCircle(center Point, radius float64) Circle {
	return Circle{center: center, radius: radius}
}

func (c *Circle) Contains(p Point) bool {
	return c.center.DistanceTo(p) <= c.radius
}

func (c *Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c *Circle) Box() Box {
	ll := Point{x: c.center.x - c.radius, y: c.center.y - c.radius}
	ur := Point{x: c.center.x + c.radius, y: c.center.y + c.radius}
	return NewBox(ll, ur)
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// ----------

type Segment struct {
	end [2]Point
}

func NewSegment(p1, p2 Point) Segment {
	end := new([2]Point)
	end[0] = p1
	end[1] = p2
	return Segment{end: *end}
}

func (s *Segment) Magnitude() float64 {
	return s.end[0].DistanceTo(s.end[1])
}

func (s *Segment) Box() Box {
	return NewBox(s.end[0], s.end[1])
}

func (s *Segment) Flip() Segment {
	return NewSegment(s.end[1], s.end[0])
}

// ----------

type Box struct {
	corner [2]Point // 0 is always lower left, 1 is always upper right
}

func NewBox(p1, p2 Point) Box {
	corner := new([2]Point)

	// normalize
	xmin, xmax := sortPair(p1.x, p2.x)
	ymin, ymax := sortPair(p1.y, p2.y)

	corner[0] = Point{x: xmin, y: ymin}
	corner[1] = Point{x: xmax, y: ymax}

	return Box{corner: *corner}
}

func (b *Box) Area() float64 {
	// normalization ensures these are positive
	dx := b.corner[1].x - b.corner[0].x
	dy := b.corner[1].y - b.corner[0].y
	return dx * dy
}

func (b *Box) Perimeter() float64 {
	// normalization ensures these are positive
	dx := b.corner[1].x - b.corner[0].x
	dy := b.corner[1].y - b.corner[0].y
	return 2*dx + 2*dy
}

func (b *Box) Contains(p Point) bool {

	// normalization ensures no need to check relative positions
	if p.x > b.corner[1].x {
		return false
	}

	if p.x < b.corner[0].x {
		return false
	}

	if p.y > b.corner[1].y {
		return false
	}

	if p.y < b.corner[0].y {
		return false
	}

	return true
}

// ----------

// TODO: Paths and Polygons

type Path struct {
	point  []Point
	closed bool
}

type Polygon Path

/**
 * Find the time that point 1 is within radius of point 2
 *
 * @param s1
 *            Starting location of point 1, units = d
 * @param s2
 *            Starting location of point 2, units = d
 * @param v1
 *            Speed vector of point 1, units = d/t
 * @param v2
 *            Speed vector of point 2, units = d/t
 * @param r
 *            Range of point 1, units = d
 * @return The time point 1 in within range of point 2, units = t. Returns NaN if never in range.
 */
func TimeIntercept(s1, s2 Point, v1, v2 Vector, radius float64) (float64, float64) {

	// locations, vector form
	// l1 = s1 + (t * v1)
	// l2 = s2 + (t * v2)

	// deltas in each dimension
	// dx = l2.x - l1.x
	// dy = l2.y - l1.y

	// range (distance equation)
	// r = sqrt (dx^2 + dy^2)
	// r^2 = dx^2 + dy^2
	// r^2 = (l2.x - l1.x)^2 + (l2.y - l1.y)^2
	// r^2 = ((s2.x + (t * v2.x)) - (s1.x + (t * v1.x)))^2 + ((s2.y + (t * v2.y)) - (s1.y + (t * v1.y)))^2

	// multiply that mess out, rearrange to the form
	// r^2 = (some big mess)t^2 + (another big mess)t + (a constant mess)
	// which is basically quadratic format
	// r^2 = at^2 +bt + c

	// coeffecients
	a := v2.y*v2.y - 2*v1.y*v2.y + v2.x*v2.x - 2*v1.x*v2.x + v1.y*v1.y + v1.x*v1.x
	b := (2*s2.y-2*s1.y)*v2.y + (2*s2.x-2*s1.x)*v2.x + (2*s1.y-2*s2.y)*v1.y + (2*s1.x-2*s2.x)*v1.x
	c := s2.y*s2.y - 2*s1.y*s2.y + s2.x*s2.x - 2*s1.x*s2.x + s1.y*s1.y + s1.x*s1.x

	// arrange that so that r^2 is just another part of the constant
	// at^2 +bt + c - r^2 = 0

	// then solve for t using known quadratic equation solution
	// t = (-b +/- sqrt(b^2 - 4a(c-r^2))) / 2a

	// but break into steps
	inner := b*b - 4*a*(c-radius*radius)

	// will the discriminant be imaginary?
	if inner < 0 {
		// there are no solutions; p1 and p2 are never within r of each other
		return math.NaN(), math.NaN()
	}

	// if inner == 0, there is one solution, but for floating point that is unlikely to be so exact
	// in any case, returning 2 equal results is still correct

	discr := math.Sqrt(inner)

	sol1 := (-b + discr) / (2 * a)
	sol2 := (-b - discr) / (2 * a)

	return sortPair(sol1, sol2)
}

func sortPair(n, m float64) (float64, float64) {
	if n <= m {
		return n, m
	}

	return m, n
}
