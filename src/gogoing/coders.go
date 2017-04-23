package gogoing

import (
	"io"
)

type Encoder interface {
	Encode(v *Event) (data []byte, err error)
}

type Decoder interface {
	Decode(data []byte) (v *Event, err error)
}


type DefaultEncoder struct {
	w io.Writer
}

func NewDefaultEncoder(w io.Writer) *DefaultEncoder {
	return &DefaultEncoder{w: w}
}

func (e *DefaultEncoder) Encode(v *Event) (data []byte, err error) {
	return
}

type DefaultDecoder struct {
	r io.Reader
}

func NewDefaultDecoder(r io.Reader) *DefaultDecoder {
	return &DefaultDecoder{r : r}
}

func (d *DefaultDecoder) Decode(data []byte) (v *Event, err error) {
	return
}