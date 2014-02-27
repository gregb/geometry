package geometry

import (
	"database/sql"
	"fmt"
	_ "github.com/gregb/pq"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var createTable = `
CREATE TABLE IF NOT EXISTS geotest
(
  id integer NOT NULL,
  t timestamp without time zone NOT NULL,
  p point,
  s lseg,
  b box,
  c circle,
  CONSTRAINT geotest_pkey PRIMARY KEY (id, t)
)`

var testTime = time.Now()
var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost:5432/pqgotest?sslmode=disable")

	//pq.TrafficLogging = true

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTable)

	if err != nil {
		panic(err)
	}
}

func TestExpectFloats(t *testing.T) {

	Convey("Given a set of []float64", t, func() {

		f1 := []float64{1}
		f2 := []float64{1, 2}
		f3 := []float64{1, 2, 3}
		f4 := []float64{1, 2, 3, 4}
		f5 := []float64{1, 2, 3, 4, 5}
		f6 := []float64{1, 2, 3, 4, 5, 6}

		Convey("Positive expectations should be met exactly", func() {
			r1, e1 := expectFloats(f1, 1)
			r2, e2 := expectFloats(f1, 2)
			r3, e3 := expectFloats(f2, 2)
			r4, e4 := expectFloats(f3, 1)
			r5, e5 := expectFloats(f3, 3)

			So(r1, ShouldResemble, f1)
			So(e1, ShouldBeNil)

			So(r2, ShouldBeNil)
			So(e2, ShouldNotBeNil)

			So(r3, ShouldResemble, f2)
			So(e3, ShouldBeNil)

			So(r4, ShouldBeNil)
			So(e4, ShouldNotBeNil)

			So(r5, ShouldResemble, f3)
			So(e5, ShouldBeNil)
		})

		Convey("Negative expectations should be met in multiples", func() {
			r1, e1 := expectFloats(f1, -1)
			r2, e2 := expectFloats(f6, -1)
			r3, e3 := expectFloats(f2, -2)
			r4, e4 := expectFloats(f3, -2)
			r5, e5 := expectFloats(f4, -2)
			r6, e6 := expectFloats(f5, -2)
			r7, e7 := expectFloats(f6, -2)

			So(r1, ShouldResemble, f1)
			So(e1, ShouldBeNil)

			So(r2, ShouldResemble, f6)
			So(e2, ShouldBeNil)

			So(r3, ShouldResemble, f2)
			So(e3, ShouldBeNil)

			So(r4, ShouldBeNil)
			So(e4, ShouldNotBeNil)

			So(r5, ShouldResemble, f4)
			So(e5, ShouldBeNil)

			So(r6, ShouldBeNil)
			So(e6, ShouldNotBeNil)

			So(r7, ShouldResemble, f6)
			So(e7, ShouldBeNil)
		})

		Convey("Non []float64 should return an error", func() {
			s := "this isn't a float slice"

			r, e := expectFloats(s, 2)
			So(r, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
	})
}

func testRoundtrip(t *testing.T, testId int, column string, value, returned interface{}) {
	const insertTemplate = "INSERT INTO geotest (id, t, %s) VALUES ($1, $2, $3)"
	const selectTemplate = "SELECT %s FROM geotest WHERE id = $1 AND t = $2"

	insertSql := fmt.Sprintf(insertTemplate, column)
	selectSql := fmt.Sprintf(selectTemplate, column)

	res, err := db.Exec(insertSql, testId, testTime, value)
	So(err, ShouldBeNil)
	aff, _ := res.RowsAffected()
	So(aff, ShouldEqual, 1)

	row := db.QueryRow(selectSql, testId, testTime)
	err = row.Scan(returned)
	So(err, ShouldBeNil)

}

func TestPointRoundtrip(t *testing.T) {
	Convey("Given a postgres table with geometric datatypes", t, func() {

		p1 := NewPoint(0, 0)
		p2 := NewPoint(-1231, 3242.832)

		Convey("Test that points can be written, read, and match", func() {

			var r1, r2 Point

			testRoundtrip(t, 1, "p", p1, &r1)
			So(r1, ShouldResemble, p1)

			testRoundtrip(t, 2, "p", p2, &r2)
			So(r2, ShouldResemble, p2)
		})
	})
}

func TestVectorRoundtrip(t *testing.T) {
	Convey("Given a postgres table with geometric datatypes", t, func() {

		// should be able to handle 15 significant digits
		v1 := NewVector(1234.56789, -9876.54321)
		v2 := NewVector(123456789012345, 0.123456789012345)

		Convey("Test that vectors can be written, read, and match", func() {

			var r1, r2 Vector

			testRoundtrip(t, 3, "p", v1, &r1)
			So(r1, ShouldResemble, v1)

			testRoundtrip(t, 4, "p", v2, &r2)
			So(r2, ShouldResemble, v2)
		})
	})
}

func TestSegmentRoundtrip(t *testing.T) {
	Convey("Given a postgres table with geometric datatypes", t, func() {

		// should be able to handle 15 significant digits
		s := NewSegment(Point{1234.56789, -9876.54321}, Point{123456789012345, 0.123456789012345})

		Convey("Test that segments can be written, read, and match", func() {
			var r Segment
			testRoundtrip(t, 5, "s", s, &r)
			So(r, ShouldResemble, s)
		})

	})
}
