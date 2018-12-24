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
	"os"
	"time"
)

type writer struct {
	s string
	e *execution
}

// Write implements io.Writer.Write.
func (w *writer) Write(p []byte) (int, error) {
	entry := &entry{
		TextPayload: string(p),
		Severity:    string(w.s),
		Time:        time.Now(),
		ExecutionID: string(w.e.id),

		wg: &w.e.wg,
	}

	done := loggingCtx.addEntry(entry)
	if done == nil {
		return os.Stderr.Write(entry.consoleOutput())
	}

	w.e.wg.Add(1)
	go func() {
		<-done
		w.e.wg.Done()
	}()

	return len(entry.TextPayload), nil
}
