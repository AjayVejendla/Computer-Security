package main

import (
	"os"
)

const a, c int = 1103515245, 12345

const chunkSize = 16

var previousBlock []byte = make([]byte, chunkSize)

//var fileBuffer []byte = make([]byte, chunkSize)

var keystreamBuffer []byte = make([]byte, chunkSize)

var outputBuffer []byte = make([]byte, chunkSize)

var fileBuffer []byte = make([]byte, chunkSize)

var x int

func main() {
	flags := os.Args[1:]

	ciphertext, _ := os.Open(flags[1])
	outputtext, _ := os.Create(flags[2])

	//Defer file close
	defer ciphertext.Close()
	defer outputtext.Close()

	//Get initial X
	x = int(hash(flags[0]) % 256)
	x = keystreamGenerate(x)

	//Load IV

	loadBlock()
	copy(previousBlock, keystreamBuffer)

	var filePosition int64 = 0

	fileInfo, _ := ciphertext.Stat()
	endPosition := fileInfo.Size()

	for filePosition != endPosition {

		filePosition = filePosition + 16
		loadBlock()
		_, _ = ciphertext.Read(outputBuffer)
		copy(fileBuffer, outputBuffer)
		//XOR block with keystream
		for i := range outputBuffer {

			outputBuffer[i] = outputBuffer[i] ^ keystreamBuffer[i]

		}

		//Reverse shuffling
		for i := 15; i >= 0; i-- {

			first := keystreamBuffer[i] & 15
			second := (keystreamBuffer[i] >> 4) & 15

			outputBuffer[first], outputBuffer[second] = outputBuffer[second], outputBuffer[first]

		}

		// XOR Previous block with recently read block
		for i := range outputBuffer {

			outputBuffer[i] = previousBlock[i] ^ outputBuffer[i]

		}

		if filePosition == endPosition {
			outputBuffer = outputBuffer[0 : 16-outputBuffer[15]]
		}

		outputtext.Write(outputBuffer)
		copy(previousBlock, fileBuffer)

	}

	//Remove padding

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
