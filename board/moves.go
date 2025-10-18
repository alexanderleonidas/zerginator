package board

import (
	"fmt"
	"time"
	"zerginator/bitoperations"
	"zerginator/globals"
)

// Moves struct holds an array of moves and a count of the number of moves
type Moves struct {
	Moves [256]uint64
	Count int
}

// MoveRecord struct holds the information needed to undo a move
type MoveRecord struct {
	move            uint64
	enPassantSquare int
	sideToMove      int
	occupancies     [3]uint64
	hashKey         uint64
}

// MoveStack array that holds the move records for undoing moves
var MoveStack []MoveRecord

// AddMove adds a move to the move list
func (m *Moves) AddMove(move uint64) {
	m.Moves[m.Count] = move
	m.Count++
}

// PrintMove prints a move to the console
func PrintMove(move uint64) {
	fmt.Printf("%s%s%s", globals.SquareToCoord[GetMoveSource(move)], globals.SquareToCoord[GetMoveTarget(move)], globals.PromotedPieces[GetMovePromotedPiece(move)])
}

// PrintMoveList prints the move list to the console
func PrintMoveList(moveList *Moves) {
	if moveList.Count == 0 {
		fmt.Printf("\n\tNo moves available\n")
		return
	}
	fmt.Printf("\n\tmove	piece	capture	doublePawnPush	enPassant\n")
	for index := 0; index < moveList.Count; index++ {
		move := moveList.Moves[index]
		captured := ""
		if GetMoveCapturedPiece(move) <= globals.BlackPawn && GetMoveCapturedPiece(move) >= globals.WhitePawn {
			captured = globals.UnicodePieces[GetMoveCapturedPiece(move)]
		} else {
			captured = "-"
		}
		fmt.Printf("\t")
		PrintMove(move)
		fmt.Printf("\t%s\t\t%s\t\t%d\t\t\t\t%d\n", globals.UnicodePieces[GetMovePiece(move)], captured, GetMoveDoublePawnPush(move), GetMoveEnPassant(move))
	}
	fmt.Printf("\n\tTotal moves: %d\n", moveList.Count)
}

// EncodeMove encodes a move into a 64-bit unsigned integer
func EncodeMove(source int, target int, piece int, promotedPiece int, capturedPiece int, doublePawnPush int, enPassant int) uint64 {
	/*
	   These are the move elements that we need to encode in a binary representation:

	   	Binary representation				Description						Hexadecimal
	   0000 0000 0000 0000 0011 1111		source square (6 bits)			0x3f
	   0000 0000 0000 1111 1100 0000		target square (6 bits)			0xfc0
	   0000 0000 0111 0000 0000 0000		piece type (3 bits)				0x7000
	   0000 0011 1000 0000 0000 0000		promoted piece (3 bits)			0x38000
	   0001 1100 0000 0000 0000 0000		captured piece (3 bit)			0x1C0000
	   0010 0000 0000 0000 0000 0000		double pawn push flag (1 bit)	0x200000
	   0100 0000 0000 0000 0000 0000		en passant flag (1 bit)			0x400000
	*/
	return uint64(source | target<<6 | piece<<12 | promotedPiece<<15 | capturedPiece<<18 | doublePawnPush<<21 | enPassant<<22)
}

// DecodeMove decodes a move and prints it to the console
func DecodeMove(move uint64) {
	fmt.Printf("%s ", globals.UnicodePieces[GetMovePiece(move)])
	fmt.Printf("%s", globals.SquareToCoord[GetMoveSource(move)])
	fmt.Printf("%s", globals.SquareToCoord[GetMoveTarget(move)])
	if GetMovePromotedPiece(move) != 0 {
		fmt.Printf("%s ", globals.UnicodePieces[GetMovePromotedPiece(move)])
	}
	if GetMoveCapturedPiece(move) <= globals.BlackPawn {
		fmt.Printf(" x %s\n", globals.UnicodePieces[GetMoveCapturedPiece(move)])
	}
}

// GetMoveSource returns the source square of the move
func GetMoveSource(move uint64) int {
	return int(move & 0x3f)
}

// GetMoveTarget returns the target square of the move
func GetMoveTarget(move uint64) int {
	return int((move & 0xfc0) >> 6)
}

// GetMovePiece returns the piece type of the move
func GetMovePiece(move uint64) int {
	return int((move & 0x7000) >> 12)
}

// GetMovePromotedPiece returns the promoted piece of the move
func GetMovePromotedPiece(move uint64) int {
	return int((move & 0x38000) >> 15)
}

// GetMoveCapturedPiece returns the captured piece of the move
func GetMoveCapturedPiece(move uint64) int {
	return int((move & 0x1C0000) >> 18)
}

// GetMoveDoublePawnPush returns the double pawn push flag of the move
func GetMoveDoublePawnPush(move uint64) int {
	return int((move & 0x200000) >> 21)
}

// GetMoveEnPassant returns the en passant flag of the move
func GetMoveEnPassant(move uint64) int {
	return int((move & 0x400000) >> 22)
}

// GenerateMoves generates all the possible moves for the current board state
func GenerateMoves(moveList *Moves) {
	var sourceSquare, targetSquare int
	var bitboard, attacks uint64
	moveList.Count = 0

	for piece := globals.WhitePawn; piece <= globals.BlackPawn; piece++ {
		bitboard = globals.Bitboards[piece]
		if globals.SideToMove == globals.WHITE {
			// generate moves for white pawns
			if piece == globals.WhitePawn {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					targetSquare = sourceSquare - 5
					// generate quiet pawn moves
					if !(targetSquare < globals.A8) && bitoperations.GetBit(globals.Occupancies[globals.BOTH], targetSquare) == 0 {
						// pawn promotion territory
						if sourceSquare >= globals.A7 && sourceSquare <= globals.E7 {
							for i := globals.WhiteKnight; i <= globals.WhiteKing; i++ {
								moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, i, globals.NoPiece, 0, 0))
							}
						} else { // pawn move
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
							if sourceSquare >= globals.A2 && sourceSquare <= globals.E2 && bitoperations.GetBit(globals.Occupancies[globals.BOTH], targetSquare-5) == 0 {
								moveList.AddMove(EncodeMove(sourceSquare, targetSquare-5, piece, globals.NoPiece, globals.NoPiece, 1, 0))
							}
						}
					}
					attacks = globals.PawnAttacks[globals.SideToMove][sourceSquare] & globals.Occupancies[globals.BLACK]
					// generate pawn captures
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						if sourceSquare >= globals.A7 && sourceSquare <= globals.E7 {
							for i := globals.WhiteKnight; i <= globals.WhiteKing; i++ {
								moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, i, globals.BlackPawn, 0, 0))
							}
						} else {
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.BlackPawn, 0, 0))
						}
						bitoperations.PopBit(&attacks, targetSquare)
					}
					// pop the least significant bit in the bitboard
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			} else if piece == globals.WhiteKnight {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					attacks = globals.KnightAttacks[sourceSquare] & ^globals.Occupancies[globals.WHITE]
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						// knight quiet move
						if bitoperations.GetBit(globals.Occupancies[globals.BLACK], targetSquare) == 0 {
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
						} else {
							// knight capture
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.BlackPawn, 0, 0))
						}
						bitoperations.PopBit(&attacks, targetSquare)
					}
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			} else if piece == globals.WhiteKing {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					attacks = globals.KingAttacks[sourceSquare] & ^globals.Occupancies[globals.WHITE]
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						// king quiet move
						if bitoperations.GetBit(globals.Occupancies[globals.BLACK], targetSquare) == 0 {
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
						} else {
							// king capture
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.BlackPawn, 0, 0))
						}
						bitoperations.PopBit(&attacks, targetSquare)
					}
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			} else if piece == globals.WhiteBishop {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					attacks = GetBishopAttacks(sourceSquare, globals.Occupancies[globals.BOTH]) & ^globals.Occupancies[globals.WHITE]
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						// knight quiet move
						if bitoperations.GetBit(globals.Occupancies[globals.BLACK], targetSquare) == 0 {
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
						} else {
							// knight capture
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.BlackPawn, 0, 0))
						}
						bitoperations.PopBit(&attacks, targetSquare)
					}
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			} else if piece == globals.WhiteRook {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					attacks = GetRookAttacks(sourceSquare, globals.Occupancies[globals.BOTH]) & ^globals.Occupancies[globals.WHITE]
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						// knight quiet move
						if bitoperations.GetBit(globals.Occupancies[globals.BLACK], targetSquare) == 0 {
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
						} else {
							// knight capture
							moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.BlackPawn, 0, 0))
						}
						bitoperations.PopBit(&attacks, targetSquare)
					}
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			}
		} else {
			// generate moves for black pawns
			if piece == globals.BlackPawn {
				for bitboard != 0 {
					sourceSquare = bitoperations.GetLeastSignificantBitIndex(bitboard)
					targetSquare = sourceSquare + 5
					// generate quiet pawn moves
					if !(targetSquare > globals.E1) && bitoperations.GetBit(globals.Occupancies[globals.BOTH], targetSquare) == 0 {
						// pawn move
						moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, globals.NoPiece, 0, 0))
					}
					attacks = globals.PawnAttacks[globals.SideToMove][sourceSquare] & globals.Occupancies[globals.WHITE]
					// generate pawn captures
					for attacks != 0 {
						targetSquare = bitoperations.GetLeastSignificantBitIndex(attacks)
						capturedPiece := 7 // initialize to an invalid piece
						// loop over the opposite sides pieces
						for p := globals.WhitePawn; p <= globals.WhiteKing; p++ {
							if bitoperations.GetBit(globals.Bitboards[p], targetSquare) == 1 {
								capturedPiece = p
								break
							}
						}
						moveList.AddMove(EncodeMove(sourceSquare, targetSquare, piece, globals.NoPiece, capturedPiece, 0, 0))
						bitoperations.PopBit(&attacks, targetSquare)
					}
					if globals.EnPassantSquare != globals.NoSquare {
						enPassantAttacks := globals.PawnAttacks[globals.SideToMove][sourceSquare] & (1 << globals.EnPassantSquare)
						if enPassantAttacks != 0 && ((1<<(globals.EnPassantSquare-5))&globals.Bitboards[globals.WhitePawn]) != 0 {
							targetEnPassantSquare := bitoperations.GetLeastSignificantBitIndex(enPassantAttacks)
							moveList.AddMove(EncodeMove(sourceSquare, targetEnPassantSquare, piece, globals.NoPiece, globals.WhitePawn, 0, 1))
						}
					}
					// pop the least significant bit in the bitboard
					bitoperations.PopBit(&bitboard, sourceSquare)
				}
			}
		}
	}
}

func MakeMove(move uint64, moveFlag int) int {
	// quiet moves
	if moveFlag == globals.AllMoves {
		// preserve the current state for undoing moves
		sourceSquare := GetMoveSource(move)
		targetSquare := GetMoveTarget(move)
		piece := GetMovePiece(move)
		promotedPiece := GetMovePromotedPiece(move)
		doublePawnPush := GetMoveDoublePawnPush(move)
		capturedPiece := GetMoveCapturedPiece(move)
		enPassant := GetMoveEnPassant(move)
		// record the move in the stack
		MoveStack = append(MoveStack, MoveRecord{move, globals.EnPassantSquare, globals.SideToMove, globals.Occupancies, globals.HashKey})

		// move the piece
		bitoperations.PopBit(&globals.Bitboards[piece], sourceSquare)
		bitoperations.SetBit(&globals.Bitboards[piece], targetSquare)

		// update the hash key
		globals.HashKey ^= PieceKeys[piece][sourceSquare] // remove piece from source square
		globals.HashKey ^= PieceKeys[piece][targetSquare] // add piece to target square

		//if there is a captured piece, remove it from the board
		if capturedPiece != globals.NoPiece {
			if bitoperations.GetBit(globals.Bitboards[capturedPiece], targetSquare) == 1 {
				bitoperations.PopBit(&globals.Bitboards[capturedPiece], targetSquare)
				// remove captured piece from hash key
				globals.HashKey ^= PieceKeys[capturedPiece][targetSquare]
			}
		}
		// if there is a promotion, remove the piece from the board and add the promoted piece
		if promotedPiece <= globals.WhiteKing && promotedPiece >= globals.WhiteKnight {
			bitoperations.PopBit(&globals.Bitboards[piece], targetSquare)
			globals.HashKey ^= PieceKeys[piece][targetSquare]
			bitoperations.SetBit(&globals.Bitboards[promotedPiece], targetSquare)
			globals.HashKey ^= PieceKeys[promotedPiece][targetSquare]
		}
		// if there is an en passant capture, remove the pawn from the board
		if enPassant != 0 {
			// only black side can perform en passant capture
			if globals.SideToMove == globals.BLACK {
				bitoperations.PopBit(&globals.Bitboards[globals.WhitePawn], targetSquare-5)
				globals.HashKey ^= PieceKeys[globals.WhitePawn][targetSquare-5]
			}
		}
		// hash en passant if available (remove enpassant square from hash key)
		if globals.EnPassantSquare != globals.NoSquare {
			globals.HashKey ^= EnPassantKeys[globals.EnPassantSquare]
		}
		globals.EnPassantSquare = globals.NoSquare
		if doublePawnPush != 0 && globals.SideToMove == globals.WHITE {
			globals.EnPassantSquare = targetSquare + 5
			// hash the en passant square
			globals.HashKey ^= EnPassantKeys[targetSquare+5]
		}
		// update occupancy bitboards
		globals.Occupancies = [3]uint64{0, 0, 0}
		for i := globals.WhitePawn; i <= globals.WhiteKing; i++ {
			globals.Occupancies[globals.WHITE] |= globals.Bitboards[i]
		}
		globals.Occupancies[globals.BLACK] |= globals.Bitboards[globals.BlackPawn]
		globals.Occupancies[globals.BOTH] |= globals.Occupancies[globals.WHITE] | globals.Occupancies[globals.BLACK]
		globals.SideToMove ^= 1
		globals.HashKey ^= SideKey // hash the side

		//// debugging for hash keys
		//hasFromScratch := GeneratePositionKey()
		//// If the hash keys do not match the incremental hash, interrupt execution
		//if hasFromScratch != globals.HashKey {
		//	fmt.Printf("\n\nMake Move!\n")
		//	fmt.Printf("move: ")
		//	PrintMove(move)
		//	PrintBoard()
		//	fmt.Printf("Hash keys should be: %x\n", hasFromScratch)
		//	globals.Scanner.Scan()
		//}
		return 1 // move made successfully
	} else {
		// capture moves
		capturedPiece := GetMoveCapturedPiece(move)
		if capturedPiece <= globals.BlackPawn && capturedPiece >= globals.WhitePawn {
			return MakeMove(move, globals.AllMoves)
		} else {
			return 0
		}
	}
}

func UnMakeMove() {
	// pop the last move from the stack
	rec := MoveStack[len(MoveStack)-1]
	MoveStack = MoveStack[:len(MoveStack)-1]
	move := rec.move
	capturedPiece := GetMoveCapturedPiece(move)
	sourceSquare := GetMoveSource(move)
	targetSquare := GetMoveTarget(move)
	piece := GetMovePiece(move)
	promotedPiece := GetMovePromotedPiece(move)
	enPassant := GetMoveEnPassant(move)

	// restore game variables
	globals.SideToMove = rec.sideToMove
	globals.EnPassantSquare = rec.enPassantSquare
	globals.Occupancies = rec.occupancies
	globals.HashKey = rec.hashKey

	// undo promotion or normal move
	if promotedPiece <= globals.WhiteKing && promotedPiece >= globals.WhiteKnight {
		bitoperations.PopBit(&globals.Bitboards[promotedPiece], targetSquare)
		bitoperations.SetBit(&globals.Bitboards[piece], sourceSquare)
	} else {
		bitoperations.PopBit(&globals.Bitboards[piece], targetSquare)
		bitoperations.SetBit(&globals.Bitboards[piece], sourceSquare)
	}

	// restore captured piece (normal capture)
	if capturedPiece != globals.NoPiece {
		// For en passant, restore pawn on correct square
		if enPassant != 0 {
			// only black side can perform en passant capture
			if globals.SideToMove == globals.BLACK {
				bitoperations.SetBit(&globals.Bitboards[globals.WhitePawn], targetSquare-5)
			}
		} else {
			bitoperations.SetBit(&globals.Bitboards[capturedPiece], targetSquare)
		}
	}

	//hasFromScratch := GeneratePositionKey()
	//// If the hash keys do not match the incremental hash, interrupt execution
	//if hasFromScratch != globals.HashKey {
	//	fmt.Printf("\n\nUnMake Move!\n")
	//	fmt.Printf("move: ")
	//	PrintMove(move)
	//	PrintBoard()
	//	fmt.Printf("Hash keys should be: %x\n", hasFromScratch)
	//	globals.Scanner.Scan()
	//}
}

// PerftDriver is a recursive procedure that counts the number of leaf nodes in the move tree up to the given depth
func PerftDriver(depth int) {
	globals.NodesVisited++
	if depth <= 0 {
		// count the nodes
		globals.LeafNodesVisited++
		return
	}
	moveList := Moves{}
	GenerateMoves(&moveList)
	for i := 0; i < moveList.Count; i++ {
		//b, o, s, e := CopyBoard()
		if MakeMove(moveList.Moves[i], globals.AllMoves) == 0 {
			continue // skip to next move if it is illegal
		}
		PerftDriver(depth - 1)
		//RestoreBoard(b, o, s, e)
		UnMakeMove()
	}
}

// PerftTest performs a performance test of the move generation and move making functions
func PerftTest(depth int) {
	globals.NodesVisited = 0
	globals.LeafNodesVisited = 0
	fmt.Println("\t--- Performance test ---")
	fmt.Println("\tMove \tNodes")
	moveList := Moves{}
	GenerateMoves(&moveList)
	startTime := time.Now()
	for i := 0; i < moveList.Count; i++ {
		//b, o, s, e := CopyBoard()
		fmt.Printf("\t")
		PrintMove(moveList.Moves[i])
		if MakeMove(moveList.Moves[i], globals.AllMoves) == 0 {
			continue // skip to next move if it is illegal
		}
		cumulativeNodes := globals.LeafNodesVisited
		PerftDriver(depth - 1)
		oldNodes := globals.LeafNodesVisited - cumulativeNodes
		//RestoreBoard(b, o, s, e)
		UnMakeMove()
		fmt.Printf("\t%d", oldNodes)
		fmt.Println()
	}
	elapsed := time.Since(startTime)
	fmt.Printf("\tDepth: %d\n", depth)
	fmt.Printf("\tLeaf Nodes: %d\n", globals.LeafNodesVisited)
	fmt.Printf("\tTotal Nodes: %d\n", globals.NodesVisited)
	fmt.Printf("\tTime: %s\n", elapsed)
}
