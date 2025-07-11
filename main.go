package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: cube <options> -flags")
		os.Exit(1)
	}

	cubes := Cubes{}

	switch os.Args[1] {
	case "play":

	case "new":
		cubes.NewCube()
		cubes.move(0, 'F')
	case "scramble":

	case "solve":

	case "print":
	}

	cubes.print()
}
