package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BUILDING CUSTOM ERROR TYPE

type errKind int

// creating enumerated type
const (
	_ errKind = iota
	noHeader
	cantReadHeader
	invalidHdrType
	invalidChkLength
)

type WaveError struct {
	kind  errKind
	value int
	err   error
}

func (e WaveError) Error() string {
	switch e.kind {
	case noHeader:
		return "no header (file too short?)"
	case cantReadHeader:
		return fmt.Sprintf("can't read header[%d]: %s", e.value, e.err.Error())
	case invalidHdrType:
		return "invalid header type"
	case invalidChkLength:
		return fmt.Sprintf("invalid chunk lengthL %d", e.value)
	default:
		return "error"
	}
}

// returns error with a particular value (e.g., header type)
func (e WaveError) with(val int) WaveError {
	e1 := e
	e1.value = val
	return e1
}

// from returns an error with a particular location and underlying error (e.g., from the standard library)
func (e WaveError) from(pos int, err error) WaveError {
	e1 := e
	e1.value = pos
	e1.err = err
	return e1
}

// we can have prototype errors we can return or customize
var (
	HeaderMissing      = WaveError{kind: noHeader}
	HeaderReadFailed   = WaveError{kind: cantReadHeader}
	InvalidHeaderType  = WaveError{kind: invalidHdrType}
	InvalidChunkLength = WaveError{kind: invalidChkLength}
)

type Header struct {
	TotalLength uint32
	riff        string
}

// example of using errors
func DecodeHeader(b []byte) (*Header, []byte, error) {
	var err error
	var pos int

	header := Header{TotalLength: uint32(len(b))}
	buf := bytes.NewReader(b)

	if len(b) < 10 {
		return &header, nil, HeaderMissing // if you don't have enough bytes, you can return one type of error
	}

	if err = binary.Read(buf, binary.BigEndian, &header.riff); err != nil { // here we start trying to decode parts of the input
		return &header, nil, HeaderReadFailed.from(pos, err) // we may get an error that is more complicated
	} // using original error variable as prototype and modifying to put particular details of the error occurring here
	// then we return that new error (of same type) with more details

	return &header, b, nil
}

// wrapped errors
// starting with Go 1.13, we can wrap one error in another
type HAL9009 struct {
	victim string
	err    error
}

// the %w format verb can wrap one error in another --> can allow for error chains
// top-level error --> intermediate error --> original error
func (h HAL9009) OpenPodBayDoors() error {
	if h.err != nil {
		return fmt.Errorf("I'm sorry %s, I can't: %w", h.victim, h.err)
	}
	return nil
}

// wrapping errors means we can also unwrap errors; here we are just returning the embedded error
func (w *WaveError) Unwrap() error {
	return w.err
}

// errors.Is can be used to walk down error chain and compare with error variable to see if a certain error is underlying cause
// Example:
// ...
// if audio, err = DecodeWaveFile(fn); err != nil {
// 	if errors.Is(err, os.ErrPermission) { // much safer way of checking for a certain error
// 		// report security violation
// 	}
// }

// can add Is to custom error, so check if a certain error is of type WaveError when errors.Is is called
// will be called using a WaveError type receiver
func (w *WaveError) Is(t error) bool {
	e, ok := t.(*WaveError) // reflection here
	if !ok {
		return false
	}
	return e.kind == w.kind
}

// errors.As looks for an error type not a value -< we can get an error of an underlying type if it's in the chain
// Example:
// ... code below is basically attempting to extract PathError from error
// if audio, err = DecodeWaveFile(fn); err != nil {
// 	var e os.PathError // a struct
// 	if errors.As(err, &e) { // much safer way of checking for a certain error
// 		// let's just pass back the underlying file error
// 		return e
// 	}
// }

// panic and recover in Go --> not really ideal to do, but can be useful in unit testing (sometimes)
func abc() {
	panic("omg")
}

// recovery of panic only works inside defer
func main() {
	defer func() {
		if p := recover(); p != nil {
			// can't really do much else than print p
			fmt.Println("recover:", p)
		}
	}()
	abc()
}
