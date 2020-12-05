package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

var fileHash string
var isValid bool = true
var numZeroes int = 0
var expectedBits int = 0
var initialHash string = ""
var pow string = ""
var leadingBits string = ""
var hashFinal string = ""

func main() {
	if len(os.Args) != 3 {
		println("Error: Incorrect number of command line args")
		return

	}

	hash := sha256.New()
	f, err := os.Open(os.Args[2])
	defer f.Close()
	headersLoaded := loadHeaders()

	if !headersLoaded {
		return
	}

	//Check if all fields are filled
	//Only checks for the fields we need

	if (initialHash == "") || (pow == "") || (leadingBits == "") || (hashFinal == "") {
		println("One or more values in header file is missing or incorrectly labelled")
		return
	}

	if err != nil {

		fmt.Println(err)
		return

	}

	//Copy file into to hash object, and update hash or return error
	if _, err := io.Copy(hash, f); err != nil {

		println("The file could not be hashed.")
		return

	}

	//Run tests
	expectedBits, err = strconv.Atoi(leadingBits)
	if err != nil {
		println("Field Leading-bits is not a valid integer")
	}

	fileHash = hex.EncodeToString(hash.Sum(nil))

	if !(fileHash == initialHash) {
		isValid = false
		println("Hash of file does not match Header")
	}

	if !validProof(pow, expectedBits, fileHash) {
		isValid = false
		println("Hash of file hash with proof of work is not valid")
	}

	if numZeroes != expectedBits {
		isValid = false
		println("The number of leading bits listed is not equal to the actual number of leading bits")
	}

	//Print if test passed
	if isValid {
		println("Pass")
	}

}

func validProof(proof string, numBits int, fileHash string) (isPowValid bool) {
	rawHash := sha256.Sum256([]byte(proof + fileHash))
	numBytes := int(math.Ceil(float64(numBits) / 8))
	slicedHash := rawHash[:]
	if !(hashFinal == hex.EncodeToString(slicedHash)) {
		isValid = false
	}

	numZeroes = 0
	prevBlock := 8
	for i := 0; i < numBytes; i++ {
		if prevBlock == 8 {
			prevBlock = bits.LeadingZeros8(uint8(rawHash[i]))
			numZeroes = numZeroes + prevBlock
		} else {
			return false
		}
	}
	if numZeroes <= numBits {
		return true
	}
	return false

}

func loadHeaders() (success bool) {
	var currentLine string

	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {

		fmt.Println(err)
		return false

	}
	fileRead := bufio.NewReader(f)

	currentLine, err = fileRead.ReadString('\n')

	//For each line remove whitespace and split by colon
	for err == nil {

		currentLine, err = fileRead.ReadString('\n')
		currentLine = strings.TrimSpace(currentLine)
		splitLine := strings.Split(currentLine, ":")

		//If the field name matches a specific name, load the value into the respective variable
		if len(splitLine) == 2 {
			splitLine[1] = strings.TrimSpace(splitLine[1])
			switch fieldName := splitLine[0]; fieldName {

			case "Initial-hash":
				initialHash = splitLine[1]
			case "Proof-of-work":
				pow = splitLine[1]
			case "Hash":
				hashFinal = splitLine[1]
			case "Leading-bits":
				leadingBits = splitLine[1]
			}

		}

	}

	//If file reading it terminated unexpectedly, print error and return
	if err != io.EOF {

		fmt.Println(err)
		return false

	}
	return true
}
