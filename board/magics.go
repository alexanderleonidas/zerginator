package board

import (
	"fmt"
	"zerginator/bitoperations"
	"zerginator/globals"
)

/*
This file contains the magic number generation algorithm for the bishop and rook pieces.
*/

// state is a pseudo random number state
var randomState uint32 = 1804289383

// GetRandomUInt32 generates a 32-bit pseudo legal number
func GetRandomUInt32() uint32 {
	number := randomState
	// XOR-shift algorithm
	number ^= number << 13
	number ^= number >> 17
	number ^= number << 5
	randomState = number
	return number
}

// GetRandomUInt64 generates a 64-bit pseudo legal number
func GetRandomUInt64() uint64 {
	// define 4 random numbers
	random1 := uint64(GetRandomUInt32()) & 0xFFFF // slicing the first 16 bits of the random number
	random2 := uint64(GetRandomUInt32()) & 0xFFFF
	random3 := uint64(GetRandomUInt32()) & 0xFFFF
	random4 := uint64(GetRandomUInt32()) & 0xFFFF
	return random1 | (random2 << 16) | (random3 << 32) | (random4 << 48)
}

// GenerateMagicNumber generates a magic number
func GenerateMagicNumber() uint64 {
	return GetRandomUInt64() & GetRandomUInt64() & GetRandomUInt64()
}

// FindMagicNumber finds an appropriate magic number
func FindMagicNumber(square int, relevantBits int, isBishop int) uint64 {
	var occupancies [4096]uint64
	var attacks [4096]uint64
	var attackMask uint64

	if isBishop == 1 {
		attackMask = MaskBishopAttacks(square)
	} else {
		attackMask = MaskRookAttacks(square)
	}

	occupancyIndices := 1 << relevantBits

	// precompute occupancies & attacks
	for i := 0; i < occupancyIndices; i++ {
		occupancies[i] = SetOccupancy(i, relevantBits, attackMask)
		if isBishop == 1 {
			attacks[i] = BishopAttacksOnTheFly(square, occupancies[i])
		} else {
			attacks[i] = RookAttacksOnTheFly(square, occupancies[i])
		}
	}

	// search for a magic number
	for i := 0; i < 100000000; i++ {
		magicNumber := GenerateMagicNumber()

		// heuristic filter, different from the original algorithm
		if bitoperations.CountBits((attackMask*magicNumber)&0xFE00000000000000) < 5 {
			continue
		}

		usedAttacks := make([]uint64, 4096) // reset for each candidate
		fail := 0

		for index := 0; index < occupancyIndices; index++ {
			magicIndex := int((occupancies[index] * magicNumber) >> (64 - relevantBits))

			if usedAttacks[magicIndex] == 0 {
				usedAttacks[magicIndex] = attacks[index]
			} else if usedAttacks[magicIndex] != attacks[index] {
				fail = 1
				break
			}
		}

		if fail == 0 {
			return magicNumber
		}
	}

	fmt.Println("Magic number search failed")
	return 0
}

// InitMagicNumbers initializes the magic numbers for bishops and rooks
func InitMagicNumbers() {
	//fmt.Println("Bishop magic numbers:")
	for square := 0; square < 40; square++ {
		// init Rook magic numbers
		fmt.Printf("0x%x,\n", FindMagicNumber(square, globals.BishopRelevantOccupancyCount[square], globals.BISHOP))
		globals.BishopMagicNumbers[square] = FindMagicNumber(square, globals.BishopRelevantOccupancyCount[square], globals.BISHOP)
	}
	fmt.Println()
	//fmt.Println("Rook magic numbers:")
	for square := 0; square < 40; square++ {
		// init Rook magic numbers
		fmt.Printf("0x%x,\n", FindMagicNumber(square, globals.RookRelevantOccupancyCount[square], globals.ROOK))
		globals.RookMagicNumbers[square] = FindMagicNumber(square, globals.RookRelevantOccupancyCount[square], globals.ROOK)
	}
}
