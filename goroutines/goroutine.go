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
	"encoding/json"
	"strconv"
	"strings"

	"github.com/thediveo/testily/chans"
	"github.com/thediveo/testily/goroutines/stacks"
)

type State uint

const (
	NoIdea State = iota

	Runnable
	Running
	Syscall
	Dead
	Copystack
	Leaked
	Preempted
	DeadExtra

	WaitGCAssistMarking
	WaitIOWait
	WaitChanReceiveNilChan
	WaitChanSendNilChan
	WaitDumpingHeap
	WaitGarbageCollection
	WaitGarbageCollectionScan
	WaitPanicWait
	WaitSelect
	WaitSelectNoCases
	WaitGCAssistWait
	WaitGCSweepWait
	WaitGCScavengeWait
	WaitChanReceive
	WaitChanSend
	WaitFinalizer
	WaitForceGCIdle
	WaitUpdateGOMAXPROCSIdle
	WaitSemacquire
	WaitSleep
	WaitSyncCondWait
	WaitSyncMutexLock
	WaitSyncRWMutexRLock
	WaitSyncRWMutexLock
	WaitSyncWaitGroupWait
	WaitTraceReaderBlocked
	WaitForGCCycle
	WaitGCWorkerIdle
	WaitGCWorkerActive
	WaitPreempted
	WaitDebugCall
	WaitGCMarkTermination
	WaitStoppingTheWorld
	WaitFlushProcCaches
	WaitTraceGoroutineStatus
	WaitTraceProcStatus
	WaitPageTraceFlush
	WaitCoroutine
	WaitGCWeakToStrongWait
	WaitSynctestRun
	WaitSynctestWait
	WaitSynctestChanReceive
	WaitSynctestChanSend
	WaitSynctestSelect
	WaitSynctestWaitGroupWait
	WaitCleanupWait
)

// Go routine status strings, including waiting reasons that are represented as
// go routine states.
//
// See also [gStatusStrings] and [waitReasonStrings] for the set of status
// strings defined as of Go 1.26.5.
//
// [gStatusStrings]: https://cs.opensource.google/go/go/+/refs/tags/go1.26.5:src/runtime/traceback.go;l=1202
// [waitReasonStrings]: https://cs.opensource.google/go/go/+/master:src/runtime/runtime2.go;l=1273
var stateDictionary = map[State]string{
	NoIdea: "???",

	Runnable:  "runnable",
	Running:   "running",
	Syscall:   "syscall",
	Dead:      "dead",
	Copystack: "copystack",
	Leaked:    "leaked",
	Preempted: "preempted",
	DeadExtra: "waiting for cgo callback",

	WaitGCAssistMarking:       "GC assist marking",
	WaitIOWait:                "IO wait",
	WaitChanReceiveNilChan:    "chan receive (nil chan)",
	WaitChanSendNilChan:       "chan send (nil chan)",
	WaitDumpingHeap:           "dumping heap",
	WaitGarbageCollection:     "garbage collection",
	WaitGarbageCollectionScan: "garbage collection scan",
	WaitPanicWait:             "panicwait",
	WaitSelect:                "select",
	WaitSelectNoCases:         "select (no cases)",
	WaitGCAssistWait:          "GC assist wait",
	WaitGCSweepWait:           "GC sweep wait",
	WaitGCScavengeWait:        "GC scavenge wait",
	WaitChanReceive:           "chan receive",
	WaitChanSend:              "chan send",
	WaitFinalizer:             "finalizer wait",
	WaitForceGCIdle:           "force gc (idle)",
	WaitUpdateGOMAXPROCSIdle:  "GOMAXPROCS updater (idle)",
	WaitSemacquire:            "semacquire",
	WaitSleep:                 "sleep",
	WaitSyncCondWait:          "sync.Cond.Wait",
	WaitSyncMutexLock:         "sync.Mutex.Lock",
	WaitSyncRWMutexRLock:      "sync.RWMutex.RLock",
	WaitSyncRWMutexLock:       "sync.RWMutex.Lock",
	WaitSyncWaitGroupWait:     "sync.WaitGroup.Wait",
	WaitTraceReaderBlocked:    "trace reader (blocked)",
	WaitForGCCycle:            "wait for GC cycle",
	WaitGCWorkerIdle:          "GC worker (idle)",
	WaitGCWorkerActive:        "GC worker (active)",
	WaitPreempted:             "preempted",
	WaitDebugCall:             "debug call",
	WaitGCMarkTermination:     "GC mark termination",
	WaitStoppingTheWorld:      "stopping the world",
	WaitFlushProcCaches:       "flushing proc caches",
	WaitTraceGoroutineStatus:  "trace goroutine status",
	WaitTraceProcStatus:       "trace proc status",
	WaitPageTraceFlush:        "page trace flush",
	WaitCoroutine:             "coroutine",
	WaitGCWeakToStrongWait:    "GC weak to strong wait",
	WaitSynctestRun:           "synctest.Run",
	WaitSynctestWait:          "synctest.Wait",
	WaitSynctestChanReceive:   "chan receive (durable)",
	WaitSynctestChanSend:      "chan send (durable)",
	WaitSynctestSelect:        "select (durable)",
	WaitSynctestWaitGroupWait: "sync.WaitGroup.Wait (durable)",
	WaitCleanupWait:           "cleanup wait",
}

// Goroutine describes various properties of a particular go routine.
type Goroutine struct {
	ID               uint64 // unique go routine ID ("goid" in Go's runtime parlance).
	State            State  // go routine state, such as "running", et cetera.
	Waiting          bool
	Scan             bool
	Durable          bool
	Leaked           bool
	ThreadLocked     bool
	SynctestBubbleID uint64
	Labels           map[string]string
}

var stateByName = func() map[string]State {
	m := make(map[string]State, len(stateDictionary))
	for state, text := range stateDictionary {
		m[text] = state
	}
	return m
}()

// All returns information about all go routines.
func All() []Goroutine {
	var gos []Goroutine
	for stack := range stacks.All() {
		header, _, ok := bytes.Cut(stack, []byte("\n"))
		if !ok {
			continue
		}
		g := new(string(header))
		if g.ID == 0 {
			continue
		}
		gos = append(gos, g)
	}
	return gos
}

// Current returns information about the caller's go routine.
func Current() Goroutine {
	header, _, ok := bytes.Cut(stacks.Current(), []byte("\n"))
	if !ok {
		return Goroutine{}
	}
	return new(string(header))
}

// ByID returns the details of go routine with the specified ID, otherwise zero
// value details.
func ByID(id uint64) Goroutine {
	for _, g := range All() {
		if g.ID != 0 && g.ID == id {
			return g
		}
	}
	return Goroutine{}
}

// New runs fn in a new go routine and immediately returns the new go routine's
// details.
func New(fn func()) Goroutine {
	ch := make(chan Goroutine)
	go func() {
		defer close(ch)
		ch <- Current()
		fn()
	}()
	return <-ch
}

// NewBlocked creates a new go routine which then blocks until the returned
// unblock function is called and only then calls fn.
func NewBlocked(fn func()) (g Goroutine, unblock func()) {
	gch := make(chan Goroutine)
	unblockch, unblock := chans.Make[chans.Nothing]()
	go func() {
		defer close(gch)
		gch <- Current()
		<-unblockch
		fn()
	}()
	g = <-gch
	return g, unblock
}

// new parses the specified go routine stack dump header line, returning
// Goroutine details; it returns a zero value Goroutine if the header line
// cannot be parsed at least in part.
func new(line string) Goroutine {
	// nota bene (https://cs.opensource.google/go/go/+/refs/tags/go1.26.5:src/runtime/traceback.go;l=1215):
	line = strings.TrimSuffix(line, ":\n")

	// First, parse the "goroutine" intro as well as the following go routine
	// ID.
	line, ok := strings.CutPrefix(line, "goroutine ")
	if !ok {
		return Goroutine{}
	}

	idstr, line, ok := strings.Cut(line, " ")
	if !ok {
		return Goroutine{}
	}
	gid, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil || gid == 0 {
		return Goroutine{}
	}

	// We skip any runtime information that might follow and go straight for the
	// "[...]" bracketed part.
	_, line, ok = strings.Cut(line, "[")
	if !ok {
		return Goroutine{}
	}
	// ...and now for the nasty part: the status string ... that can contain
	// multiple words separated by spaces. The end of the status string is
	// determined by:
	//  - an opening paranthesis "(",
	//  - a comma "," (followed by a space),
	//  - a space followed by "labels:",
	//  - the ending "]".
	// Now, "]" might pop up in quoted label keys and values, so we should check
	// only for it if we can't find anything else earlier.
	idx := 0
gobble:
	for ; idx < len(line); idx++ {
		switch line[idx] {
		case ',', '(', ']':
			break gobble
		case ' ':
			switch {
			case strings.HasPrefix(line[idx:], " labels:"):
				break gobble
			case strings.HasPrefix(line[idx:], " ("):
				break gobble
			}
		}
	}
	state, ok := stateByName[line[:idx]]
	if !ok {
		return Goroutine{}
	}
	line = line[idx:]

	var leaked, scan, durable, locked bool
	var synctestBubbleid uint64
	var labels map[string]string
insidebrackets:
	for {
		if strings.HasPrefix(line, " (") {
			var bracketed string
			bracketed, line, ok = strings.Cut(line[2:], ")")
			if !ok {
				// no idea what's going on, so stopping parsing and going with what we picked up so far.
				break insidebrackets
			}
			switch bracketed {
			case "leaked":
				leaked = true
			case "scan":
				scan = true
			case "durable":
				durable = true
			default:
				// silently swallow, we don't know what it means.
			}
			continue insidebrackets
		}

		if strings.HasPrefix(line, ", ") {
			const (
				lockedToThread = "locked to thread"
				synctestBubble = "synctest bubble "
			)
			line = line[2:]
			switch {
			case strings.HasPrefix(line, lockedToThread):
				locked = true
				line = line[len(lockedToThread):]
				continue
			case strings.HasPrefix(line, synctestBubble):
				line = line[len(synctestBubble):]
				idx := 0
				for ; idx < len(line); idx++ {
					if line[idx] < '0' || line[idx] > '9' {
						break
					}
				}
				synctestBubbleid, err = strconv.ParseUint(line[:idx], 10, 64)
				if err != nil {
					// no idea what's going on, so stopping parsing and going with what we picked up so far.
					break insidebrackets
				}
				line = line[idx:]
				continue
			case len(line) > 0 && line[0] >= '0' && line[0] <= '9':
				_, line, ok = strings.Cut(line, "minutes")
				if !ok {
					// no idea what's going on, so stopping parsing and going with what we picked up so far.
					break insidebrackets
				}
				// we just swallow
				continue
			}
			break insidebrackets
		}

		if strings.HasPrefix(line, " labels:") {
			labels = parseLabels(line)
			break insidebrackets
		}

		// no idea what's going on, so stopping parsing and going with what we picked up so far.
		break insidebrackets
	}

	return Goroutine{
		ID:               gid,
		State:            state,
		Waiting:          state >= WaitGCAssistMarking,
		Scan:             scan,
		Durable:          durable,
		Leaked:           leaked,
		ThreadLocked:     locked,
		SynctestBubbleID: synctestBubbleid,
		Labels:           labels,
	}
}

// parseLabels returns the key-value labels attached to a go routine from its
// stackdump "goroutine" line.
func parseLabels(line string) map[string]string {
	// according to runtime.goroutineheader
	// (https://cs.opensource.google/go/go/+/refs/tags/go1.26.5:src/runtime/traceback.go;l=1215)
	// there aren't any "labels:" strings before the trailing optional labels
	// key-value map; thus, we can get away with a simple forward search and
	// cut.
	_, labelsJSON, ok := strings.Cut(line, " labels:")
	if !ok {
		return nil
	}
	// Following must be a JSON-like labels key-value map. In order to be open
	// to further extensions of the "goroutine" statement/line, we do not use
	// json.Unmarshal but instead a json.Decoder that reads only one JSON value,
	// the labels map/object and then stops, ignoring anything else that might
	// be following.
	m := map[string]string{}
	if err := json.NewDecoder(strings.NewReader(labelsJSON)).Decode(&m); err != nil {
		return nil
	}
	return m
}
