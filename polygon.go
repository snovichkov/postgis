// Package postgis contains implementation of postgis primitives.
package postgis

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"io"
)

// PolygonS type.
type PolygonS struct {
	SRID  int32
	Lines []Line
}

// Type implements Geometry interface.
func (p *PolygonS) Type() uint32 {
	return 0x20000003
}

// Read implement Geometry interface.
func (p *PolygonS) Read(reader io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(reader, order, &p.SRID); err != nil {
		return err
	}

	var wkbLinesCount uint32
	if err = binary.Read(reader, order, &wkbLinesCount); err != nil {
		return err
	}

	p.Lines = make([]Line, wkbLinesCount)

	for i := uint32(0); i < wkbLinesCount; i++ {
		if err = p.Lines[i].Read(reader, order); err != nil {
			return err
		}
	}

	return nil
}

// Write implement Geometry interface.
func (p *PolygonS) Write(writer io.Writer, order binary.ByteOrder) (err error) {
	if err = binary.Write(writer, order, p.SRID); err != nil {
		return err
	}

	if err = binary.Write(writer, order, uint32(len(p.Lines))); err != nil {
		return err
	}

	for _, l := range p.Lines {
		if err = l.Write(writer, order); err != nil {
			return err
		}
	}

	return nil
}

// Value implements the types.Valuer interface.
func (p PolygonS) Value() (_ driver.Value, err error) {
	var buf = bytes.NewBuffer(nil)
	if err = buf.WriteByte(wkbNDR); err != nil {
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
func (p *PolygonS) Scan(src interface{}) (err error) {
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

		if wkbByteOrder, err = reader.ReadByte(); err != nil {
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

		return p.Read(reader, byteOrder)
	default:
		return ErrUnsupportedSourceDataType
	}
}
