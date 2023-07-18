// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package lpconsole_test

import (
	"context"
	"testing"

	"github.com/thorstenrie/lpconsole"
)

func TestApp(t *testing.T) {
	lpconsole.Usage(&lpconsole.Help{App: "Test", Description: "A test app", Version: "1.0.0"})
	ctx := context.Background()
	lpconsole.Run(ctx)
}
