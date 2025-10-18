package ai

import (
	"fmt"
	"time"
	"zerginator/board"
	"zerginator/globals"
)

// MaxPly is the maximum ply depth
const MaxPly int = 64

// FollowPV indicates if we are following the principal variation
var FollowPV bool

// ScorePV indicates if we score of the principal variation node
var ScorePV bool

// PVLength is the length of the principal variation
var PVLength [MaxPly]int

// PVTable is the principal variation table [ply][ply]
var PVTable [MaxPly][MaxPly]uint64

// FullDepthMoves is the number of moves to search at full depth
const FullDepthMoves int = 4

// ReductionLimit is the maximum depth reduction for late move reductions
const ReductionLimit int = 3

// BestMove is the best move found so far
var BestMove uint64

// SearchPosition performs a search to find the best move for the current position
func SearchPosition(depth int) {
	// clear the helper data
	FollowPV = false
	ScorePV = false
	globals.Ply = 0
	globals.NodesVisited = -1 // -1 to not count the root node
	//globals.Stopped = false
	KillerMoves = [2][64]uint64{}
	HistoryHeuristic = [6][40]uint64{}
	PVTable = [MaxPly][MaxPly]uint64{}
	PVLength = [MaxPly]int{}

	alpha := -50000
	beta := 50000
	startTime := time.Now()
	// Iterative Deepening
	for d := 1; d <= depth; d++ {
		//if globals.Stopped {
		//	break // break if time is up
		//}
		FollowPV = true
		value := -negamax(d, alpha, beta)
		// Aspiration Windows
		if value <= alpha || value >= beta {
			// we are outside the window, so try again with a full window
			value = -negamax(d, -50000, 50000)
		}
		alpha = value - 50
		beta = value + 50
		elapsed := time.Since(startTime)
		fmt.Printf("\ninfo score cp %d depth %d nodes %d time %dms pv ", value, d, globals.NodesVisited, elapsed.Milliseconds())
		for i := 0; i < PVLength[0]; i++ {
			board.PrintMove(PVTable[0][i])
			fmt.Printf(" ")
		}
	}
	BestMove = PVTable[0][0]
	moveList := board.Moves{}
	board.GenerateMoves(&moveList)
	OrderMoves(&moveList, 0)
	isLegal := false
	for i := 0; i < moveList.Count; i++ {
		if moveList.Moves[i] == BestMove {
			isLegal = true
			break
		}
	}
	if !isLegal && moveList.Count > 0 {
		// Fallback: pick the first legal move
		BestMove = moveList.Moves[0]
	}
	fmt.Printf("\nbestmove: ")
	board.PrintMove(BestMove)
	fmt.Println()
}

// negamax performs a search to the given depth with alpha-beta pruning
func negamax(depth int, alpha int, beta int) int {
	PVLength[globals.Ply] = globals.Ply
	var bestMove uint64 = 0 // best move found so far to store in TT
	score := ProbeTranspositionTable(&bestMove, depth, alpha, beta)
	hashFlag := HashFlagAlpha

	if globals.Ply != 0 && IsRepetition() {
		return 0 // return 0 if we found a repetition
	}
	if globals.Ply != 0 && score != noHashEntry && beta-alpha > 1 {
		return score
	}

	if score != noHashEntry && globals.Ply >= 1 {
		return score
	}
	if depth <= 0 {
		// run quiescence search here to avoid the horizon effect
		return quiescence(alpha, beta)
	}
	// check for maximum ply
	if globals.Ply > MaxPly-1 || board.IsTerminalPosition() {
		// we are too deep in the search tree
		return EvaluatePosition()
	}
	globals.NodesVisited++
	legalMoves := 0
	/* Null Move Pruning using reduced depth search.
	This asks, "If I do nothing here, can the opponent do anything?" We give the opponent a free try, and if our
	position is so good that we exceed beta, we can assume that we would exceed beta if we searched all our moves */
	if depth >= 3 && globals.Ply != 0 {
		a, b, c, d, e := board.CopyBoard()
		if globals.EnPassantSquare != globals.NoSquare {
			globals.HashKey ^= board.EnPassantKeys[globals.EnPassantSquare]
		}
		globals.EnPassantSquare = globals.NoSquare // remove en passant square
		globals.SideToMove ^= 1                    // switch side to move, giving the opponent a free move
		globals.HashKey ^= board.SideKey
		globals.Ply++
		globals.RepetitionIndex++
		globals.RepetitionTable[globals.RepetitionIndex] = globals.HashKey
		score = -negamax(depth-1-2, -beta, -beta+1) // null move search with d-1-R, R=2
		board.RestoreBoard(a, b, c, d, e)
		globals.Ply--
		globals.RepetitionIndex--
		//if globals.Stopped {
		//	return 0 // return 0 if time is up
		//}
		if score >= beta {
			return beta
		}
	}
	// generate all the children of the current position
	children := board.Moves{}
	board.GenerateMoves(&children)
	if FollowPV {
		EnablePVScore(&children)
	}
	// order the children by score
	OrderMoves(&children, bestMove)
	movesSearched := 0
	value := -100000
	for i := 0; i < children.Count; i++ {
		globals.Ply++
		globals.RepetitionIndex++
		globals.RepetitionTable[globals.RepetitionIndex] = globals.HashKey
		move := children.Moves[i]
		// make the move and check if it is legal
		if board.MakeMove(move, globals.AllMoves) == 0 {
			globals.Ply--
			globals.RepetitionIndex--
			continue // skip illegal moves
		}
		legalMoves++
		// Late Move Reductions
		if movesSearched == 0 {
			// if this is the first move, search it with a full window
			score = -negamax(depth-1, -beta, -alpha)
		} else {
			// condition to consider late move reductions
			if movesSearched >= FullDepthMoves && depth >= ReductionLimit {
				/* When doing our late move reductions, we hope that the moves we are reducing depths for
				would never produce a beta-cutoff */
				score = -negamax(depth-2, -alpha-1, -alpha)
			} else {
				score = alpha + 1
			}
			// Principle Variation Search
			if score > alpha {
				/* Once we find a move with a score between alpha and beta, the rest of the children are
				searched with a window (alpha, alpha+1) in the aim to prove they are no better. */
				score = -negamax(depth-1, -alpha-1, -alpha)
				if score > alpha && score < beta {
					/* If we find out that the algorithm is wrong, and that a later move is better than the
					first PV move, we re-search that move with the full window.*/
					score = -negamax(depth-1, -beta, -alpha)
				}
			}
		}
		value = max(value, score)
		board.UnMakeMove()
		//if globals.Stopped {
		//	return 0 // return 0 if time is up
		//}
		globals.Ply--
		globals.RepetitionIndex--
		movesSearched++
		// found a better move
		if value > alpha {
			hashFlag = HashFlagExact
			bestMove = move
			// on quiet moves, update the history heuristic
			if board.GetMoveCapturedPiece(move) == globals.NoPiece {
				HistoryHeuristic[board.GetMovePiece(move)][board.GetMoveTarget(move)] += uint64(depth) * uint64(depth)
			}
			alpha = value
			PVTable[globals.Ply][globals.Ply] = move // store best move
			for nextPly := globals.Ply + 1; nextPly < PVLength[globals.Ply+1]; nextPly++ {
				// copy move from deeper ply to current ply
				PVTable[globals.Ply][nextPly] = PVTable[globals.Ply+1][nextPly]
			}
			PVLength[globals.Ply] = PVLength[globals.Ply+1]

			// beta cutoff
			if beta <= alpha {
				RecordHash(bestMove, depth, value, HashFlagBeta)
				if board.GetMoveCapturedPiece(move) == globals.NoPiece {
					// store killer move
					KillerMoves[1][globals.Ply] = KillerMoves[0][globals.Ply]
					KillerMoves[0][globals.Ply] = move
				}
				return beta
			}
		}
	}
	if legalMoves == 0 {
		// if the current player cannot move, the game ends in a draw
		return 0
	}
	RecordHash(bestMove, depth, value, hashFlag)
	return value
}

// quiescence performs a quiescence search to avoid the horizon effect
func quiescence(alpha int, beta int) int {
	globals.NodesVisited++
	// check for maximum ply
	if globals.Ply > MaxPly-1 || board.IsTerminalPosition() {
		// we are too deep in the search tree
		return EvaluatePosition()
	}
	evaluation := EvaluatePosition()
	if evaluation >= beta {
		return beta
	}
	if evaluation > alpha {
		alpha = evaluation
	}
	children := board.Moves{}
	board.GenerateMoves(&children)
	OrderMoves(&children, 0)
	for i := 0; i < children.Count; i++ {
		globals.Ply++
		globals.RepetitionIndex++
		globals.RepetitionTable[globals.RepetitionIndex] = globals.HashKey
		// make the move and check if it is legal
		if board.MakeMove(children.Moves[i], globals.OnlyCaptures) == 0 {
			globals.Ply--
			globals.RepetitionIndex--
			continue // skip illegal moves
		}
		alpha = max(alpha, -quiescence(-beta, -alpha))
		board.UnMakeMove()
		//if globals.Stopped {
		//	return 0 // return 0 if time is up
		//}
		globals.Ply--
		globals.RepetitionIndex--
		if beta <= alpha {
			return beta
		}
	}
	return alpha
}

func IsRepetition() bool {
	for i := 0; i < globals.RepetitionIndex; i++ {
		// if we found the hash key same with a current
		if globals.RepetitionTable[i] == globals.HashKey {
			return true
		}
	}
	// if no repetition found
	return false
}
