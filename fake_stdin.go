package lpconsole

import (
	"bufio"
	"os"
	"time"

	"github.com/thorstenrie/tserr"
)

type fakeStdin struct {
	in, r, w, o *os.File
	e           error
	d           time.Duration
}

var (
	f = &fakeStdin{o: os.Stdin}
)

/*
func infile(in *os.File) {
	Close()
	r, w, err := os.Pipe()
	if err != nil {
		r.Close()
		w.Close()
		return
	}
	fs.in, fs.r, fs.w = in, r, w
	os.Stdin = fs.r
}*/

func Restore() error {
	if f.r != nil {
		f.e = f.r.Close()
	}
	if f.w != nil {
		f.e = f.w.Close()
	}
	if f.in != nil {
		f.e = f.in.Close()
	}
	f.w, f.r, f.in = nil, nil, nil
	os.Stdin = f.o
	return f.e
}

func Delay(d time.Duration) error {
	if d < 0 {
		return tserr.Higher(&tserr.HigherArgs{Var: "d", Actual: int64(d), LowerBound: 0})
	}
	f.d = d
	return nil
}

func Err() error {
	return f.e
}

func Stdin(in *os.File) error {
	Restore()
	f.r, f.w, f.e = os.Pipe()
	if (f.e != nil) || (f.w == nil) || (f.r == nil) {
		Restore()
		return tserr.NotAvailable(&tserr.NotAvailableArgs{S: "os.Pipe", Err: f.e})
	}
	f.in = in
	os.Stdin = f.r
	go write()
	return nil
}

func write() {
	s := bufio.NewScanner(f.in)
	for s.Scan() {
		_, err := f.w.WriteString(s.Text() + "\n")
		if err != nil {
			return
		}
		time.Sleep(f.d)
	}
	f.w.Close()
}
