package clock

import (
	"fmt"
	"sync"
	"time"
	"zerginator/globals"
)

// MaxTimeComputerMove is the maximum time a computer move can take.
const MaxTimeComputerMove = 15 * time.Second

// Clock manages two-player time controls.
type Clock struct {
	mu        sync.Mutex
	remaining time.Duration
	running   bool
	lastStart time.Time
}

// NewClock creates a new clock with the given time in minutes.
func NewClock(minutes int) *Clock {
	return &Clock{remaining: time.Duration(minutes) * time.Minute}
}

// Start starts the clock.
func (c *Clock) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return
	}
	c.lastStart = time.Now()
	c.running = true
}

// Stop stops the clock.
func (c *Clock) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return
	}
	elapsed := time.Since(c.lastStart)
	c.remaining -= elapsed
	c.running = false
	if c.remaining <= 0 {
		c.remaining = 0
	}
}

// TimeLeft returns the remaining time on the clock.
func (c *Clock) TimeLeft() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return c.remaining - time.Since(c.lastStart)
	}
	return c.remaining
}

// IsExpired returns true if the clock has expired.
func (c *Clock) IsExpired() bool {
	return c.TimeLeft() <= 0
}

// GameClock manages the clocks for both players and handles turn switching.
type GameClock struct {
	White *Clock
	Black *Clock
	turn  int
}

// NewGameClock creates a new GameClock with default time settings.
func NewGameClock() *GameClock {
	return &GameClock{
		White: NewClock(10),
		Black: NewClock(10),
		turn:  0,
	}
}

// SwitchTurn switches the turn between players and starts/stops the appropriate clock.
func (g *GameClock) SwitchTurn() {
	if g.turn == globals.WHITE {
		g.White.Stop()
		g.Black.Start()
		g.turn = globals.BLACK
	} else {
		g.Black.Stop()
		g.White.Start()
		g.turn = globals.WHITE
	}
}

// Status prints the remaining time for both players.
func (g *GameClock) Status() {
	fmt.Printf("\tWhite: %v | Black: %v\n", g.White.TimeLeft().Round(time.Second), g.Black.TimeLeft().Round(time.Second))
}
