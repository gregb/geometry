package geometry

import (
	"launchpad.net/gocheck"
	"math"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct{}

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestPoint(c *gocheck.C) {

	p2 := NewPoint(3, 4)
	p3 := NewPoint(3, 4)
	p4 := NewPoint(4, 5)
	v1 := NewVector(3, 4)
	v2 := NewVector(-3, -4)
	v3 := NewVector(1, 1)

	c.Check(p2 == p3, gocheck.Equals, true)
	c.Check(ORIGIN.DistanceTo(p2), gocheck.Equals, float64(5))

	c.Check(ORIGIN, gocheck.Not(gocheck.Equals), V_ZERO)
	c.Check(ORIGIN.Translate(V_ZERO), gocheck.Equals, ORIGIN)
	c.Check(p2.Translate(v3), gocheck.Equals, p4)

	c.Check(ORIGIN.VectorTo(p2), gocheck.Equals, v1)
	c.Check(p2.VectorTo(ORIGIN), gocheck.Equals, v2)
}

func (s *Suite) TestVector(c *gocheck.C) {

	v1 := NewVector(3, 4)
	v2 := NewVector(-3, -4)
	v3 := NewVector(1, 1)

	c.Check(v1.Scale(-1), gocheck.Equals, v2)
	c.Check(V_BASIS_Y.Scale(1), gocheck.Equals, V_BASIS_Y)

	c.Check(V_ZERO.Magnitude(), gocheck.Equals, float64(0))
	c.Check(v1.Magnitude(), gocheck.Equals, float64(5))
	c.Check(v2.Magnitude(), gocheck.Equals, float64(5))
	c.Check(v3.Magnitude(), gocheck.Equals, math.Sqrt(2))

	c.Check(V_BASIS_X.Unit(), gocheck.Equals, V_BASIS_X)
	dumb := v1.Unit()
	c.Check(dumb.Magnitude(), gocheck.Equals, float64(1))

	c.Check(v1.Plus(v2), gocheck.Equals, V_ZERO)

	c.Check(V_ZERO.Plus(v3, V_BASIS_X, V_BASIS_Y, V_BASIS_X, V_BASIS_Y, V_BASIS_Y), gocheck.Equals, v1)

	c.Check(V_BASIS_X.Angle(), gocheck.Equals, float64(0))
	c.Check(V_BASIS_Y.Angle(), gocheck.Equals, math.Pi/2)
	c.Check(v3.Angle(), gocheck.Equals, math.Pi/4)
}

func (s *Suite) TestSegment(c *gocheck.C) {

	p1 := NewPoint(-3, 4)
	p2 := NewPoint(-4, 3)
	v1 := NewVector(-1, -1)
	s1 := NewSegment(ORIGIN, p1)
	s2 := NewSegment(p1, p2)
	s3 := NewSegment(p2, p1)
	s4 := p1.SegmentOf(v1)
	b1 := NewBox(ORIGIN, p1)

	c.Check(s1.Magnitude(), gocheck.Equals, float64(5))
	c.Check(s2.Magnitude(), gocheck.Equals, math.Sqrt(2))

	c.Check(s3.Flip(), gocheck.Equals, s2)
	c.Check(s4, gocheck.Equals, s2)
	c.Check(s1.Box(), gocheck.Equals, b1)

}

func (s *Suite) TestCircle(c *gocheck.C) {

	p1 := NewPoint(-3, 4)
	p2 := NewPoint(-4, 3)
	p3 := NewPoint(-1, -1)
	p4 := NewPoint(-8, -1)
	p5 := NewPoint(2, 9)
	b1 := NewBox(p4, p5)
	c1 := NewCircle(p1, 5)

	c.Check(c1.Contains(ORIGIN), gocheck.Equals, true)
	c.Check(c1.Contains(p1), gocheck.Equals, true)
	c.Check(c1.Contains(p2), gocheck.Equals, true)
	c.Check(c1.Contains(p3), gocheck.Equals, false)

	c.Check(c1.Area(), gocheck.Equals, math.Pi*25)
	c.Check(c1.Perimeter(), gocheck.Equals, math.Pi*10)
	c.Check(c1.Box(), gocheck.Equals, b1)
}

func (s *Suite) TestBox(c *gocheck.C) {
	p1 := NewPoint(-2, 2)
	p2 := NewPoint(2, 2)
	p3 := NewPoint(2, -2)
	p4 := NewPoint(-2, -2)
	p5 := NewPoint(-3, 1)
	p6 := NewPoint(3, 1)
	p7 := NewPoint(3, -1)
	p8 := NewPoint(-3, -1)

	b1 := NewBox(p1, p3)
	b2 := NewBox(p3, p1)
	b3 := NewBox(p2, p4)
	b4 := NewBox(p4, p2) // canonical

	c.Check(b1, gocheck.Equals, b4)
	c.Check(b2, gocheck.Equals, b4)
	c.Check(b3, gocheck.Equals, b4)
	c.Check(b4, gocheck.Equals, b4)

	c.Check(b4.Area(), gocheck.Equals, float64(16))
	c.Check(b4.Perimeter(), gocheck.Equals, float64(16))

	c.Check(b4.Contains(p1), gocheck.Equals, true)
	c.Check(b4.Contains(p2), gocheck.Equals, true)
	c.Check(b4.Contains(p3), gocheck.Equals, true)
	c.Check(b4.Contains(p4), gocheck.Equals, true)
	c.Check(b4.Contains(ORIGIN), gocheck.Equals, true)

	c.Check(b4.Contains(p5), gocheck.Equals, false)
	c.Check(b4.Contains(p6), gocheck.Equals, false)
	c.Check(b4.Contains(p7), gocheck.Equals, false)
	c.Check(b4.Contains(p8), gocheck.Equals, false)
}

func checkIntercept(s1, s2 Point, v1, v2 Vector, radius float64) bool {

	// get first intercept time
	t1, _ := TimeIntercept(s1, s2, v1, v2, radius)

	// calc locations of each point at t1
	i1 := s1.Translate(v1.Scale(t1))
	i2 := s2.Translate(v2.Scale(t1))

	// how far apart are they?
	calcR := i1.DistanceTo(i2)

	// should be equal to radius
	// but allow for floating point weirdness
	if math.Abs(calcR-radius) < 0.00001 {
		return true
	}

	return false
}
