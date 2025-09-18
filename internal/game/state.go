package game

import (
	"time"
)

// GameMode represents the current mode of the game
type GameMode int

const (
	AttractMode GameMode = iota
	Playing
	GameOver
	HighScore
)

// String returns the string representation of the game mode
func (gm GameMode) String() string {
	switch gm {
	case AttractMode:
		return "AttractMode"
	case Playing:
		return "Playing"
	case GameOver:
		return "GameOver"
	case HighScore:
		return "HighScore"
	default:
		return "Unknown"
	}
}

// GameState represents the complete state of the game
type GameState struct {
	// Game mode and flow
	Mode        GameMode
	Paused      bool
	GameStarted bool
	GameEnded   bool

	// Player state
	Player      *PlayerShip
	Lives       int
	Score       int
	HighScore   int

	// Game entities
	Invaders    []*Invader
	Bullets     []*Bullet
	UFO         *UFO
	Barriers    [][]bool // 2D array representing barrier blocks

	// Game timing
	Wave         int
	WaveCleared  bool
	LastUpdate   time.Time
	DeltaTime    float64

	// Game world dimensions
	ScreenWidth  int
	ScreenHeight int

	// Game timing constants (in seconds)
	FixedDeltaTime float64 // 1/20 = 0.05 for 20Hz updates

	// Input state
	InputState   *InputState
}

// InputState tracks the current input state
type InputState struct {
	LeftPressed  bool
	RightPressed bool
	FirePressed  bool
	FireJustPressed bool
	PauseJustPressed bool
}

// NewGameState creates a new game state with default values
func NewGameState(screenWidth, screenHeight int) *GameState {
	return &GameState{
		Mode:           AttractMode,
		Lives:          3,
		Score:          0,
		HighScore:      0,
		Wave:           1,
		ScreenWidth:    screenWidth,
		ScreenHeight:   screenHeight,
		FixedDeltaTime: 1.0 / 20.0, // 20Hz update rate
		InputState:     &InputState{},
		LastUpdate:     time.Now(),
	}
}

// InitializeNewGame sets up a fresh game state for starting a new game
func (gs *GameState) InitializeNewGame() {
	gs.Mode = Playing
	gs.Paused = false
	gs.GameStarted = true
	gs.GameEnded = false
	gs.Lives = 3
	gs.Score = 0
	gs.Wave = 1
	gs.WaveCleared = false

	// Initialize player
	gs.Player = NewPlayerShip(float64(gs.ScreenWidth/2), float64(gs.ScreenHeight-40))

	// Initialize invaders
	gs.initializeInvaders()

	// Clear bullets and UFO
	gs.Bullets = []*Bullet{}
	gs.UFO = nil

	// Initialize barriers
	gs.initializeBarriers()

	// Reset input state
	gs.InputState = &InputState{}
	gs.LastUpdate = time.Now()
}

// ResetToAttractMode resets the game state to attract mode
func (gs *GameState) ResetToAttractMode() {
	gs.Mode = AttractMode
	gs.Paused = false
	gs.GameStarted = false
	gs.GameEnded = false
	gs.Player = nil
	gs.Invaders = []*Invader{}
	gs.Bullets = []*Bullet{}
	gs.UFO = nil
	gs.InputState = &InputState{}
}

// GameOver transitions the game to game over state
func (gs *GameState) GameOver() {
	gs.Mode = GameOver
	gs.GameEnded = true

	// Update high score if necessary
	if gs.Score > gs.HighScore {
		gs.HighScore = gs.Score
	}
}

// NextWave advances to the next wave
func (gs *GameState) NextWave() {
	gs.Wave++
	gs.WaveCleared = false
	gs.initializeInvaders()

	// Reset player position
	if gs.Player != nil {
		gs.Player.Position.X = float64(gs.ScreenWidth / 2)
		gs.Player.Position.Y = float64(gs.ScreenHeight - 40)
		gs.Player.Velocity.X = 0
	}

	// Clear player bullets (but keep enemy bullets for challenge)
	newBullets := []*Bullet{}
	for _, bullet := range gs.Bullets {
		if !bullet.IsPlayerBullet {
			newBullets = append(newBullets, bullet)
		}
	}
	gs.Bullets = newBullets
}

// initializeInvaders creates the initial invader formation
func (gs *GameState) initializeInvaders() {
	gs.Invaders = []*Invader{}

	// Grid configuration
	const rows = 5
	const cols = 11
	const spacingX = 40
	const spacingY = 30
	const startX = 100
	const startY = 80

	for row := 0; row < rows; row++ {
		var invaderType InvaderType
		var points int

		// Different invader types by row
		switch row {
		case 0:
			invaderType = InvaderTypeSmall
			points = 30
		case 1, 2:
			invaderType = InvaderTypeMedium
			points = 20
		case 3, 4:
			invaderType = InvaderTypeLarge
			points = 10
		}

		for col := 0; col < cols; col++ {
			x := float64(startX + col*spacingX)
			y := float64(startY + row*spacingY)

			invader := NewInvader(invaderType, x, y, points)
			gs.Invaders = append(gs.Invaders, invader)
		}
	}
}

// initializeBarriers creates the defensive barriers
func (gs *GameState) initializeBarriers() {
	// Simple barrier implementation - 4 barriers across the screen
	const barrierCount = 4
	const barrierWidth = 22
	const barrierHeight = 16

	// Initialize barriers array
	gs.Barriers = make([][]bool, barrierCount*barrierWidth)
	for i := range gs.Barriers {
		gs.Barriers[i] = make([]bool, barrierHeight)
	}

	for barrier := 0; barrier < barrierCount; barrier++ {
		startX := barrier * barrierWidth

		// Fill in barrier blocks
		for x := 0; x < barrierWidth; x++ {
			for y := 0; y < barrierHeight; y++ {
				// Create a simple rectangular barrier with some gaps
				if y < 3 || y > barrierHeight-4 || x < 3 || x > barrierWidth-4 {
					continue // Leave edges open
				}
				if y > 8 && y < 12 && x > 8 && x < 14 {
					continue // Leave center gap
				}

				if startX+x < len(gs.Barriers) && y < len(gs.Barriers[0]) {
					gs.Barriers[startX+x][y] = true
				}
			}
		}
	}
}

// IsWaveCleared checks if all invaders have been destroyed
func (gs *GameState) IsWaveCleared() bool {
	return len(gs.Invaders) == 0
}

// GetLiveInvaderCount returns the number of remaining invaders
func (gs *GameState) GetLiveInvaderCount() int {
	count := 0
	for _, invader := range gs.Invaders {
		if invader.Alive {
			count++
		}
	}
	return count
}

// AddScore adds points to the player's score
func (gs *GameState) AddScore(points int) {
	gs.Score += points
	if gs.Score > gs.HighScore {
		gs.HighScore = gs.Score
	}
}

// LoseLife removes a life from the player
func (gs *GameState) LoseLife() {
	gs.Lives--
	if gs.Lives <= 0 {
		gs.GameOver()
	}
}

// TogglePause toggles the pause state of the game
func (gs *GameState) TogglePause() {
	if gs.Mode == Playing {
		gs.Paused = !gs.Paused
	}
}

// UpdateDeltaTime calculates and updates the delta time since last update
func (gs *GameState) UpdateDeltaTime() {
	now := time.Now()
	gs.DeltaTime = now.Sub(gs.LastUpdate).Seconds()
	gs.LastUpdate = now
}