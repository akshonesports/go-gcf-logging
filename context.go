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
	"fmt"
	"os"
	"sync"
)

type loggingContext struct {
	initOnce sync.Once

	queueMutex   sync.Mutex
	queue        chan *batch
	currentBatch *batch
}

// startNewBatch prepares a new batch.
//
// Note: startNewBatch is not thread safe.
func (c *loggingContext) startNewBatch() *batch {
	c.currentBatch = &batch{
		ready: make(chan struct{}),
		done:  make(chan struct{}),
	}
	c.queue <- c.currentBatch
	return c.currentBatch
}

func (c *loggingContext) addEntry(entry *entry) chan struct{} {
	if c.queue == nil {
		return nil
	}

	c.queueMutex.Lock()
	defer c.queueMutex.Unlock()

	// Start a new batch if the current one would grow too much.
	if len(c.currentBatch.Entries) > 0 &&
		(len(c.currentBatch.Entries)+1 > maxLogBatchEntries ||
			c.currentBatch.length+len(entry.TextPayload) > maxLogBatchLength) {
		c.startNewBatch()
	}

	c.currentBatch.addEntry(entry)

	return c.currentBatch.ready
}

func (c *loggingContext) startReportWorker() {
	for logBatch := range c.queue {
		<-logBatch.ready

		c.queueMutex.Lock()
		if logBatch == c.currentBatch {
			c.startNewBatch()
		}
		c.queueMutex.Unlock()

		if err := logBatch.report(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			killInstance()
		}
	}
}

func (c *loggingContext) initialize() {
	c.initOnce.Do(func() {
		c.queue = make(chan *batch, 5)
		c.startNewBatch()
		go c.startReportWorker()
	})
}

var loggingCtx loggingContext
