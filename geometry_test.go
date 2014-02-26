package geometry

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestPoint(t *testing.T) {

	Convey("Given a set of points and vectors", t, func() {
		p2 := NewPoint(3, 4)
		p3 := NewPoint(3, 4)
		p4 := NewPoint(4, 5)
		v1 := NewVector(3, 4)
		v2 := NewVector(-3, -4)
		v3 := NewVector(1, 1)

		Convey("Two points at the same place should equal each other", func() {
			So(p2, ShouldResemble, p3)
		})

		Convey("Distance should calculate correctly", func() {
			So(Origin.DistanceTo(p2), ShouldEqual, float64(5))
		})

		Convey("A point and vector should not be equal, even if their components match", func() {
			So(Origin, ShouldNotEqual, ZeroVector)
		})

		Convey("Translated points should be where expected", func() {
			So(Origin.Translate(ZeroVector), ShouldResemble, Origin)
			So(p2.Translate(v3), ShouldResemble, p4)
		})

		Convey("Vectors calculaed between points should be as expected", func() {
			So(Origin.VectorTo(p2), ShouldResemble, v1)
			So(p2.VectorTo(Origin), ShouldResemble, v2)
		})

		Convey("Values should extract correctly", func() {
			x, y := p4.Values()
			So(x, ShouldEqual, 4)
			So(y, ShouldEqual, 5)
		})
	})
}

func TestVector(t *testing.T) {

	Convey("Given a set of vectors", t, func() {
		v1 := NewVector(3, 4)
		v2 := NewVector(-3, -4)
		v3 := NewVector(1, 1)
		v4 := NewVector(2, 3)

		Convey("Scaling should produce the correct result", func() {
			So(v1.Scale(-1), ShouldResemble, v2)
			So(BasisY.Scale(1), ShouldResemble, BasisY)
		})

		Convey("Magnitude should be calculated correctly", func() {
			So(ZeroVector.Magnitude(), ShouldEqual, float64(0))
			So(v1.Magnitude(), ShouldEqual, float64(5))
			So(v2.Magnitude(), ShouldEqual, float64(5))
			So(v3.Magnitude(), ShouldEqual, math.Sqrt(2))
		})

		Convey("Angle should be calculated correctly", func() {
			So(BasisX.Angle(), ShouldEqual, float64(0))
			So(BasisY.Angle(), ShouldEqual, math.Pi/2)
			So(v3.Angle(), ShouldEqual, math.Pi/4)
		})

		Convey("Calculated unit vectors should have magnitude 1 and the same angle as the original", func() {
			So(BasisX.Unit(), ShouldResemble, BasisX)
			dumb := v1.Unit()
			So(dumb.Magnitude(), ShouldEqual, float64(1))
			So(dumb.Angle(), ShouldAlmostEqual, v1.Angle())
		})

		Convey("Addition should produce the correct result", func() {
			So(v1.Plus(v2), ShouldResemble, ZeroVector)
			So(ZeroVector.Plus(v3, BasisX, BasisY, BasisX, BasisY, BasisY), ShouldResemble, v1)
		})

		Convey("Dot product should be calculated correctly", func() {
			So(v1.Dot(v2), ShouldEqual, (-3*3)+(-4*4)) // -25
			So(v1.Dot(v3), ShouldEqual, 3+4)           // 7
			So(v1.Dot(v4), ShouldEqual, (3*2)+(4*3))   // 18
		})

		Convey("Cross product should be calculated correctly", func() {
			So(v1.CrossZ(v2), ShouldEqual, (3*-4)-(4*-3)) // 0
			So(v1.CrossZ(v3), ShouldEqual, 3*1-4*1)       // -1
			So(v1.CrossZ(v4), ShouldEqual, (3*3)-(4*2))   // 1
		})

		Convey("As as segment should produce the expeced segment", func() {
			s := NewSegment(Origin, Point{2, 3})
			So(v4.AsSegment(), ShouldResemble, s)

		})

		Convey("Values should extract correctly", func() {
			x, y := v1.Values()
			So(x, ShouldEqual, 3)
			So(y, ShouldEqual, 4)
		})
	})
}

func TestSegment(t *testing.T) {

	Convey("Given a set of segments", t, func() {

		p1 := NewPoint(-3, 4)
		p2 := NewPoint(-4, 3)
		v1 := NewVector(-1, -1)
		s1 := NewSegment(Origin, p1)
		s2 := NewSegment(p1, p2)
		s3 := NewSegment(p2, p1)
		s4 := p1.SegmentTo(v1)
		b1 := NewBox(Origin, p1)

		So(s1.Magnitude(), ShouldEqual, float64(5))
		So(s2.Magnitude(), ShouldEqual, math.Sqrt(2))

		So(s3.Flip(), ShouldResemble, s2)
		So(s4, ShouldResemble, s2)
		So(s1.AsBox(), ShouldResemble, b1)
	})
}

func TestCircle(t *testing.T) {

	Convey("Given a circle", t, func() {

		p1 := NewPoint(-3, 4)
		p2 := NewPoint(-4, 3)
		p3 := NewPoint(-1, -1)
		p4 := NewPoint(-8, -1)
		p5 := NewPoint(2, 9)
		b1 := NewBox(p4, p5)
		c1 := NewCircle(p1, 5)

		Convey("It should contain the points expected", func() {
			So(c1.Contains(Origin), ShouldEqual, true)
			So(c1.Contains(p1), ShouldEqual, true)
			So(c1.Contains(p2), ShouldEqual, true)
			So(c1.Contains(p3), ShouldEqual, false)
		})

		Convey("Perimeter and area should be calculated correctly", func() {
			So(c1.Area(), ShouldEqual, math.Pi*25)
			So(c1.Perimeter(), ShouldEqual, math.Pi*10)
		})

		Convey("Its bounding box should be correct", func() {
			So(c1.Box(), ShouldResemble, b1)
		})

		Convey("Values should extract correctly", func() {
			x, y, r := c1.Values()
			So(x, ShouldEqual, -3)
			So(y, ShouldEqual, 4)
			So(r, ShouldEqual, 5)
		})
	})
}

func TestBox(t *testing.T) {

	Convey("Given two sets of identical boxes", t, func() {
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

		Convey("Similar boxes should be constructed canonically", func() {
			So(b1, ShouldResemble, b4)
			So(b2, ShouldResemble, b4)
			So(b3, ShouldResemble, b4)
			So(b4, ShouldResemble, b4)
		})

		Convey("Perimeter and area should be calculated correctly", func() {
			So(b4.Area(), ShouldEqual, float64(16))
			So(b4.Perimeter(), ShouldEqual, float64(16))
		})

		Convey("It should contain the points expected", func() {
			So(b4.Contains(p1), ShouldEqual, true)
			So(b4.Contains(p2), ShouldEqual, true)
			So(b4.Contains(p3), ShouldEqual, true)
			So(b4.Contains(p4), ShouldEqual, true)
			So(b4.Contains(Origin), ShouldEqual, true)

			So(b4.Contains(p5), ShouldEqual, false)
			So(b4.Contains(p6), ShouldEqual, false)
			So(b4.Contains(p7), ShouldEqual, false)
			So(b4.Contains(p8), ShouldEqual, false)
		})
	})
}

func TestIntercept(t *testing.T) {

	Convey("Given points and velocity vectors", t, func() {

		Convey("Converging paths should intercept correctly", func() {
			s1 := Point{0, 0}
			s2 := Point{8, 0}
			v1 := Vector{1, 1}
			v2 := Vector{-1, 1}

			t1a, t2a := TimeIntercept(s1, s2, v1, v2, 2)
			So(t1a, ShouldAlmostEqual, 3)
			So(t2a, ShouldAlmostEqual, 5)
		})

		Convey("Diverging paths should intercept correctly", func() {
			s1 := Point{0, 0}
			s2 := Point{8, 0}
			v1 := Vector{1, 1}
			v2 := Vector{-1, 1}

			// Starting at the same point, diverging and out of reach eventually
			// Theoretically they intercepted in the past
			t1b, t2b := TimeIntercept(s1, s1, v2, v1, 2)
			So(t1b, ShouldAlmostEqual, -1)
			So(t2b, ShouldAlmostEqual, 1)

			// Starting at different points, and diverging.
			// Theoretically they intercepted in the past,
			// and already diverged in the past
			// and will not come in range in the future
			t1c, t2c := TimeIntercept(s1, s2, v2, v1, 2)
			So(t1c, ShouldAlmostEqual, -5)
			So(t2c, ShouldAlmostEqual, -3)
		})

		Convey("Parallel paths should intercept correctly", func() {
			s00 := Point{0, 0}
			s80 := Point{8, 0}
			s_2_4 := Point{-2, -4}

			s10 := Point{1, 0}

			v11 := Vector{1, 1}
			v22 := Vector{2, 2}

			// Starting out of range, same velocity
			// Always out of range (in the past and future)
			t1d, t2d := TimeIntercept(s00, s80, v11, v11, 2)
			So(math.IsNaN(t1d), ShouldBeTrue)
			So(math.IsNaN(t2d), ShouldBeTrue)

			// Starting out of range, different velocity (but still parallel)
			// Always out of range (in the past and future)
			// [Slightly different code path than previous test]
			t1e, t2e := TimeIntercept(s00, s80, v11, v22, 2)
			So(math.IsNaN(t1e), ShouldBeTrue)
			So(math.IsNaN(t2e), ShouldBeTrue)

			// Starting in range, same velocity
			// Never out of range (past or future)
			t1f, t2f := TimeIntercept(s00, s10, v11, v11, 2)
			So(math.IsInf(t1f, -1), ShouldBeTrue)
			So(math.IsInf(t2f, 1), ShouldBeTrue)

			// Starting out of range, different velocity (but still parallel)
			// Catches up, then outpaces
			t1g, t2g := TimeIntercept(s00, s_2_4, v11, v22, 2)
			So(t1g, ShouldAlmostEqual, 2)
			So(t2g, ShouldAlmostEqual, 4)
		})
	})
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
