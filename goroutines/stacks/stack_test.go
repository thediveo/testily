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
	"slices"
	"sync"

	. "github.com/thediveo/testily/chans"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("go routine stack dumps", func() {

	When("iterating over goroutine stack dumps", func() {

		It("yields individual stack dumps", func() {
			dumps := slices.Collect(All())
			Expect(len(dumps)).To(BeNumerically(">", 1))
			Expect(dumps).To(HaveEach(Not(HaveSuffix("\n\n"))))
			Expect(dumps).To(HaveEach(HavePrefix("goroutine ")))
		})

		It("aborts iterating over stack dumps", func() {
			count := 0
			for range All() {
				count++
				break
			}
			Expect(count).To(Equal(1))
		})

	})

	When("retrieving the current go routine's stack dump", func() {

		It("ends in a single newline", func() {
			dump := Current()
			Expect(dump).NotTo(BeEmpty())
			Expect(dump).To(HaveSuffix("\n"))
			Expect(dump).NotTo(HaveSuffix("\n\n"))
		})

		It("returns rather long stack traces", func() {

			ch, done := Make[Nothing]()
			defer done()

			var wg sync.WaitGroup

			var dumpfn func(level int)
			dumpfn = func(level int) {
				level--
				if level > 0 {
					dumpfn(level)
					return
				}
				wg.Done()
				<-ch
			}
			const times = 100
			wg.Add(times)
			for range times {
				go dumpfn(1_000)
			}
			wg.Wait()
			var buff bytes.Buffer
			for dump := range All() {
				buff.Write(dump)
			}
			dumps := buff.Bytes()
			GinkgoWriter.Printf("dump size: %d\n", len(dumps))
			Expect(len(dumps)).To(BeNumerically(">", initialStackBufferSize))
			Expect(dumps).To(HaveSuffix("\n"), "unexpected truncation")
		})

	})

})
