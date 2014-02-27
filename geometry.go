package geometry

import (
	"math"
)

// A Shape is an enclosed 2D area.
// Circles, Boxes, and Polygons are all shapes.
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

// Point is a point on the 2D plane.
// It is represented in the postgres database by the <point> type.
// Points are immutable.  Use Value() to inspect contents.
type Point struct {
	x, y float64
}

// The Origin is the zero point on the 2D plane (0,0)
var Origin = NewPoint(0, 0)

// NewPoint returns a new point at the given coordinates.
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func distanceBetween(p1, p2 Point) float64 {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceTo returns the euclidean distance between this point and another.
func (p Point) DistanceTo(other Point) float64 {
	return distanceBetween(p, other)
}

// VectorTo computes a vector pointing from this point to anoter point.
func (p Point) VectorTo(other Point) Vector {
	dx := other.x - p.x
	dy := other.y - p.y
	return Vector{x: dx, y: dy}
}

// Translate returns a new point translated by the given vector.
func (p Point) Translate(v Vector) Point {
	return Point{x: p.x + v.x, y: p.y + v.y}
}

// AsSegment returns a segment from this point to a point as translated by the
// given vector.
func (p Point) SegmentTo(v Vector) Segment {
	return NewSegment(p, Point{x: p.x + v.x, y: p.y + v.y})
}

// Value returns the coordinate values of the point.
func (p Point) Values() (x, y float64) {
	return p.x, p.y
}

// ----------

// Vector is a vector on the 2D plane.
// It is represented in the postgres database by the <point> type.
// Vectors are immutable.  Use Value() to inspect contents.
type Vector Point

// The Zero vector. [0,0]
var ZeroVector = NewVector(0, 0)

// The X basis vector [1,0]
var BasisX = NewVector(1, 0)

// The Y basis vector [0,1]
var BasisY = NewVector(0, 1)

// NewVector creates a new vector in the 2D plane.
func NewVector(x, y float64) Vector {
	return Vector{x: x, y: y}
}

// Magnitude returns the magnitude (length) of this vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

// Unit returns a new vector in the same direction as this vector, but with
// a magnitude of 1.
func (v Vector) Unit() Vector {
	m := v.Magnitude()
	return Vector{x: v.x / m, y: v.y / m}
}

// Scale returns a new vector scaled by the multiplier n.
func (v Vector) Scale(n float64) Vector {
	return Vector{v.x * n, v.y * n}
}

// Angle computes the angle of the vector
func (v Vector) Angle() float64 {
	return math.Atan2(v.y, v.x)
}

// Plus returns a new vector whose value is the sum of this vector, plus
// all the vectors passed as parameters.
func (v Vector) Plus(vs ...Vector) Vector {
	dx := v.x
	dy := v.y

	for _, v2 := range vs {
		dx += v2.x
		dy += v2.y
	}

	return Vector{x: dx, y: dy}
}

// Dot returns the dot product of the vector with another vector
func (v Vector) Dot(v2 Vector) float64 {
	return v.x*v2.x + v.y*v2.y
}

// CrossZ returns the Z component of the cross product of this vector with
// another vector.  Because the inputs are 2D vectors in the X/Y plane, and
// their implied Z components are zero, their cross product is colinear with
// the Z axis (X and Y components being zero), and not representable by a
// Vector in this package.  This method returns only the Z component, whose
// direction can be determined by its sign.
func (v Vector) CrossZ(v2 Vector) float64 {
	return v.x*v2.y - v.y*v2.x
}

// AsSegment returns a segment from the origin to the vector endpoint.
// To obtain a segment from a point other than the origin, use
// Point.AsSegment()
func (v Vector) AsSegment() Segment {
	return NewSegment(Origin, Point(v))
}

// Value returns the coordinate values of the vector.
func (v Vector) Values() (x, y float64) {
	return v.x, v.y
}

// ----------

// Circle is a circle on the 2D plane.
// It is represented in the postgres database by the <circle> type.
// Implements the Shape interface.
type Circle struct {
	center Point
	radius float64
}

// NewCircle returns a circle centered on the given point, and with the
// given radius.
func NewCircle(center Point, radius float64) Circle {
	return Circle{center: center, radius: radius}
}

// Contains returns true if the point is on or inside the circle.
func (c Circle) Contains(p Point) bool {
	return c.center.DistanceTo(p) <= c.radius
}

// Area returns the area of the circle.
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// Box returns a box which exactly encloses the circle.
// The box is tangent to the circle at its sides' midpoints, and shares
// a center with the circle.
func (c Circle) Box() Box {
	ll := Point{x: c.center.x - c.radius, y: c.center.y - c.radius}
	ur := Point{x: c.center.x + c.radius, y: c.center.y + c.radius}
	return NewBox(ll, ur)
}

// Perimeter returns the perimeter of the circle
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// Value returns the coordinate values of the circle.
func (c Circle) Values() (x, y, radius float64) {
	return c.center.x, c.center.y, c.radius
}

// ----------

// A Segment is a line segment defined by two points.
// It is represented in the postgres database by the <lseg> type.
type Segment [2]Point

// NewSegment returns a segment connecting the given points.
func NewSegment(p1, p2 Point) Segment {
	return Segment([2]Point{p1, p2})
}

// Magnitude returns the length of the segment.
func (s Segment) Magnitude() float64 {
	return s[0].DistanceTo(s[1])
}

// Box returns the segment as a box whose corners are the segment's endpoints.
func (s Segment) AsBox() Box {
	return NewBox(s[0], s[1])
}

// Flip returns a new segment whose endpoints are exchanged
func (s Segment) Flip() Segment {
	return NewSegment(s[1], s[0])
}

// ----------

// A Box is a rectangle on the 2D plane.
// A box's size and location is determined by two diagonally opposite corners.
// The normal form of a box is for the first corner to be at the upper right
// (coordinates with the largest X and Y values), and the second corner to be
// the at the lower left (as defined by a standard X/Y plane with values
// decreasing to the left and down).  Using the NewBox function guarantees
// a normal box.
type Box [2]Point

func NewBox(p1, p2 Point) Box {
	b := new(Box)

	// normalize
	xmin, xmax := sortPair(p1.x, p2.x)
	ymin, ymax := sortPair(p1.y, p2.y)

	b[0] = Point{x: xmax, y: ymax}
	b[1] = Point{x: xmin, y: ymin}

	return *b
}

// Area calculates the area of the box.
func (b Box) Area() float64 {
	// normalization ensures these are positive
	dx := b[0].x - b[1].x
	dy := b[0].y - b[1].y
	return dx * dy
}

// Perimeter calculates the perimeter of the box.
func (b Box) Perimeter() float64 {
	// normalization ensures these are positive
	dx := b[0].x - b[1].x
	dy := b[0].y - b[1].y
	return 2*dx + 2*dy
}

// Contains returns whether the given point is on or inside the box.
func (b Box) Contains(p Point) bool {

	// normalization ensures no need to check relative positions
	if p.x > b[0].x {
		return false
	}

	if p.x < b[1].x {
		return false
	}

	if p.y > b[0].y {
		return false
	}

	if p.y < b[1].y {
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

// TimeIntercept computes the interception time of two moving points.
// Two points, s1 and s2, which have respective velocities of v1 and v2, may
// intercept at two times, returned by this function.  Interception is defined
// as being with 'radius' distance of each other.  Go get a true intersection,
// use a radius of zero.  Interception times may be in the past or future,
// depending on starting locations and velocities.
// If the velocities are parallel, two degenerate cases may arise.  Points may
// never be in range, in which case NaN is returned for both time, or always
// in range, in which case -Inf and +Inf are returned.
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

	// if a==0 the end result will be a divide by zero, so just bail now
	if a == 0 {
		// Geometric meaning: The points have the same velocities, so
		// they will either ALWAYS be within r of each other, or NEVER
		if s1.DistanceTo(s2) <= radius {
			// All times are solutions
			// The points have been, and always will be, <= r from each other
			return math.Inf(-1), math.Inf(1)
		}

		// No times are solution
		return math.NaN(), math.NaN()
	}

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
