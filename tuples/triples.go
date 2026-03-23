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

// Tuple3 represents a 3-tuple of comparable values.
type Tuple3[T1, T2, T3 comparable] struct {
	e1 T1
	e2 T2
	e3 T3
}

// Triple is an alias for a 3-tuple.
type Triple[T1, T2, T3 comparable] = Tuple3[T1, T2, T3]

// PackTriple packs three values into a 3-tuple.
func PackTriple[T1, T2, T3 comparable](v1 T1, v2 T2, v3 T3) Triple[T1, T2, T3] {
	return Pack3(v1, v2, v3)
}

// Pack three values into a 3-tuple.
func Pack3[T1, T2, T3 comparable](v1 T1, v2 T2, v3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{
		e1: v1,
		e2: v2,
		e3: v3,
	}
}

// Unpack the values of the elements of a 3-tuple.
func (t Tuple3[T1, T2, T3]) Unpack() (T1, T2, T3) {
	return t.e1, t.e2, t.e3
}
