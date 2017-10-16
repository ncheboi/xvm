package main

import (
	"fmt"

	"./xvm"
)

func main() {
	_, err := xvm.NewRun()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("hello")
	}
}
