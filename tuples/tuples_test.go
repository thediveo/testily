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

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tuples", func() {

	It("packs and unpacks pairs", func() {
		pair := PackPair(42, "abc")
		e1, e2 := pair.Unpack()
		Expect(e1).To(Equal(42))
		Expect(e2).To(Equal("abc"))
		Expect(pair).To(Equal(PackPair(42, "abc")))
	})

	It("packs and unpacks triples", func() {
		triple := PackTriple(42, "abc", 66.6)
		e1, e2, e3 := triple.Unpack()
		Expect(e1).To(Equal(42))
		Expect(e2).To(Equal("abc"))
		Expect(e3).To(Equal(66.6))
	})

})
