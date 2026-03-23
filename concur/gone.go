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

package concur

import "github.com/thediveo/testily/nothing"

// Nothing should actually not be stuttering.
type Nothing = nothing.Nothing

// CloseWhenGone calls the passed function in its own new goroutine, closing the
// returned (buffered) channel after the called function has returned or
// panicked.
func CloseWhenGone(fn func()) chan Nothing { return closeWhenGone(fn, nil) }

func closeWhenGone(fn func(), recoverfn func(r any)) chan Nothing {
	ch := make(chan Nothing, 1)
	go func() {
		defer func() {
			close(ch)
			if recoverfn != nil {
				recoverfn(recover())
			}
		}()
		fn()
	}()
	return ch
}

// PassWhenGone calls the passed function in its own new goroutine, passing its
// result via the returned (buffered) channel that then gets closed. If the
// called function panics, the returned channel is closed without producing any
// result.
func PassWhenGone[T any](fn func() T) chan T {
	ch := make(chan T, 1)
	go func() {
		defer close(ch)
		ch <- fn()
	}()
	return ch
}
