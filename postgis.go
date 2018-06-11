// Package postgis contains implementation of postgis primitives.
package postgis

import (
	"encoding/binary"
	"errors"
	"io"
)

// Geometry interface.
type Geometry interface {
	Type() uint32
	Read(reader io.Reader, order binary.ByteOrder) error
	Write(writer io.Writer, order binary.ByteOrder) error
}

const (
	wkbXDR byte = 0
	wkbNDR byte = 1
)

var (
	// ErrUnexpectedGeometryType triggered when detect unexpected geometry type.
	ErrUnexpectedGeometryType = errors.New("unexpected geometry type")

	// ErrUnsupportedByteOrder triggered when byte order unsupported.
	ErrUnsupportedByteOrder = errors.New("unsupported byte order")

	// ErrUnsupportedSourceDataType triggered when source have invalid data type.
	ErrUnsupportedSourceDataType = errors.New("unexpected source data type")
)
