package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"strconv"
	"time"
)

var fileHash string

var numBits int

var numIterations int

var numZeroes int = 0

func main() {
	if len(os.Args) != 3 {
		println("Error: Incorrect number of command line args")
		return
	}
	hash := sha256.New()
	f, err := os.Open(os.Args[2])
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return

	}

	numBits, err = strconv.Atoi(os.Args[1])

	if err != nil {

		println("Error: Num bits is not a valid integer")
		return

	}

	//copy file into hash and update or return error
	if _, err := io.Copy(hash, f); err != nil {

		println("Error: Could not hash file.")
		return

	}

	fileHash = hex.EncodeToString(hash.Sum(nil))

	var prefix []byte

	prefix = nil
	count := 0
	start := time.Now()

	//String generation loop
	for prefix == nil {
		count++
		prefix = make([]byte, count)
		prefix = generateStrings(prefix, count-1)

	}
	end := time.Now()
	elapsed := end.Sub(start)
	rawHash := sha256.Sum256([]byte(string(prefix) + fileHash))
	combinedHash := rawHash[:]

	//Print results
	fmt.Println("File: ", os.Args[2])
	fmt.Println("Initial-hash: ", fileHash)
	fmt.Println("Proof-of-work: ", string(prefix))
	fmt.Println("Hash: ", hex.EncodeToString(combinedHash))
	fmt.Println("Leading-bits: ", numZeroes)
	fmt.Println("Iterations: ", numIterations)
	fmt.Println("Compute-time: ", elapsed.Seconds())
}

func validProof(proof string, numBits int, fileHash string) (isValid bool) {
	//Get hash of proof + hash of file
	rawHash := sha256.Sum256([]byte(proof + fileHash))

	//Figure out how many leading bytes of the hash we actually need to check
	numBytes := int(math.Ceil(float64(numBits) / 8))
	/**combinedHash := rawHash[:]
	fmt.Println(proof)
	fmt.Println(hex.EncodeToString(combinedHash))
	fmt.Println(rawHash)
	**/

	numZeroes = 0
	//Makes sure leading 0s are consecuvitve
	prevBlock := 8

	//Calculate leading 0s for byte, and add to total
	for i := 0; i < numBytes; i++ {
		if prevBlock == 8 {
			prevBlock = bits.LeadingZeros8(uint8(rawHash[i]))
			numZeroes = numZeroes + prevBlock
		} else {
			return false
		}
	}

	//If the number of leading 0s is equal to or greater than the puzzle length, POW is valid
	if numZeroes >= numBits {
		return true
	}

	//Otherwise return false
	return false

}

//Makes more sense to do this process iteratively. Probably possible to optimize as tail-call recusion
//but I'm just not sure how.
func generateStrings(prefix []byte, index int) (validString []byte) {
	//Loop through possible ascii values
	for i := 33; i < 128; i++ {
		numIterations++
		//Ignore double quotes, single quotes, and colons
		//Could possibly optimize by looping through a list of valid ascii characters
		//Instead of using the numeric values and having to check for invalid characters every time
		if (i != 34) && (i != 39) && (i != 58) {
			prefix[index] = byte(i)
		}

		//If valid proof, return string
		if validProof(string(prefix), numBits, fileHash) {

			return prefix
		}

		//If not valid, recursive call to iterate to next possible string
		if index > 0 {
			validString = generateStrings(prefix, index-1)
		}

		//Base case
		if validString != nil {
			return validString
		}
	}
	return nil
}
