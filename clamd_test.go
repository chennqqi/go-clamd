package clamd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var (
	clamdSocket = "/tmp/clamd.socket"
	tmpdir      = "/tmp"
)

func init() {
	if s := os.Getenv("CLAMD_SOCKET"); s != "" {
		clamdSocket = s
	}
	if s := os.Getenv("CLAMD_TEMP"); s != "" {
		tmpdir = s
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
			t.Fatal(i, ":", err)
		}
		r := <-ch
		if r.Status != RES_FOUND {
			t.Fatal("status", i, ":", r)
		}
	}
}

func TestScanFileMemoryLeak(t *testing.T) {
	c := NewClamd(clamdSocket)
	const maxLoop = 20000

	f, err := os.CreateTemp(tmpdir, "eicar_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(EICAR)
	f.Close()
	defer os.Remove(f.Name())

	fmt.Println("test", f.Name())
	for i := 0; i < maxLoop; i++ {
		ch, err := c.ScanFile(f.Name())
		if err != nil {
			t.Fatal("ScanFile", err)
		}
		r := <-ch
		if r == nil {
			t.Fatal("result of", i, ":", f.Name(), "is nil")
		}
		if r.Status != RES_FOUND {
			t.Fatal("status", i, ":", f.Name(), r)
		}
	}
}
