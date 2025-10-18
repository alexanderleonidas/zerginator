package gui

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
	"zerginator/ai"
	"zerginator/bitoperations"
	"zerginator/board"
	"zerginator/clock"

	//"zerginator/clock"
	"zerginator/globals"
	"zerginator/uci"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	tileSize     = 80
	boardWidth   = 5
	boardHeight  = 8
	panelHeight  = 110
	ScreenWidth  = tileSize * boardWidth
	ScreenHeight = tileSize*boardHeight + panelHeight
	panelY       = tileSize * boardHeight
	ctrlBtnW     = 150
	ctrlBtnH     = 40
)

const (
	stateMenu = iota
	statePlaying
	stateReset
	stateGameOver
	statePromotion
)

// reusable UI images and scaled pieces
var (
	lightSquareImg *ebiten.Image
	darkSquareImg  *ebiten.Image
	checkboxOnImg  *ebiten.Image
	checkboxOffImg *ebiten.Image
	buttonImg      *ebiten.Image
	scaledPieceImg map[int]*ebiten.Image
)

// loadPNG loads a PNG image from the specified path and returns it as an *ebiten.Image
func loadPNG(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// InitImages initializes the piece images and the reusable UI assets.
// Scales piece images to tileSize once to avoid per-frame scaling.
func InitImages() error {
	scaledPieceImg = make(map[int]*ebiten.Image)

	// create reusable square images
	lightSquareImg = ebiten.NewImage(tileSize, tileSize)
	lightSquareImg.Fill(color.RGBA{R: 240, G: 217, B: 181, A: 255})
	darkSquareImg = ebiten.NewImage(tileSize, tileSize)
	darkSquareImg.Fill(color.RGBA{R: 181, G: 136, B: 99, A: 255})

	// checkbox images
	checkboxOnImg = ebiten.NewImage(20, 20)
	checkboxOnImg.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
	checkboxOffImg = ebiten.NewImage(20, 20)
	checkboxOffImg.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255})

	// button image
	buttonImg = ebiten.NewImage(200, 60)
	buttonImg.Fill(color.RGBA{R: 50, G: 100, B: 200, A: 255})

	// load raw images
	paths := map[int]string{
		globals.WhitePawn:   "images/white_pawn.png",
		globals.WhiteRook:   "images/white_rook.png",
		globals.WhiteKnight: "images/white_knight.png",
		globals.WhiteBishop: "images/white_bishop.png",
		globals.WhiteKing:   "images/white_king.png",
		globals.BlackPawn:   "images/black_pawn.png",
	}

	for k, p := range paths {
		img, err := loadPNG(p)
		if err != nil {
			return err
		}
		// scale into a tile-size image once
		dst := ebiten.NewImage(tileSize, tileSize)
		op := &ebiten.DrawImageOptions{}
		scaleX := float64(tileSize) / float64(img.Bounds().Dx())
		scaleY := float64(tileSize) / float64(img.Bounds().Dy())
		op.GeoM.Scale(scaleX, scaleY)
		dst.DrawImage(img, op)
		scaledPieceImg[k] = dst
	}

	return nil
}

func formatSeconds(s time.Duration) string {
	if s < 0 {
		s = 0
	}
	s = s / time.Second
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

// Game represents a game state.
type Game struct {
	selectedSource  int
	state           int
	pvp             bool
	pvc             bool
	cvc             bool
	playerPlays     int
	movesMade       int
	pieceOptions    []int
	bottomSelection []int
	winner          int
	clock           *clock.GameClock
}

// numClicks tracks the number of clicks (0, 1, or 2)
var numClicks int

// pendingPromotionFrom is the square from which the piece is promoted from
var pendingPromotionFrom int = globals.NoSquare

// pendingPromotionTo is the square to which the piece is promoted to
var pendingPromotionTo int = globals.NoSquare

// Update updates the game state.
func (g *Game) Update() error {
	switch g.state {
	// Menu state: handle menu interactions
	case stateMenu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()

			// checkbox area for "Play vs Player" (small box)
			bw := 200
			cbx := (ScreenWidth - bw) / 2
			cby := (ScreenHeight / 2) - 100
			if x >= cbx && x <= cbx+20 && y >= cby && y <= cby+20 {
				g.pvp = !g.pvp
				if g.pvp {
					g.pvc = false
					g.cvc = false
				}
				log.Printf("Play Vs Player toggled: %v\n", g.pvp)
				return nil
			}

			// Checkbox area for "Play vs Computer"
			bw2 := 200
			cbx2 := (ScreenWidth - bw2) / 2
			cby2 := (ScreenHeight / 2) - 60
			if x >= cbx2 && x <= cbx2+20 && y >= cby2 && y <= cby2+20 {
				g.pvc = !g.pvc
				if g.pvc {
					g.pvp = false
					g.cvc = false
				}
				log.Printf("Play Vs Computer toggled: %v\n", g.pvc)
				return nil
			}

			// Checkbox area for "Computer vs Computer"
			bw3 := 200
			cbx3 := (ScreenWidth - bw3) / 2
			cby3 := (ScreenHeight / 2) - 20
			if x >= cbx3 && x <= cbx3+20 && y >= cby3 && y <= cby3+20 {
				g.cvc = !g.cvc
				if g.cvc {
					g.pvc = false
					g.pvp = false
				}
				log.Printf("Computer Vs Computer toggled: %v\n", g.cvc)
				return nil
			}

			// side selection buttons
			btw4W, btw4H := 50, 50
			btw4X := (ScreenWidth-btw4W)/2 - 100
			btw4Y := (ScreenHeight-btw4H)/2 + 100
			if x >= btw4X && x <= btw4X+btw4W && y >= btw4Y && y <= btw4Y+btw4H {
				if g.playerPlays == globals.WHITE {
					g.playerPlays = -1
				} else {
					g.playerPlays = globals.WHITE
					log.Println("Player chooses WHITE side to start")
				}
			}

			btw5W, btw5H := 50, 50
			btw5X := (ScreenWidth-btw5W)/2 + 0
			btw5Y := (ScreenHeight-btw5H)/2 + 100
			if x >= btw5X && x <= btw5X+btw5W && y >= btw5Y && y <= btw5Y+btw5H {
				if g.playerPlays == globals.BLACK {
					g.playerPlays = -1
				} else {
					g.playerPlays = globals.BLACK
					log.Println("Player chooses BLACK side to start")
				}
			}

			btw6W, btw6H := 50, 50
			btw6X := (ScreenWidth-btw6W)/2 + 100
			btw6Y := (ScreenHeight-btw6H)/2 + 100
			if x >= btw6X && x <= btw6X+btw6W && y >= btw6Y && y <= btw6Y+btw6H {
				if g.playerPlays == 3 {
					g.playerPlays = -1
				} else {
					g.playerPlays = 3
					log.Println("Player chooses RANDOM side to start")
				}
			}

			// check if the player pressed the button for organising the bottom row
			btnX := (ScreenWidth - 200) / 2
			btnY := (ScreenHeight-60)/2 + 200
			if x >= btnX && x <= btnX+ctrlBtnW && y >= btnY && y <= btnY+ctrlBtnH {
				if g.movesMade == 0 {
					g.state = stateReset
				} else {
					log.Println("Cannot reset pieces: game already in progress")
				}
			}

			// Start button area (centered)
			btnW, btnH := 200, 60
			bx := (ScreenWidth - 200) / 2
			by := (ScreenHeight-60)/2 + 300
			if x >= bx && x <= bx+btnW && y >= by && y <= by+btnH {
				if (g.pvp && !g.pvc && !g.cvc) || (!g.pvp && g.pvc && !g.cvc) || (!g.pvp && !g.pvc && g.cvc) && g.playerPlays != -1 {
					// Apply settings
					g.state = statePlaying
					g.movesMade = 0
					g.selectedSource = globals.NoSquare
					numClicks = 0
					g.clock = clock.NewGameClock()
					log.Println("Game started (from menu)")
					board.ParseFEN(board.GetStartPosFEN())
					board.PrintBoard()
				}
			}
		}
	// Playing state: handle game interactions
	case statePlaying:
		// Check Win conditions
		// Black side wins if one of its pawns is in the 1st rank
		for sq := globals.A1; sq <= globals.E1; sq++ {
			if bitoperations.GetBit(globals.Bitboards[globals.BlackPawn], sq) == 1 {
				g.winner = globals.BLACK
				g.state = stateGameOver
				return nil
			}
		}
		// Black side wins if all white pieces are captured
		if globals.Occupancies[globals.WHITE] == 0 {
			g.winner = globals.BLACK
			g.state = stateGameOver
			return nil
		}
		// White side wins if all black pieces are captured
		if globals.Bitboards[globals.BlackPawn] == 0 {
			g.winner = globals.WHITE
			g.state = stateGameOver
			return nil
		}
		tempMovesList := board.Moves{}
		board.GenerateMoves(&tempMovesList)
		if tempMovesList.Count == 0 {
			g.winner = 3
			g.state = stateGameOver
			return nil
		}
		if g.clock != nil {
			if g.clock.White.IsExpired() {
				log.Println("White side clock expired")
				g.winner = globals.BLACK
				g.state = stateGameOver
				break
			}
			if g.clock.Black.IsExpired() {
				log.Println("Black side clock expired")
				g.winner = globals.WHITE
				g.state = stateGameOver
				break
			}
		}
		// Detect a single mouse click event
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			file := x / tileSize
			rank := y / tileSize
			square := rank*boardWidth + file

			// check if the player clicked the "Reset Game" button
			btn1X := (ScreenWidth-100)/2 - 140
			btn1Y := panelY + (panelHeight-ctrlBtnH)/2 - 10
			if x >= btn1X && x <= btn1X+100 && y >= btn1Y && y <= btn1Y+ctrlBtnH {
				board.ParseFEN(board.GetStartPosFEN())
				g.movesMade = 0
				g.selectedSource = globals.NoSquare
				g.clock = clock.NewGameClock()
			}

			// check if the player clicked the "Undo Move" button
			btn2X := (ScreenWidth-100)/2 + 140
			btn2Y := panelY + (panelHeight-ctrlBtnH)/2 - 10
			if x >= btn2X && x <= btn2X+100 && y >= btn2Y && y <= btn2Y+ctrlBtnH {
				if g.movesMade > 0 {
					board.UnMakeMove()
					g.movesMade--
				}
			}

			// check if the player clicked the "Go to menu" button
			btn3X := (ScreenWidth - 120) / 2
			btn3Y := panelY + (panelHeight-40)/2 + 25
			if x >= btn3X && x <= btn3X+120 && y >= btn3Y && y <= btn3Y+40 {
				g.state = stateMenu
				g.selectedSource = globals.NoSquare
				g.movesMade = 0
				g.cvc = false
				g.pvp = false
				g.pvc = false
			}

			if g.pvp || (g.pvc && globals.SideToMove == g.playerPlays) && square < 40 {
				if bitoperations.GetBit(globals.Occupancies[globals.BOTH], square) == 1 && g.selectedSource == globals.NoSquare && numClicks == 0 {
					g.selectedSource = square
					log.Printf("Clicked Source square: %s\n", globals.SquareToCoord[square])
					numClicks++
				} else if g.selectedSource != globals.NoSquare && numClicks == 1 && square < 40 {
					numClicks = 0
					log.Printf("Clicked Target square: %s\n", globals.SquareToCoord[square])
					if bitoperations.GetBit(globals.Bitboards[globals.WhitePawn], g.selectedSource) == 1 && square >= globals.A8 && square <= globals.E8 {
						pendingPromotionFrom = g.selectedSource
						pendingPromotionTo = square
						g.state = statePromotion
						return nil
					}
					moveString := globals.SquareToCoord[g.selectedSource] + globals.SquareToCoord[square]
					move := uci.ParseMove(moveString)
					if move != 0 {
						if board.MakeMove(move, globals.AllMoves) == 1 {
							g.clock.SwitchTurn()
							g.movesMade++
							board.PrintBoard()
							g.clock.Status()
						}
					}
					// reset for next move
					g.selectedSource = globals.NoSquare
				}
			}
		} else {
			// Computer to make move if it's its turn
			if (g.pvc && globals.SideToMove != g.playerPlays) || g.cvc {
				uci.ParseGo("go depth 13")
				time.Sleep(1 * time.Second)
				if board.MakeMove(ai.BestMove, globals.AllMoves) == 1 {
					g.clock.SwitchTurn()
					g.movesMade++
					board.PrintBoard()
					g.clock.Status()
				}
			}
		}
	// Reset state: handle reset interactions
	case stateReset:
		// initialize options and current selections when entering the reset state
		if g.pieceOptions == nil {
			g.pieceOptions = []int{-1, globals.WhiteKnight, globals.WhiteBishop, globals.WhiteRook, globals.WhiteKing}
		}
		if g.bottomSelection == nil || len(g.bottomSelection) == 0 {
			count := 5
			g.bottomSelection = make([]int, count)
			for i := 0; i < count; i++ {
				g.bottomSelection[i] = -1
				sq := globals.A1 // A1..E1 -> squares 0..4
				for piece := range scaledPieceImg {
					if bitoperations.GetBit(globals.Bitboards[piece], sq) == 1 {
						g.bottomSelection[i] = piece
						break
					}
				}
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			btnW, btnH := 200, 60
			bx := (ScreenWidth - btnW) / 2
			by := (ScreenHeight-btnH)/2 + 200
			if x >= bx && x <= bx+btnW && y >= by && y <= by+btnH {
				// make sure all the chosen pieces are unique
				unique := true
				seen := make(map[int]struct{})
				for _, val := range g.bottomSelection {
					if _, exists := seen[val]; exists {
						unique = false
					}
					seen[val] = struct{}{}
				}
				if unique {
					// apply the new position
					bottomRow := globals.ConvertConstantsToString[g.bottomSelection[0]] + globals.ConvertConstantsToString[g.bottomSelection[1]] + globals.ConvertConstantsToString[g.bottomSelection[2]] + globals.ConvertConstantsToString[g.bottomSelection[3]] + globals.ConvertConstantsToString[g.bottomSelection[4]]
					fen := "ppppp/ppppp/ppppp/5/5/5/PPPPP/" + bottomRow + " w -"
					board.ParseFEN(fen)
					g.state = statePlaying
					g.state = statePlaying
					g.movesMade = 0
					g.selectedSource = globals.NoSquare
					g.clock = clock.NewGameClock()
					log.Println("Game started (from menu)")
				} else {
					log.Println("Cannot set position: duplicate pieces selected")
				}
			} else {
				windW, windH := 60, 60
				count := 5
				gap := (ScreenWidth - count*windW) / (count + 1)
				baseY := (ScreenHeight-windH)/2 + 100
				for i := 0; i < count; i++ {
					bx = gap + i*(windW+gap)
					by = baseY
					if x >= bx && x <= bx+windW && y >= by && y <= by+windH {
						// cycle through pieceOptions
						current := g.bottomSelection[i]
						idx := -1
						for j, opt := range g.pieceOptions {
							if opt == current {
								idx = j
								break
							}
						}
						var nextIdx int
						if idx == -1 || idx+1 >= len(g.pieceOptions) {
							nextIdx = 0
						} else {
							nextIdx = idx + 1
						}
						g.bottomSelection[i] = g.pieceOptions[nextIdx]
						return nil
					}
				}
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = stateMenu
			return nil
		}
	case stateGameOver:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			log.Println("Game Over, total moves made: ", g.movesMade)
			// reset all state variables
			g.pvp = false
			g.pvc = false
			g.cvc = false
			g.pieceOptions = make([]int, 0)
			g.bottomSelection = make([]int, 0)
			g.winner = 0
			g.selectedSource = globals.NoSquare
			g.movesMade = 0
			g.playerPlays = 0
			g.state = stateMenu
			return nil
		}
	case statePromotion:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			windW, windH := 60, 60
			count := 4
			gap := (ScreenWidth - count*windW) / (count + 1)
			baseY := (ScreenHeight-windH)/2 + 100

			// require a pending promotion to proceed
			if pendingPromotionFrom == globals.NoSquare || pendingPromotionTo == globals.NoSquare {
				// nothing to promote
				g.state = statePlaying
				return nil
			}

			for i := range globals.PromotedPieces {
				bx := gap + (i-1)*(windW+gap)
				by := baseY
				if x >= bx && x <= bx+windW && y >= by && y <= by+windH {
					// map clicked index to promotion piece and uci promotion char
					moveString := globals.SquareToCoord[pendingPromotionFrom] + globals.SquareToCoord[pendingPromotionTo] + globals.PromotedPieces[i]
					move := uci.ParseMove(moveString)
					if move != 0 {
						if board.MakeMove(move, globals.AllMoves) == 1 {
							g.movesMade++
							board.PrintBoard()
						}
					}
					// clear pending promotion and reset selection
					pendingPromotionFrom = globals.NoSquare
					pendingPromotionTo = globals.NoSquare
					g.state = statePlaying
					g.selectedSource = globals.NoSquare
					numClicks = 0
					return nil
				}
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// cancel promotion choice
			pendingPromotionFrom = globals.NoSquare
			pendingPromotionTo = globals.NoSquare
			g.state = statePlaying
			return nil
		}
	}
	return nil
}

// Draw draws the game screen using pre-created images to avoid per-frame allocations.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == stateMenu {
		screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})

		ebitenutil.DebugPrintAt(screen, "Zerginator - Menu", ScreenWidth/2-80, ScreenHeight/2-200)
		ebitenutil.DebugPrintAt(screen, "Please check only one box!", ScreenWidth/2-80, ScreenHeight/2-140)

		// draw checkboxes using pre-made images
		bw := 200
		cbx := (ScreenWidth - bw) / 2
		cby := (ScreenHeight / 2) - 100
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cbx), float64(cby))
		if g.pvp {
			screen.DrawImage(checkboxOnImg, op)
		} else {
			screen.DrawImage(checkboxOffImg, op)
		}
		ebitenutil.DebugPrintAt(screen, "Play vs Player (click box)", cbx+28, cby-2)

		bw2 := 200
		cbx2 := (ScreenWidth - bw2) / 2
		cby2 := (ScreenHeight / 2) - 60
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(cbx2), float64(cby2))
		if g.pvc {
			screen.DrawImage(checkboxOnImg, op2)
		} else {
			screen.DrawImage(checkboxOffImg, op2)
		}
		ebitenutil.DebugPrintAt(screen, "Play vs Computer (click box)", cbx2+28, cby2-2)

		bw3 := 200
		cbx3 := (ScreenWidth - bw3) / 2
		cby3 := (ScreenHeight / 2) - 20
		op3 := &ebiten.DrawImageOptions{}
		op3.GeoM.Translate(float64(cbx3), float64(cby3))
		if g.cvc {
			screen.DrawImage(checkboxOnImg, op3)
		} else {
			screen.DrawImage(checkboxOffImg, op3)
		}
		ebitenutil.DebugPrintAt(screen, "Computer vs Computer (click box)", cbx3+28, cby3-2)

		ebitenutil.DebugPrintAt(screen, "Player, please choose a side to start!", ScreenWidth/2-110, ScreenHeight/2+40)

		// side checkboxes (use same 50x50 image but positioned)
		btw4W, btw4H := 50, 50
		btw4X := (ScreenWidth-btw4W)/2 - 100
		btw4Y := (ScreenHeight-btw4H)/2 + 100
		op4 := &ebiten.DrawImageOptions{}
		op4.GeoM.Translate(float64(btw4X), float64(btw4Y))
		if g.playerPlays == globals.WHITE {
			box4 := ebiten.NewImage(btw4W, btw4H)
			box4.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
			screen.DrawImage(box4, op4)
		} else {
			box4 := ebiten.NewImage(btw4W, btw4H)
			box4.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255})
			screen.DrawImage(box4, op4)
		}
		ebitenutil.DebugPrintAt(screen, "WHITE", btw4X+10, btw4Y+18)

		btw5W, btw5H := 50, 50
		btw5X := (ScreenWidth-btw5W)/2 + 0
		btw5Y := (ScreenHeight-btw5H)/2 + 100
		op5 := &ebiten.DrawImageOptions{}
		op5.GeoM.Translate(float64(btw5X), float64(btw5Y))
		if g.playerPlays == globals.BLACK {
			box5 := ebiten.NewImage(btw5W, btw5H)
			box5.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
			screen.DrawImage(box5, op5)
		} else {
			box5 := ebiten.NewImage(btw5W, btw5H)
			box5.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255})
			screen.DrawImage(box5, op5)
		}
		ebitenutil.DebugPrintAt(screen, "BLACK", btw5X+10, btw5Y+18)

		btw6W, btw6H := 50, 50
		btw6X := (ScreenWidth-btw6W)/2 + 100
		btw6Y := (ScreenHeight-btw6H)/2 + 100
		op6 := &ebiten.DrawImageOptions{}
		op6.GeoM.Translate(float64(btw6X), float64(btw6Y))
		if g.playerPlays == 3 {
			box6 := ebiten.NewImage(btw6W, btw6H)
			box6.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
			screen.DrawImage(box6, op6)
		} else {
			box6 := ebiten.NewImage(btw6W, btw6H)
			box6.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255})
			screen.DrawImage(box6, op6)
		}
		ebitenutil.DebugPrintAt(screen, "RANDOM", btw6X+8, btw6Y+18)

		// Choose bottom row pieces for WHITE
		btnX := (ScreenWidth - 200) / 2
		btnY := (ScreenHeight-60)/2 + 200
		btnOp := &ebiten.DrawImageOptions{}
		btnOp.GeoM.Translate(float64(btnX), float64(btnY))
		screen.DrawImage(buttonImg, btnOp)
		ebitenutil.DebugPrintAt(screen, "Choose Bottom Row", btnX+50, btnY+22)

		ebitenutil.DebugPrintAt(screen, "OR", btnX+90, btnY+75)

		// Start button
		startBtnX := (ScreenWidth - 200) / 2
		startBtnY := (ScreenHeight-60)/2 + 300
		opStartBtn := &ebiten.DrawImageOptions{}
		opStartBtn.GeoM.Translate(float64(startBtnX), float64(startBtnY))
		screen.DrawImage(buttonImg, opStartBtn)
		ebitenutil.DebugPrintAt(screen, "Start Game", startBtnX+60, startBtnY+22)

		return
	} else if g.state == stateReset {
		screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
		ebitenutil.DebugPrintAt(screen, "Choose White Bottom Pieces", ScreenWidth/2-90, ScreenHeight/2-150)
		ebitenutil.DebugPrintAt(screen, "Click on the Squares to change the pieces", ScreenWidth/2-100, ScreenHeight/2-100)
		ebitenutil.DebugPrintAt(screen, "Press escape to go back (you will loose the order if you do)", ScreenWidth/2-170, ScreenHeight/2-50)

		windW, windH := 60, 60
		count := 5
		gap := (ScreenWidth - count*windW) / (count + 1)
		baseY := (ScreenHeight-windH)/2 + 100
		labels := []string{"A1", "B1", "C1", "D1", "E1"}

		for i := 0; i < count; i++ {
			x := gap + i*(windW+gap)
			opTemp := &ebiten.DrawImageOptions{}
			opTemp.GeoM.Translate(float64(x), float64(baseY))

			box := ebiten.NewImage(windW, windH)
			// draw background box
			box.Fill(color.RGBA{R: 220, G: 220, B: 220, A: 255})
			screen.DrawImage(box, opTemp)

			// draw piece image if set, scale to fit the widget
			if g.bottomSelection != nil && i < len(g.bottomSelection) && g.bottomSelection[i] != -1 {
				piece := g.bottomSelection[i]
				if img := scaledPieceImg[piece]; img != nil {
					imgOp := &ebiten.DrawImageOptions{}
					scale := float64(windW) / float64(tileSize)
					imgOp.GeoM.Scale(scale, scale)
					// center scaled image in the widget
					imgOp.GeoM.Translate(float64(x), float64(baseY))
					screen.DrawImage(img, imgOp)
				}
			}

			ebitenutil.DebugPrintAt(screen, labels[i], x, baseY-20)
		}

		btnX := (ScreenWidth - 200) / 2
		btnY := (ScreenHeight-60)/2 + 200
		opBtn := &ebiten.DrawImageOptions{}
		opBtn.GeoM.Translate(float64(btnX), float64(btnY))
		screen.DrawImage(buttonImg, opBtn)
		ebitenutil.DebugPrintAt(screen, "Set Position and Start", btnX+40, btnY+22)

		return
	} else if g.state == stateGameOver {
		screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
		var winnerText string
		if g.winner == globals.WHITE {
			winnerText = "White Wins!"
		} else if g.winner == globals.BLACK {
			winnerText = "Black Wins!"
		} else {
			winnerText = "Draw!"
		}
		ebitenutil.DebugPrintAt(screen, winnerText, ScreenWidth/2-60, ScreenHeight/2-60)
		//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time left: %d seconds", g.timeLeft), ScreenWidth/2-80, ScreenHeight/2-20)
		ebitenutil.DebugPrintAt(screen, "Click to return to menu", ScreenWidth/2-80, ScreenHeight/2+20)
		return
	} else if g.state == statePromotion {
		screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
		ebitenutil.DebugPrintAt(screen, "Choose White Pawn Promotion", ScreenWidth/2-90, ScreenHeight/2-150)
		ebitenutil.DebugPrintAt(screen, "Click on the Square to choose the promoted piece", ScreenWidth/2-150, ScreenHeight/2-100)
		ebitenutil.DebugPrintAt(screen, "Press escape to go back", ScreenWidth/2-100, ScreenHeight/2-50)

		windW, windH := 60, 60
		count := 4
		gap := (ScreenWidth - count*windW) / (count + 1)
		baseY := (ScreenHeight-windH)/2 + 100

		for i := globals.WhiteKnight; i <= globals.WhiteKing; i++ {
			x := gap + (i-1)*(windW+gap)
			opTemp := &ebiten.DrawImageOptions{}
			opTemp.GeoM.Translate(float64(x), float64(baseY))

			box := ebiten.NewImage(windW, windH)
			// draw background box
			box.Fill(color.RGBA{R: 220, G: 220, B: 220, A: 255})
			screen.DrawImage(box, opTemp)

			if img := scaledPieceImg[i]; img != nil {
				imgOp := &ebiten.DrawImageOptions{}
				scale := float64(windW) / float64(tileSize)
				imgOp.GeoM.Scale(scale, scale)
				// translate using the widget position so the image is drawn at the correct spot
				imgOp.GeoM.Translate(float64(x), float64(baseY))
				screen.DrawImage(img, imgOp)
			}
		}
		return
	}

	// StatePlaying: draw tiles using pre-made square images and draw pieces using scaledPieceImgs
	for rank := 0; rank < boardHeight; rank++ {
		for file := 0; file < boardWidth; file++ {
			x, y := file*tileSize, rank*tileSize
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			if (rank+file)%2 == 0 {
				screen.DrawImage(darkSquareImg, op)
			} else {
				screen.DrawImage(lightSquareImg, op)
			}
		}
	}

	// draw pieces by iterating piece types then squares (fewer allocations than per-tile inner loop)
	for piece, img := range scaledPieceImg {
		if img == nil {
			continue
		}
		for sq := 0; sq < boardWidth*boardHeight; sq++ {
			if bitoperations.GetBit(globals.Bitboards[piece], sq) == 1 {
				file := sq % boardWidth
				rank := sq / boardWidth
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(file*tileSize), float64(rank*tileSize))
				screen.DrawImage(img, op)
			}
		}
	}

	// panel background
	panel := ebiten.NewImage(ScreenWidth, panelHeight)
	panel.Fill(color.RGBA{R: 40, G: 40, B: 44, A: 230})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(panelY))
	screen.DrawImage(panel, op)

	// reset board button
	btn1X := (ScreenWidth-100)/2 - 140
	btn1Y := panelY + (panelHeight-ctrlBtnH)/2 - 20
	btn1 := ebiten.NewImage(100, ctrlBtnH)
	btn1.Fill(color.RGBA{R: 70, G: 120, B: 200, A: 255})
	btn1Op := &ebiten.DrawImageOptions{}
	btn1Op.GeoM.Translate(float64(btn1X), float64(btn1Y))
	screen.DrawImage(btn1, btn1Op)

	// button label: \(Reset Game\)
	ebitenutil.DebugPrintAt(screen, "New Game", btn1X+18, btn1Y+12)

	if g.clock != nil {
		// draw on the left and right sides of the panel
		ebitenutil.DebugPrintAt(screen, "White: "+formatSeconds(g.clock.White.TimeLeft()), 125, panelY+15)
		ebitenutil.DebugPrintAt(screen, "Black: "+formatSeconds(g.clock.Black.TimeLeft()), ScreenWidth-195, panelY+15)
	}

	// unmake move button
	btn2X := (ScreenWidth-100)/2 + 140
	btn2Y := panelY + (panelHeight-ctrlBtnH)/2 - 20
	btn2 := ebiten.NewImage(100, ctrlBtnH)
	btn2.Fill(color.RGBA{R: 70, G: 120, B: 200, A: 255})
	btn2Op := &ebiten.DrawImageOptions{}
	btn2Op.GeoM.Translate(float64(btn2X), float64(btn2Y))
	screen.DrawImage(btn2, btn2Op)

	// button label: \(Reset Game\)
	ebitenutil.DebugPrintAt(screen, "Undo move", btn2X+18, btn2Y+12)

	// go back to the menu button
	btn3X := (ScreenWidth - 120) / 2
	btn3Y := panelY + (panelHeight-40)/2 + 25
	btn3 := ebiten.NewImage(120, 40)
	btn3.Fill(color.RGBA{R: 70, G: 120, B: 200, A: 255})
	btn3Op := &ebiten.DrawImageOptions{}
	btn3Op.GeoM.Translate(float64(btn3X), float64(btn3Y))
	screen.DrawImage(btn3, btn3Op)

	ebitenutil.DebugPrintAt(screen, "Go back to menu", btn3X+18, btn3Y+12)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
