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
	"bytes"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	TextPayload string
	Severity    string
	Time        time.Time
	ExecutionID string

	wg *sync.WaitGroup
}

func (e *entry) consoleOutput() []byte {
	var logBuf bytes.Buffer
	fmt.Fprintf(&logBuf, "[%s]", e.Severity[:1])
	if e.ExecutionID != "" {
		fmt.Fprintf(&logBuf, "[%s]", e.ExecutionID)
	}

	logBuf.WriteByte(' ')
	logBuf.WriteString(e.TextPayload)
	if len(e.TextPayload) == 0 || e.TextPayload[len(e.TextPayload)-1] != '\n' {
		logBuf.WriteByte('\n')
	}

	return logBuf.Bytes()
}

type batch struct {
	Entries []*entry

	length int

	ready chan struct{}
	done  chan struct{}
}

// addEntry adds a log entry to the batch.
//
// Note: addEntry is not thread safe.
func (b *batch) addEntry(entry *entry) {
	if b.Entries == nil {
		close(b.ready)
	}

	b.Entries = append(b.Entries, entry)
	b.length += len(entry.TextPayload)
}

func (b *batch) report() error {
	if len(b.Entries) == 0 {
		return nil
	}

	defer close(b.done)

	if err := postToSupervisor("/_ah/log", b, supervisorLogTimeout); err != nil {
		return err
	}

	return nil
}