package geometry

import (
	_ "encoding/json"
	"fmt"
)

const (
	POINT_AS_ARRAY  string  = "[%f,%f]"
	POINT_AS_OBJECT string  = "{\"x\":%f,\"y\":%f}"
	CIRCLE_AS_ARRAY string  = "[%f,%f,%f]"
	CIRCLE_AS_OBJECT string = "{\"c\":{\"x\":%f,\"y\":%f},\"r\":%f}"
)

var JsonMarshalFormat string = POINT_AS_ARRAY

// Implements json.Marshaller interface
func (p *Point) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(JsonMarshalFormat, p.x, p.y)
	return []byte(json), nil
}

// Implements json.Marshaller interface
func (v *Vector) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(JsonMarshalFormat, v.x, v.y)
	return []byte(json), nil
}

// Implements json.Marshaller interface
func (c *Circle) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(JsonMarshalFormat, c.center.x, c.center.y, c.radius)
	return []byte(json), nil
}
