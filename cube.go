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

func getChar(side int) byte {
	switch side {
	case 0:
		return 'W'
	case 1:
		return 'O'
	case 2:
		return 'G'
	case 3:
		return 'R'
	case 4:
		return 'B'
	case 5:
		return 'Y'
	default:
		return '?'
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

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// fallback or default width
		return 80
	}
	return width
}

func (cubes *Cubes) loadCubeStrings(index int) [6]string {

	var strings [6]string
	strings[0] = fmt.Sprintf("   %c%c      ",
		getChar((*cubes)[index].left[1].colours[(2+cubes.getPiece(index, 1).orientation)%3]),
		getChar((*cubes)[index].right[1].colours[(1+cubes.getPiece(index, 5).orientation)%3]))
	strings[1] = fmt.Sprintf("   %c%c      ",
		getChar((*cubes)[index].left[0].colours[(1+cubes.getPiece(index, 0).orientation)%3]),
		getChar((*cubes)[index].right[0].colours[(2+cubes.getPiece(index, 4).orientation)%3]))
	strings[2] = fmt.Sprintf("%c%c %c%c %c%c %c%c",
		getChar((*cubes)[index].left[1].colours[(0+cubes.getPiece(index, 1).orientation)%3]),
		getChar((*cubes)[index].left[0].colours[(0+cubes.getPiece(index, 0).orientation)%3]),
		getChar((*cubes)[index].left[0].colours[(2+cubes.getPiece(index, 0).orientation)%3]),
		getChar((*cubes)[index].right[0].colours[(1+cubes.getPiece(index, 4).orientation)%3]),
		getChar((*cubes)[index].right[0].colours[(0+cubes.getPiece(index, 4).orientation)%3]),
		getChar((*cubes)[index].right[1].colours[(0+cubes.getPiece(index, 5).orientation)%3]),
		getChar((*cubes)[index].right[1].colours[(2+cubes.getPiece(index, 5).orientation)%3]),
		getChar((*cubes)[index].left[1].colours[(1+cubes.getPiece(index, 1).orientation)%3]))
	strings[3] = fmt.Sprintf("%c%c %c%c %c%c %c%c",
		getChar((*cubes)[index].left[2].colours[(0+cubes.getPiece(index, 2).orientation)%3]),
		getChar((*cubes)[index].left[3].colours[(0+cubes.getPiece(index, 3).orientation)%3]),
		getChar((*cubes)[index].left[3].colours[(1+cubes.getPiece(index, 3).orientation)%3]),
		getChar((*cubes)[index].right[3].colours[(2+cubes.getPiece(index, 7).orientation)%3]),
		getChar((*cubes)[index].right[3].colours[(0+cubes.getPiece(index, 7).orientation)%3]),
		getChar((*cubes)[index].right[2].colours[(0+cubes.getPiece(index, 6).orientation)%3]),
		getChar((*cubes)[index].right[2].colours[(1+cubes.getPiece(index, 6).orientation)%3]),
		getChar((*cubes)[index].left[2].colours[(2+cubes.getPiece(index, 2).orientation)%3]))
	strings[4] = fmt.Sprintf("   %c%c      ",
		getChar((*cubes)[index].left[3].colours[(2+cubes.getPiece(index, 3).orientation)%3]),
		getChar((*cubes)[index].right[3].colours[(1+cubes.getPiece(index, 7).orientation)%3]))
	strings[5] = fmt.Sprintf("   %c%c      ",
		getChar((*cubes)[index].left[2].colours[(1+cubes.getPiece(index, 2).orientation)%3]),
		getChar((*cubes)[index].right[2].colours[(2+cubes.getPiece(index, 6).orientation)%3]))

	return strings
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

func (cubes *Cubes) rotatePiece(cube, piece, turns int) {
	p := cubes.getPiece(cube, piece)
	p.orientation = (p.orientation + turns) % 3
}

func (cubes *Cubes) move(cube int, move byte) {
	switch move {
	case 'F':
		cubes.swapPieces(cube, 0, 4)
		cubes.swapPieces(cube, 0, 7)
		cubes.swapPieces(cube, 0, 3)
		cubes.rotatePiece(cube, 0, 2)
		cubes.rotatePiece(cube, 4, 1)
		cubes.rotatePiece(cube, 3, 1)
		cubes.rotatePiece(cube, 7, 2)
	case 'B':
	case 'U':
	case 'D':
	case 'R':
	case 'L':
	}
}

func (cubes *Cubes) print() {

	var cubeStrings [][6]string

	for i := range *cubes {
		cubeStrings = append(cubeStrings, cubes.loadCubeStrings(i))
	}

	// 6 x 8 per cube
	//   RR
	//   RR
	// WWGGYYBB
	// WWGGYYBB
	//   OO
	//   OO

	for i, _ := range *cubes {
		for row := 0; row < 6; row++ {
			fmt.Println(cubeStrings[i][row])
		}
	}

}
