package main

import "fmt"
import "git.ff02.de/display/fetchers"

func main() {
	fmt.Println("Hello, World!")

	_ = fetchers.KVV()
	_ = fetchers.DWD()
}
