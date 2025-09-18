package game

import (
	"math"
	"time"
)

// Engine handles the core game loop and logic
type Engine struct {
	state           *GameState
	lastUFOTime     time.Time
	gameStartTime   time.Time
	invaderMoveTimer float64
	invaderDropTimer float64

	// Invader movement parameters
	invaderMoveSpeed     float64
	invaderDropDistance  float64
	invaderMoveInterval  float64
	baseInvaderSpeed     float64

	// Timing accumulators for fixed timestep
	accumulator float64
}

// NewEngine creates a new game engine
func NewEngine(screenWidth, screenHeight int) *Engine {
	return &Engine{
		state:                NewGameState(screenWidth, screenHeight),
		lastUFOTime:          time.Now(),
		gameStartTime:        time.Now(),
		baseInvaderSpeed:     1.0,  // base speed multiplier
		invaderDropDistance:  20.0, // pixels to drop down
		invaderMoveInterval:  1.0,  // seconds between horizontal moves
	}
}

// GetState returns the current game state
func (e *Engine) GetState() *GameState {
	return e.state
}

// StartNewGame initializes a new game
func (e *Engine) StartNewGame() {
	e.state.InitializeNewGame()
	e.gameStartTime = time.Now()
	e.lastUFOTime = time.Now()
	e.resetInvaderMovement()
}

// ProcessAnalogInput processes analog input for camera control
func (e *Engine) ProcessAnalogInput(analogX float64, firePressed, fireJustPressed, pauseJustPressed bool) {
	// Handle mode-specific input
	switch e.state.Mode {
	case AttractMode:
		if fireJustPressed || pauseJustPressed {
			e.StartNewGame()
		}
	case Playing:
		if pauseJustPressed {
			e.state.TogglePause()
		}
		if !e.state.Paused && e.state.Player != nil && e.state.Player.Alive {
			// Direct position control based on analog input
			// Map analogX (-1 to 1) to screen position
			centerX := float64(e.state.ScreenWidth) / 2
			maxOffset := float64(e.state.ScreenWidth) / 2 - 30 // Keep ship on screen

			// Set player position directly based on head position
			targetX := centerX + (analogX * maxOffset)

			// Smooth the movement slightly
			currentX := e.state.Player.Position.X
			e.state.Player.Position.X = currentX*0.3 + targetX*0.7

			// Keep within bounds
			if e.state.Player.Position.X < 30 {
				e.state.Player.Position.X = 30
			}
			if e.state.Player.Position.X > float64(e.state.ScreenWidth)-30 {
				e.state.Player.Position.X = float64(e.state.ScreenWidth) - 30
			}

			// Handle shooting
			if firePressed {
				bullet := e.state.Player.TryShoot()
				if bullet != nil {
					e.state.Bullets = append(e.state.Bullets, bullet)
				}
			}
		}
	case GameOver, HighScore:
		if fireJustPressed || pauseJustPressed {
			e.state.Mode = AttractMode
		}
	}
}

// ProcessInput processes input events and updates input state
func (e *Engine) ProcessInput(leftPressed, rightPressed, firePressed, fireJustPressed, pauseJustPressed bool) {
	input := e.state.InputState

	// Update input state
	input.LeftPressed = leftPressed
	input.RightPressed = rightPressed
	input.FirePressed = firePressed
	input.FireJustPressed = fireJustPressed
	input.PauseJustPressed = pauseJustPressed

	// Handle mode-specific input
	switch e.state.Mode {
	case AttractMode:
		if fireJustPressed {
			e.StartNewGame()
		}
	case Playing:
		if pauseJustPressed {
			e.state.TogglePause()
		}
		if !e.state.Paused {
			e.processPlayingInput()
		}
	case GameOver:
		if fireJustPressed {
			e.state.ResetToAttractMode()
		}
	case HighScore:
		// Handle high score input if needed
		if fireJustPressed {
			e.state.ResetToAttractMode()
		}
	}
}

// processPlayingInput handles input during gameplay
func (e *Engine) processPlayingInput() {
	input := e.state.InputState
	player := e.state.Player

	if player == nil || !player.Alive {
		return
	}

	// Handle movement
	player.ApplyInput(input.LeftPressed, input.RightPressed, e.state.FixedDeltaTime)

	// Handle shooting
	if input.FireJustPressed {
		bullet := player.TryShoot()
		if bullet != nil {
			e.state.Bullets = append(e.state.Bullets, bullet)
		}
	}
}

// Update runs a fixed timestep update loop
func (e *Engine) Update(deltaTime float64) {
	// Update delta time in state for reference
	e.state.DeltaTime = deltaTime

	// Don't update if paused
	if e.state.Paused {
		return
	}

	// Accumulate time for fixed timestep updates
	e.accumulator += deltaTime

	// Run fixed timestep updates
	for e.accumulator >= e.state.FixedDeltaTime {
		e.fixedUpdate(e.state.FixedDeltaTime)
		e.accumulator -= e.state.FixedDeltaTime
	}
}

// fixedUpdate performs updates at a fixed timestep (20Hz)
func (e *Engine) fixedUpdate(deltaTime float64) {
	switch e.state.Mode {
	case AttractMode:
		e.updateAttractMode(deltaTime)
	case Playing:
		e.updatePlaying(deltaTime)
	case GameOver:
		e.updateGameOver(deltaTime)
	case HighScore:
		e.updateHighScore(deltaTime)
	}
}

// updateAttractMode handles attract mode updates
func (e *Engine) updateAttractMode(deltaTime float64) {
	// Simple attract mode - could show demo gameplay or scrolling text
	// For now, just wait for player input
}

// updatePlaying handles the main gameplay updates
func (e *Engine) updatePlaying(deltaTime float64) {
	// Update player
	if e.state.Player != nil {
		e.state.Player.Update(deltaTime, float64(e.state.ScreenWidth))
	}

	// Update invaders
	e.updateInvaders(deltaTime)

	// Update bullets
	e.updateBullets(deltaTime)

	// Update UFO
	e.updateUFO(deltaTime)

	// Handle collisions
	e.handleCollisions()

	// Check win/lose conditions
	e.checkGameConditions()

	// Spawn UFO occasionally
	e.maybeSpawnUFO()
}

// updateGameOver handles game over state
func (e *Engine) updateGameOver(deltaTime float64) {
	// Game over screen logic - could have animations or effects
}

// updateHighScore handles high score display
func (e *Engine) updateHighScore(deltaTime float64) {
	// High score screen logic
}

// updateInvaders updates all invaders and handles formation movement
func (e *Engine) updateInvaders(deltaTime float64) {
	liveInvaders := []*Invader{}

	// Update individual invaders
	for _, invader := range e.state.Invaders {
		if !invader.Alive {
			continue
		}

		invader.Update(deltaTime)
		liveInvaders = append(liveInvaders, invader)

		// Handle invader shooting
		if bullet := invader.TryShoot(deltaTime); bullet != nil {
			e.state.Bullets = append(e.state.Bullets, bullet)
		}
	}

	e.state.Invaders = liveInvaders

	// Handle formation movement
	e.updateInvaderFormation(deltaTime)
}

// updateInvaderFormation handles the classic invader formation movement
func (e *Engine) updateInvaderFormation(deltaTime float64) {
	if len(e.state.Invaders) == 0 {
		return
	}

	e.invaderMoveTimer += deltaTime

	// Calculate movement speed based on remaining invaders (fewer = faster)
	invaderCount := float64(len(e.state.Invaders))
	speedMultiplier := e.baseInvaderSpeed * (55.0 / (invaderCount + 5.0))
	currentMoveInterval := e.invaderMoveInterval / speedMultiplier

	if e.invaderMoveTimer >= currentMoveInterval {
		e.invaderMoveTimer = 0

		// Find the bounds of the formation
		leftmost, rightmost := e.findInvaderBounds()

		// Determine if we need to drop down and reverse direction
		shouldDrop := false
		direction := e.state.Invaders[0].Direction

		if direction > 0 && rightmost >= float64(e.state.ScreenWidth-20) {
			shouldDrop = true
			direction = -1
		} else if direction < 0 && leftmost <= 20 {
			shouldDrop = true
			direction = 1
		}

		// Move all invaders
		moveDistance := 10.0 * float64(direction)

		for _, invader := range e.state.Invaders {
			invader.Direction = direction

			if shouldDrop {
				invader.Move(0, e.invaderDropDistance)
			} else {
				invader.Move(moveDistance, 0)
			}
		}

		// Check if invaders reached the bottom
		e.checkInvaderReachBottom()
	}
}

// findInvaderBounds finds the leftmost and rightmost invader positions
func (e *Engine) findInvaderBounds() (leftmost, rightmost float64) {
	if len(e.state.Invaders) == 0 {
		return 0, 0
	}

	leftmost = e.state.Invaders[0].Position.X
	rightmost = e.state.Invaders[0].Position.X

	for _, invader := range e.state.Invaders {
		if invader.Position.X < leftmost {
			leftmost = invader.Position.X
		}
		if invader.Position.X > rightmost {
			rightmost = invader.Position.X
		}
	}

	return leftmost, rightmost
}

// checkInvaderReachBottom checks if any invader has reached the bottom
func (e *Engine) checkInvaderReachBottom() {
	bottomLine := float64(e.state.ScreenHeight - 100) // Line above player area

	for _, invader := range e.state.Invaders {
		if invader.Position.Y >= bottomLine {
			// Game over - invaders reached the bottom
			e.state.GameOver()
			return
		}
	}
}

// updateBullets updates all bullets and removes dead ones
func (e *Engine) updateBullets(deltaTime float64) {
	liveBullets := []*Bullet{}

	for _, bullet := range e.state.Bullets {
		if !bullet.Alive {
			continue
		}

		bullet.Update(deltaTime, float64(e.state.ScreenWidth), float64(e.state.ScreenHeight))

		if bullet.Alive {
			liveBullets = append(liveBullets, bullet)
		}
	}

	e.state.Bullets = liveBullets
}

// updateUFO updates the UFO if it exists
func (e *Engine) updateUFO(deltaTime float64) {
	if e.state.UFO != nil {
		e.state.UFO.Update(deltaTime, float64(e.state.ScreenWidth))
		if !e.state.UFO.Alive {
			e.state.UFO = nil
		}
	}
}

// maybeSpawnUFO spawns a UFO occasionally
func (e *Engine) maybeSpawnUFO() {
	if e.state.UFO != nil {
		return // UFO already exists
	}

	gameTime := time.Since(e.gameStartTime).Seconds()
	if ShouldSpawnUFO(e.lastUFOTime, gameTime) {
		e.lastUFOTime = time.Now()

		// Spawn from random side
		var startX float64
		var direction int

		if math.Mod(gameTime, 2.0) < 1.0 {
			// Spawn from left
			startX = -50
			direction = 1
		} else {
			// Spawn from right
			startX = float64(e.state.ScreenWidth) + 50
			direction = -1
		}

		e.state.UFO = NewUFO(startX, 50, direction)
	}
}

// handleCollisions handles all collision detection and responses
func (e *Engine) handleCollisions() {
	// Player bullets vs invaders
	e.handlePlayerBulletCollisions()

	// Player bullets vs UFO
	e.handlePlayerBulletUFOCollisions()

	// Enemy bullets vs player
	e.handleEnemyBulletCollisions()

	// Player vs enemy bullets (already handled above)
	// Bullets vs barriers would go here if implemented
}

// handlePlayerBulletCollisions handles collisions between player bullets and invaders
func (e *Engine) handlePlayerBulletCollisions() {
	for _, bullet := range e.state.Bullets {
		if !bullet.Alive || !bullet.IsPlayerBullet {
			continue
		}

		for _, invader := range e.state.Invaders {
			if !invader.Alive {
				continue
			}

			if bullet.Bounds.Intersects(invader.Bounds) {
				// Collision detected
				bullet.Alive = false
				invader.Alive = false
				e.state.AddScore(invader.Points)
				break // Bullet can only hit one invader
			}
		}
	}
}

// handlePlayerBulletUFOCollisions handles collisions between player bullets and UFO
func (e *Engine) handlePlayerBulletUFOCollisions() {
	if e.state.UFO == nil || !e.state.UFO.Alive {
		return
	}

	for _, bullet := range e.state.Bullets {
		if !bullet.Alive || !bullet.IsPlayerBullet {
			continue
		}

		if bullet.Bounds.Intersects(e.state.UFO.Bounds) {
			// Collision detected
			bullet.Alive = false
			e.state.UFO.Alive = false
			e.state.AddScore(e.state.UFO.Points)
			break // Bullet hits UFO
		}
	}
}

// handleEnemyBulletCollisions handles collisions between enemy bullets and player
func (e *Engine) handleEnemyBulletCollisions() {
	if e.state.Player == nil || !e.state.Player.Alive {
		return
	}

	for _, bullet := range e.state.Bullets {
		if !bullet.Alive || bullet.IsPlayerBullet {
			continue
		}

		if bullet.Bounds.Intersects(e.state.Player.Bounds) {
			// Player hit by enemy bullet
			bullet.Alive = false
			e.state.Player.Alive = false
			e.state.LoseLife()

			// Respawn player if lives remaining
			if e.state.Lives > 0 {
				e.respawnPlayer()
			}
			break
		}
	}
}

// respawnPlayer respawns the player after a brief delay
func (e *Engine) respawnPlayer() {
	// For now, respawn immediately at starting position
	e.state.Player = NewPlayerShip(float64(e.state.ScreenWidth/2), float64(e.state.ScreenHeight-40))

	// Clear enemy bullets for fairness
	playerBullets := []*Bullet{}
	for _, bullet := range e.state.Bullets {
		if bullet.IsPlayerBullet {
			playerBullets = append(playerBullets, bullet)
		}
	}
	e.state.Bullets = playerBullets
}

// checkGameConditions checks for win/lose conditions
func (e *Engine) checkGameConditions() {
	// Check if wave is cleared
	if e.state.IsWaveCleared() && !e.state.WaveCleared {
		e.state.WaveCleared = true
		// Start next wave after a brief delay
		// For now, immediately start next wave
		e.state.NextWave()
		e.resetInvaderMovement()
	}
}

// resetInvaderMovement resets invader movement timing
func (e *Engine) resetInvaderMovement() {
	e.invaderMoveTimer = 0
	e.invaderDropTimer = 0
}