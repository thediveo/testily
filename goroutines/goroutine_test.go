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

package goroutines

import (
	"bytes"
	"context"
	"runtime/pprof"
	"time"

	. "github.com/thediveo/testily/chans"
	"github.com/thediveo/testily/godebug"
	"github.com/thediveo/testily/goroutines/stacks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("parsing stack dumps for go routine details", Ordered, func() {

	It("creates a new go routine which then blocks", func() {
		done, bedone := Make[Nothing]()
		g, unblock := NewBlocked(func() {
			bedone()
		})
		Eventually(ByID).WithArguments(g.ID).
			Within(2*time.Second).ProbeEvery(1*time.Millisecond).
			Should(HaveField("State", WaitChanReceive),
				"new go routine never blocked")
		unblock()
		Eventually(done).Within(1 * time.Second).ProbeEvery(1 * time.Millisecond).
			Should(BeClosed())
	})

	Context("goroutines", func() {

		DescribeTable("goroutine header lines",
			func(line string, expected Goroutine) {
				Expect(new(line)).To(Equal(expected))
			},
			Entry(nil, "", Goroutine{}),
			Entry(nil, "goroutine", Goroutine{}),
			Entry(nil, "goroutine ", Goroutine{}),
			Entry(nil, "goroutine abc", Goroutine{}),
			Entry(nil, "goroutine 0 ", Goroutine{}),
			Entry(nil, "goroutine 123abc ", Goroutine{}),
			Entry(nil, "goroutine 123", Goroutine{}),
			Entry(nil, "goroutine 123", Goroutine{}),
			Entry(nil, "goroutine 123 foo", Goroutine{}),
			Entry(nil, "goroutine 123 [", Goroutine{}),
			Entry(nil, "goroutine 123 [the quick brown fox", Goroutine{}),
			Entry(nil, "goroutine 123 [the quick brown fox]", Goroutine{}),
			Entry(nil, "goroutine 123 [chan receive]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
			}),
			Entry(nil, "goroutine 123 [chan receive (]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked) (scan) (durable)]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
				Leaked:  true,
				Scan:    true,
				Durable: true,
			}),
			Entry(nil, "goroutine 123 [chan receive, locked to thread]", Goroutine{
				ID:           123,
				State:        WaitChanReceive,
				Waiting:      true,
				ThreadLocked: true,
			}),
			Entry(nil, "goroutine 123 [chan receive, 10]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked), synctest bubble abc]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
				Leaked:  true,
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked), foobar]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
				Leaked:  true,
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked), 10 minutes, locked to thread, synctest bubble 666]", Goroutine{
				ID:               123,
				State:            WaitChanReceive,
				Waiting:          true,
				Leaked:           true,
				ThreadLocked:     true,
				SynctestBubbleID: 666,
			}),
			Entry(nil, "goroutine 123 [chan receive labels:{\"foo\":\"bar\"}]", Goroutine{
				ID:      123,
				State:   WaitChanReceive,
				Waiting: true,
				Labels: map[string]string{
					"foo": "bar",
				},
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked), locked to thread labels:{\"foo\":\"bar\"}]", Goroutine{
				ID:           123,
				State:        WaitChanReceive,
				Waiting:      true,
				Leaked:       true,
				ThreadLocked: true,
				Labels: map[string]string{
					"foo": "bar",
				},
			}),
			Entry(nil, "goroutine 123 [chan receive (leaked), locked to thread labels:{\"foo\"]", Goroutine{
				ID:           123,
				State:        WaitChanReceive,
				Waiting:      true,
				Leaked:       true,
				ThreadLocked: true,
			}),
		)

		It("returns details about all go routines", func() {
			gos := All()
			Expect(gos).NotTo(BeEmpty())
			Expect(gos).To(HaveEach(HaveField("ID", Not(BeZero()))))
		})

		It("returns details about the caller's go routine", func() {
			Expect(Current()).To(HaveField("ID", Not(BeZero())))
		})

		It("returns details about a particular go routine", func() {
			g, bedone := NewBlocked(func() {})
			defer bedone()

			Expect(g.ID).NotTo(BeZero())
			Expect(g).NotTo(Equal(Current()))
			Eventually(ByID).WithArguments(g.ID).
				ProbeEvery(2 * time.Second).ProbeEvery(10 * time.Millisecond).
				Should(HaveField("State", WaitChanReceive))
		})

	})

	Context("goroutine labels", func() {

		It("picks up goroutine labels", func() {
			if godebug.Settings()["tracebacklabels"] != "1" {
				Skip("stack dumps with labels require go toolchain 1.26+")
			}

			ch, done := Make[Nothing]()
			defer done()

			go func() {
				pprof.Do(context.Background(), pprof.Labels("foo", "barz"), func(context.Context) {
					<-ch // wait...
				})
			}()

			Eventually(func() map[string]string {
				GinkgoWriter.Println("----")
				for dump := range stacks.All() {
					line, _, _ := bytes.Cut(dump, []byte("\n"))
					GinkgoWriter.Println(string(line))
					labels := parseLabels(string(line))
					if labels != nil {
						return labels
					}
				}
				return nil
			}).Within(1 * time.Second).ProbeEvery(100 * time.Millisecond).
				Should(HaveKeyWithValue("foo", "barz"))
		})

	})

})
