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

package chans

import (
	"sync"

	"github.com/thediveo/testily/nothing"
)

// CloseableChannel is any channel transporting values of type T that can be
// closed. Please note that read channels cannot be closed, only bidirectional
// and send channels.
type CloseableChannel[T any] interface {
	~chan T | ~chan<- T
}

// CloseOnce returns a function that ensures to close the passed
// CloseableChannel only once.
func CloseOnce[T any, C CloseableChannel[T]](ch C) func() {
	return sync.OnceFunc(func() { close(ch) })
}

// Nothing shouldn't be stuttering.
type Nothing = nothing.Nothing

// Make returns a channel for values of type T together with a closer function
// that when called multiple times will close the returned channel only at its
// first call.
//
//	ch, closer := Make[foo.Sprockets]()
//
// In case the channel should act as a signaller without transporting values
// other than the zero value after closing, use
//
//	ch, closer := Make[Nothing]()
func Make[T any]() (ch chan T, closer func()) {
	ch = make(chan T)
	return ch, CloseOnce(ch)
}
