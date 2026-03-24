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

package concur_test

import (
	"fmt"
	"time"

	. "github.com/thediveo/testily/concur"
	. "github.com/thediveo/testily/tuples"
)

func ExampleCloseWhenGone() {
	ch := CloseWhenGone(func() { time.Sleep(10 * time.Millisecond) })
	select {
	case <-ch:
		fmt.Printf("done\n")
	case <-time.After(1 * time.Second):
		fmt.Printf("error: timeout\n")
	}
	// Output: done
}

func ExamplePassWhenGone() {
	ch := PassWhenGone(
		func() Pair[string, int] {
			return PackPair("meaning of life", 42)
		})
	select {
	case qa := <-ch:
		q, a := qa.Unpack()
		fmt.Printf("question %v, answer %v", q, a)
	case <-time.After(1 * time.Second):
		fmt.Printf("error: timeout\n")
	}
	// Output: question meaning of life, answer 42
}
