package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Sides struct {
	colours   map[int]string
	opposites map[int]int
}

type CornerPiece struct {
	// clockwise
	colours     [3]int
	orientation int
}

// 2x2
// defining left and right slices
type Cube struct {
	left  [4]*CornerPiece
	right [4]*CornerPiece
}

type Cubes []Cube

type DisplayConfig struct {
	mode string
}

type PlayConfig struct {
	running bool
}

func getChar(side int, config DisplayConfig) string {
	switch side {
	case 0:
		if config.mode == "ansi" {
			return "\033[47m  \033[0m"
		}
		return "W"
	case 1:
		if config.mode == "ansi" {
			return "\033[48;5;208m  \033[0m"
		}
		return "O"
	case 2:
		if config.mode == "ansi" {
			return "\033[42m  \033[0m"
		}
		return "G"
	case 3:
		if config.mode == "ansi" {
			return "\033[41m  \033[0m"
		}
		return "R"
	case 4:
		if config.mode == "ansi" {
			return "\033[44m  \033[0m"
		}
		return "B"
	case 5:
		if config.mode == "ansi" {
			return "\033[43m  \033[0m"
		}
		return "Y"
	default:
		return "?"
	}
}

func NewSides() *Sides {
	return &Sides{
		colours: map[int]string{
			0: "white",
			1: "orange",
			2: "green",
			3: "red",
			4: "blue",
			5: "yellow",
		},
		opposites: map[int]int{
			0: 5,
			1: 3,
			2: 4,
			3: 1,
			4: 2,
			5: 0,
		},
	}
}

func (s *Sides) OppositeColour(side int) (string, bool) {
	opp, ok := s.opposites[side]
	if !ok {
		return "", false
	}
	colour, ok := s.colours[opp]
	return colour, ok
}

func pieceColours(piece int) (int, int, int) {
	switch piece {
	case 0:
		return 1, 0, 2
	case 1:
		return 1, 4, 0
	case 2:
		return 1, 5, 4
	case 3:
		return 1, 2, 5
	case 4:
		return 3, 2, 0
	case 5:
		return 3, 0, 4
	case 6:
		return 3, 4, 5
	case 7:
		return 3, 5, 2
	default:
		panic("invalid corner index")
	}
}

func (cubes *Cubes) NewCube() Cube {
	cube := Cube{}
	// left side
	for leftPiece := 0; leftPiece < 4; leftPiece++ {
		s1, s2, s3 := pieceColours(leftPiece)
		cube.left[leftPiece] = &CornerPiece{colours: [3]int{s1, s2, s3}, orientation: 0}
	}
	// right side
	for rightPiece := 0; rightPiece < 4; rightPiece++ {
		s1, s2, s3 := pieceColours(rightPiece + 4)
		cube.right[rightPiece] = &CornerPiece{colours: [3]int{s1, s2, s3}, orientation: 0}
	}
	*cubes = append(*cubes, cube)
	return cube
}

func (cubes *Cubes) getPiece(cube, piece int) *CornerPiece {
	if piece < 4 && piece >= 0 {
		return (*cubes)[cube].left[piece]
	} else if piece < 8 {
		return (*cubes)[cube].right[piece-4]
	} else {
		fmt.Println("error: invalid piece index")
		return nil
	}
}

func (cubes *Cubes) swapPieces(cube, piece1, piece2 int) {
	p1 := cubes.getPiece(cube, piece1)
	p2 := cubes.getPiece(cube, piece2)
	if p1 == nil || p2 == nil {
		fmt.Println("error: invalid piece pointer")
		return
	}
	*p1, *p2 = *p2, *p1
}

func (cubes *Cubes) rotatePiece(cube, piece, turns int, clockwise bool) {
	p := cubes.getPiece(cube, piece)
	direction := 1
	if clockwise {
		direction = -1
	}
	p.orientation = (p.orientation + (turns * direction))
	if p.orientation < 0 {
		p.orientation += 3
	}
	p.orientation = p.orientation % 3
}

func (cubes *Cubes) move(cube int, move byte) {
	switch move {
	case 'h':
		// F
		cubes.swapPieces(cube, 0, 4)
		cubes.swapPieces(cube, 0, 7)
		cubes.swapPieces(cube, 0, 3)
		cubes.rotatePiece(cube, 0, 1, true)
		cubes.rotatePiece(cube, 4, 1, false)
		cubes.rotatePiece(cube, 7, 1, true)
		cubes.rotatePiece(cube, 3, 1, false)
	case 'g':
		// F'
		cubes.swapPieces(cube, 0, 3)
		cubes.swapPieces(cube, 0, 7)
		cubes.swapPieces(cube, 0, 4)
		cubes.rotatePiece(cube, 0, 1, true)
		cubes.rotatePiece(cube, 4, 1, false)
		cubes.rotatePiece(cube, 7, 1, true)
		cubes.rotatePiece(cube, 3, 1, false)
	case 'a':
		// B
		cubes.swapPieces(cube, 1, 2)
		cubes.swapPieces(cube, 1, 6)
		cubes.swapPieces(cube, 1, 5)
		cubes.rotatePiece(cube, 1, 1, false)
		cubes.rotatePiece(cube, 2, 1, true)
		cubes.rotatePiece(cube, 6, 1, false)
		cubes.rotatePiece(cube, 5, 1, true)
	case ';':
		// B'
		cubes.swapPieces(cube, 1, 5)
		cubes.swapPieces(cube, 1, 6)
		cubes.swapPieces(cube, 1, 2)
		cubes.rotatePiece(cube, 1, 1, false)
		cubes.rotatePiece(cube, 2, 1, true)
		cubes.rotatePiece(cube, 6, 1, false)
		cubes.rotatePiece(cube, 5, 1, true)
	case 'i':
		// U
		cubes.swapPieces(cube, 0, 1)
		cubes.swapPieces(cube, 0, 5)
		cubes.swapPieces(cube, 0, 4)
		cubes.rotatePiece(cube, 0, 1, false)
		cubes.rotatePiece(cube, 1, 1, true)
		cubes.rotatePiece(cube, 5, 1, false)
		cubes.rotatePiece(cube, 4, 1, true)
	case 'r':
		// U'
		cubes.swapPieces(cube, 0, 4)
		cubes.swapPieces(cube, 0, 5)
		cubes.swapPieces(cube, 0, 1)
		cubes.rotatePiece(cube, 0, 1, false)
		cubes.rotatePiece(cube, 1, 1, true)
		cubes.rotatePiece(cube, 5, 1, false)
		cubes.rotatePiece(cube, 4, 1, true)
	case 's':
		// D
		cubes.swapPieces(cube, 3, 7)
		cubes.swapPieces(cube, 3, 6)
		cubes.swapPieces(cube, 3, 2)
		cubes.rotatePiece(cube, 3, 1, true)
		cubes.rotatePiece(cube, 7, 1, false)
		cubes.rotatePiece(cube, 6, 1, true)
		cubes.rotatePiece(cube, 2, 1, false)
	case 'l':
		// D'
		cubes.swapPieces(cube, 3, 2)
		cubes.swapPieces(cube, 3, 6)
		cubes.swapPieces(cube, 3, 7)
		cubes.rotatePiece(cube, 3, 1, true)
		cubes.rotatePiece(cube, 7, 1, false)
		cubes.rotatePiece(cube, 6, 1, true)
		cubes.rotatePiece(cube, 2, 1, false)
	case 'o':
		// R
		cubes.swapPieces(cube, 4, 5)
		cubes.swapPieces(cube, 4, 6)
		cubes.swapPieces(cube, 4, 7)
	case 'j':
		// R'
		cubes.swapPieces(cube, 4, 7)
		cubes.swapPieces(cube, 4, 6)
		cubes.swapPieces(cube, 4, 5)
	case 'f':
		// L
		cubes.swapPieces(cube, 0, 3)
		cubes.swapPieces(cube, 0, 2)
		cubes.swapPieces(cube, 0, 1)
	case 'e':
		// L'
		cubes.swapPieces(cube, 0, 1)
		cubes.swapPieces(cube, 0, 2)
		cubes.swapPieces(cube, 0, 3)
	}
}

func (cubes *Cubes) loadCubeStrings(index int, config DisplayConfig) [8]string {

	var strings [8]string
	strings[0] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[1].colours[(1+cubes.getPiece(index, 1).orientation)%3], config),
		getChar((*cubes)[index].right[1].colours[(2+cubes.getPiece(index, 5).orientation)%3], config))
	strings[1] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[1].colours[(2+cubes.getPiece(index, 1).orientation)%3], config),
		getChar((*cubes)[index].right[1].colours[(1+cubes.getPiece(index, 5).orientation)%3], config))
	strings[2] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[0].colours[(1+cubes.getPiece(index, 0).orientation)%3], config),
		getChar((*cubes)[index].right[0].colours[(2+cubes.getPiece(index, 4).orientation)%3], config))
	strings[3] = fmt.Sprintf("%s%s%s%s%s%s",
		getChar((*cubes)[index].left[1].colours[(0+cubes.getPiece(index, 1).orientation)%3], config),
		getChar((*cubes)[index].left[0].colours[(0+cubes.getPiece(index, 0).orientation)%3], config),
		getChar((*cubes)[index].left[0].colours[(2+cubes.getPiece(index, 0).orientation)%3], config),
		getChar((*cubes)[index].right[0].colours[(1+cubes.getPiece(index, 4).orientation)%3], config),
		getChar((*cubes)[index].right[0].colours[(0+cubes.getPiece(index, 4).orientation)%3], config),
		getChar((*cubes)[index].right[1].colours[(0+cubes.getPiece(index, 5).orientation)%3], config))
	strings[4] = fmt.Sprintf("%s%s%s%s%s%s",
		getChar((*cubes)[index].left[2].colours[(0+cubes.getPiece(index, 2).orientation)%3], config),
		getChar((*cubes)[index].left[3].colours[(0+cubes.getPiece(index, 3).orientation)%3], config),
		getChar((*cubes)[index].left[3].colours[(1+cubes.getPiece(index, 3).orientation)%3], config),
		getChar((*cubes)[index].right[3].colours[(2+cubes.getPiece(index, 7).orientation)%3], config),
		getChar((*cubes)[index].right[3].colours[(0+cubes.getPiece(index, 7).orientation)%3], config),
		getChar((*cubes)[index].right[2].colours[(0+cubes.getPiece(index, 6).orientation)%3], config))
	strings[5] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[3].colours[(2+cubes.getPiece(index, 3).orientation)%3], config),
		getChar((*cubes)[index].right[3].colours[(1+cubes.getPiece(index, 7).orientation)%3], config))
	strings[6] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[2].colours[(1+cubes.getPiece(index, 2).orientation)%3], config),
		getChar((*cubes)[index].right[2].colours[(2+cubes.getPiece(index, 6).orientation)%3], config))
	strings[7] = fmt.Sprintf("    %s%s    ",
		getChar((*cubes)[index].left[2].colours[(2+cubes.getPiece(index, 2).orientation)%3], config),
		getChar((*cubes)[index].right[2].colours[(1+cubes.getPiece(index, 6).orientation)%3], config))

	return strings
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// fallback and default width
		return 80
	}
	return width
}

func (cubes *Cubes) print(config DisplayConfig) {

	var cubeStrings [][8]string

	for i := range *cubes {
		cubeStrings = append(cubeStrings, cubes.loadCubeStrings(i, config))
	}

	fmt.Print("\r\n")
	for i := range *cubes {
		for row := 0; row < len(cubeStrings[i]); row++ {
			fmt.Print(cubeStrings[i][row], "\r\n")
		}
	}
	fmt.Print("\r\n")
}

func showOptions() {
	fmt.Print("\r\nF  - [h]   U  - [i]   R  - [o]   L  - [f]   B  - [a]   D  - [s]")
	fmt.Print("\r\nF' - [g]   U' - [r]   R' - [j]   L' - [e]   B' - [;]   D' - [l]")
	fmt.Print("\r\n\nquit - [q]\r\n")
}

func deleteLines(count int) {
	for i := 0; i < count; i++ {
		fmt.Print("\033[1A\033[2K")
	}
}

func disableCursorBlink() {
	fmt.Print("\033[?12l")
}

func enableCursorBlink() {
	fmt.Print("\033[?12h")
}

func disableCursor() {
	fmt.Print("\033[?25l")
}

func enableCursor() {
	fmt.Print("\033[?25h")
}

func (cubes *Cubes) Play(config *PlayConfig, displayConfig *DisplayConfig) {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	disableCursorBlink()
	disableCursor()
	defer enableCursorBlink()
	defer enableCursor()

	input := make([]byte, 1)
	dynamicLineCount := 10
	// termWidth := getTerminalWidth()
	// cubeWidth := 10
	errorCount := 0
	config.running = true
	showOptions()
	cubes.print(*displayConfig)
	for config.running {
		// read input
		_, err := os.Stdin.Read(input)
		if err != nil {
			fmt.Print("input error: ", err)
			errorCount++
			continue
		}
		deleteLines(dynamicLineCount)
		errorCount = 0

		// update cube
		switch input[0] {
		case 'h':
			cubes.move(0, 'h')
		case 'j':
			cubes.move(0, 'j')
		case 'i':
			cubes.move(0, 'i')
		case 'o':
			cubes.move(0, 'o')
		case 'f':
			cubes.move(0, 'f')
		case 'l':
			cubes.move(0, 'l')
		case ';':
			cubes.move(0, ';')
		case 'g':
			cubes.move(0, 'g')
		case 'r':
			cubes.move(0, 'r')
		case 'e':
			cubes.move(0, 'e')
		case 's':
			cubes.move(0, 's')
		case 'a':
			cubes.move(0, 'a')
		case ' ':
			// timer
		case 'q':
			config.running = false
		}

		// render
		cubes.print(*displayConfig)
	}
}
