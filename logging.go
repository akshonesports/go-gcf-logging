// Copyright 2018 Akshon Media Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/functions/metadata"
)

type contextKey struct{}

type execution struct {
	wg sync.WaitGroup
	id string
}

func WithContext(ctx context.Context) (context.Context, error) {
	var e execution

	md, err := metadata.FromContext(ctx)
	if err == nil {
		e.id = md.EventID
	}

	return context.WithValue(ctx, contextKey{}, &e), nil
}

func Wait(ctx context.Context) {
	e, ok := ctx.Value(contextKey{}).(*execution)
	if !ok {
		return
	}

	e.wg.Wait()
}

type Severity string

const (
	Default   Severity = "DEFAULT"
	Debug     Severity = "DEBUG"
	Info      Severity = "INFO"
	Notice    Severity = "NOTICE"
	Warning   Severity = "WARNING"
	Error     Severity = "ERROR"
	Critical  Severity = "CRITICAL"
	Alert     Severity = "ALERT"
	Emergency Severity = "EMERGENCY"
)

func Logger(ctx context.Context, severity Severity, flags int) *log.Logger {
	w := writer{s: string(severity)}

	e, ok := ctx.Value(contextKey{}).(*execution)
	if ok {
		w.e = e
	}

	return log.New(&w, "", flags)
}
