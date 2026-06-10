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

package close

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

// Once returns a function that ensures to close the passed CloseableChannel
// only once.
func Once[T any, C CloseableChannel[T]](ch C) func() {
	return sync.OnceFunc(func() { close(ch) })
}

// Nothing shouldn't be stuttering.
type Nothing = nothing.Nothing

// NewChan returns a channel for values of type T together with a closer
// function that when called multiple times will close the returned channel only
// at its first call.
//
//	ch, closer := NewChan[foo.Sprockets]()
//
// In case the channel should act as a signaller without
// transporting values other than the zero value after closing, use
//
//	ch, closer := NewChan[Nothing]()
func NewChan[T any]() (ch chan T, closer func()) {
	ch = make(chan T)
	return ch, Once(ch)
}
