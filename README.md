# `testily`

![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)

A small collection of weird test convenience helpers with odd names, to aid in
DRY. Preferably dot-imported as a simple DSL to make your unit tests more
concise. A nice companion to [Ginkgo](https://github.com/onsi/gomega).

- package `tuples`: type-safe tuples packing and unpacking: `Pair`, `Triple`.
- package `concur`: calling functions on new goroutines and signalling when
  they've finished, optionally channeling their results back: `CloseWhenDone`,
  `PassWhenDone`.
- package `nothing`: `interface{}` has `any`, we have `Nothing` for `struct{}`.

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
