package clamd

import (
	"io/ioutil"
	"bytes"
	"io"
	"os"
	"testing"
)

var (
	clamdSocket = "/tmp/clamd.socket"
)

func init() {
	if s := os.Getenv("CLAMD_SOCKET"); s != "" {
		socket = clamdSocket
	}
}

func TestScanStreamMemoryLeak(t *testing.T) {
	c := NewClamd(clamdSocket)
	const maxLoop = 20000

	for i := 0; i < maxLoop; i++ {
		abortCh := make(chan bool)
		buf := bytes.NewBuffer(EICAR)
		ch, err := c.ScanStream(buf, abortCh)
		if err != nil {
			t.Fatal(err)
		}
		r := <-ch
		if r.Status!=RES_OK {
			t.Fatal("status", r)
		}
	}
}

func TestScanFileMemoryLeak(t *testing.T) {
	c := NewClamd(clamdSocket)
	const maxsocket = 20000
	f,err := os.CreateTemp(".", "test-*"))
	if err!=nil {
		t.Fatal(err)
	}
	f.Write(EICAR)
	f.Close()
	defer os.Remove(f.Name())
	
	for i := 0; ; i++ {
		ch,err := c.ScanFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
				r := <-ch
		if r.Status!=RES_OK {
			t.Fatal("status", r)
		}
	}
}
