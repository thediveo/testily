// Copyright 2026 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stacks

import (
	"bytes"
	"iter"
	"runtime"
	"sync/atomic"
)

// All iterates over the output of multiple go routine stack dumps, yielding
// each individual go routine stack dump. The yielded stack dump comes from the
// underlying full stack dump in order to avoid unnecessary copy operations.
//
// The first element (if any) is always the stack dump of the calling go
// routine.
func All() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		dump := stacks(true)
		const dumpSeparator = "\n\n"

		for start := 0; ; {
			idx := bytes.Index(dump[start:], []byte(dumpSeparator))
			if idx < 0 {
				if start < len(dump) {
					yield(dump[start:]) // ...we will next terminate anyway.
				}
				return
			}
			// nota bene: keep a single final line end.
			if !yield(dump[start : start+idx+1]) {
				return
			}
			start += idx + len(dumpSeparator)
		}
	}
}

// Current returns the stack dump for the calling go routine.
func Current() []byte {
	return stacks(false)
}

var stackBufferSize atomic.Int64

const initialStackBufferSize = 64 * 1024

func init() {
	stackBufferSize.Store(initialStackBufferSize)
}

// stacks returns stack trace information for the current go routine if all is
// false, or for all go routines if all is true. It is a wrapper around
// [runtime.Stack] in order to hide the stack trace information buffer
// allocation shenanigans.
func stacks(all bool) []byte {
	for size := stackBufferSize.Load(); ; size *= 2 {
		stacks := make([]byte, size)
		n := runtime.Stack(stacks, all)
		if int64(n) >= size {
			continue
		}
		for {
			prevsize := stackBufferSize.Load()
			if size <= prevsize || stackBufferSize.CompareAndSwap(prevsize, size) {
				break
			}
		}
		return stacks[:n]
	}
}
