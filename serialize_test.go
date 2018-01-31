package serialize

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

/* Test types */

type testStruct struct {
	A uint32
	B int64
	C []uint16
	D string
	E []uint16 `compressed:"true"`
	F map[uint8]string

	G testSubStruct
	H testInterface
}

type testSubStruct struct {
	A byte
}

type testInterface struct {
	A int64
}

func (t *testInterface) Deserialize(r io.Reader) error {
	return Read(r, &t.A)
}

func (t *testInterface) Serialize(w io.Writer) error {
	return Write(w, &t.A)
}

/* Test data */
var source = &testStruct{1, 2, []uint16{3, 4, 5, 6}, "test", compressData, map[uint8]string{7: "seven", 8: "eight"}, testSubStruct{9}, testInterface{10}}
var testData, _ = hex.DecodeString("0000000100000000000000020000000400030004000500060000000474657374000000ddecc1411500000405b0af8ea072882a83fbde52b36900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0eb020000ffff000000020700000005736576656e0800000005656967687409000000000000000a")
var compressData = make([]uint16, 100000)

func init() {
	for i := range compressData {
		compressData[i] = 42
	}
}

/* Tests */

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Write(buf, source); err != nil {
		t.Error(err)
		return
	}

	result := buf.Bytes()
	if len(result) != len(testData) {
		fmt.Printf("Failing TestWrite, result \"%s\" length (%d) does not match test data length (%d)\n", hex.EncodeToString(result), len(result), len(testData))
		t.FailNow()
	}

	for i := 0; i < len(testData); i++ {
		if result[i] != testData[i] {
			fmt.Printf("Failing TestWrite, hex output \"%s\" does not match test data slice\n", hex.EncodeToString(result))
			t.FailNow()
		}
	}
}

func TestRead(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(testData)

	tst := new(testStruct)
	if err := Read(buf, tst); err != nil {
		t.Error(err)
		return
	}

	compare(t, "A", tst.A, source.A)
	compare(t, "B", tst.B, source.B)
	compare(t, "len(C)", len(tst.C), len(source.C))
	for i := range source.C {
		compare(t, fmt.Sprintf("C[%d]", i), tst.C[i], source.C[i])
	}
	compare(t, "D", tst.D, source.D)
	compare(t, "len(E)", len(tst.E), len(source.E))
	for i := range source.E {
		compare(t, fmt.Sprintf("E[%d]", i), tst.E[i], source.E[i])
	}
	for k := range source.F {
		compare(t, fmt.Sprintf("F[%d]", k), tst.F[k], source.F[k])
	}
	compare(t, "G (sub-struct)", tst.G.A, source.G.A)
	compare(t, "H (interface)", tst.H.A, source.H.A)
}

func compare(t *testing.T, field string, value1, value2 interface{}) {
	if value1 != value2 {
		fmt.Printf("Failing TestRead, decoded data field '%s' with value '%v' does not match '%v'", field, value1, value2)
		t.FailNow()
	}
}
