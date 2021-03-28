# unsafeless

Transmuting types in Go without using `unsafe` or `reflect`. See [unsafeless.go2](unsafeless.go2) for an explanation
and [unsafeless\_test.go2](unsafeless_test.go2) for example usage.

This concept doesn't actually need Go 2, but the generics make it easier to use? Maybe? The Go 2 might be even nicer if
type inference worked properly.

You can either test this by getting the `dev.go2go` toolchain
(https://go.googlesource.com/go/+/refs/heads/dev.go2go/README.go2go.md), or by copy-pasting it into the go2go
playground (https://go2goplay.golang.org).
