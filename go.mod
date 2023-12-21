module git.ff02.de/display

// Need to pin Go itself to 1.17, as the Kindle runs on kernel 2.6.31 and Go >= 1.18 requires a linux kernel >= 2.6.32
// https://go-review.googlesource.com/c/go/+/346789
go 1.17

require golang.org/x/image v0.14.0

require golang.org/x/text v0.14.0 // indirect
