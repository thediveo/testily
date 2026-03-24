# `testily`

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/testily)
[![License](https://img.shields.io/github/license/thediveo/testily)](https://img.shields.io/github/license/thediveo/testily)
![build and test](https://github.com/thediveo/testily/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
![goroutines](https://img.shields.io/badge/go%20routines-not%20leaking-success)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/testily)](https://goreportcard.com/report/github.com/thediveo/testily)
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)

A small collection of weird test convenience helpers with odd names, to aid in
DRY. Preferably dot-imported as a simple DSL to make your unit tests more
concise. A nice companion to [Ginkgo](https://github.com/onsi/gomega).

- package [`tuples`](https://pkg.go.dev/github.com/thediveo/testily/tuples):
  type-safe tuples packing and unpacking:
  [`Pair`](https://pkg.go.dev/github.com/thediveo/testily/tuples#Pair),
  [`Triple`](https://pkg.go.dev/github.com/thediveo/testily/tuples#Triple). Use
  `PackPair(...)`, `Unpack()`, et cetera.
- package [`concur`](https://pkg.go.dev/github.com/thediveo/testily/concur):
  calling functions on new goroutines and signalling when they've finished,
  optionally channeling their results back:
  [`CloseWhenDone`](https://pkg.go.dev/github.com/thediveo/testily/concur#CloseWhenGone),
  [`PassWhenDone`](https://pkg.go.dev/github.com/thediveo/testily/concur#PassWhenGone).
- package [`nothing`](https://pkg.go.dev/github.com/thediveo/testily/nothing):
  `interface{}` has `any`, we have `Nothing` for `struct{}`.

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/testily`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `testily.code-workspace` and off you go...

## Supported Go Versions

`lxkns` supports versions of Go that are noted by the [Go release
policy](https://golang.org/doc/devel/release.html#policy), that is, major
versions _N_ and _N_-1 (where _N_ is the current major version).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright and License

`testily` is Copyright 2026 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
