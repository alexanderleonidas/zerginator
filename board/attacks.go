package board

import (
	"zerginator/bitoperations"
	"zerginator/globals"
)

// MaskPawnAttacks returns the bitboard of all the pawn attacks on the given square
func MaskPawnAttacks(side int, square int) uint64 {
	// results/attacks bitboard
	var attacks uint64 = 0
	// piece bitboard
	var bitboard uint64 = 0

	// set piece on board
	bitoperations.SetBit(&bitboard, square)

	if side == globals.WHITE {
		// handle off-board captures and generate attacks
		if (bitboard>>4)&globals.NotAFile != 0 {
			attacks |= bitboard >> 4
		}
		if (bitboard>>6)&globals.NotEFile != 0 {
			attacks |= bitboard >> 6
		}
	} else {
		if (bitboard<<4)&globals.NotEFile != 0 {
			attacks |= bitboard << 4
		}
		if (bitboard<<6)&globals.NotAFile != 0 {
			attacks |= bitboard << 6
		}
	}
	return attacks & globals.Bitboard40Mask
}

// MaskKnightAttacks returns the bitboard of all the knight attacks on the given square
func MaskKnightAttacks(square int) uint64 {
	var attacks uint64 = 0
	var bitboard uint64 = 0
	bitoperations.SetBit(&bitboard, square)
	if (bitboard>>11)&globals.NotEFile != 0 {
		attacks |= bitboard >> 11
	}
	if (bitboard>>9)&globals.NotAFile != 0 {
		attacks |= bitboard >> 9
	}
	if (bitboard>>7)&globals.NotDEFile != 0 {
		attacks |= bitboard >> 7
	}
	if (bitboard>>3)&globals.NotABFile != 0 {
		attacks |= bitboard >> 3
	}
	if (bitboard<<11)&globals.NotAFile != 0 {
		attacks |= bitboard << 11
	}
	if (bitboard<<9)&globals.NotEFile != 0 {
		attacks |= bitboard << 9
	}
	if (bitboard<<7)&globals.NotABFile != 0 {
		attacks |= bitboard << 7
	}
	if (bitboard<<3)&globals.NotDEFile != 0 {
		attacks |= bitboard << 3
	}
	return attacks & globals.Bitboard40Mask
}

// MaskKingAttacks returns the bitboard of all the king attacks on the given square
func MaskKingAttacks(square int) uint64 {
	var attacks uint64 = 0
	var bitboard uint64 = 0
	bitoperations.SetBit(&bitboard, square)
	if (bitboard>>1)&globals.NotEFile != 0 {
		attacks |= bitboard >> 1
	}
	if (bitboard>>4)&globals.NotAFile != 0 {
		attacks |= bitboard >> 4
	}
	if (bitboard >> 5) != 0 {
		attacks |= bitboard >> 5
	}
	if (bitboard>>6)&globals.NotEFile != 0 {
		attacks |= bitboard >> 6
	}
	if (bitboard<<1)&globals.NotAFile != 0 {
		attacks |= bitboard << 1
	}
	if (bitboard<<4)&globals.NotEFile != 0 {
		attacks |= bitboard << 4
	}
	if (bitboard << 5) != 0 {
		attacks |= bitboard << 5
	}
	if (bitboard<<6)&globals.NotAFile != 0 {
		attacks |= bitboard << 6
	}
	return attacks & globals.Bitboard40Mask
}

// InitLeapersAttacks initializes the pawn, king, and knight attacks tables
func InitLeapersAttacks() {
	for square := 0; square < 40; square++ {
		globals.PawnAttacks[globals.WHITE][square] = MaskPawnAttacks(globals.WHITE, square)
		globals.PawnAttacks[globals.BLACK][square] = MaskPawnAttacks(globals.BLACK, square)
		globals.KnightAttacks[square] = MaskKnightAttacks(square)
		globals.KingAttacks[square] = MaskKingAttacks(square)
	}
}

// MaskBishopAttacks returns the bitboard of all the bishop attacks on the given square
func MaskBishopAttacks(square int) uint64 {
	var attacks uint64 = 0
	var targetRank, targetFile = square / 5, square % 5
	// mask the relevant bishop occupancy bits
	for r, f := targetRank+1, targetFile+1; r <= 6 && f <= 3; r, f = r+1, f+1 {
		attacks |= 1 << (r*5 + f)
	}
	for r, f := targetRank-1, targetFile+1; r >= 1 && f <= 3; r, f = r-1, f+1 {
		attacks |= 1 << (r*5 + f)
	}
	for r, f := targetRank+1, targetFile-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		attacks |= 1 << (r*5 + f)
	}
	for r, f := targetRank-1, targetFile-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		attacks |= 1 << (r*5 + f)
	}
	return attacks & globals.Bitboard40Mask
}

// MaskRookAttacks returns the bitboard of all the rook attacks on the given square
func MaskRookAttacks(square int) uint64 {
	var attacks uint64 = 0
	var targetRank, targetFile = square / 5, square % 5
	// mask the relevant rook occupancy bits
	for r := targetRank + 1; r <= 6; r++ {
		attacks |= 1 << (r*5 + targetFile)
	}
	for r := targetRank - 1; r >= 1; r-- {
		attacks |= 1 << (r*5 + targetFile)
	}
	for f := targetFile + 1; f <= 3; f++ {
		attacks |= 1 << (targetRank*5 + f)
	}
	for f := targetFile - 1; f >= 1; f-- {
		attacks |= 1 << (targetRank*5 + f)
	}
	return attacks & globals.Bitboard40Mask
}

// BishopAttacksOnTheFly generates bishop attacks taking into account a piece block
func BishopAttacksOnTheFly(square int, block uint64) uint64 {
	var attacks uint64 = 0
	var targetRank, targetFile = square / 5, square % 5
	// generate bishop attacks
	for r, f := targetRank+1, targetFile+1; r <= 7 && f <= 4; r, f = r+1, f+1 {
		attacks |= 1 << (r*5 + f)
		if (1<<(r*5+f))&block != 0 {
			break
		}
	}
	for r, f := targetRank-1, targetFile+1; r >= 0 && f <= 4; r, f = r-1, f+1 {
		attacks |= 1 << (r*5 + f)
		if (1<<(r*5+f))&block != 0 {
			break
		}
	}
	for r, f := targetRank+1, targetFile-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= 1 << (r*5 + f)
		if (1<<(r*5+f))&block != 0 {
			break
		}
	}
	for r, f := targetRank-1, targetFile-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= 1 << (r*5 + f)
		if (1<<(r*5+f))&block != 0 {
			break
		}
	}
	return attacks & globals.Bitboard40Mask
}

// RookAttacksOnTheFly generates rook attacks taking into account a piece block
func RookAttacksOnTheFly(square int, block uint64) uint64 {
	var attacks uint64 = 0
	var targetRank, targetFile = square / 5, square % 5
	// mask the relevant rook occupancy bits
	for r := targetRank + 1; r <= 7; r++ {
		attacks |= 1 << (r*5 + targetFile)
		if (1<<(r*5+targetFile))&block != 0 {
			break
		}
	}
	for r := targetRank - 1; r >= 0; r-- {
		attacks |= 1 << (r*5 + targetFile)
		if (1<<(r*5+targetFile))&block != 0 {
			break
		}
	}
	for f := targetFile + 1; f <= 4; f++ {
		attacks |= 1 << (targetRank*5 + f)
		if (1<<(targetRank*5+f))&block != 0 {
			break
		}
	}
	for f := targetFile - 1; f >= 0; f-- {
		attacks |= 1 << (targetRank*5 + f)
		if (1<<(targetRank*5+f))&block != 0 {
			break
		}
	}

	return attacks & globals.Bitboard40Mask
}

// SetOccupancy sets the occupancy of a square on the board and returns the occupancy
func SetOccupancy(index int, maskBitCount int, attackMask uint64) uint64 {
	/*
		This function creates the occupancy bitboard for a given attack pattern. The index parameter serves
		as a mask to determine which squares of the attack pattern are occupied. For example, if the attack
		pattern consists of 4 squares, and the index parameter is 0b0001, then the first square of the attack
		pattern is considered blocked. Likewise, if index=0b0101 that means that the first and third squares are
		blocked.
	*/
	var occupancy uint64
	for i := 0; i < maskBitCount; i++ {
		square := bitoperations.GetLeastSignificantBitIndex(attackMask)
		// pop the least significant bit in the attack_mask
		bitoperations.PopBit(&attackMask, square)
		if (index & (1 << i)) != 0 {
			occupancy |= 1 << square
		}
	}
	return occupancy & globals.Bitboard40Mask
}

// InitSlidersAttacks initializes the sliders attacks tables
func InitSlidersAttacks(isBishop int) {
	// init bishop and rook attacks
	for square := 0; square < 40; square++ {
		globals.BishopMasks[square] = MaskBishopAttacks(square)
		globals.RookMasks[square] = MaskRookAttacks(square)
		var attackMask uint64
		if isBishop == globals.BISHOP {
			attackMask = globals.BishopMasks[square]

		} else {
			attackMask = globals.RookMasks[square]
		}
		relevantBitCount := bitoperations.CountBits(attackMask)
		occupancyIndices := 1 << relevantBitCount
		for index := 0; index < occupancyIndices; index++ {
			if isBishop == globals.BISHOP {
				occupancy := SetOccupancy(index, relevantBitCount, attackMask)
				magicIndex := (occupancy * globals.BishopMagicNumbers[square]) >> (64 - globals.BishopRelevantOccupancyCount[square])
				globals.BishopAttacks[square][magicIndex] = BishopAttacksOnTheFly(square, occupancy)
			} else {
				occupancy := SetOccupancy(index, relevantBitCount, attackMask)
				magicIndex := (occupancy * globals.RookMagicNumbers[square]) >> (64 - globals.RookRelevantOccupancyCount[square])
				globals.RookAttacks[square][magicIndex] = RookAttacksOnTheFly(square, occupancy)
			}
		}
	}
}

// GetBishopAttacks returns the bitboard of all the bishop attacks on the given square
func GetBishopAttacks(square int, occupancy uint64) uint64 {
	// Here we get the bishop attacks for the current board occupancies
	occupancy &= globals.BishopMasks[square]
	occupancy *= globals.BishopMagicNumbers[square]
	occupancy >>= 64 - globals.BishopRelevantOccupancyCount[square]
	return globals.BishopAttacks[square][occupancy]
}

// GetRookAttacks returns the bitboard of all the rook attacks on the given square
func GetRookAttacks(square int, occupancy uint64) uint64 {
	// Here we get the bishop attacks for the current board occupancies
	occupancy &= globals.RookMasks[square]
	occupancy *= globals.RookMagicNumbers[square]
	occupancy >>= 64 - globals.RookRelevantOccupancyCount[square]
	return globals.RookAttacks[square][occupancy]
}
