package main

// #include <stdio.h>
// #include <errno.h>
import "C"
import "time"

//export HelloWorld
func HelloWorld() {
	for {
		print("Hello World\n")
		time.Sleep(1 * time.Second)
	}
}

func main() {}
