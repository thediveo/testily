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

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getting functions done", func() {

	Context("CloseWhenGone", func() {

		It("closes the returned channel after the called function returns", func() {
			done := CloseWhenGone(func() {})
			Eventually(done).Should(BeClosed())
		})

		It("closes the returned channel when the called function panics", func() {
			panik := make(chan any, 1)
			done := closeWhenGone(func() { panic(42) }, func(r any) { panik <- r })
			Eventually(done).Should(BeClosed())
			Eventually(panik).Should(Receive(Equal(42)))
		})

	})

	Context("PassWhenGone", func() {

		It("passed the result of the concurrent function", func() {
			awaitResult := PassWhenGone(func() int { return 42 })
			Eventually(awaitResult).Should(Receive(Equal(42)))
		})

	})

})
