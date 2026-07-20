// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"context"
	"net/http"
	"time"

	"github.com/bborbe/log"
)

func CreateSetLoglevelHandler(ctx context.Context) http.Handler {
	return log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute))
}
