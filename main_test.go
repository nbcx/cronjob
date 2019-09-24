package main

import (
	"fmt"
	"os"
)

func TsetMain() {
	fmt.Printf("process %d\n", os.Getpid())
	ejob := NewEJobApplication()
	ejob.run()
}

