package binary

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

// Pack ...
func Pack(v interface{}) (output []byte, err error) {
	var buffer bytes.Buffer
	payload, _ := Marshall(v)

	encoder := NewEncoder(&buffer)
	encoder.writeUVariant(uint32(len(payload)))
	output = append(buffer.Bytes(), payload...)

	return
}

// Marshall encode the payload to binary format
func Marshall(v interface{}) (output []byte, err error) {
	var buffer bytes.Buffer
	buffer.Grow(64)

	// Encode and set the buffer if successful
	if err = MarshallTo(v, &buffer); err == nil {
		output = buffer.Bytes()
	}

	return
}

// MarshallTo ...
func MarshallTo(v interface{}, dst io.Writer) (err error) {
	encoder := NewEncoder(dst)
	tp := reflect.ValueOf(v)
	if tp.Kind() == reflect.Struct {
		encoder.writeStructValue(tp)
	} else {
		panic("not impl")
	}
	return
}

// Encoder ...
type Encoder struct {
	buffer [8]byte
	out    io.Writer
	err    error
}

// NewEncoder ...
func NewEncoder(out io.Writer) *Encoder {
	return &Encoder{
		out: out,
	}
}

// Write ...
func (e *Encoder) Write(p []byte) {
	if e.err == nil {
		_, e.err = e.out.Write(p)
	}
}

// WriteValue ...
func (e *Encoder) writeValue(tp reflect.Value) {
	switch fieldType := tp.Kind(); fieldType {
	case reflect.Struct:
		e.writeStructValue(tp)
	case reflect.Uint8:
		e.writeUint8(uint8(tp.Uint()))
	case reflect.Uint16:
		e.writeUint16(uint16(tp.Uint()))
	case reflect.Uint32:
		e.writeUint32(uint32(tp.Uint()))
	case reflect.Uint64:
		e.writeUint64(uint64(tp.Uint()))
	case reflect.Int8:
		e.writeUint8(uint8(tp.Int()))
	case reflect.Int16:
		e.writeUint16(uint16(tp.Int()))
	case reflect.Int32:
		e.writeUint32(uint32(tp.Int()))
	case reflect.Int64:
		e.writeUint64(uint64(tp.Int()))
	case reflect.Bool:
		e.writeBool(tp.Bool())
	case reflect.Array:
		// todo
	case reflect.Map:
		e.writeUVariant(uint32(tp.Len()))
		for _, key := range tp.MapKeys() {
			mapped := tp.MapIndex(key)
			e.writeValue(key)
			e.writeValue(mapped)
		}
	case reflect.Slice:
		e.writeUVariant(uint32(tp.Len()))
		for i := 0; i < tp.Len(); i++ {
			e.writeValue(tp.Index(i))
		}
	case reflect.String:
		e.writeString(tp.String())
	default:
		panic("do not support!")
	}
}

func (e *Encoder) writeStructValue(tp reflect.Value) {
	for i := 0; i < tp.NumField(); i++ {
		e.writeValue(tp.Field(i))
	}
}

func (e *Encoder) writeString(v string) {
	e.writeUVariant(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		e.Write([]byte{v[i]})
	}
}

func (e *Encoder) writeBool(v bool) {
	if v {
		e.writeUint8(1)
	} else {
		e.writeUint8(0)
	}
}

func (e *Encoder) writeUint8(v uint8) {
	e.Write([]byte{v})
}

func (e *Encoder) writeUint16(v uint16) {
	binary.LittleEndian.PutUint16(e.buffer[:], v)

	e.Write(e.buffer[:2])
}

func (e *Encoder) writeUint32(v uint32) {
	binary.LittleEndian.PutUint32(e.buffer[:], v)

	e.Write(e.buffer[:4])
}

func (e *Encoder) writeUint64(v uint64) {
	binary.LittleEndian.PutUint64(e.buffer[:], v)

	e.Write(e.buffer[:8])
}

func (e *Encoder) writeUVariant(v uint32) {
	if v > 0x8000 {
		e.writeUint16(uint16((v & 0x7fff) | 0x8000))
		e.writeUint8(uint8(v >> 15))
	} else {
		e.writeUint16(uint16(v))
	}

	binary.LittleEndian.PutUint32(e.buffer[:], v)

	e.Write(e.buffer[:4])
}
