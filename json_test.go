package geometry

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMarshalPoint(t *testing.T) {
	Convey("Given a point or vector", t, func() {

		p1 := Point{1, 2}
		p2 := Point{-1234, 5678}
		v3 := Vector{0.23423, 32.123122256}
		v4 := Vector{-8451394857194, 0.00000003}

		Convey("Test that generated correct Json for array format", func() {

			Options.Point = Array
			Options.Vector = Array

			b1, e1 := p1.MarshalJSON()
			So(e1, ShouldBeNil)
			So(string(b1), ShouldEqual, "[1,2]")

			b2, e2 := p2.MarshalJSON()
			So(e2, ShouldBeNil)
			So(string(b2), ShouldEqual, "[-1234,5678]")

			b3, e3 := v3.MarshalJSON()
			So(e3, ShouldBeNil)
			So(string(b3), ShouldEqual, "[0.23423,32.123122256]")

			b4, e4 := v4.MarshalJSON()
			So(e4, ShouldBeNil)
			So(string(b4), ShouldEqual, "[-8.451394857194e+12,3e-08]")

		})

		Convey("Test that generated correct Json for object format", func() {

			Options.Point = Object
			Options.Vector = Object

			b1, e1 := p1.MarshalJSON()
			So(e1, ShouldBeNil)
			So(string(b1), ShouldEqual, `{"x":1,"y":2}`)

			b2, e2 := p2.MarshalJSON()
			So(e2, ShouldBeNil)
			So(string(b2), ShouldEqual, `{"x":-1234,"y":5678}`)

			b3, e3 := v3.MarshalJSON()
			So(e3, ShouldBeNil)
			So(string(b3), ShouldEqual, `{"x":0.23423,"y":32.123122256}`)

			b4, e4 := v4.MarshalJSON()
			So(e4, ShouldBeNil)
			So(string(b4), ShouldEqual, `{"x":-8.451394857194e+12,"y":3e-08}`)

		})
	})

	Convey("Given a segment or box", t, func() {

		s1 := Segment{Point{1, 2}, Point{-1234, 5678}}
		b2 := Box{Point{8451394857194, 32.123122256}, Point{0.23423, 0.00000003}}

		Convey("Test that generated correct Json for array format", func() {

			Options.Point = Array
			Options.Box = Array
			Options.Segment = Array

			r1, e1 := s1.MarshalJSON()
			So(e1, ShouldBeNil)
			So(string(r1), ShouldEqual, "[1,2,-1234,5678]")

			r2, e2 := b2.MarshalJSON()
			So(e2, ShouldBeNil)
			So(string(r2), ShouldEqual, "[8.451394857194e+12,32.123122256,0.23423,3e-08]")

		})

		Convey("Test that generated correct Json for object format", func() {

			Options.Point = Object
			Options.Box = Object
			Options.Segment = Object

			r1, e1 := s1.MarshalJSON()
			So(e1, ShouldBeNil)
			So(string(r1), ShouldEqual, `{"0":{"x":1,"y":2},"1":{"x":-1234,"y":5678}}`)

			r2, e2 := b2.MarshalJSON()
			So(e2, ShouldBeNil)
			So(string(r2), ShouldEqual, `{"0":{"x":8.451394857194e+12,"y":32.123122256},"1":{"x":0.23423,"y":3e-08}}`)

		})

			Convey("Test that generated correct Json for compound format", func() {

					Options.Point = Object
					Options.Box = Compound
					Options.Segment = Compound

					r1, e1 := s1.MarshalJSON()
					So(e1, ShouldBeNil)
					So(string(r1), ShouldEqual, `[{"x":1,"y":2},{"x":-1234,"y":5678}]`)

					r2, e2 := b2.MarshalJSON()
					So(e2, ShouldBeNil)
					So(string(r2), ShouldEqual, `[{"x":8.451394857194e+12,"y":32.123122256},{"x":0.23423,"y":3e-08}]`)

					Options.Point = Array
					Options.Box = Compound
					Options.Segment = Compound

					r1, e1 = s1.MarshalJSON()
					So(e1, ShouldBeNil)
					So(string(r1), ShouldEqual, `[[1,2],[-1234,5678]]`)

					r2, e2 = b2.MarshalJSON()
					So(e2, ShouldBeNil)
					So(string(r2), ShouldEqual, `[[8.451394857194e+12,32.123122256],[0.23423,3e-08]]`)

				})
		})
}
