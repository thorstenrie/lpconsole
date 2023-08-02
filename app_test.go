// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package lpconsole_test

import (
	"context"
	"fmt"
)

func stop(ctx context.Context, args []string) error {
	fmt.Println("Stopping application")
	return nil
}

func foo(ctx context.Context, args []string) error {
	for _, s := range args {
		fmt.Println(s)
	}
	return nil
}
