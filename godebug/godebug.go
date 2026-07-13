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

package godebug

import (
	"os"
	"runtime/debug"
	"slices"
	"strings"
)

// Settings returns GODEBUG= settings.
func Settings() map[string]string {
	var godebugValue string
	if bi, ok := debug.ReadBuildInfo(); ok {
		idx := slices.IndexFunc(bi.Settings, func(e debug.BuildSetting) bool {
			return e.Key == "DefaultGODEBUG"
		})
		if idx >= 0 {
			godebugValue = bi.Settings[idx].Value
		}
	}
	settings := make(map[string]string)
	for setting := range strings.SplitSeq(godebugValue, ",") {
		key, val, ok := strings.Cut(strings.TrimSpace(setting), "=")
		if ok {
			settings[key] = val
		}
	}
	for setting := range strings.SplitSeq(os.Getenv("GODEBUG"), ",") {
		key, val, ok := strings.Cut(strings.TrimSpace(setting), "=")
		if ok {
			settings[key] = val
		}
	}
	return settings
}
