package serialize

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"sync"
)

var mapIndex = sync.Map{}

func getMapHandlerFromType(t reflect.Type) readWriter {
	infoV, found := mapIndex.Load(t.String())
	if found {
		return infoV.(readWriter)
	}

	t.Key()
	info := &mapReadWriter{
		mapType: t,

		keyType:    t.Key(),
		keyHandler: getTypeHandler(t.Key()),

		valueType:    t.Elem(),
		valueHandler: getTypeHandler(t.Elem()),
	}
	mapIndex.Store(t.String(), info)

	return info
}

type mapReadWriter struct {
	variable

	mapType                  reflect.Type
	keyType, valueType       reflect.Type
	keyHandler, valueHandler readWriter
}

func (s *mapReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}

	ul := binary.BigEndian.Uint32(b)
	if ul > math.MaxInt32 {
		return fmt.Errorf("transmitted map size too large (%d>%d)", ul, math.MaxInt32)
	}
	l := int(ul)

	mp := reflect.MakeMapWithSize(s.mapType, l)

	for i := 0; i < l; i++ {
		key := reflect.Indirect(reflect.New(s.keyType))
		val := reflect.Indirect(reflect.New(s.valueType))

		handleVariableReader(r, s.keyHandler, key)
		handleVariableReader(r, s.valueHandler, val)

		mp.SetMapIndex(key, val)
	}

	v.Set(mp)
	return nil
}

func (s *mapReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v.Len()))

	if _, err := w.Write(b); err != nil {
		return err
	}

	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)

		handleVariableWriter(w, s.keyHandler, key)
		handleVariableWriter(w, s.valueHandler, val)
	}

	return nil
}
