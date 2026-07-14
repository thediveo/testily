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

package godebug_test

import (
	"os"
	"strings"

	"github.com/thediveo/testily/godebug"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("godebug settings", func() {

	It("fetches (some) godebug settings", func() {
		Expect(os.Setenv("GODEBUG", strings.Join(
			append(strings.Split(os.Getenv("GODEBUG"), ","), "frobatz=123"), ","))).
			To(Succeed())
		godebugs := godebug.Settings()
		Expect(godebugs).NotTo(BeEmpty())
		Expect(godebugs).To(HaveKeyWithValue("frobatz", "123"))
	})

})
