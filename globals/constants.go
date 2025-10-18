package globals

import (
	"bufio"
	"os"
)

/* - - - - - - - - - - - - - - - - - - - - - - -
|											   |
				FEN CONSTANTS
|											   |
- - - - - - - - - - - - - - - - - - - - - - - */

// FenEmpty is the FEN string for an empty board
const FenEmpty string = "5/5/5/5/5/5/5/5 w -"

// FenStartWhiteBottomRow holds the different starting positions with white pieces at the bottom row
var FenStartWhiteBottomRow = [120]string{
	"RNK1B", "NRK1B", "KRN1B", "RKN1B", "NKR1B", "KNR1B", "1NRKB", "N1RKB", "R1NKB", "1RNKB",
	"NR1KB", "RN1KB", "RK1NB", "KR1NB", "1RKNB", "R1KNB", "K1RNB", "1KRNB", "1KNRB", "K1NRB",
	"N1KRB", "1NKRB", "KN1RB", "NK1RB", "BK1RN", "KB1RN", "1BKRN", "B1KRN", "K1BRN", "1KBRN",
	"RKB1N", "KRB1N", "BRK1N", "RBK1N", "KBR1N", "BKR1N", "B1RKN", "1BRKN", "RB1KN", "BR1KN",
	"1RBKN", "R1BKN", "R1KBN", "1RKBN", "KR1BN", "RK1BN", "1KRBN", "K1RBN", "N1RBK", "1NRBK",
	"RN1BK", "NR1BK", "1RNBK", "R1NBK", "B1NRK", "1BNRK", "NB1RK", "BN1RK", "1NBRK", "N1BRK",
	"NRB1K", "RNB1K", "BNR1K", "NBR1K", "RBN1K", "BRN1K", "BR1NK", "RB1NK", "1BRNK", "B1RNK",
	"R1BNK", "1RBNK", "KRBN1", "RKBN1", "BKRN1", "KBRN1", "RBKN1", "BRKN1", "NRKB1", "RNKB1",
	"KNRB1", "NKRB1", "RKNB1", "KRNB1", "KBNR1", "BKNR1", "NKBR1", "KNBR1", "BNKR1", "NBKR1",
	"NBRK1", "BNRK1", "RNBK1", "NRBK1", "BRNK1", "RBNK1", "1BNKR", "B1NKR", "N1BKR", "1NBKR",
	"BN1KR", "NB1KR", "KB1NR", "BK1NR", "1KBNR", "K1BNR", "B1KNR", "1BKNR", "1NKBR", "N1KBR",
	"K1NBR", "1KNBR", "NK1BR", "KN1BR", "KNB1R", "NKB1R", "BKN1R", "KBN1R", "NBK1R", "BNK1R"}

// FenDebugStartPosition is the FEN string for the starting position
const FenDebugStartPosition string = "ppppp/ppppp/ppppp/5/5/5/PPPPP/RNK1B w -"

// ppppp/ppppp/ppppp/5/5/5/PPPPP/1BNKR w -
const FenDebug2 string = "ppppp/ppp1p/p2p1/Ppppp/1P3/1RN1P/2PPB/2K2 w -"
const FenDebug3 string = "ppppp/ppp1p/p2p1/PppNp/1P1P1/1R2P/2P1B/2K2 b d3"
const FenDebug4 string = "1p3/2P11/PR3/3Pp/Pp3/2B2/12p1/1NP1K b a3"

/* - - - - - - - - - - - - - - - - - - - - - - -
|											   |
				BOARD CONSTANTS
|											   |
- - - - - - - - - - - - - - - - - - - - - - - */

// Side to move constants
const (
	WHITE = iota
	BLACK
	BOTH
)

// Bishop and Rooks constants for generating masks
const (
	ROOK = iota
	BISHOP
)

// Piece constants
const (
	WhitePawn = iota
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteKing
	BlackPawn
	NoPiece
)

// Move type constants
const (
	AllMoves = iota
	OnlyCaptures
)

// AsciiPieces is a constant that holds the ascii representation of each piece
var AsciiPieces = [6]string{
	" P", // WHITE_PAWN
	" N", // WHITE_KNIGHT
	" B", // WHITE_BISHOP
	" R", // WHITE_ROOK
	" K", // WHITE_KING
	" p"} // BlackPawn

// UnicodePieces is a constant that holds the Unicode representation of each piece, make sure you have a Dark theme
var UnicodePieces = [7]string{"♟", "♞", "♝", "♜", "♚", "♙"}

// ♜ ♞ ♝ ♛ ♚ ♟︎
// ♙ ♖ ♘ ♗ ♕ ♔

// ConvertAsciiToConstants is a map that converts the ascii representation of each piece to its constant
var ConvertAsciiToConstants = map[rune]int{
	'P': WhitePawn,
	'N': WhiteKnight,
	'B': WhiteBishop,
	'R': WhiteRook,
	'K': WhiteKing,
	'p': BlackPawn}

// ConvertConstantsToString is a map that converts the piece constant to its ascii representation
var ConvertConstantsToString = map[int]string{
	WhitePawn:   "P",
	WhiteKnight: "N",
	WhiteBishop: "B",
	WhiteRook:   "R",
	WhiteKing:   "K",
	BlackPawn:   "p",
	-1:          "1"}

// PromotedPieces holds the map for promoted pieces constants
var PromotedPieces = map[int]string{
	WhiteKnight: "N",
	WhiteBishop: "B",
	WhiteRook:   "R",
	WhiteKing:   "K"}

/*
constant representing each square on the 5x8 board
a8, b8, c8, d8, e8,
a7, b7, c7, d7, e7,
a6, b6, c6, d6, e6,
a5, b5, c5, d5, e5,
a4, b4, c4, d4, e4,
a3, b3, c3, d3, e3,
a2, b2, c2, d2, e2,
a1, b1, c1, d1, e1
*/
const (
	A8 = iota
	B8
	C8
	D8
	E8
	A7
	B7
	C7
	D7
	E7
	A6
	B6
	C6
	D6
	E6
	A5
	B5
	C5
	D5
	E5
	A4
	B4
	C4
	D4
	E4
	A3
	B3
	C3
	D3
	E3
	A2
	B2
	C2
	D2
	E2
	A1
	B1
	C1
	D1
	E1
	NoSquare
)

// MirrorSquare is a constant that holds the mirrored square for each square on the board
var MirrorSquare = [40]int{
	A1, B1, C1, D1, E1,
	A2, B2, C2, D2, E2,
	A3, B3, C3, D3, E3,
	A4, B4, C4, D4, E4,
	A5, B5, C5, D5, E5,
	A6, B6, C6, D6, E6,
	A7, B7, C7, D7, E7,
	A8, B8, C8, D8, E8}

// SquareToCoord is a constant that holds the algebraic notation for each square on the board
var SquareToCoord = [40]string{
	"a8", "b8", "c8", "d8", "e8",
	"a7", "b7", "c7", "d7", "e7",
	"a6", "b6", "c6", "d6", "e6",
	"a5", "b5", "c5", "d5", "e5",
	"a4", "b4", "c4", "d4", "e4",
	"a3", "b3", "c3", "d3", "e3",
	"a2", "b2", "c2", "d2", "e2",
	"a1", "b1", "c1", "d1", "e1"}

/*
NotAFile is the bitboard for the file A with all ones except for the first file which is 0
this can be created like this and run in the main function once to get the constant value:
8   0  1  1  1  1
7   0  1  1  1  1
6   0  1  1  1  1
5   0  1  1  1  1
4   0  1  1  1  1
3   0  1  1  1  1
2   0  1  1  1  1
1   0  1  1  1  1

	A  B  C  D  E
*/
const NotAFile uint64 = 1064043510750
const NotEFile uint64 = 532021755375
const NotDEFile uint64 = 248276819175
const NotABFile uint64 = 993107276700

// Bitboard40Mask is a constant that masks the last 24 bits of a 64-bit integer
const Bitboard40Mask uint64 = (1 << 40) - 1

// SideToMove is a constant that holds the side to move
var SideToMove int

// EnPassantSquare is a constant that holds the square of the en passant target
var EnPassantSquare int = NoSquare

// Bitboards holds the bitboards for each piece: black has only pawns, white has pawns, knight, bishop and rook, king
var Bitboards [6]uint64

// Occupancies hold the occupancy of each square
var Occupancies [3]uint64

// HashKey holds the hash key for the current position
var HashKey uint64

// RepetitionTable holds the repetition table for the current position
var RepetitionTable [10000]uint64

// RepetitionIndex holds the repetition index for the current position
var RepetitionIndex int

// Ply is the current ply in the search tree
var Ply int

/* - - - - - - - - - - - - - - - - - - - - - - -
|											   |
				ATTACK CONSTANTS
|											   |
- - - - - - - - - - - - - - - - - - - - - - - */

// PawnAttacks Here we create a way to form the pawn attack table, which is a 2d array [side_to_move][square]
var PawnAttacks [2][40]uint64

// KnightAttacks is a table of all the knight attacks on the bitboard
var KnightAttacks [40]uint64

// KingAttacks is a table of all the king attacks on the bitboard
var KingAttacks [40]uint64

// BishopMasks is a table of all the bishop masks on the bitboard
var BishopMasks [40]uint64

// RookMasks is a table of all the rook masks on the bitboard
var RookMasks [40]uint64

// BishopAttacks is a table of all the bishop attacks on the bitboard
var BishopAttacks [40][512]uint64

// RookAttacks is a table of all the rook attacks on the bitboard
var RookAttacks [40][4096]uint64

// BishopRelevantOccupancyCount are the relevant occupancy bit count for every square on the board
var BishopRelevantOccupancyCount = [40]int{
	3, 2, 2, 2, 3,
	3, 2, 2, 2, 3,
	4, 3, 4, 3, 4,
	5, 4, 4, 4, 5,
	5, 4, 4, 4, 5,
	4, 3, 4, 3, 4,
	3, 2, 2, 2, 3,
	3, 2, 2, 2, 3}

// RookRelevantOccupancyCount are the relevant occupancy bit count for every square on the board
var RookRelevantOccupancyCount = [40]int{
	9, 8, 8, 8, 9,
	8, 7, 7, 7, 8,
	8, 7, 7, 7, 8,
	8, 7, 7, 7, 8,
	8, 7, 7, 7, 8,
	8, 7, 7, 7, 8,
	8, 7, 7, 7, 8,
	9, 8, 8, 8, 9}

// RookMagicNumbers RooKMagicNumbers is a table of all the rook magic numbers
var RookMagicNumbers = [40]uint64{
	0x400804082801400,
	0xa0081028000084d,
	0x240800c0400420a0,
	0x402004040000092,
	0x800a04400204082,
	0x20480100000890,
	0x10201008000030,
	0x20404040482000,
	0x340200404420001,
	0x2160048090020201,
	0x421002004400040,
	0x80801008000806,
	0x4421020040041044,
	0xc80844020000600,
	0x4002480900000,
	0x8020081100600000,
	0x100482010000001,
	0x4080204104042000,
	0x8140100804000000,
	0x8408200120880300,
	0x40100240110a000c,
	0x120804081001904,
	0x20204080160000,
	0xc0408040400101,
	0x420044020001020,
	0x40040210000a00b,
	0x8050101080800,
	0x882081080022010,
	0x208200420190800,
	0xc008448800800185,
	0x4010480a0000080,
	0x2100810500004c1,
	0x2008080810400620,
	0x431010480040000,
	0x8048a40000080,
	0x8804002800002,
	0x10010080408000,
	0x10020202040a004,
	0x8020200400800001,
	0x11048082002040}

// BishopMagicNumbers is a table of all the bishop magic numbers
var BishopMagicNumbers = [40]uint64{
	0x25c080804200101,
	0x895284800000b20,
	0x24a040000000442,
	0x30c000080408000,
	0x48200004002020,
	0x8110888004010,
	0x1404640080008810,
	0xccc01088b5106,
	0x6822020000500,
	0x80b1428200040040,
	0x110222100083402,
	0x10010b501000001,
	0x80108900002020,
	0x481143100060020,
	0x4084314800000000,
	0x80c8104240048100,
	0x84008100021000,
	0x44008a05000100,
	0x480c910801c210,
	0xe440801080008602,
	0x208201010800004,
	0x109089020020100,
	0x200208400540a80,
	0x2088408240206206,
	0x424080811414000,
	0x30c20c40440000,
	0x6000108091020900,
	0x4010822040018008,
	0x3212888242040002,
	0x809001002034c,
	0x60008d2020043210,
	0x830404b40c350241,
	0x8801854800842,
	0x8021320600810,
	0x14e0200041004,
	0x501259044204a0,
	0x84a4400840,
	0x9c0028130025001,
	0x80000108c001010,
	0x4329c40856000904}

/* - - - - - - - - - - - - - - - - - - - - - - -
|											   |
				EVALUATION CONSTANTS
|											   |
- - - - - - - - - - - - - - - - - - - - - - - */

// MaterialValues holds the material values for each piece type
var MaterialValues = map[int]int{
	WhitePawn:   100,
	WhiteKnight: 300,
	WhiteBishop: 350,
	WhiteRook:   500,
	WhiteKing:   400,
	BlackPawn:   -100}

// PawnPositionalValues holds the positional values for pawns on the board
var PawnPositionalValues = [40]int{
	90, 90, 90, 90, 90,
	30, 30, 40, 30, 30,
	20, 20, 30, 20, 20,
	10, 10, 20, 10, 10,
	5, 10, 20, 10, 5,
	0, 0, 5, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0}

// KnightPositionalValues holds the positional values for knights on the board
var KnightPositionalValues = [40]int{
	-5, 0, 0, 0, -5,
	-5, 0, 5, 0, -5,
	-5, 5, 10, 5, -5,
	-5, 20, 30, 20, -5,
	-5, 20, 30, 20, -5,
	-5, 5, 10, 5, -5,
	-5, 0, 0, 0, -5,
	-5, -5, -5, -5, -5}

// BishopPositionalValues holds the positional values for bishop on the board
var BishopPositionalValues = [40]int{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 10, 10, 10, 0,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 5, 5, 5, 0,
	0, 10, 10, 10, 0,
	-10, -10, -10, -10, -10}

// RookPositionalValues holds the positional values for rook on the board
var RookPositionalValues = [40]int{
	50, 50, 50, 50, 50,
	50, 50, 50, 50, 50,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 0, 20, 0, 0}

// KingPositionalValues holds the positional values for king on the board
var KingPositionalValues = [40]int{
	0, 0, 0, 0, 0,
	0, 5, 5, 5, 0,
	0, 5, 10, 5, 0,
	0, 10, 20, 10, 0,
	0, 10, 20, 10, 0,
	0, 5, 10, 5, 0,
	0, 5, -5, 5, 0,
	-15, -15, -15, -15, -15}

// GetRankFromSquare holds the rank for each square on the board
var GetRankFromSquare = [40]int{
	7, 7, 7, 7, 7,
	6, 6, 6, 6, 6,
	5, 5, 5, 5, 5,
	4, 4, 4, 4, 4,
	3, 3, 3, 3, 3,
	2, 2, 2, 2, 2,
	1, 1, 1, 1, 1,
	0, 0, 0, 0, 0}

// DoublePawnPenalty holds the doubled pawn penalty
var DoublePawnPenalty = -10

// IsolatedPawnPenalty holds the isolated pawn penalty
const IsolatedPawnPenalty = -10

// SemiOpenFileScore holds the semi-open file score
const SemiOpenFileScore = 10

// OpenFileScore holds the open file score
const OpenFileScore = 15

// PassedPawnBonus holds the passed pawn bonus
var PassedPawnBonus = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

// FileMasks holds the file masks for each file on the board
var FileMasks [40]uint64

// RankMasks holds the rank masks for each rank on the board
var RankMasks [40]uint64

// IsolatedMasks holds the isolated pawn masks for each file on the board
var IsolatedMasks [40]uint64

// WhitePassedMasks holds the passed white pawn masks for each file on the board
var WhitePassedMasks [40]uint64

// BlackPassedMasks holds the passed black pawn masks for each file on the board
var BlackPassedMasks [40]uint64

/* - - - - - - - - - - - - - - - - - - - - - - -
|											   |
				OTHER CONSTANTS
|											   |
- - - - - - - - - - - - - - - - - - - - - - - */

// Scanner is used to read input from the console
var Scanner = bufio.NewScanner(os.Stdin)

// NodesVisited counts the number of nodes visited during perft tests
var NodesVisited int

// LeafNodesVisited counts the number of leaf nodes visited during perft tests
var LeafNodesVisited int

// TotalNodesEveryPlyFromStart holds the total number of nodes at each depth from 1 to 15 from the initial position
var TotalNodesEveryPlyFromStart = [15]int{13, 78, 986, 6530, 88362, 654080, 9370535, 75313428, 1126242061, 9632104995, 149071825095, 12, 13, 14, 15}

// 148471/149071825095

// LeafNodesEveryPlyFromStart holds the total number of leaf nodes at each depth from 1 to 15 from the initial position
var LeafNodesEveryPlyFromStart = [15]int{13, 65, 908, 5544, 81832, 565718, 8716455, 65942893, 1050928633, 8505862934, 139439720100, 12, 13, 14, 15}
