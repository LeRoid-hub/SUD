package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Version = "0.1.1"

type Sudoku struct {
	Title  string
	Author string
	Rules  string
	Board  [][]int
}

func (sudoku *Sudoku) toBytes() []byte {
	var bytes []byte
	bytes = append(bytes, []byte(sudoku.Title)...)
	bytes = append(bytes, 0)
	bytes = append(bytes, []byte(sudoku.Author)...)
	bytes = append(bytes, 0)
	bytes = append(bytes, []byte(sudoku.Rules)...)
	bytes = append(bytes, 0)
	board := boardtoBytes(sudoku.Board)
	bytes = append(bytes, board...)
	return bytes
}

func boardtoBytes(boardraw [][]int) []byte {
	var board []byte
	var buffer byte
	for _, row := range boardraw {
		for _, cell := range row {
			if cell > 9 || cell < 0 {
				print(cell, " is not a valid cell value")
				panic("Invalid cell value")
			}

			fmt.Println(cell, buffer)

			if buffer == 0 {
				buffer = byte(cell)
			} else {
				buffer = buffer<<4 | byte(cell)
				board = append(board, buffer)
				buffer = 0
			}
		}
	}
	buffer = buffer << 4
	board = append(board, buffer)
	return board
}

func bytesToBoard(input []byte) [][]int {
	var output [][]int
	row := []int{}

	for _, cell := range input {
		lower := cell & 0x0F
		upper := cell >> 4

		row = append(row, int(upper))
		if len(row) == 9 {
			rowCopy := make([]int, len(row))
			copy(rowCopy, row)
			output = append(output, rowCopy)
			row = []int{}
		}

		if lower == 0xF {
			break
		}

		row = append(row, int(lower))
		if len(row) == 9 {
			rowCopy := make([]int, len(row))
			copy(rowCopy, row)
			output = append(output, rowCopy)
			row = []int{}
		}
	}

	return output
}

func Serialize(sudoku Sudoku) {
	var serialized []byte
	var header []byte
	var body []byte
	var footer []byte
	var length uint16

	//HEADER
	version := getVersion()
	header = append(header, byte(version>>8), byte(version))

	//BODY
	body = sudoku.toBytes()

	//FOOTER

	//LENGTH
	length = uint16(len(header) + 2 + len(body) + len(footer))
	header = append(header, byte(length>>8), byte(length))

	serialized = append(serialized, header...)
	serialized = append(serialized, body...)
	serialized = append(serialized, footer...)
	os.WriteFile("sudoku.sud", serialized, 0644)
}

func Deserialize() {
	serialized, err := os.ReadFile("sudoku.sud")
	if err != nil {
		panic(err)
	}

	printBytes(serialized)

	version := uint16(serialized[0])<<8 | uint16(serialized[1])

	if version != getVersion() {
		printBytes([]byte{byte(version >> 8), byte(version)})
		println("=/")
		printBytes([]byte{byte(getVersion() >> 8), byte(getVersion())})
		panic("Invalid version")
	}

	body := serialized[4 : len(serialized)-2]
	board := bytesToBoard(body)
	fmt.Println(board)
}

func getVersion() uint16 {
	version := strings.Split(Version, ".")

	majorInt, err := strconv.Atoi(version[0])
	if err != nil || majorInt > 15 {
		fmt.Println("Invalid major version, the major version must be an integer between 0 and 15")
		panic(err)
	}
	major := uint16(majorInt)

	minorInt, err := strconv.Atoi(version[1])
	if err != nil || minorInt > 15 {
		fmt.Println("Invalid minor version, the minor version must be an integer between 0 and 15")
		panic(err)
	}
	minor := uint16(minorInt)

	patchInt, err := strconv.Atoi(version[2])
	if err != nil || patchInt > 255 {
		fmt.Println("Invalid patch version, the patch version must be an integer between 0 and 255")
		panic(err)
	}
	patch := uint16(patchInt)

	versionbyte := major<<12 | minor<<8 | patch

	return versionbyte
}

func usedBits(n uint16) int {
	bits := 0
	for n > 0 {
		bits++
		n >>= 1
	}
	return bits
}

func printBytes(bytes []byte) {
	for _, n := range bytes {
		fmt.Printf("%08b ", n)
		fmt.Print("\n")
	}
}

func main() {
	sudoku := Sudoku{
		Title:  "Sudoku",
		Author: "John Doe",
		Rules:  "Fill the board with numbers from 1 to 9",
		Board: [][]int{
			{5, 3, 0, 0, 7, 0, 0, 0, 0},
			{6, 0, 0, 1, 9, 5, 0, 0, 0},
			{0, 9, 8, 0, 0, 0, 0, 6, 0},
			{8, 0, 0, 0, 6, 0, 0, 0, 3},
			{4, 0, 0, 8, 0, 3, 0, 0, 1},
			{7, 0, 0, 0, 2, 0, 0, 0, 6},
			{0, 6, 0, 0, 0, 0, 2, 8, 0},
			{0, 0, 0, 4, 1, 9, 0, 0, 5},
			{0, 0, 0, 0, 8, 0, 0, 7, 9},
		},
	}

	Serialize(sudoku)
	Deserialize()
}

/*
func clutter(input [][]int) []byte {
	var output []byte
	var buffer byte
	var full bool
	for _, row := range input {
		for _, cell := range row {
			if cell > 9 || cell < 0 {
				print(cell, " is not a valid cell value")
				panic("Invalid cell value")
			}

			if !full {
				buffer = byte(cell)
				full = true
			} else {
				buffer = buffer<<4 | byte(cell)
				output = append(output, buffer)
				buffer = 0
				full = false
			}
		}
	}

	if buffer != 0 {
		buffer = buffer << 4
		buffer = buffer | 0x0F
		output = append(output, buffer)
	}

	return output
}

func declutter(input []byte) {
	var output [][]int
	row := []int{}

	for _, cell := range input {
		lower := cell & 0x0F
		upper := cell >> 4

		row = append(row, int(upper))
		if len(row) == 9 {
			rowCopy := make([]int, len(row))
			copy(rowCopy, row)
			output = append(output, rowCopy)
			row = []int{}
		}

		if lower == 0xF {
			break
		}

		row = append(row, int(lower))
		if len(row) == 9 {
			rowCopy := make([]int, len(row))
			copy(rowCopy, row)
			output = append(output, rowCopy)
			row = []int{}
		}
	}

	fmt.Print("\n")
	fmt.Println(output)
}*/
