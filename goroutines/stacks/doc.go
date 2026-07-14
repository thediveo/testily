/*
Package stacks provides retrieving (textual) stack traces for all or only the
calling go routine. It is a convenience wrapper around [runtime.Stack] that
takes care of sufficiently large result buffer allocation.
*/
package stacks
