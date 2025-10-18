/*
This program was influenced by the chess programming resources available at:
https://github.com/maksimKorzh/chess_programming/blob/master/src/bbc/parse_fen/bbc.c#L429
https://github.com/bluefeversoft/vice
*/

// 5/pp1R1/5/1B3/5/4p/5/5 w -
// ppppp/ppppp/ppppp/5/5/5/PPPPP/N1KBR w -
// ppppp/ppppp/3pp/1pP2/5/5/2PPP/N1KBR w -
// ppppp/ppppp/p1ppp/1p3/2P2/1P3/P2PP/BKR1N w -
package main

import (
	"log"
	"zerginator/ai"
	"zerginator/board"
	"zerginator/globals"
	"zerginator/gui"
	"zerginator/uci"

	"github.com/hajimehoshi/ebiten/v2"
)

func initAll() {
	board.InitLeapersAttacks()
	//game.InitMagicNumbers()
	board.InitSlidersAttacks(globals.BISHOP)
	board.InitSlidersAttacks(globals.ROOK)
	board.InitRandomKeys()
	ai.ClearTranspositionTable()
	ai.InitPawnEvaluationMasks()
}

func main() {
	initAll()

	if err := gui.InitImages(); err != nil {
		log.Fatalf("failed to load images: %v", err)
	}

	debug := false
	graphics := true
	if graphics {
		ebiten.SetWindowSize(gui.ScreenWidth, gui.ScreenHeight)
		ebiten.SetWindowTitle("Zerginator 1.0")
		if err := ebiten.RunGame(&gui.Game{}); err != nil {
			log.Fatal(err)
		}
	} else if debug {
		board.ParseFEN(globals.FenDebugStartPosition)
		//board.ParseFEN("5/p1ppp/5/5/5/5/P1PPP/1R3 w -")
		board.PrintBoard()
		//fmt.Println("\tScore: ", ai.EvaluatePosition())
		uci.ParseGo("go depth 15")
		//board.PrintBoard()
		//for depth := 0; depth <= 0; depth++ {
		//	moves.PerftTest(depth)
		//}
	} else {
		uci.MainUciLoop()
	}

	/*
		Testing the UCI protocol implementation
	*/
	//board.ParseFEN("1p3/2P11/PR3/3Pp/1p3/2B2/P2p1/1NP1K w -")
	//board.PrintBoard()
	//move := uci.ParseMove("a2a4")
	//
	//if move != 0 {
	//	moves.MakeMove(move, globals.AllMoves)
	//	board.PrintBoard()
	//	moves.UnMakeMove()
	//	board.PrintBoard()
	//} else {
	//	fmt.Printf("Invalid move!\n")
	//}

	//testPos := "position startpos moves a2a4 a6a5 b2b4 b6b5 c2c4 c6c5 d2d4 d6d5 e2e4 e6e5"
	//testGo := "go depth 6"
	//game.ParsePosition(testPos)
	//game.PrintBoard()
	//game.ParseGo(testGo)
	////////////////////////////////////////////////////////////////////////////////
	/*
		Testing saving and restoring the game state as well as the Perft Driver for move generation
	*/

	//game.ParseFEN("1p3/2P11/PR3/3Pp/1p3/2B2/P2p1/1N1K w e6")
	//game.ParseFEN("5/5/5/2pP1/5/5/P4/5 w c6")
	//game.ParseFEN(game.FenStartPosition)
	//game.ParseFEN("1p3/2P11/PR3/3Pp/Pp3/2B2/12p1/1NP1K b a3")
	//
	//game.PrintBoard()
	//moveList := game.Moves{}
	//game.GenerateMoves(&moveList)
	//game.PrintMoveList(&moveList)
	//for i := 0; i < moveList.Count; i++ {
	//	fmt.Println("Making move: ", i)
	//	move := moveList.Moves[i]
	//	//b, o, s, e := game.CopyBoard()
	//	if game.MakeMove(move, game.AllMoves) == 0 {
	//		fmt.Println("Illegal Move!")
	//	}
	//	game.PrintBoard()
	//	//game.PrintBitBoard(game.Occupancies[game.BOTH])
	//	//game.Scanner.Scan()
	//	//game.RestoreBoard(b, o, s, e)
	//	game.UnMakeMove()
	//	game.PrintBoard()
	//	//game.PrintBitBoard(game.Occupancies[game.BOTH])
	//	//game.Scanner.Scan()
	//}
	//fmt.Println(moveList.Count)
	//game.PerftDriver(3)
	//fmt.Println(game.LeafNodesVisited)
	//game.PerftTest(7)
	////////////////////////////////////////////////////////////////////////////////
	/*
		Testing the encoded move representation
	*/
	// basic idea
	//var move uint64
	//move = (move | 39) << 6 // encode target square of e1
	//game.PrintBitBoard(move)
	//targetSquare := (move & 0xfc0) >> 6
	//fmt.Printf("targetSquare: %s\n", game.SquareToCoord[targetSquare])

	//move := game.EncodeMove(game.A1, game.A7, game.WhiteRook, 0, 1, 0, 0)
	//move1 := game.EncodeMove(game.D7, game.D5, game.BlackPawn, 0, 0, 1, 0)
	//var moveList game.Moves
	//moveList.AddMove(move)
	//moveList.AddMove(move1)
	//game.PrintMoveList(&moveList)

	//moveList := game.Moves{}
	//game.ParseFEN("1p3/2P11/PR3/3Pp/Pp3/2B2/12p1/1NP1K b a3")
	//game.PrintBoard()
	//game.GenerateMoves(&moveList)
	//game.PrintMoveList(&moveList)

	////////////////////////////////////////////////////////////////////////////////
	/*
		Test the bitboard generation for piece attacks
	*/
	//game.ParseFEN("1p3/2P11/PR3/3Pp/Pp3/2B2/12p1/1NP1K w a3")
	//game.PrintBoard()
	//game.GenerateMoves()
	////////////////////////////////////////////////////////////////////////////////
	/*
		How to check if a square of the board is being attacked
	*/
	// understanding the trick to get the attacked squares
	//game.ParseFEN("5/5/5/2P2/5/5/5/5 w -") // white pawn on c5
	//game.PrintBoard()
	//game.PrintBitBoard(game.PawnAttacks[game.BLACK][game.B6])
	//game.PrintBitBoard(game.Bitboards[game.WhitePawn])
	//game.PrintBitBoard(game.PawnAttacks[game.BLACK][game.B6] & game.Bitboards[game.WhitePawn])
	//fmt.Println(game.PawnAttacks[game.BLACK][game.C5] & game.Bitboards[game.WhitePawn])
	//game.PrintAttackedSquares(game.WHITE)
	//
	//game.ParseFEN("5/5/5/5/2N2/5/5/5 w -") // knight on c4
	//game.PrintBoard()
	//game.PrintBitBoard(game.KnightAttacks[game.C4])
	//game.PrintBitBoard(game.Bitboards[game.WhiteKnight])
	//game.PrintBitBoard(game.KnightAttacks[game.B6] & game.Bitboards[game.WhiteKnight])
	//game.PrintAttackedSquares(game.WHITE)
	//
	//game.ParseFEN("5/5/5/5/2R2/5/5/5 w -") // rook on c4
	//game.PrintBoard()
	//game.PrintBitBoard(game.GetRookAttacks(game.C5, game.BOTH))
	//game.PrintBitBoard(game.Bitboards[game.WhiteRook])
	//game.PrintBitBoard(game.GetRookAttacks(game.C5, game.BOTH) & game.Bitboards[game.WhiteRook])
	//game.PrintAttackedSquares(game.WHITE)
	//
	//game.ParseFEN(game.FenStartPosition)
	//game.PrintBoard()
	//game.PrintAttackedSquares(game.WHITE)
	//game.PrintAttackedSquares(game.BLACK)
	////////////////////////////////////////////////////////////////////////////////
	/*
		Test the FEN strings
	*/
	//game.ParseFEN(game.FenStartPosition)
	//game.PrintBoard()
	//game.PrintBitBoard(game.Occupancies[game.WHITE])
	//game.PrintBitBoard(game.Occupancies[game.BLACK])
	//game.PrintBitBoard(game.Occupancies[game.BOTH])
	////////////////////////////////////////////////////////////////////////////////
	/*
		Test the ASCII and Unicode pieces
	*/
	//game.SetBit(&game.Bitboards[game.WhitePawn], game.E2)
	//game.SetBit(&game.Bitboards[game.WhitePawn], game.D2)
	//game.SetBit(&game.Bitboards[game.WhitePawn], game.C2)
	//game.SetBit(&game.Bitboards[game.WhitePawn], game.B2)
	//game.SetBit(&game.Bitboards[game.WhitePawn], game.A2)
	//game.SetBit(&game.Bitboards[game.WhiteRook], game.B1)
	//game.SetBit(&game.Bitboards[game.WhiteKnight], game.D1)
	//game.SetBit(&game.Bitboards[game.WhiteBishop], game.E1)
	//game.SetBit(&game.Bitboards[game.WhiteKing], game.A1)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.E6)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.D6)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.C6)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.B6)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.A6)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.E7)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.D7)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.C7)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.B7)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.A7)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.E8)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.D8)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.C8)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.B8)
	//game.SetBit(&game.Bitboards[game.BlackPawn], game.A8)
	//
	//game.PrintBitBoard(game.Bitboards[game.BlackPawn])
	////fmt.Printf("piece: %s\n", game.AsciiPieces[game.WHITE_PAWN])
	////fmt.Printf("piece: %s\n", game.UnicodePieces[game.WHITE_PAWN])
	////fmt.Printf("piece: %s\n", game.AsciiPieces[game.ConvertAsciiToConstants['P']])
	////fmt.Printf("piece: %s\n", game.UnicodePieces[game.ConvertAsciiToConstants['P']])
	//game.SideToMove = game.WHITE
	//game.EnPassantSquare = game.C3
	//game.PrintBoard()
	//
	//fmt.Println()
	//fmt.Println()
	//for i := game.WhitePawn; i <= game.BlackPawn; i++ {
	//	fmt.Printf("Bitboard for: %s\n", game.UnicodePieces[i])
	//	game.PrintBitBoard(game.Bitboards[i])
	//}
	////////////////////////////////////////////////////////////////////////////////
	/*
		Test the bitboard to make sure the attacks are correct
	*/
	//var occupancy uint64
	//game.SetBit(&occupancy, game.C3)
	//game.SetBit(&occupancy, game.B4)
	//game.SetBit(&occupancy, game.D2)
	//game.PrintBitBoard(occupancy)
	//game.PrintBitBoard(game.GetBishopAttacks(game.D4, occupancy))
	//game.PrintBitBoard(game.GetRookAttacks(game.D4, occupancy))
	////////////////////////////////////////////////////////////////////////////////
	/*
		Get the magic numbers for the bishop and rook pieces
	*/
	// magic number routine
	//game.InitMagicNumbers()

	//fmt.Printf("%d\n")
	//game.PrintBitBoard(uint64(game.GetRandomUInt32()))
	//game.PrintBitBoard(uint64(game.GetRandomUInt32()) & 0xFFFF)
	//game.PrintBitBoard(game.GetRandomUInt64())
	//game.PrintBitBoard(game.GenerateMagicNumber())
	////////////////////////////////////////////////////////////////////////////////
	/*
		Get the number of relevant occupancy bits counts for each square on the board for a given piece attack mask.
		This is useful for magic bitboard generation, so we save it later.
	*/
	//for rank := 0; rank < 8; rank++ {
	//	for file := 0; file < 5; file++ {
	//		square := rank*5 + file
	//		//fmt.Printf(" %d,", game.CountBits(game.MaskRookAttacks(square)))
	//		fmt.Printf(" %d,", game.CountBits(game.MaskBishopAttacks(square)))
	//	}
	//	fmt.Println()
	//}
	////////////////////////////////////////////////////////////////////////////////
	/*
		Generate all the rook attack masks for a given square when there are pieces blocking the way
	*/
	//attackMask := game.MaskRookAttacks(game.A2)
	//game.PrintBitBoard(attackMask)

	// set index to 4096
	//for i := 0; i < 4096; i++ {
	//	game.PrintBitBoard(game.SetOccupancy(i, game.CountBits(attackMask), attackMask))
	//}
	////////////////////////////////////////////////////////////////////////////////
	/*
		Set a block pattern on the board and see how the bishop and rook attack generation works
	*/
	//var block uint64 = 0
	//game.SetBit(&block, game.B6)
	//game.SetBit(&block, game.D3)
	//game.SetBit(&block, game.A3)
	//game.SetBit(&block, game.C2)
	//game.SetBit(&block, game.C7)
	//game.PrintBitBoard(block)
	//fmt.Println("Bit Count: ", game.CountBits(block))
	//fmt.Println("Least significant bit index: ", game.GetLeastSignificantBitIndex(block))
	//fmt.Println("Least significant bit coordinate: ", game.SquareToCoord[game.GetLeastSignificantBitIndex(block)])
	//
	//var test uint64 = 0
	//game.SetBit(&test, game.GetLeastSignificantBitIndex(block))
	//game.PrintBitBoard(test)
	//
	//game.PrintBitBoard(game.BishopAttacksOnTheFly(game.C5, block))
	//game.PrintBitBoard(game.RookAttacksOnTheFly(game.C3, block))

	////////////////////////////////////////////////////////////////////////////////
	/*
		Play around with the attack masks for each piece
	*/
	//game.PrintBitBoard(game.MaskPawnAttacks(game.WHITE, game.B1))
	//game.PrintBitBoard(game.MaskPawnAttacks(game.BLACK, game.D8))
	//game.PrintBitBoard(game.MaskKnightAttacks(game.E8))
	//game.PrintBitBoard(game.MaskKingAttacks(game.D4))
	//game.PrintBitBoard(game.MaskBishopAttacks(game.A5))
	//game.PrintBitBoard(game.MaskRookAttacks(game.C4))

	////////////////////////////////////////////////////////////////////////////////
	/*
		Uncomment a line to choose a Leaper attack piece for which you want to see all the attack patterns
	*/
	//game.InitLeapersAttacks()
	//for square := 0; square < 40; square++ {
	//	game.PrintBitBoard(game.PawnAttacks[game.BLACK][square])
	//	game.PrintBitBoard(game.PawnAttacks[game.WHITE][square])
	//	game.PrintBitBoard(game.KnightAttacks[square])
	//	game.PrintBitBoard(game.KingAttacks[square])
	//	game.PrintBitBoard(game.MaskBishopAttacks(square))
	//	game.PrintBitBoard(game.MaskRookAttacks(square))
	//	game.PrintBitBoard(game.BishopAttacksOnTheFly(square, 0))
	//}
	////////////////////////////////////////////////////////////////////////////////
	/*
		Just a start here, play around with the bitboard by setting Bits to it
	*/
	//var board1 uint64
	//var board2 uint64 = 0b0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0001_1111
	//var board3 uint64 = 35
	//fmt.Printf("%b\n", board1)
	//fmt.Printf("%b\n", board2)
	//fmt.Printf("%b\n", board3)
	//
	//game.SetBit(&board1, game.E5)
	//game.PrintBitBoard(board1)
	//game.PopBit(&board1, game.E5)
	//game.PrintBitBoard(board1)
	//
	//game.PrintBitBoard(board2)
	//game.PrintBitBoard(board3)

	////////////////////////////////////////////////////////////////////////////////

	//for rank := 1; rank < 9; rank++ {
	//	//fmt.Printf("\"a%d\",\"b%d\",\"c%d\",\"d%d\",\"e%d\",", rank, rank, rank, rank, rank)
	//	fmt.Printf("A%d, B%d, C%d, D%d, E%d,\n", rank, rank, rank, rank, rank)
	//}

}
