/*
Package tuple mistreats structs as a poor man's idea of tuples, using Go
generics. This package is in no way intended to be a comprehensive tuple
provider. Instead, it's a bare minimum to pack, pass, and unpack tuples in
places such as channels where only a single type is expected.
*/
package tuples
