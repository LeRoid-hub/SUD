package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Version = "1.1.1"

type Sudoku struct {
	Title  string
	Author string
	Rules  string
	Board  [][]int
}

func Serialize() {
	var serialized []byte
	var header []byte
	var body []byte
	var footer []byte
	var length uint16
	var empty byte = 0

	//HEADER
	version := getVersion()
	header = append(header, byte(version>>8), byte(version))

	//BODY
	title := stringToBytes("Sudoku")
	author := stringToBytes("John Doe")
	rules := stringToBytes("Fill the board with numbers from 1 to 9")

	body = append(body, title...)
	body = append(body, empty)
	body = append(body, author...)
	body = append(body, empty)
	body = append(body, rules...)
	body = append(body, empty)

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

func stringToBytes(s string) []byte {
	return []byte(s)
}

func main() {
	Serialize()
	Deserialize()

	print(usedBits(5))
	printBytes(stringToBytes("Hello, World!"))
}
