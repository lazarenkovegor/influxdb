// Code generated by tmpl; DO NOT EDIT.
// https://github.com/benbjohnson/tmpl
//
// Source: point.gen.go.tmpl

package query

//lint:file-ignore U1000 Ignore all unused code, it's generated

import (
	"context"
	"encoding/binary"
	"io"

	internal "github.com/influxdata/influxdb/v2/influxql/query/internal"
	"google.golang.org/protobuf/proto"
)

// FloatPoint represents a point with a float64 value.
// DO NOT ADD ADDITIONAL FIELDS TO THIS STRUCT.
// See TestPoint_Fields in influxql/point_test.go for more details.
type FloatPoint struct {
	Name string
	Tags Tags

	Time  int64
	Value float64
	Aux   []interface{}

	// Total number of points that were combined into this point from an aggregate.
	// If this is zero, the point is not the result of an aggregate function.
	Aggregated uint32
	Nil        bool
}

func (v *FloatPoint) name() string { return v.Name }
func (v *FloatPoint) tags() Tags   { return v.Tags }
func (v *FloatPoint) time() int64  { return v.Time }
func (v *FloatPoint) nil() bool    { return v.Nil }
func (v *FloatPoint) value() interface{} {
	if v.Nil {
		return nil
	}
	return v.Value
}
func (v *FloatPoint) aux() []interface{} { return v.Aux }

// Clone returns a copy of v.
func (v *FloatPoint) Clone() *FloatPoint {
	if v == nil {
		return nil
	}

	other := *v
	if v.Aux != nil {
		other.Aux = make([]interface{}, len(v.Aux))
		copy(other.Aux, v.Aux)
	}

	return &other
}

// CopyTo makes a deep copy into the point.
func (v *FloatPoint) CopyTo(other *FloatPoint) {
	other.Name, other.Tags = v.Name, v.Tags
	other.Time = v.Time
	other.Value, other.Nil = v.Value, v.Nil
	if v.Aux != nil {
		if len(other.Aux) != len(v.Aux) {
			other.Aux = make([]interface{}, len(v.Aux))
		}
		copy(other.Aux, v.Aux)
	}
}

func encodeFloatPoint(p *FloatPoint) *internal.Point {
	return &internal.Point{
		Name:       proto.String(p.Name),
		Tags:       proto.String(p.Tags.ID()),
		Time:       proto.Int64(p.Time),
		Nil:        proto.Bool(p.Nil),
		Aux:        encodeAux(p.Aux),
		Aggregated: proto.Uint32(p.Aggregated),

		FloatValue: proto.Float64(p.Value),
	}
}

func decodeFloatPoint(pb *internal.Point) *FloatPoint {
	return &FloatPoint{
		Name:       pb.GetName(),
		Tags:       newTagsID(pb.GetTags()),
		Time:       pb.GetTime(),
		Nil:        pb.GetNil(),
		Aux:        decodeAux(pb.Aux),
		Aggregated: pb.GetAggregated(),
		Value:      pb.GetFloatValue(),
	}
}

// floatPoints represents a slice of points sortable by value.
type floatPoints []FloatPoint

func (a floatPoints) Len() int { return len(a) }
func (a floatPoints) Less(i, j int) bool {
	if a[i].Time != a[j].Time {
		return a[i].Time < a[j].Time
	}
	return a[i].Value < a[j].Value
}
func (a floatPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// floatPointsByValue represents a slice of points sortable by value.
type floatPointsByValue []FloatPoint

func (a floatPointsByValue) Len() int { return len(a) }

func (a floatPointsByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a floatPointsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// floatPointsByTime represents a slice of points sortable by value.
type floatPointsByTime []FloatPoint

func (a floatPointsByTime) Len() int           { return len(a) }
func (a floatPointsByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a floatPointsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// floatPointByFunc represents a slice of points sortable by a function.
type floatPointsByFunc struct {
	points []FloatPoint
	cmp    func(a, b *FloatPoint) bool
}

func (a *floatPointsByFunc) Len() int           { return len(a.points) }
func (a *floatPointsByFunc) Less(i, j int) bool { return a.cmp(&a.points[i], &a.points[j]) }
func (a *floatPointsByFunc) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }

func (a *floatPointsByFunc) Push(x interface{}) {
	a.points = append(a.points, x.(FloatPoint))
}

func (a *floatPointsByFunc) Pop() interface{} {
	p := a.points[len(a.points)-1]
	a.points = a.points[:len(a.points)-1]
	return p
}

func floatPointsSortBy(points []FloatPoint, cmp func(a, b *FloatPoint) bool) *floatPointsByFunc {
	return &floatPointsByFunc{
		points: points,
		cmp:    cmp,
	}
}

// FloatPointEncoder encodes FloatPoint points to a writer.
type FloatPointEncoder struct {
	w io.Writer
}

// NewFloatPointEncoder returns a new instance of FloatPointEncoder that writes to w.
func NewFloatPointEncoder(w io.Writer) *FloatPointEncoder {
	return &FloatPointEncoder{w: w}
}

// EncodeFloatPoint marshals and writes p to the underlying writer.
func (enc *FloatPointEncoder) EncodeFloatPoint(p *FloatPoint) error {
	// Marshal to bytes.
	buf, err := proto.Marshal(encodeFloatPoint(p))
	if err != nil {
		return err
	}

	// Write the length.
	if err := binary.Write(enc.w, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}

	// Write the encoded point.
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// FloatPointDecoder decodes FloatPoint points from a reader.
type FloatPointDecoder struct {
	r     io.Reader
	stats IteratorStats
	ctx   context.Context
}

// NewFloatPointDecoder returns a new instance of FloatPointDecoder that reads from r.
func NewFloatPointDecoder(ctx context.Context, r io.Reader) *FloatPointDecoder {
	return &FloatPointDecoder{r: r, ctx: ctx}
}

// Stats returns iterator stats embedded within the stream.
func (dec *FloatPointDecoder) Stats() IteratorStats { return dec.stats }

// DecodeFloatPoint reads from the underlying reader and unmarshals into p.
func (dec *FloatPointDecoder) DecodeFloatPoint(p *FloatPoint) error {
	for {
		// Read length.
		var sz uint32
		if err := binary.Read(dec.r, binary.BigEndian, &sz); err != nil {
			return err
		}

		// Read point data.
		buf := make([]byte, sz)
		if _, err := io.ReadFull(dec.r, buf); err != nil {
			return err
		}

		// Unmarshal into point.
		var pb internal.Point
		if err := proto.Unmarshal(buf, &pb); err != nil {
			return err
		}

		// If the point contains stats then read stats and retry.
		if pb.Stats != nil {
			dec.stats = decodeIteratorStats(pb.Stats)
			continue
		}

		// Decode into point object.
		*p = *decodeFloatPoint(&pb)

		return nil
	}
}

// IntegerPoint represents a point with a int64 value.
// DO NOT ADD ADDITIONAL FIELDS TO THIS STRUCT.
// See TestPoint_Fields in influxql/point_test.go for more details.
type IntegerPoint struct {
	Name string
	Tags Tags

	Time  int64
	Value int64
	Aux   []interface{}

	// Total number of points that were combined into this point from an aggregate.
	// If this is zero, the point is not the result of an aggregate function.
	Aggregated uint32
	Nil        bool
}

func (v *IntegerPoint) name() string { return v.Name }
func (v *IntegerPoint) tags() Tags   { return v.Tags }
func (v *IntegerPoint) time() int64  { return v.Time }
func (v *IntegerPoint) nil() bool    { return v.Nil }
func (v *IntegerPoint) value() interface{} {
	if v.Nil {
		return nil
	}
	return v.Value
}
func (v *IntegerPoint) aux() []interface{} { return v.Aux }

// Clone returns a copy of v.
func (v *IntegerPoint) Clone() *IntegerPoint {
	if v == nil {
		return nil
	}

	other := *v
	if v.Aux != nil {
		other.Aux = make([]interface{}, len(v.Aux))
		copy(other.Aux, v.Aux)
	}

	return &other
}

// CopyTo makes a deep copy into the point.
func (v *IntegerPoint) CopyTo(other *IntegerPoint) {
	other.Name, other.Tags = v.Name, v.Tags
	other.Time = v.Time
	other.Value, other.Nil = v.Value, v.Nil
	if v.Aux != nil {
		if len(other.Aux) != len(v.Aux) {
			other.Aux = make([]interface{}, len(v.Aux))
		}
		copy(other.Aux, v.Aux)
	}
}

func encodeIntegerPoint(p *IntegerPoint) *internal.Point {
	return &internal.Point{
		Name:       proto.String(p.Name),
		Tags:       proto.String(p.Tags.ID()),
		Time:       proto.Int64(p.Time),
		Nil:        proto.Bool(p.Nil),
		Aux:        encodeAux(p.Aux),
		Aggregated: proto.Uint32(p.Aggregated),

		IntegerValue: proto.Int64(p.Value),
	}
}

func decodeIntegerPoint(pb *internal.Point) *IntegerPoint {
	return &IntegerPoint{
		Name:       pb.GetName(),
		Tags:       newTagsID(pb.GetTags()),
		Time:       pb.GetTime(),
		Nil:        pb.GetNil(),
		Aux:        decodeAux(pb.Aux),
		Aggregated: pb.GetAggregated(),
		Value:      pb.GetIntegerValue(),
	}
}

// integerPoints represents a slice of points sortable by value.
type integerPoints []IntegerPoint

func (a integerPoints) Len() int { return len(a) }
func (a integerPoints) Less(i, j int) bool {
	if a[i].Time != a[j].Time {
		return a[i].Time < a[j].Time
	}
	return a[i].Value < a[j].Value
}
func (a integerPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// integerPointsByValue represents a slice of points sortable by value.
type integerPointsByValue []IntegerPoint

func (a integerPointsByValue) Len() int { return len(a) }

func (a integerPointsByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a integerPointsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// integerPointsByTime represents a slice of points sortable by value.
type integerPointsByTime []IntegerPoint

func (a integerPointsByTime) Len() int           { return len(a) }
func (a integerPointsByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a integerPointsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// integerPointByFunc represents a slice of points sortable by a function.
type integerPointsByFunc struct {
	points []IntegerPoint
	cmp    func(a, b *IntegerPoint) bool
}

func (a *integerPointsByFunc) Len() int           { return len(a.points) }
func (a *integerPointsByFunc) Less(i, j int) bool { return a.cmp(&a.points[i], &a.points[j]) }
func (a *integerPointsByFunc) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }

func (a *integerPointsByFunc) Push(x interface{}) {
	a.points = append(a.points, x.(IntegerPoint))
}

func (a *integerPointsByFunc) Pop() interface{} {
	p := a.points[len(a.points)-1]
	a.points = a.points[:len(a.points)-1]
	return p
}

func integerPointsSortBy(points []IntegerPoint, cmp func(a, b *IntegerPoint) bool) *integerPointsByFunc {
	return &integerPointsByFunc{
		points: points,
		cmp:    cmp,
	}
}

// IntegerPointEncoder encodes IntegerPoint points to a writer.
type IntegerPointEncoder struct {
	w io.Writer
}

// NewIntegerPointEncoder returns a new instance of IntegerPointEncoder that writes to w.
func NewIntegerPointEncoder(w io.Writer) *IntegerPointEncoder {
	return &IntegerPointEncoder{w: w}
}

// EncodeIntegerPoint marshals and writes p to the underlying writer.
func (enc *IntegerPointEncoder) EncodeIntegerPoint(p *IntegerPoint) error {
	// Marshal to bytes.
	buf, err := proto.Marshal(encodeIntegerPoint(p))
	if err != nil {
		return err
	}

	// Write the length.
	if err := binary.Write(enc.w, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}

	// Write the encoded point.
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// IntegerPointDecoder decodes IntegerPoint points from a reader.
type IntegerPointDecoder struct {
	r     io.Reader
	stats IteratorStats
	ctx   context.Context
}

// NewIntegerPointDecoder returns a new instance of IntegerPointDecoder that reads from r.
func NewIntegerPointDecoder(ctx context.Context, r io.Reader) *IntegerPointDecoder {
	return &IntegerPointDecoder{r: r, ctx: ctx}
}

// Stats returns iterator stats embedded within the stream.
func (dec *IntegerPointDecoder) Stats() IteratorStats { return dec.stats }

// DecodeIntegerPoint reads from the underlying reader and unmarshals into p.
func (dec *IntegerPointDecoder) DecodeIntegerPoint(p *IntegerPoint) error {
	for {
		// Read length.
		var sz uint32
		if err := binary.Read(dec.r, binary.BigEndian, &sz); err != nil {
			return err
		}

		// Read point data.
		buf := make([]byte, sz)
		if _, err := io.ReadFull(dec.r, buf); err != nil {
			return err
		}

		// Unmarshal into point.
		var pb internal.Point
		if err := proto.Unmarshal(buf, &pb); err != nil {
			return err
		}

		// If the point contains stats then read stats and retry.
		if pb.Stats != nil {
			dec.stats = decodeIteratorStats(pb.Stats)
			continue
		}

		// Decode into point object.
		*p = *decodeIntegerPoint(&pb)

		return nil
	}
}

// UnsignedPoint represents a point with a uint64 value.
// DO NOT ADD ADDITIONAL FIELDS TO THIS STRUCT.
// See TestPoint_Fields in influxql/point_test.go for more details.
type UnsignedPoint struct {
	Name string
	Tags Tags

	Time  int64
	Value uint64
	Aux   []interface{}

	// Total number of points that were combined into this point from an aggregate.
	// If this is zero, the point is not the result of an aggregate function.
	Aggregated uint32
	Nil        bool
}

func (v *UnsignedPoint) name() string { return v.Name }
func (v *UnsignedPoint) tags() Tags   { return v.Tags }
func (v *UnsignedPoint) time() int64  { return v.Time }
func (v *UnsignedPoint) nil() bool    { return v.Nil }
func (v *UnsignedPoint) value() interface{} {
	if v.Nil {
		return nil
	}
	return v.Value
}
func (v *UnsignedPoint) aux() []interface{} { return v.Aux }

// Clone returns a copy of v.
func (v *UnsignedPoint) Clone() *UnsignedPoint {
	if v == nil {
		return nil
	}

	other := *v
	if v.Aux != nil {
		other.Aux = make([]interface{}, len(v.Aux))
		copy(other.Aux, v.Aux)
	}

	return &other
}

// CopyTo makes a deep copy into the point.
func (v *UnsignedPoint) CopyTo(other *UnsignedPoint) {
	other.Name, other.Tags = v.Name, v.Tags
	other.Time = v.Time
	other.Value, other.Nil = v.Value, v.Nil
	if v.Aux != nil {
		if len(other.Aux) != len(v.Aux) {
			other.Aux = make([]interface{}, len(v.Aux))
		}
		copy(other.Aux, v.Aux)
	}
}

func encodeUnsignedPoint(p *UnsignedPoint) *internal.Point {
	return &internal.Point{
		Name:       proto.String(p.Name),
		Tags:       proto.String(p.Tags.ID()),
		Time:       proto.Int64(p.Time),
		Nil:        proto.Bool(p.Nil),
		Aux:        encodeAux(p.Aux),
		Aggregated: proto.Uint32(p.Aggregated),
	}
}

func decodeUnsignedPoint(pb *internal.Point) *UnsignedPoint {
	return &UnsignedPoint{
		Name:       pb.GetName(),
		Tags:       newTagsID(pb.GetTags()),
		Time:       pb.GetTime(),
		Nil:        pb.GetNil(),
		Aux:        decodeAux(pb.Aux),
		Aggregated: pb.GetAggregated(),
		Value:      pb.GetUnsignedValue(),
	}
}

// unsignedPoints represents a slice of points sortable by value.
type unsignedPoints []UnsignedPoint

func (a unsignedPoints) Len() int { return len(a) }
func (a unsignedPoints) Less(i, j int) bool {
	if a[i].Time != a[j].Time {
		return a[i].Time < a[j].Time
	}
	return a[i].Value < a[j].Value
}
func (a unsignedPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// unsignedPointsByValue represents a slice of points sortable by value.
type unsignedPointsByValue []UnsignedPoint

func (a unsignedPointsByValue) Len() int { return len(a) }

func (a unsignedPointsByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a unsignedPointsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// unsignedPointsByTime represents a slice of points sortable by value.
type unsignedPointsByTime []UnsignedPoint

func (a unsignedPointsByTime) Len() int           { return len(a) }
func (a unsignedPointsByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a unsignedPointsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// unsignedPointByFunc represents a slice of points sortable by a function.
type unsignedPointsByFunc struct {
	points []UnsignedPoint
	cmp    func(a, b *UnsignedPoint) bool
}

func (a *unsignedPointsByFunc) Len() int           { return len(a.points) }
func (a *unsignedPointsByFunc) Less(i, j int) bool { return a.cmp(&a.points[i], &a.points[j]) }
func (a *unsignedPointsByFunc) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }

func (a *unsignedPointsByFunc) Push(x interface{}) {
	a.points = append(a.points, x.(UnsignedPoint))
}

func (a *unsignedPointsByFunc) Pop() interface{} {
	p := a.points[len(a.points)-1]
	a.points = a.points[:len(a.points)-1]
	return p
}

func unsignedPointsSortBy(points []UnsignedPoint, cmp func(a, b *UnsignedPoint) bool) *unsignedPointsByFunc {
	return &unsignedPointsByFunc{
		points: points,
		cmp:    cmp,
	}
}

// UnsignedPointEncoder encodes UnsignedPoint points to a writer.
type UnsignedPointEncoder struct {
	w io.Writer
}

// NewUnsignedPointEncoder returns a new instance of UnsignedPointEncoder that writes to w.
func NewUnsignedPointEncoder(w io.Writer) *UnsignedPointEncoder {
	return &UnsignedPointEncoder{w: w}
}

// EncodeUnsignedPoint marshals and writes p to the underlying writer.
func (enc *UnsignedPointEncoder) EncodeUnsignedPoint(p *UnsignedPoint) error {
	// Marshal to bytes.
	buf, err := proto.Marshal(encodeUnsignedPoint(p))
	if err != nil {
		return err
	}

	// Write the length.
	if err := binary.Write(enc.w, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}

	// Write the encoded point.
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// UnsignedPointDecoder decodes UnsignedPoint points from a reader.
type UnsignedPointDecoder struct {
	r     io.Reader
	stats IteratorStats
	ctx   context.Context
}

// NewUnsignedPointDecoder returns a new instance of UnsignedPointDecoder that reads from r.
func NewUnsignedPointDecoder(ctx context.Context, r io.Reader) *UnsignedPointDecoder {
	return &UnsignedPointDecoder{r: r, ctx: ctx}
}

// Stats returns iterator stats embedded within the stream.
func (dec *UnsignedPointDecoder) Stats() IteratorStats { return dec.stats }

// DecodeUnsignedPoint reads from the underlying reader and unmarshals into p.
func (dec *UnsignedPointDecoder) DecodeUnsignedPoint(p *UnsignedPoint) error {
	for {
		// Read length.
		var sz uint32
		if err := binary.Read(dec.r, binary.BigEndian, &sz); err != nil {
			return err
		}

		// Read point data.
		buf := make([]byte, sz)
		if _, err := io.ReadFull(dec.r, buf); err != nil {
			return err
		}

		// Unmarshal into point.
		var pb internal.Point
		if err := proto.Unmarshal(buf, &pb); err != nil {
			return err
		}

		// If the point contains stats then read stats and retry.
		if pb.Stats != nil {
			dec.stats = decodeIteratorStats(pb.Stats)
			continue
		}

		// Decode into point object.
		*p = *decodeUnsignedPoint(&pb)

		return nil
	}
}

// StringPoint represents a point with a string value.
// DO NOT ADD ADDITIONAL FIELDS TO THIS STRUCT.
// See TestPoint_Fields in influxql/point_test.go for more details.
type StringPoint struct {
	Name string
	Tags Tags

	Time  int64
	Value string
	Aux   []interface{}

	// Total number of points that were combined into this point from an aggregate.
	// If this is zero, the point is not the result of an aggregate function.
	Aggregated uint32
	Nil        bool
}

func (v *StringPoint) name() string { return v.Name }
func (v *StringPoint) tags() Tags   { return v.Tags }
func (v *StringPoint) time() int64  { return v.Time }
func (v *StringPoint) nil() bool    { return v.Nil }
func (v *StringPoint) value() interface{} {
	if v.Nil {
		return nil
	}
	return v.Value
}
func (v *StringPoint) aux() []interface{} { return v.Aux }

// Clone returns a copy of v.
func (v *StringPoint) Clone() *StringPoint {
	if v == nil {
		return nil
	}

	other := *v
	if v.Aux != nil {
		other.Aux = make([]interface{}, len(v.Aux))
		copy(other.Aux, v.Aux)
	}

	return &other
}

// CopyTo makes a deep copy into the point.
func (v *StringPoint) CopyTo(other *StringPoint) {
	other.Name, other.Tags = v.Name, v.Tags
	other.Time = v.Time
	other.Value, other.Nil = v.Value, v.Nil
	if v.Aux != nil {
		if len(other.Aux) != len(v.Aux) {
			other.Aux = make([]interface{}, len(v.Aux))
		}
		copy(other.Aux, v.Aux)
	}
}

func encodeStringPoint(p *StringPoint) *internal.Point {
	return &internal.Point{
		Name:       proto.String(p.Name),
		Tags:       proto.String(p.Tags.ID()),
		Time:       proto.Int64(p.Time),
		Nil:        proto.Bool(p.Nil),
		Aux:        encodeAux(p.Aux),
		Aggregated: proto.Uint32(p.Aggregated),

		StringValue: proto.String(p.Value),
	}
}

func decodeStringPoint(pb *internal.Point) *StringPoint {
	return &StringPoint{
		Name:       pb.GetName(),
		Tags:       newTagsID(pb.GetTags()),
		Time:       pb.GetTime(),
		Nil:        pb.GetNil(),
		Aux:        decodeAux(pb.Aux),
		Aggregated: pb.GetAggregated(),
		Value:      pb.GetStringValue(),
	}
}

// stringPoints represents a slice of points sortable by value.
type stringPoints []StringPoint

func (a stringPoints) Len() int { return len(a) }
func (a stringPoints) Less(i, j int) bool {
	if a[i].Time != a[j].Time {
		return a[i].Time < a[j].Time
	}
	return a[i].Value < a[j].Value
}
func (a stringPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// stringPointsByValue represents a slice of points sortable by value.
type stringPointsByValue []StringPoint

func (a stringPointsByValue) Len() int { return len(a) }

func (a stringPointsByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a stringPointsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// stringPointsByTime represents a slice of points sortable by value.
type stringPointsByTime []StringPoint

func (a stringPointsByTime) Len() int           { return len(a) }
func (a stringPointsByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a stringPointsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// stringPointByFunc represents a slice of points sortable by a function.
type stringPointsByFunc struct {
	points []StringPoint
	cmp    func(a, b *StringPoint) bool
}

func (a *stringPointsByFunc) Len() int           { return len(a.points) }
func (a *stringPointsByFunc) Less(i, j int) bool { return a.cmp(&a.points[i], &a.points[j]) }
func (a *stringPointsByFunc) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }

func (a *stringPointsByFunc) Push(x interface{}) {
	a.points = append(a.points, x.(StringPoint))
}

func (a *stringPointsByFunc) Pop() interface{} {
	p := a.points[len(a.points)-1]
	a.points = a.points[:len(a.points)-1]
	return p
}

func stringPointsSortBy(points []StringPoint, cmp func(a, b *StringPoint) bool) *stringPointsByFunc {
	return &stringPointsByFunc{
		points: points,
		cmp:    cmp,
	}
}

// StringPointEncoder encodes StringPoint points to a writer.
type StringPointEncoder struct {
	w io.Writer
}

// NewStringPointEncoder returns a new instance of StringPointEncoder that writes to w.
func NewStringPointEncoder(w io.Writer) *StringPointEncoder {
	return &StringPointEncoder{w: w}
}

// EncodeStringPoint marshals and writes p to the underlying writer.
func (enc *StringPointEncoder) EncodeStringPoint(p *StringPoint) error {
	// Marshal to bytes.
	buf, err := proto.Marshal(encodeStringPoint(p))
	if err != nil {
		return err
	}

	// Write the length.
	if err := binary.Write(enc.w, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}

	// Write the encoded point.
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// StringPointDecoder decodes StringPoint points from a reader.
type StringPointDecoder struct {
	r     io.Reader
	stats IteratorStats
	ctx   context.Context
}

// NewStringPointDecoder returns a new instance of StringPointDecoder that reads from r.
func NewStringPointDecoder(ctx context.Context, r io.Reader) *StringPointDecoder {
	return &StringPointDecoder{r: r, ctx: ctx}
}

// Stats returns iterator stats embedded within the stream.
func (dec *StringPointDecoder) Stats() IteratorStats { return dec.stats }

// DecodeStringPoint reads from the underlying reader and unmarshals into p.
func (dec *StringPointDecoder) DecodeStringPoint(p *StringPoint) error {
	for {
		// Read length.
		var sz uint32
		if err := binary.Read(dec.r, binary.BigEndian, &sz); err != nil {
			return err
		}

		// Read point data.
		buf := make([]byte, sz)
		if _, err := io.ReadFull(dec.r, buf); err != nil {
			return err
		}

		// Unmarshal into point.
		var pb internal.Point
		if err := proto.Unmarshal(buf, &pb); err != nil {
			return err
		}

		// If the point contains stats then read stats and retry.
		if pb.Stats != nil {
			dec.stats = decodeIteratorStats(pb.Stats)
			continue
		}

		// Decode into point object.
		*p = *decodeStringPoint(&pb)

		return nil
	}
}

// BooleanPoint represents a point with a bool value.
// DO NOT ADD ADDITIONAL FIELDS TO THIS STRUCT.
// See TestPoint_Fields in influxql/point_test.go for more details.
type BooleanPoint struct {
	Name string
	Tags Tags

	Time  int64
	Value bool
	Aux   []interface{}

	// Total number of points that were combined into this point from an aggregate.
	// If this is zero, the point is not the result of an aggregate function.
	Aggregated uint32
	Nil        bool
}

func (v *BooleanPoint) name() string { return v.Name }
func (v *BooleanPoint) tags() Tags   { return v.Tags }
func (v *BooleanPoint) time() int64  { return v.Time }
func (v *BooleanPoint) nil() bool    { return v.Nil }
func (v *BooleanPoint) value() interface{} {
	if v.Nil {
		return nil
	}
	return v.Value
}
func (v *BooleanPoint) aux() []interface{} { return v.Aux }

// Clone returns a copy of v.
func (v *BooleanPoint) Clone() *BooleanPoint {
	if v == nil {
		return nil
	}

	other := *v
	if v.Aux != nil {
		other.Aux = make([]interface{}, len(v.Aux))
		copy(other.Aux, v.Aux)
	}

	return &other
}

// CopyTo makes a deep copy into the point.
func (v *BooleanPoint) CopyTo(other *BooleanPoint) {
	other.Name, other.Tags = v.Name, v.Tags
	other.Time = v.Time
	other.Value, other.Nil = v.Value, v.Nil
	if v.Aux != nil {
		if len(other.Aux) != len(v.Aux) {
			other.Aux = make([]interface{}, len(v.Aux))
		}
		copy(other.Aux, v.Aux)
	}
}

func encodeBooleanPoint(p *BooleanPoint) *internal.Point {
	return &internal.Point{
		Name:       proto.String(p.Name),
		Tags:       proto.String(p.Tags.ID()),
		Time:       proto.Int64(p.Time),
		Nil:        proto.Bool(p.Nil),
		Aux:        encodeAux(p.Aux),
		Aggregated: proto.Uint32(p.Aggregated),

		BooleanValue: proto.Bool(p.Value),
	}
}

func decodeBooleanPoint(pb *internal.Point) *BooleanPoint {
	return &BooleanPoint{
		Name:       pb.GetName(),
		Tags:       newTagsID(pb.GetTags()),
		Time:       pb.GetTime(),
		Nil:        pb.GetNil(),
		Aux:        decodeAux(pb.Aux),
		Aggregated: pb.GetAggregated(),
		Value:      pb.GetBooleanValue(),
	}
}

// booleanPoints represents a slice of points sortable by value.
type booleanPoints []BooleanPoint

func (a booleanPoints) Len() int { return len(a) }
func (a booleanPoints) Less(i, j int) bool {
	if a[i].Time != a[j].Time {
		return a[i].Time < a[j].Time
	}
	return !a[i].Value
}
func (a booleanPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// booleanPointsByValue represents a slice of points sortable by value.
type booleanPointsByValue []BooleanPoint

func (a booleanPointsByValue) Len() int { return len(a) }

func (a booleanPointsByValue) Less(i, j int) bool { return !a[i].Value }

func (a booleanPointsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// booleanPointsByTime represents a slice of points sortable by value.
type booleanPointsByTime []BooleanPoint

func (a booleanPointsByTime) Len() int           { return len(a) }
func (a booleanPointsByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a booleanPointsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// booleanPointByFunc represents a slice of points sortable by a function.
type booleanPointsByFunc struct {
	points []BooleanPoint
	cmp    func(a, b *BooleanPoint) bool
}

func (a *booleanPointsByFunc) Len() int           { return len(a.points) }
func (a *booleanPointsByFunc) Less(i, j int) bool { return a.cmp(&a.points[i], &a.points[j]) }
func (a *booleanPointsByFunc) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }

func (a *booleanPointsByFunc) Push(x interface{}) {
	a.points = append(a.points, x.(BooleanPoint))
}

func (a *booleanPointsByFunc) Pop() interface{} {
	p := a.points[len(a.points)-1]
	a.points = a.points[:len(a.points)-1]
	return p
}

func booleanPointsSortBy(points []BooleanPoint, cmp func(a, b *BooleanPoint) bool) *booleanPointsByFunc {
	return &booleanPointsByFunc{
		points: points,
		cmp:    cmp,
	}
}

// BooleanPointEncoder encodes BooleanPoint points to a writer.
type BooleanPointEncoder struct {
	w io.Writer
}

// NewBooleanPointEncoder returns a new instance of BooleanPointEncoder that writes to w.
func NewBooleanPointEncoder(w io.Writer) *BooleanPointEncoder {
	return &BooleanPointEncoder{w: w}
}

// EncodeBooleanPoint marshals and writes p to the underlying writer.
func (enc *BooleanPointEncoder) EncodeBooleanPoint(p *BooleanPoint) error {
	// Marshal to bytes.
	buf, err := proto.Marshal(encodeBooleanPoint(p))
	if err != nil {
		return err
	}

	// Write the length.
	if err := binary.Write(enc.w, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}

	// Write the encoded point.
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// BooleanPointDecoder decodes BooleanPoint points from a reader.
type BooleanPointDecoder struct {
	r     io.Reader
	stats IteratorStats
	ctx   context.Context
}

// NewBooleanPointDecoder returns a new instance of BooleanPointDecoder that reads from r.
func NewBooleanPointDecoder(ctx context.Context, r io.Reader) *BooleanPointDecoder {
	return &BooleanPointDecoder{r: r, ctx: ctx}
}

// Stats returns iterator stats embedded within the stream.
func (dec *BooleanPointDecoder) Stats() IteratorStats { return dec.stats }

// DecodeBooleanPoint reads from the underlying reader and unmarshals into p.
func (dec *BooleanPointDecoder) DecodeBooleanPoint(p *BooleanPoint) error {
	for {
		// Read length.
		var sz uint32
		if err := binary.Read(dec.r, binary.BigEndian, &sz); err != nil {
			return err
		}

		// Read point data.
		buf := make([]byte, sz)
		if _, err := io.ReadFull(dec.r, buf); err != nil {
			return err
		}

		// Unmarshal into point.
		var pb internal.Point
		if err := proto.Unmarshal(buf, &pb); err != nil {
			return err
		}

		// If the point contains stats then read stats and retry.
		if pb.Stats != nil {
			dec.stats = decodeIteratorStats(pb.Stats)
			continue
		}

		// Decode into point object.
		*p = *decodeBooleanPoint(&pb)

		return nil
	}
}
