package postgis

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygonS_Type(t *testing.T) {
	var p PolygonS
	assert.EqualValues(t, 0x20000003, p.Type())
}

func TestPolygonS_Value(t *testing.T) {
	var testCases = []struct {
		value  PolygonS
		error  error
		expect string
	}{{
		value: PolygonS{
			4326,
			[]Line{{
				Points: []Point{
					{75.15, 29.53},
					{77, 29},
					{77.6, 29.5},
					{75.15, 29.53},
				},
			}},
		},
		expect: "0103000020E610000001000000040000009A99999999C9524048E17A14AE873D4000000000004053400000000000003D4066666666666653400000000000803D409A99999999C9524048E17A14AE873D40",
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

func TestPolygonS_Scan(t *testing.T) {
	var testCases = []struct {
		value string
		error error
	}{{
		value: "0103000020E61000000100000004000000C38366D7BD413E40BADBF5D214FD4D40E7FBA9F1D2413E40E2067C7E18FD4D4052B81E85EB413E40D237691A14FD4D40C38366D7BD413E40BADBF5D214FD4D40",
	}}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", index), func(t *testing.T) {
			var (
				p   PolygonS
				err = p.Scan([]byte(testCase.value))
			)

			assert.EqualValues(t, testCase.error, err)
			t.Logf("%v", p)
			//assert.EqualValues(t, testCase.expect, strings.ToUpper(encoded))
		})
	}
}
