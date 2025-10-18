package board

import (
	"fmt"
	"math/rand"
	"zerginator/bitoperations"
	"zerginator/globals"
)

// PrintBitBoard prints the bitboard to the console
func PrintBitBoard(bitBoard uint64) {
	// To access this function in main.go, it must be capitalised as only these are exported
	fmt.Println()
	fmt.Println()
	// loop over the board ranks
	for rank := 0; rank < 8; rank++ {
		// loop over the board files
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			if file == 0 {
				fmt.Printf("\t%d ", 8-rank)
			}
			fmt.Printf("  %d", bitoperations.GetBit(bitBoard, square))
		}
		fmt.Println()
	}
	fmt.Printf("\t    A  B  C  D  E\n")
	fmt.Println("\tBitboard: ", bitBoard)
	fmt.Println()
}

// PrintBoard combines the bitboards of the pieces and prints them to the console
func PrintBoard() {
	fmt.Println()
	fmt.Println()
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			if file == 0 {
				fmt.Printf("\t%d ", 8-rank)
			}
			piece := -1
			for bbPiece := 0; bbPiece <= globals.BlackPawn; bbPiece++ {
				if bitoperations.GetBit(globals.Bitboards[bbPiece], square) == 1 {
					piece = bbPiece
				}
			}
			if piece == -1 {
				fmt.Printf("  .")
			} else {
				fmt.Printf("  %s", globals.UnicodePieces[piece])
				//fmt.Printf(" %s", AsciiPieces[piece])
			}
		}
		fmt.Println()
	}
	fmt.Printf("\t    A  B  C  D  E\n\n")
	if globals.SideToMove == globals.WHITE {
		fmt.Println("\tSide to move: White")
	} else {
		fmt.Println("\tSide to move: Black")
	}
	if globals.EnPassantSquare != globals.NoSquare {
		fmt.Println("\tEn-passant square:", globals.SquareToCoord[globals.EnPassantSquare])
	} else {
		fmt.Println("\tEn-passant square: None")
	}
	fmt.Printf("\tHash key: %x\n", globals.HashKey)
	//fmt.Println()
}

// ParseFEN parses the FEN string and populates the bitboards and state variables
func ParseFEN(fen string) {
	/*
		This parses a custom FEN string and populates the bitboards and state variables.
		The FEN string only contains the piece placement data, the side to move and en-passant square.
		For example, a FEN string where only a white pawn is on a1 would be "5/5/5/5/5/5/5/P4 - -"
	*/
	// reset board
	for i := 0; i <= globals.BlackPawn; i++ {
		globals.Bitboards[i] = 0
	}
	for i := 0; i < len(globals.Occupancies); i++ {
		globals.Occupancies[i] = 0
	}
	globals.SideToMove = globals.WHITE
	globals.EnPassantSquare = globals.NoSquare
	globals.RepetitionIndex = 0
	for i := 0; i < len(globals.RepetitionTable); i++ {
		globals.RepetitionTable[i] = 0
	}

	idx := 0 // index in fen string
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 5; file++ {
			if idx >= len(fen) {
				break
			}
			ch := rune(fen[idx])
			square := rank*5 + file
			// piece placement
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				piece := globals.ConvertAsciiToConstants[ch]
				bitoperations.SetBit(&globals.Bitboards[piece], square)
				idx++
				continue
			}
			// digit → empty squares
			if ch >= '0' && ch <= '9' {
				offset := int(ch - '0')
				file += offset - 1 // -1 because of the loop increment
				idx++
				continue
			}
		}
		// skip rank separator
		if idx < len(fen) && fen[idx] == '/' {
			idx++
		}
	}
	// side to move
	idx++
	if fen[idx] == 'w' {
		globals.SideToMove = globals.WHITE
	} else if fen[idx] == 'b' {
		globals.SideToMove = globals.BLACK
	}
	idx += 2
	if fen[idx] != '-' {
		file := int(fen[idx] - 'a')
		rank := 8 - int(fen[idx+1]-'0')
		globals.EnPassantSquare = rank*5 + file
	} else {
		globals.EnPassantSquare = globals.NoSquare
	}
	//fmt.Printf("'%s", fen[idx:])

	// init the occupancy bitboards
	for piece := globals.WhitePawn; piece <= globals.WhiteKing; piece++ {
		globals.Occupancies[globals.WHITE] |= globals.Bitboards[piece]
	}
	globals.Occupancies[globals.BLACK] |= globals.Bitboards[globals.BlackPawn]
	globals.Occupancies[globals.BOTH] |= globals.Occupancies[globals.WHITE] | globals.Occupancies[globals.BLACK]

	// init the hash key
	globals.HashKey = GeneratePositionKey()
}

// GetStartPosFEN returns a random starting position FEN string
func GetStartPosFEN() string {
	return "ppppp/ppppp/ppppp/5/5/5/PPPPP/" + globals.FenStartWhiteBottomRow[rand.Intn(len(globals.FenStartWhiteBottomRow))] + " w -"
}

// PrintAttackedSquares prints the attacked squares of the given side to the console
func PrintAttackedSquares(side int) {
	fmt.Println()
	fmt.Println()
	for rank := 0; rank < 8; rank++ {
		// loop over the board files
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			if file == 0 {
				fmt.Printf("\t%d ", 8-rank)
			}

			fmt.Printf("  %d", IsSquareAttacked(square, side))
		}
		fmt.Println()
	}
	fmt.Printf("\t    A  B  C  D  E\n")
	fmt.Println()
}

// IsSquareAttacked returns 1 if the given square is attacked by a piece of the given side, 0 otherwise
func IsSquareAttacked(square int, side int) int {

	/*
		Here we use another trick to check if a square is attacked by a piece of the given side.
		We make a bitwise AND between the black pawn attack mask on a given square and the white pawn bitboard.
		More formally, we take the intersection of the set of squares that are attacked by the black pawn on the
		given square and the set of squares that are occupied by the white pawn.
		For example, say there is a white pawn on c5, we know it will attack squares b6 and d6. A black pawn would
		have to be on b6 to attack c5. You can see that the attacks are antisymmetric. If we take the intersection
		of the set of squares attacked by the black pawn on b6 and the set of squares occupied by the white pawn,
		we get a non-empty set, which means that b6 is attacked by a white pawn. We do the same for d6.
	*/
	if side == globals.WHITE {
		if (globals.PawnAttacks[globals.BLACK][square]&globals.Bitboards[globals.WhitePawn]) != 0 ||
			(globals.KnightAttacks[square]&globals.Bitboards[globals.WhiteKnight]) != 0 ||
			(globals.KingAttacks[square]&globals.Bitboards[globals.WhiteKing]) != 0 ||
			(GetBishopAttacks(square, globals.Occupancies[globals.BOTH])&globals.Bitboards[globals.WhiteBishop]) != 0 ||
			(GetRookAttacks(square, globals.Occupancies[globals.BOTH])&globals.Bitboards[globals.WhiteRook]) != 0 {
			return 1
		}
	} else if side == globals.BLACK {
		if (globals.PawnAttacks[globals.WHITE][square] & globals.Bitboards[globals.BlackPawn]) != 0 {
			return 1
		}
	}
	return 0
}

// CopyBoard returns a copy of the current board state
func CopyBoard() ([6]uint64, [3]uint64, int, int, uint64) {
	var BitboardsCopy [6]uint64
	var OccupanciesCopy [3]uint64
	var SideToMoveCopy, EnPassantSquareCopy int
	var HashKeyCopy uint64
	copy(BitboardsCopy[:], globals.Bitboards[:])
	copy(OccupanciesCopy[:], globals.Occupancies[:])
	SideToMoveCopy = globals.SideToMove
	EnPassantSquareCopy = globals.EnPassantSquare
	HashKeyCopy = globals.HashKey
	return BitboardsCopy, OccupanciesCopy, SideToMoveCopy, EnPassantSquareCopy, HashKeyCopy
}

// RestoreBoard restores the board state from a copy
func RestoreBoard(bitboards [6]uint64, occupancies [3]uint64, sideToMove int, enPassantSquare int, hashKey uint64) {
	copy(globals.Bitboards[:], bitboards[:])
	copy(globals.Occupancies[:], occupancies[:])
	globals.SideToMove = sideToMove
	globals.EnPassantSquare = enPassantSquare
	globals.HashKey = hashKey
}

func IsTerminalPosition() bool {
	// Black side wins if a black paw reaches white bottom row
	for square := globals.A1; square <= globals.E1; square++ {
		if bitoperations.GetBit(globals.Bitboards[globals.BlackPawn], square) == 1 {
			return true
		}
	}
	// Black side wins if all white pieces are captured
	if globals.Occupancies[globals.WHITE] == 0 {
		return true
	}
	// White side wins if all white pieces are captured
	if globals.Occupancies[globals.BLACK] == 0 {
		return true
	}

	return false
}
