// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package lpconsole_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/thorstenrie/lpconsole"
	"github.com/thorstenrie/tsfio"
)

func TestApp(t *testing.T) {
	lpconsole.Usage(&lpconsole.Help{App: "Test", Description: "A test app", Version: "1.0.0"})
	lpconsole.Add(&lpconsole.Command{Key: "foo", Help: "Foo", Function: foo})
	lpconsole.Add(&lpconsole.Command{Key: "stop", Help: "Stop", Function: stop})
	lpconsole.Exit("stop")
	f, _ := tsfio.OpenFile(tsfio.Filename("testdata/stdin.txt"))
	lpconsole.SetInput(f)
	ctx := context.Background()
	lpconsole.Run(ctx)
}

func TestFakeStdin(t *testing.T) {
	fs, _ := tsfio.OpenFile(tsfio.Filename("testdata/stdin.txt"))
	lpconsole.Delay(time.Second)
	lpconsole.Stdin(fs)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	lpconsole.Restore()
}
