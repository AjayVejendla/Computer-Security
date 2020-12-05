package main

import (
	"os"
)

const a, c int = 1103515245, 12345

//Number of bytes to read at once
const chunkSize = 16

//Create a buffer to partially read the file into
var fileBuffer []byte = make([]byte, chunkSize)

//Create a buffer for file writing
//Doubles as temp buffer for the previous block
var outputBuffer []byte = make([]byte, chunkSize)

//Hold 16 bits of keystream to use
var keystreamBuffer []byte = make([]byte, chunkSize)

//Temp value for keystream
var x int

func main() {

	//Array for command line arguments
	flags := os.Args[1:]

	//Open the file specified in the first command line arg
	//File to read from and XOR contents of
	plaintext, err1 := os.Open(flags[1])
	ciphertext, _ := os.Create(flags[2])

	//Close files when main exits
	//defer ciphertext.Close()
	defer plaintext.Close()

	//Get initial X
	x = int(hash(flags[0]) % 256)
	x = keystreamGenerate(x)

	//Load IV as previous block
	loadBlock()
	copy(outputBuffer, keystreamBuffer)

	//Var to store how many bytes were read with each chunk
	//Will be 128 unless last chunk is read
	var numBytes int = 16
	for err1 == nil && numBytes == chunkSize {

		loadBlock()
		numBytes, err1 = plaintext.Read(fileBuffer)

		//Write out the contents of the ouput buffer
		//Only write up to the number of bytes read in the last chunk to avoid writing too many bytes

		//Loop to pad bytes

		for i := range fileBuffer[numBytes:] {
			fileBuffer[i+numBytes] = byte(chunkSize - numBytes)
		}
		// XOR Previous block with recently read block
		for i := range outputBuffer {

			outputBuffer[i] = fileBuffer[i] ^ outputBuffer[i]

		}

		//Shuffle block using keystream

		for i := range keystreamBuffer {

			first := keystreamBuffer[i] & 15
			second := (keystreamBuffer[i] >> 4) & 15

			outputBuffer[first], outputBuffer[second] = outputBuffer[second], outputBuffer[first]

		}

		//xor result with keystream
		for i := range outputBuffer {

			outputBuffer[i] = outputBuffer[i] ^ keystreamBuffer[i]

		}

		ciphertext.Write(outputBuffer)
	}

}

func hash(input string) uint64 {

	var hash uint64

	for _, c := range input {
		hash = uint64(c) + (hash << 6) + (hash << 16) - hash
	}

	return hash
}

//Don't actually need to encapsulate as function, potentially slower because a,c are created on every call
//But speed difference should be negligible, doing this for readability
func keystreamGenerate(x_n int) int {

	return (((x_n * a) + c) % 256)

}

func loadBlock() {

	//Load 16 bytes from keystream
	for i := range keystreamBuffer {

		keystreamBuffer[i] = byte(x)
		x = keystreamGenerate(x)

	}

}
