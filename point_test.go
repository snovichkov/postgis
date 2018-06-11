package postgis

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointS_Type(t *testing.T) {
	var p PointS
	assert.EqualValues(t, 0x20000001, p.Type())
}

func TestPointS_Scan(t *testing.T) {
	var testCases = []struct {
		value  interface{}
		error  error
		expect PointS
	}{{
		value: nil,
	}, {
		value:  []uint8("0101000020E61000006FBBD05CA7AF404079E57ADB4C874140"),
		expect: PointS{4326, 33.372295, 35.057033},
	}, {
		value: []byte("0102000020E61000006FBBD05CA7AF404079E57ADB4C874140"),
		error: ErrUnexpectedGeometryType,
	}, {
		value: []byte("0201000020E61000006FBBD05CA7AF404079E57ADB4C874140"),
		error: ErrUnsupportedByteOrder,
	}, {
		value: 1,
		error: ErrUnsupportedSourceDataType,
	}}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", index), func(t *testing.T) {
			var (
				p   PointS
				err = p.Scan(testCase.value)
			)

			assert.EqualValues(t, testCase.error, err)
			assert.EqualValues(t, testCase.expect.SRID, p.SRID)
			assert.EqualValues(t, testCase.expect.X, p.X)
			assert.EqualValues(t, testCase.expect.Y, p.Y)
		})
	}
}

func TestPointS_Value(t *testing.T) {
	var testCases = []struct {
		value  PointS
		error  error
		expect string
	}{{
		value:  PointS{4326, 33.372295, 35.057033},
		expect: "0101000020E61000006FBBD05CA7AF404079E57ADB4C874140",
	}}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", index), func(t *testing.T) {
			var (
				raw, err = testCase.value.Value()
				encoded  = hex.EncodeToString(raw.([]byte))
			)

			assert.EqualValues(t, testCase.error, err)
			assert.EqualValues(t, testCase.expect, strings.ToUpper(encoded))
		})
	}
}
