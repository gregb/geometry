package geometry

import (
	_ "encoding/json"
	"strconv"
)

type FormatFlag byte

const Array FormatFlag = 10
const Compound FormatFlag = 11
const Object FormatFlag = 12

type JsonOptions struct {
	Point   FormatFlag
	Vector  FormatFlag
	Segment FormatFlag
	Box     FormatFlag
	Circle  FormatFlag
}

var DefaultJsonOptions = JsonOptions{
	Point:   Array,
	Vector:  Array,
	Segment: Compound,
	Box:     Compound,
	Circle:  Compound,
}

var Options = DefaultJsonOptions

func appendPoint(b []byte, p Point, style FormatFlag) []byte {

	switch style {
	case Array, Compound:
		b = append(b, '[')
		b = strconv.AppendFloat(b, p.x, 'g', -1, 64)
		b = append(b, ',')
		b = strconv.AppendFloat(b, p.y, 'g', -1, 64)
		b = append(b, ']')
	case Object:
		b = append(b, '{', '"', 'x', '"', ':')
		b = strconv.AppendFloat(b, p.x, 'g', -1, 64)
		b = append(b, ',', '"', 'y', '"', ':')
		b = strconv.AppendFloat(b, p.y, 'g', -1, 64)
		b = append(b, '}')
	}

	return b
}

func appendSegmentOrBox(b []byte, ps [2]Point, style FormatFlag) []byte {

	switch style {
	case Array:
		b = append(b, '[')
		b = strconv.AppendFloat(b, ps[0].x, 'g', -1, 64)
		b = append(b, ',')
		b = strconv.AppendFloat(b, ps[0].y, 'g', -1, 64)
		b = append(b, ',')
		b = strconv.AppendFloat(b, ps[1].x, 'g', -1, 64)
		b = append(b, ',')
		b = strconv.AppendFloat(b, ps[1].y, 'g', -1, 64)
		b = append(b, ']')
	case Compound:
		b = append(b, '[')
		b = appendPoint(b, ps[0], Options.Point)
		b = append(b, ',')
		b = appendPoint(b, ps[1], Options.Point)
		b = append(b, ']')
	case Object:
		b = append(b, '{', '"', '0', '"', ':')
		b = appendPoint(b, ps[0], Options.Point)
		b = append(b, ',', '"', '1', '"', ':')
		b = appendPoint(b, ps[1], Options.Point)
		b = append(b, '}')
	}

	return b
}

func appendCircle(b []byte, c Circle, style FormatFlag) []byte {

	switch style {
	case Array:
		b = append(b, '[')
		b = appendPoint(b, c.center, Array)
		b = append(b, ',')
		b = strconv.AppendFloat(b, c.radius, 'g', -1, 64)
		b = append(b, ']')
	case Compound:
		b = append(b, '[')
		b = appendPoint(b, c.center, Options.Point)
		b = append(b, ',')
		b = strconv.AppendFloat(b, c.radius, 'g', -1, 64)
		b = append(b, ']')
	case Object:
		b = append(b, '{', '"', 'c', '"', ':')
		b = appendPoint(b, c.center, Options.Point)
		b = append(b, ',', '"', 'r', '"', ':')
		b = strconv.AppendFloat(b, c.radius, 'g', -1, 64)
		b = append(b, '}')
	}

	return b
}

// Implements json.Marshaller interface
func (p Point) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 8)
	b = appendPoint(b, p, Options.Point)
	return b, nil
}

// Implements json.Marshaller interface
func (v Vector) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 8)
	p := Point(v)
	b = appendPoint(b, p, Options.Vector)
	return b, nil
}

// Implements json.Marshaller interface
func (b Box) MarshalJSON() ([]byte, error) {
	bytes := make([]byte, 0, 8)
	bytes = appendSegmentOrBox(bytes, b, Options.Box)
	return bytes, nil
}

// Implements json.Marshaller interface
func (s Segment) MarshalJSON() ([]byte, error) {
	bytes := make([]byte, 0, 8)
	bytes = appendSegmentOrBox(bytes, s, Options.Segment)
	return bytes, nil
}

// Implements json.Marshaller interface
func (c Circle) MarshalJSON() ([]byte, error) {
	bytes := make([]byte, 0, 8)
	bytes = appendCircle(bytes, c, Options.Circle)
	return bytes, nil
}
