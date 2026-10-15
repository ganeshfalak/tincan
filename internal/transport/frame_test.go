package transport

import (
	"bytes"
	"testing"
)

func TestWriteFrameHello(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrame(&buf, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	want := []byte{0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o'}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("got %v", buf.Bytes())
	}
}

func TestFrameRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrame(&buf, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFrame(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestReadFrameTooLarge(t *testing.T) {
	// 1048577 is MaxFrame + 1, written as four big-endian bytes.
	_, err := ReadFrame(bytes.NewReader([]byte{0x00, 0x10, 0x00, 0x01}))
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadFrameTruncated(t *testing.T) {
	_, err := ReadFrame(bytes.NewReader([]byte{0, 0, 0, 5, 'h', 'e'}))
	if err == nil {
		t.Fatal("expected an error")
	}
}
