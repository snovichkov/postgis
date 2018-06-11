// Package postgis contains implementation of postgis primitives.
package postgis

import (
	"encoding/binary"
	"io"
)

// Line type.
type Line struct {
	Points []Point
}

// Type implements Geometry interface.
func (*Line) Type() uint32 {
	return 2
}

// Read implements
func (l *Line) Read(reader io.Reader, order binary.ByteOrder) (err error) {
	var count int32
	if err = binary.Read(reader, order, &count); err != nil {
		return err
	}

	l.Points = make([]Point, count)

	for i := int32(0); i < count; i++ {
		if err = l.Points[i].Read(reader, order); err != nil {
			return err
		}
	}

	return nil
}

// Write implements Geometry interface.
func (l *Line) Write(writer io.Writer, order binary.ByteOrder) (err error) {
	if err = binary.Write(writer, order, int32(len(l.Points))); err != nil {
		return err
	}

	for _, p := range l.Points {
		if err = p.Write(writer, order); err != nil {
			return err
		}
	}

	return nil
}
