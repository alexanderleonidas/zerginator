package ai

import (
	"zerginator/bitoperations"
	"zerginator/board"
	"zerginator/globals"
)

// EvaluatePosition evaluates the current board position and returns a score
func EvaluatePosition() int {
	score := 0
	for p := globals.WhitePawn; p <= globals.BlackPawn; p++ {
		bitboard := globals.Bitboards[p]
		for bitboard != 0 {
			piece := p
			square := bitoperations.GetLeastSignificantBitIndex(bitboard)
			// add the material value of the piece to the score
			score += globals.MaterialValues[piece]

			// add the positional value of the piece to the score
			switch {
			case p == globals.WhitePawn:
				// positional score
				score += globals.PawnPositionalValues[square]
				// double pawn penalty
				doublePawns := bitoperations.CountBits(globals.Bitboards[globals.WhitePawn] & globals.FileMasks[square])
				if doublePawns > 1 {
					score += globals.DoublePawnPenalty * doublePawns
				}
				// isolated pawn penalty
				if globals.Bitboards[globals.WhitePawn]&globals.IsolatedMasks[square] == 0 {
					score += globals.IsolatedPawnPenalty
				}
				// passed pawn bonus
				if globals.WhitePassedMasks[square]&globals.Bitboards[globals.BlackPawn] == 0 {
					score += globals.PassedPawnBonus[globals.GetRankFromSquare[square]]
				}
			case p == globals.WhiteKnight:
				score += globals.KnightPositionalValues[square]
			case p == globals.WhiteBishop:
				// positional score
				score += globals.BishopPositionalValues[square]
				// mobility score
				score += bitoperations.CountBits(board.GetBishopAttacks(square, globals.Occupancies[globals.BOTH]))
			case p == globals.WhiteRook:
				// positional score
				score += globals.RookPositionalValues[square]
				// semi-open file score
				if globals.Bitboards[globals.WhitePawn]&globals.FileMasks[square] == 0 {
					score += globals.SemiOpenFileScore
				}
				// open file score
				if (globals.Bitboards[globals.WhitePawn]|globals.Bitboards[globals.BlackPawn])&globals.FileMasks[square] == 0 {
					score += globals.OpenFileScore
				}
			case p == globals.WhiteKing:
				score += globals.KingPositionalValues[square]
			case p == globals.BlackPawn:
				// positional score
				score -= globals.PawnPositionalValues[globals.MirrorSquare[square]]
				// double pawn penalty
				doublePawns := bitoperations.CountBits(globals.Bitboards[globals.BlackPawn] & globals.FileMasks[square])
				if doublePawns > 1 {
					score -= globals.DoublePawnPenalty * doublePawns
				}
				// isolated pawn penalty
				if globals.Bitboards[globals.BlackPawn]&globals.IsolatedMasks[square] == 0 {
					score -= globals.IsolatedPawnPenalty
				}
				// win condition for black
				if square >= globals.A1 && square <= globals.E1 {
					score -= 50000
				}
				// passed pawn bonus
				if globals.BlackPassedMasks[square]&globals.Bitboards[globals.WhitePawn] == 0 {
					score -= globals.PassedPawnBonus[globals.GetRankFromSquare[globals.MirrorSquare[square]]]
				}
			}
			// pop the least significant bit
			bitoperations.PopBit(&bitboard, square)
		}
	}
	// no white pieces so black wins
	if globals.Occupancies[globals.WHITE] == 0 {
		score -= 50000
	}
	// no black pieces so white wins
	if globals.Occupancies[globals.BLACK] == 0 {
		score += 50000
	}
	if globals.SideToMove == globals.BLACK {
		score *= -1
	}
	return score
}

// SetFileRankMask returns the file and rank masks
func SetFileRankMask(fileNumber int, rankNumber int) uint64 {
	var mask uint64
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			if -1 != fileNumber && file == fileNumber {
				bitoperations.SetBit(&mask, square)
			} else if -1 != rankNumber && rank == rankNumber {
				bitoperations.SetBit(&mask, square)
			}
		}
	}

	return mask
}

// InitPawnEvaluationMasks initializes the pawn evaluation masks
func InitPawnEvaluationMasks() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			// set the file masks
			globals.FileMasks[square] |= SetFileRankMask(file, -1)
			// set the rank masks
			globals.RankMasks[square] |= SetFileRankMask(-1, rank)
			// set the isolated masks
			globals.IsolatedMasks[square] |= SetFileRankMask(file-1, -1)
			globals.IsolatedMasks[square] |= SetFileRankMask(file+1, -1)
		}
	}

	// set the passed masks after setting the other masks
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 5; file++ {
			square := rank*5 + file
			// set the black passed masks
			globals.BlackPassedMasks[square] |= SetFileRankMask(file-1, rank)
			globals.BlackPassedMasks[square] |= SetFileRankMask(file, rank)
			globals.BlackPassedMasks[square] |= SetFileRankMask(file+1, rank)
			for i := 0; i < rank+1; i++ {
				globals.BlackPassedMasks[square] &= ^globals.RankMasks[i*5+file]
			}
			// set the white passed masks
			globals.WhitePassedMasks[square] |= SetFileRankMask(file-1, -1)
			globals.WhitePassedMasks[square] |= SetFileRankMask(file, -1)
			globals.WhitePassedMasks[square] |= SetFileRankMask(file+1, -1)
			for i := 0; i < (8 - rank); i++ {
				globals.WhitePassedMasks[square] &= ^globals.RankMasks[(7-i)*5+file]
			}
		}
	}
}
