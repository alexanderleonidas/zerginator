# Zerginator

Zerginator is a chess engine written in Go. It provides a UCI-compatible headless engine and an optional graphical front-end using `github.com/hajimehoshi/ebiten/v2`. The project focuses on a compact, efficient bitboard-based implementation with tools for testing (perft), evaluation and move-generation.

## Key features
- UCI protocol support (`uci` package) — interactive engine mode and `go depth N` parsing.
- Optional GUI using `ebiten` (`gui` package).
- Bitboard representation for board state.
- Fast move generation and legality checking (`board`, `moves` concepts).
- Perft and debugging helpers for correctness testing.
- FEN parsing and position setup.
- Zobrist hashing (random keys) and a transposition table (`ai` / `board`).
- Pawn evaluation masks and evaluation scaffolding (`ai`).
- Compact encoded move representation and make/unmake move routines.

## Implemented algorithms & techniques
- Bitboards for board representation and fast bitwise operations.
- Leaper attack tables (pawn, knight, king) precomputed at init.
- Sliding attack generation for rook/bishop on-the-fly (masking & occupancy).
- Zobrist hashing (random keys) for position keys.
- Transposition table (hash table) for cached search results.
- Packed move encoding (single integer) for efficient move lists.
- FEN parsing and position setup for testing and UCI.
- Perft driver for move-generation verification.

## Planned / experimental features
- Magic bitboards / magic number generator (commented in code for future use).
- Additional search enhancements (tuning move ordering, quiescence, iterative deepening, etc.).

## Project structure (high level)
- `main.go` — program entry, init routines and mode selection (GUI / UCI / debug).
- `board` — bitboard generation, attack masks, move generation helpers.
- `ai` — evaluation, transposition table and search-related helpers.
- `uci` — UCI protocol parsing and main engine loop.
- `gui` — Ebiten-based graphical front-end and image loading.
- `globals` — shared constants and configuration.

## External libraries & tools
- Go (modules) — language and build system.
- `github.com/hajimehoshi/ebiten/v2` — 2D game library used for the GUI.
- (Optional) `mingw-w64` — when cross-compiling with `cgo` to Windows.
- Recommended dev environment: GoLand (project tested there).

## Build / cross-compile (examples)
- Native build (current OS):
  - `go build -o zerginator ./...`
- Cross-compile for Windows (pure Go, 64-bit):
  - `GOOS=windows GOARCH=amd64 go build -o zerginator.exe ./cmd/zerginator` \
    or if `main` is in repo root: `GOOS=windows GOARCH=amd64 go build -o zerginator.exe .`
- If `cgo` is required, install `mingw-w64` and set `CC`:
  - `brew install mingw-w64`  
  - `CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o zerginator.exe .`

## Usage notes
- GUI mode uses Ebiten windowing; headless mode runs the UCI loop.
- Use the UCI `go depth N` command to trigger depth-limited search from the UCI interface.

## Credits
- Project written in Go. GUI powered by `github.com/hajimehoshi/ebiten/v2`.
- Engine ideas influenced by classic chess programming resources and sample code referenced in comments.
