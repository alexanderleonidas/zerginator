package bitoperations

import (
	"zerginator/globals"
)

// GetBit returns the bit at the given square
func GetBit(bitBoard uint64, square int) uint64 {
	/*
		Here we shift the bitBoard to the right by 1 for a shift count of square index
		and then use the bitwise AND operator to check if the least significant bit is a 1 or 0.
		Shifting the bitBoard to the right by 1 is equivalent to dividing the bitBoard by 2 and vice versa.
	*/
	return (bitBoard >> square) & 1
}

// SetBit sets the bit at the given square to 1
func SetBit(bitBoard *uint64, square int) {
	*bitBoard |= 1 << square
	*bitBoard &= globals.Bitboard40Mask
}

// PopBit sets the bit at the given square to 0
func PopBit(bitBoard *uint64, square int) {
	/*
		It is important to note that if we use this multiple times, the bit will be flipped
		multiple times. Therefore, we have to check if the bit is already set to 1 before flipping it
	*/
	if GetBit(*bitBoard, square) == 1 {
		*bitBoard ^= 1 << square
	} else {
		return
	}
}

// CountBits returns the number of bits set to 1 in the given bitBoard
func CountBits(bitBoard uint64) int {
	/*
		Here we use a trick to count the number of bits set to 1 in a bitBoard.
		We repeatedly turn off the rightmost bit in the bitBoard and increment a counter
		until the bitBoard is 0.
	*/
	count := 0
	for bitBoard != 0 {
		count++
		// reset the least significant bit
		bitBoard &= bitBoard - 1
	}
	return count
}

// GetLeastSignificantBitIndex returns the index of the least significant 1st bit set to 1
func GetLeastSignificantBitIndex(bitBoard uint64) int {
	/*
		To get the index of the least significant 1st bit set to 1, we use a trick to isolate the least significant 1-bit
		and then add trailing ones up to that bit. We then use the CountBits function to get the index.
	*/
	if bitBoard != 0 {
		// add trailing ones up to the least significant 1-bit
		bitBoard = (bitBoard & -bitBoard) - 1
		return CountBits(bitBoard)
	} else {
		return -1
	}
}
