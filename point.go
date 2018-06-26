// Package postgis contains implementation of postgis primitives.
package postgis

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"io"
)

type (
	// Point type.
	Point struct {
		X, Y float64
	}

	// PointS type.
	PointS struct {
		SRID int32
		X, Y float64
	}
)

// Type implement Geometry interface.
func (*Point) Type() uint32 {
	return 1
}

// Read implement geometry interface.
func (p *Point) Read(reader io.Reader, order binary.ByteOrder) (err error) {
	return binary.Read(reader, order, p)
}

// Write implement geometry interface.
func (p *Point) Write(writer io.Writer, order binary.ByteOrder) (err error) {
	return binary.Write(writer, order, *p)
}

// Type implement Geometry interface.
func (*PointS) Type() uint32 {
	return 0x20000001
}

// Write implement geometry interface.
func (p *PointS) Read(reader io.Reader, order binary.ByteOrder) (err error) {
	return binary.Read(reader, order, p)
}

// Write implement geometry interface.
func (p *PointS) Write(writer io.Writer, order binary.ByteOrder) (err error) {
	return binary.Write(writer, order, *p)
}

// Value implements the types.Valuer interface.
func (p PointS) Value() (_ driver.Value, err error) {
	var buf = bytes.NewBuffer(nil)
	if err = binary.Write(buf, binary.LittleEndian, wkbNDR); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, p.Type()); err != nil {
		return nil, err
	}

	if err = p.Write(buf, binary.LittleEndian); err != nil {
		return nil, err
	}

	return hex.EncodeToString(
		buf.Bytes(),
	), nil
}

// Scan implements the sql.Scanner interface.
func (p *PointS) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		var raw = make([]byte, hex.DecodedLen(len(v)))
		if _, err = hex.Decode(raw, v); err != nil {
			return err
		}

		var (
			reader       = bytes.NewReader(raw)
			byteOrder    binary.ByteOrder
			wkbByteOrder byte
		)

		if binary.Read(reader, binary.LittleEndian, &wkbByteOrder); err != nil {
			return err
		}

		switch wkbByteOrder {
		case wkbXDR:
			byteOrder = binary.BigEndian
		case wkbNDR:
			byteOrder = binary.LittleEndian
		default:
			return ErrUnsupportedByteOrder
		}

		var wkbGeometryType uint32
		if err = binary.Read(reader, byteOrder, &wkbGeometryType); err != nil {
			return err
		}

		if wkbGeometryType != p.Type() {
			return ErrUnexpectedGeometryType
		}

		return binary.Read(reader, byteOrder, p)
	default:
		return ErrUnsupportedSourceDataType
	}
}
