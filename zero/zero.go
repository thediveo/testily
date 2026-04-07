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

/*
Package zero provides types zero values.
*/
package zero

// Value returns a zero value of the specified type. This is the preferred form
// avoiding stuttering when importing normally, as opposed to dot-importing.
func Value[T any]() T {
	var zero T
	return zero
}

// Zero returns a zero value of the specified type. This is the preferred form
// when dot-importing.
func Zero[T any]() T { return Value[T]() }
