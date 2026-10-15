package transport

import (
	"net"
	"testing"
)

func TestTCPRoundTrip(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errc := make(chan error, 1)
	gotc := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errc <- err
			return
		}
		defer conn.Close()
		buf, err := ReadFrame(conn)
		if err != nil {
			errc <- err
			return
		}
		gotc <- string(buf)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := WriteFrame(conn, []byte("hello")); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-errc:
		t.Fatal(err)
	case got := <-gotc:
		if got != "hello" {
			t.Fatalf("got %q", got)
		}
	}
}
