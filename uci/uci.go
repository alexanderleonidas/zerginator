package uci

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"zerginator/ai"
	"zerginator/board"
	"zerginator/clock"
	"zerginator/globals"
)

// TimeKeeper is the game clock
var TimeKeeper *clock.GameClock

// ParseMove takes a move in string format (e.g. "a2a4", "b7b8Q") and converts it to the internal move representation
func ParseMove(moveString string) uint64 {
	moveList := board.Moves{}
	board.GenerateMoves(&moveList)
	sourceSquare := (moveString[0] - 'a') + (8-(moveString[1]-'0'))*5
	targetSquare := (moveString[2] - 'a') + (8-(moveString[3]-'0'))*5
	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i]
		// make sure the source and target squares are available within the move list
		if int(sourceSquare) == board.GetMoveSource(move) && int(targetSquare) == board.GetMoveTarget(move) {
			promotedPiece := board.GetMovePromotedPiece(move)
			// if there is a promoted piece, make sure it matches the move string
			if promotedPiece <= globals.WhiteKing && promotedPiece >= globals.WhiteKnight {
				if len(moveString) == 5 {
					if promotedPiece == globals.WhiteKnight && moveString[4] == 'N' {
						return move
					} else if promotedPiece == globals.WhiteBishop && moveString[4] == 'B' {
						return move
					} else if promotedPiece == globals.WhiteKing && moveString[4] == 'K' {
						return move
					} else if promotedPiece == globals.WhiteRook && moveString[4] == 'R' {
						return move
					}
				}
				continue // continue loop if no legal move found
			}
			return move // legal move
		}
	}
	return 0 // return illegal move
}

// ParsePosition sets up the board position based on the UCI "position" command
func ParsePosition(command string) {
	/*
		This procedure sets up the board position based on the UCI "position" command.
		It handles both the "startpos" and "fen" options to initialize the board state.
		Examples of valid commands:
		- "position startpos"
		- "position startpos moves e2e4 e4e5 d2d4 b8c6"
		- "position fen ppppp/ppp1p/p2p1/Ppppp/1P3/1RN1P/2PPB/2K2 w -"
		- "position fen ppppp/ppp1p/p2p1/Ppppp/1P3/1RN1P/2PPB/2K2 w - moves e2e4 e4e5 d2d4 b8c6"
		- "position moves e2e4 e4e5 d2d4 b8c6"
		- "position undo"
	*/
	currentChar := 9
	command = command[currentChar:] // remove "position "
	if strings.HasPrefix(command, "startpos") {
		// initialize the board to the starting position
		board.ParseFEN(board.GetStartPosFEN())
	} else if strings.HasPrefix(command, "fen") {
		// initialize the board to the given FEN
		currentChar = strings.Index(command, "fen")
		if currentChar == -1 {
			board.ParseFEN(board.GetStartPosFEN())
		} else {
			currentChar += 4 // skip "fen "
			board.ParseFEN(command[currentChar:])
		}
	}
	// check for moves
	currentChar = strings.Index(command, "moves")
	if currentChar != -1 {
		currentChar += 6 // skip "moves "
		// loop over all moves in the command
		for _, move := range strings.Split(command[currentChar:], " ") {
			move = strings.TrimSpace(move)
			parsedMove := ParseMove(move)
			if parsedMove == 0 {
				break
			}
			globals.RepetitionIndex++
			globals.RepetitionTable[globals.RepetitionIndex] = globals.HashKey
			board.MakeMove(parsedMove, globals.AllMoves)
		}
	}

	// check for undo move
	currentChar = strings.Index(command, "undo")
	if currentChar != -1 {
		if len(board.MoveStack) > 0 {
			board.UnMakeMove()
		}
	}
	board.PrintBoard()
}

func ParseGo(command string) {
	/*
		This procedure parses the UCI "go" command to make the engine search for the best move. An example
		command is "go depth 6".
	*/
	depth := -1
	// look for "depth" in the command
	if idx := strings.Index(command, "depth"); idx != -1 {
		// parse the number after "depth "
		valStr := strings.TrimSpace(command[idx+6:])
		if n, err := strconv.Atoi(valStr); err == nil {
			depth = n
		}
	} else {
		// default depth
		depth = 13
	}
	// search position
	ai.SearchPosition(depth)
}

// MainUciLoop is the main loop that handles UCI commands
func MainUciLoop() {
	var input string
	// main loop
	fmt.Println("Zerginator 1.0")
	for globals.Scanner.Scan() {
		// reset input
		input = strings.TrimSpace(globals.Scanner.Text())
		if len(input) == 0 || input == "" || input[0] == '\n' {
			continue
		}
		switch {
		case strings.HasPrefix(input, "isready"):
			// engine is ready
			fmt.Printf("readyok\n")
			continue
		case strings.HasPrefix(input, "position"):
			ParsePosition(input)
			ai.ClearTranspositionTable()
		case strings.HasPrefix(input, "ucinewgame"):
			ParsePosition("position startpos")
			ai.ClearTranspositionTable()
		case strings.HasPrefix(input, "go"):
			ParseGo(input)
		case strings.HasPrefix(input, "uci"):
			fmt.Println("ID name: Zerginator 1.0")
			fmt.Println("uciok")
		case strings.HasPrefix(input, "startime"):
			TimeKeeper = clock.NewGameClock()
		case strings.HasPrefix(input, "quit"):
			return
		}
	}
	if err := globals.Scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
	}
}
