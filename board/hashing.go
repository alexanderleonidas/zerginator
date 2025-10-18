package board

import (
	"zerginator/bitoperations"
	"zerginator/globals"
)

// PieceKeys is the hash keys for each piece on each square [piece][square]
var PieceKeys [6][40]uint64

// EnPassantKeys is the hash keys for each en passant square
var EnPassantKeys [40]uint64

// SideKey is the hash key for the side to move
var SideKey uint64

// InitRandomKeys initializes the hash keys
func InitRandomKeys() {
	randomState = 1804289383
	for piece := globals.WhitePawn; piece <= globals.BlackPawn; piece++ {
		for square := 0; square < 40; square++ {
			PieceKeys[piece][square] = GetRandomUInt64()
			//fmt.Printf("%x\n", PieceKeys[piece][square])
		}
	}
	for square := 0; square < 40; square++ {
		EnPassantKeys[square] = GetRandomUInt64()
	}
	SideKey = GetRandomUInt64()
}

// GeneratePositionKey generates the hash key for the current position
func GeneratePositionKey() uint64 {
	var finalKey uint64
	var bitboard uint64
	for piece := globals.WhitePawn; piece <= globals.BlackPawn; piece++ {
		bitboard = globals.Bitboards[piece]
		for bitboard != 0 {
			square := bitoperations.GetLeastSignificantBitIndex(bitboard)
			// hash piece on square
			finalKey ^= PieceKeys[piece][square]
			bitoperations.PopBit(&bitboard, square)
		}
	}
	if globals.EnPassantSquare != globals.NoSquare {
		// hash en passant square
		finalKey ^= EnPassantKeys[globals.EnPassantSquare]
	}
	// hash the side only if it is black to move
	if globals.SideToMove == globals.BLACK {
		finalKey ^= SideKey
	}
	return finalKey
}
