package geometry

// info from http://www.postgresql.org/docs/9.2/static/datatype-geometric.html

import (
	"database/sql/driver"
	"fmt"
)

// Checks that the number of floats returned by the sql driver matches expectations.
// Src is expected to be a []float64, but the typecast is done here to consolidate error
// checking so that each type's Scan() method does not have to do it itself.
func expectFloats(src interface{}, expected int) ([]float64, error) {

	floats, ok := src.([]float64)

	if !ok {
		return nil, fmt.Errorf("Expected []float64 from driver, got %T instead", src)
	}

	// if positive, expect exactly that number
	if expected > 0 {
		if len(floats) != expected {
			return nil, fmt.Errorf("Expected %d floats while parsing geometry, but got %d instead", expected, len(floats))
		}
	} else {
		// otherwise, any multiple of |expected| is ok
		// if expected == -1, then ANY amount is ok
		extra := len(floats) % (-expected)
		if extra != 0 {
			return nil, fmt.Errorf("Expected a multiple of %d floats while parsing geometry, but got %d instead", -expected, len(floats))
		}
	}

	return floats, nil
}

// assert that types implement driver.Valuer is a pq.Encoder
var _ driver.Valuer = Point{}
var _ driver.Valuer = Vector{}
var _ driver.Valuer = Segment{}
var _ driver.Valuer = Circle{}
var _ driver.Valuer = Box{}

// ----------

func (p *Point) Scan(src interface{}) error {
	floats, err := expectFloats(src, 2)

	if err != nil {
		return fmt.Errorf("Error while parsing data for Point: %s", err)
	}

	p.x = floats[0]
	p.y = floats[1]

	return nil
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", p.x, p.y), nil
}

// ----------

func (v *Vector) Scan(src interface{}) error {
	floats, err := expectFloats(src, 2)

	if err != nil {
		return fmt.Errorf("Error while parsing data for Vector: %s", err)
	}

	v.x = floats[0]
	v.y = floats[1]

	return nil
}

func (v Vector) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", v.x, v.y), nil
}

// ----------

func (s *Segment) Scan(src interface{}) error {
	floats, err := expectFloats(src, 4)

	if err != nil {
		return fmt.Errorf("Error while parsing data for Segment: %s", err)
	}

	s[0].x = floats[0]
	s[0].y = floats[1]
	s[1].x = floats[2]
	s[1].y = floats[3]

	return nil
}

func (s Segment) Value() (driver.Value, error) {
	return fmt.Sprintf("[(%f,%f),(%f,%f)]", s[0].x, s[0].y, s[1].x, s[1].y), nil
}

// ----------

func (b *Box) Scan(src interface{}) error {
	floats, err := expectFloats(src, 4)

	if err != nil {
		return fmt.Errorf("Error while parsing data for Box: %s", err)
	}

	b[0].x = floats[0]
	b[0].y = floats[1]
	b[1].x = floats[2]
	b[1].y = floats[3]

	return nil
}

func (b Box) Value() (driver.Value, error) {
	return fmt.Sprintf("((%f,%f),(%f,%f))", b[0].x, b[0].y, b[1].x, b[1].y), nil
}

// ----------

func (c *Circle) Scan(src interface{}) error {
	floats, err := expectFloats(src, 3)

	if err != nil {
		return fmt.Errorf("Error while parsing data for Circle: %s", err)
	}

	c.center.x = floats[0]
	c.center.y = floats[1]
	c.radius = floats[2]

	return nil
}

func (c Circle) Value() (driver.Value, error) {
	return fmt.Sprintf("<(%f,%f),%f>", c.center.x, c.center.y, c.radius), nil
}
