package ikea

import (
	"io"
	"reflect"
)

func getPointerHandlerFromType(t reflect.Type) readWriter {
	e := t.Elem()
	return &pointerWrapper{getTypeHandler(e), e}
}

type pointerWrapper struct {
	readWriter
	typ reflect.Type
}

func (p *pointerWrapper) isFixed() bool {
	return p.readWriter.isFixed()
}

func (p *pointerWrapper) vLength(v reflect.Value) (int, error) {
	return p.readWriter.(variableReadWriter).vLength(v.Elem())
}

func (p *pointerWrapper) readVariable(r io.Reader, v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(p.typ))
	}
	return p.readWriter.(variableReadWriter).readVariable(r, v.Elem())
}

func (p *pointerWrapper) writeVariable(w io.Writer, v reflect.Value) error {
	s := reflect.New(p.typ)
	v.Set(s)
	return p.readWriter.(variableReadWriter).writeVariable(w, s.Elem())
}

func (p *pointerWrapper) length() int {
	return p.readWriter.(fixedReadWriter).length()
}

func (p *pointerWrapper) readFixed(b []byte, v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.New(p.typ))
	}
	p.readWriter.(fixedReadWriter).readFixed(b, v.Elem())
}

func (p *pointerWrapper) writeFixed(b []byte, v reflect.Value) {
	s := reflect.New(p.typ)
	v.Set(s)
	p.readWriter.(fixedReadWriter).writeFixed(b, s.Elem())
}
