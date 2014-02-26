package geometry

import (
	"database/sql"
	_ "github.com/gregb/pq"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func setupRoundtripTest() *sql.DB {
	db, err := sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost:5432/pqgotest?sslmode=disable")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TEMP TABLE geotest (id int, p point, s lseg, b box, c circle);")

	if err != nil {
		panic(err)
	}

	return db
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

func TestPointRoundtrip(t *testing.T) {
	Convey("Given a postgres table with geometric datatypes", t, func() {
		db := setupRoundtripTest()

		p1 := NewPoint(0, 0)
		p2 := NewPoint(-1231, 3242.832)

		Convey("Test that values can be written", func() {
			r1, e1 := db.Exec("INSERT INTO geotest (id, p) VALUES (1, $1)", p1)
			So(e1, ShouldBeNil)

			r2, e2 := db.Exec("INSERT INTO geotest (id, p) VALUES (2, $1)", p2)
			So(e2, ShouldBeNil)

			a1, _ := r1.RowsAffected()
			a2, _ := r2.RowsAffected()

			So(a1, ShouldEqual, 1)
			So(a2, ShouldEqual, 1)
		})
		Convey("Test that values can be read back, and match written values", func() {

			var r1, r2 Point

			row1 := db.QueryRow("SELECT p FROM geotest WHERE id = $1", 1)
			err := row1.Scan(&r1)
			So(err, ShouldBeNil)
			So(r1, ShouldResemble, p1)

			row2 := db.QueryRow("SELECT p FROM geotest WHERE id = $1", 2)
			err = row2.Scan(&r2)
			So(err, ShouldBeNil)
			So(r2, ShouldResemble, p2)

		})
	})
}
