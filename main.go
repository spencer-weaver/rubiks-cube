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

	displayConfig := DisplayConfig{
		mode: "ansi",
	}

	playConfig := PlayConfig{}

	cubes := Cubes{}

	switch os.Args[1] {
	case "play":
		cubes.NewCube()
		cubes.Play(&playConfig, &displayConfig)
	case "load":

	case "new":
		cubes.NewCube()
	case "scramble":

	case "solve":

	case "print":
		cubes.print(displayConfig)
	}
}
