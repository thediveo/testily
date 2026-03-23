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

package tuples

// Tuple2 represents a 2-tuple of comparable values.
type Tuple2[T1, T2 comparable] struct {
	e1 T1
	e2 T2
}

// Pair is an alias for a 2-tuple.
type Pair[T1, T2 comparable] = Tuple2[T1, T2]

// PackPair packs two values into a 2-tuple.
func PackPair[T1, T2 comparable](v1 T1, v2 T2) Pair[T1, T2] { return Pack2(v1, v2) }

// Pack2 packs two values into a 2-tuple.
func Pack2[T1, T2 comparable](v1 T1, v2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{
		e1: v1,
		e2: v2,
	}
}

// Unpack the two values of the elements of a 2-tuple.
func (t Tuple2[T1, T2]) Unpack() (T1, T2) {
	return t.e1, t.e2
}
