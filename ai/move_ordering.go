package ai

import (
	"fmt"
	"zerginator/board"
	"zerginator/globals"
)

/*
	In this file an implementation of move ordering to enhance the search efficiency of the negamax
	search. This will help create beta cutoffs earlier in the search tree. This also holds the killer
	move and history heuristic tables. The transposition table is also implemented here.

	This incorporates the most valuable victim - least valuable attacker (MVV-LVA) heuristic
		(Victim)|  Pawn | Knight | Bishop | Rook | King
	(Attacker)	|
		Pawn	|	104		204		304		404		504
		Knight	|	103		203		303		403		503
		Bishop	|	102		202		302		402		502
		Rook	|	101		201		301		401		501
		King	|	100		200		300		400		500

	A typical ordering of moves would be:
	1. TT entry
	2. PV move
	3. Captures (ordered by MVV-LVA)
	4. 1st Killer move
	5. 2nd Killer move
	6. History heuristic moves
	7. Unsorted moves
*/

var MvvLvaScores = [6][10]int{
	{104, 204, 304, 404, 504, 104, 204, 304, 404, 504},
	{103, 203, 303, 403, 503, 103, 203, 203, 403, 503},
	{102, 202, 302, 402, 502, 102, 303, 302, 402, 502},
	{101, 201, 301, 401, 501, 101, 401, 401, 401, 501},
	{100, 200, 300, 400, 500, 104, 503, 502, 400, 500},
	{104, 204, 304, 404, 504, 104, 202, 303, 404, 504}}

// KillerMoves stores the killer moves with [id][ply]
var KillerMoves [2][64]uint64

// HistoryHeuristic stores the history heuristic scores for moves with [piece][square]
var HistoryHeuristic [6][40]uint64

func ScoreMove(move uint64) int {
	// score the principle variation move highest if we are following the PV
	if ScorePV {
		if PVTable[0][globals.Ply] == move {
			ScorePV = false
			return 20000
		}
	}
	if board.GetMoveCapturedPiece(move) >= globals.WhitePawn && board.GetMoveCapturedPiece(move) <= globals.BlackPawn {
		// return the MVV-LVA score [source_square][target_piece]
		return MvvLvaScores[board.GetMovePiece(move)][board.GetMoveCapturedPiece(move)] + 10000
	} else {
		// score quiet moves
		if KillerMoves[0][globals.Ply] == move {
			return 9000
		} else if KillerMoves[1][globals.Ply] == move {
			return 8000
		} else {
			return int(HistoryHeuristic[board.GetMovePiece(move)][board.GetMoveTarget(move)])
		}
	}
}

// EnablePVScore enables the principal variation scoring
func EnablePVScore(moveList *board.Moves) {
	FollowPV = false
	for i := 0; i < moveList.Count; i++ {
		if PVTable[0][globals.Ply] == moveList.Moves[i] {
			ScorePV = true
			FollowPV = true
		}
	}
}

func PrintMoveScores(moveList *board.Moves) {
	fmt.Printf("\n\tMove | Score\n")
	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i]
		board.PrintMove(move)
		fmt.Printf(" | %d\n", ScoreMove(move))
	}
}

func OrderMoves(moveList *board.Moves, bestMove uint64) {
	var moveScores []int
	for i := 0; i < moveList.Count; i++ {
		if bestMove == moveList.Moves[i] {
			moveScores = append(moveScores, 30000)
		} else {
			moveScores = append(moveScores, ScoreMove(moveList.Moves[i]))
		}
	}
	// simple bubble sort, as move lists are short
	for i := 0; i < moveList.Count; i++ {
		swapped := false
		for j := 0; j < moveList.Count-i-1; j++ {
			if moveScores[j] < moveScores[j+1] {
				// swap scores
				moveScores[j], moveScores[j+1] = moveScores[j+1], moveScores[j]
				// swap moves
				moveList.Moves[j], moveList.Moves[j+1] = moveList.Moves[j+1], moveList.Moves[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

// noHashEntry indicates that there is no valid entry in the hash table
var noHashEntry = 10000000

// hash table size
var hashSize = 0x400000 // 4 million entries

// Transposition Table flags
const (
	HashFlagExact = iota
	HashFlagAlpha
	HashFlagBeta
)

// taggedHashEntry represents an entry in the transposition table
type taggedHashEntry struct {
	key      uint64 // position hash key
	depth    int    // current search depth
	flags    int    // node flag: score>=beta, score<=alpha, score>alpha
	value    int    // score for the position
	bestMove uint64 // best move for the position
}

// transpositionTable is the transposition table
var transpositionTable = make([]taggedHashEntry, hashSize)

// ClearTranspositionTable clears the transposition table
func ClearTranspositionTable() {
	for i := 0; i < hashSize; i++ {
		// reset entry
		transpositionTable[i].key = 0
		transpositionTable[i].depth = 0
		transpositionTable[i].flags = 0
		transpositionTable[i].value = 0
	}
}

// ProbeTranspositionTable checks the TT and returns either a stored value/bound or noHashEntry.
// It ensures the caller's bestMove pointer is updated when a TT entry (even shallow) exists.
func ProbeTranspositionTable(bestMove *uint64, depth int, alpha int, beta int) int {
	//Create a pointer to point to the entry in the transposition table based on the current board hash key
	hashEntry := &transpositionTable[globals.HashKey%uint64(hashSize)]
	if hashEntry.key == globals.HashKey {
		// provide the best move if it is not nil
		if bestMove != nil {
			*bestMove = hashEntry.bestMove
		}
		// only use stored value if it was searched to at least the requested depth
		if hashEntry.depth >= depth {
			switch hashEntry.flags {
			case HashFlagExact:
				return hashEntry.value
			case HashFlagAlpha:
				if hashEntry.value <= alpha {
					return alpha
				}
			case HashFlagBeta:
				if hashEntry.value >= beta {
					return beta
				}
			}
		}
	}
	return noHashEntry
}

// RecordHash records the hash entry in the transposition table
func RecordHash(bestMove uint64, depth int, value int, hashFlag int) {
	hashEntry := &transpositionTable[globals.HashKey%uint64(hashSize)]
	hashEntry.key = globals.HashKey
	hashEntry.depth = depth
	hashEntry.flags = hashFlag
	hashEntry.value = value
	hashEntry.bestMove = bestMove
}
