package main

import (
	"context"
	"fmt"
	"log"
	"syscall/js"
	"time"

	"github.com/jonasrmichel/bobn/internal/game"
	"github.com/jonasrmichel/bobn/internal/wasm"
)

// Game represents the main game state
type Game struct {
	canvas    js.Value
	ctx       js.Value
	width     int
	height    int
	running   bool
	lastTime  float64

	// Game components
	bridge    *wasm.JSBridge
	renderer  *wasm.Renderer
	engine    *game.Engine
	camera    *wasm.CameraController

	// Timing
	accumulator float64
	frameTime   float64

	// Camera input
	cameraX   float64
	cameraY   float64
}

// NewGame creates a new game instance
func NewGame(canvas js.Value) *Game {
	ctx := canvas.Call("getContext", "2d")
	if ctx.IsUndefined() || ctx.IsNull() {
		log.Fatal("Failed to get 2D context from canvas")
		return nil
	}

	width := canvas.Get("width").Int()
	height := canvas.Get("height").Int()

	// Initialize game components
	bridge := wasm.NewJSBridge()
	err := bridge.Initialize("gameCanvas")
	if err != nil {
		log.Printf("Failed to initialize bridge: %v", err)
		return nil
	}

	engine := game.NewEngine(width, height)
	renderer := wasm.NewRenderer(bridge, width, height)

	// Set the renderer to use the same context
	renderer.SetContext(ctx)

	// Initialize camera controller
	camera := wasm.NewCameraController()
	camera.Initialize()

	g := &Game{
		canvas:      canvas,
		ctx:         ctx,
		width:       width,
		height:      height,
		bridge:      bridge,
		engine:      engine,
		renderer:    renderer,
		camera:      camera,
		frameTime:   1000.0 / 60.0, // 60 FPS target
	}

	// Set up camera position callback
	camera.SetPositionCallback(func(x, y float64) {
		g.cameraX = x
		g.cameraY = y
	})

	return g
}

// Start begins the game loop
func (g *Game) Start() {
	log.Println("Start() called - starting game loop")
	g.running = true
	g.gameLoop()
}

// Stop ends the game loop
func (g *Game) Stop() {
	g.running = false
}

// gameLoop runs the main game loop
func (g *Game) gameLoop() {
	log.Println("gameLoop started")
	frameCount := 0

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !g.running {
			renderFrame.Release()
			return nil
		}

		currentTime := args[0].Float()
		if g.lastTime == 0 {
			g.lastTime = currentTime
		}

		deltaTime := currentTime - g.lastTime
		g.lastTime = currentTime

		frameCount++
		if frameCount == 1 {
			log.Println("First frame - calling update and render")
		}

		g.update(deltaTime)
		g.render()

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	log.Println("Calling requestAnimationFrame")
	js.Global().Call("requestAnimationFrame", renderFrame)
}

// update handles game logic updates with fixed timestep
func (g *Game) update(deltaTime float64) {
	// Fixed timestep accumulator pattern for consistent physics
	g.accumulator += deltaTime

	// Fixed update step (50ms = 20Hz)
	fixedTimeStep := 50.0
	for g.accumulator >= fixedTimeStep {
		// Get input state from bridge
		input := g.bridge.GetInputState()

		// If camera is enabled, use analog control
		if g.camera.IsEnabled() && g.engine.GetState().Mode == game.Playing {
			// Use camera position for analog control
			g.engine.ProcessAnalogInput(
				g.cameraX,  // Analog X position (-1 to 1)
				input.FirePressed,  // Only fire when Space is pressed
				input.FireJustPressed,
				input.PauseJustPressed || input.EnterJustPressed,
			)
		} else {
			// Use digital keyboard input
			leftPressed := input.LeftPressed
			rightPressed := input.RightPressed

			// Process input and update game state
			g.engine.ProcessInput(
				leftPressed,
				rightPressed,
				input.FirePressed,
				input.FireJustPressed,
				input.PauseJustPressed || input.EnterJustPressed,
			)
		}
		g.engine.Update(fixedTimeStep / 1000.0) // Convert to seconds

		g.accumulator -= fixedTimeStep
	}

	// Update UI elements in HTML
	g.updateUI()
}

// render handles drawing the game
func (g *Game) render() {
	// Use the stored context
	ctx := g.ctx

	// Check if context is valid
	if !ctx.Truthy() {
		log.Println("ERROR: Canvas context is not valid!")
		return
	}

	// Clear canvas
	ctx.Set("fillStyle", "#000000")
	ctx.Call("fillRect", 0, 0, g.width, g.height)

	if g.engine == nil {
		// Show loading message if not ready
		ctx.Set("fillStyle", "#00ff00")
		ctx.Set("font", "20px monospace")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", "ENGINE NOT INITIALIZED", g.width/2, g.height/2)
		return
	}

	// Draw the game directly
	state := g.engine.GetState()

	// Draw stars background
	ctx.Set("fillStyle", "#ffffff")
	for i := 0; i < 50; i++ {
		x := (i * 73) % g.width
		y := (i * 37) % g.height
		ctx.Call("fillRect", x, y, 2, 2)
	}

	// Draw game based on mode
	switch state.Mode {
	case game.AttractMode:
		// Draw title
		ctx.Set("fillStyle", "#00ff00")
		ctx.Set("font", "48px monospace")
		ctx.Set("textAlign", "center")
		ctx.Set("textBaseline", "middle")
		ctx.Call("fillText", "BOBN", g.width/2, 150)

		ctx.Set("font", "20px monospace")
		ctx.Set("fillStyle", "#00ffff")
		ctx.Call("fillText", "SPACE INVADERS", g.width/2, 200)

		// Instructions
		ctx.Set("font", "16px monospace")
		ctx.Set("fillStyle", "#ffff00")
		ctx.Call("fillText", "PRESS ENTER TO START", g.width/2, 300)
		ctx.Call("fillText", "USE ARROWS TO MOVE", g.width/2, 330)
		ctx.Call("fillText", "SPACE TO FIRE", g.width/2, 360)

	case game.Playing:
		// Draw player ship
		if state.Player != nil && state.Player.Alive {
			ctx.Set("fillStyle", "#00ff00")
			// Simple triangle ship
			ctx.Call("beginPath")
			ctx.Call("moveTo", state.Player.Position.X, state.Player.Position.Y)
			ctx.Call("lineTo", state.Player.Position.X-15, state.Player.Position.Y+20)
			ctx.Call("lineTo", state.Player.Position.X+15, state.Player.Position.Y+20)
			ctx.Call("closePath")
			ctx.Call("fill")
		}

		// Draw invaders
		for _, invader := range state.Invaders {
			if invader.Alive {
				ctx.Set("fillStyle", "#ff00ff")
				ctx.Call("fillRect", invader.Position.X-10, invader.Position.Y-5, 20, 10)
				// Eyes
				ctx.Set("fillStyle", "#000000")
				ctx.Call("fillRect", invader.Position.X-6, invader.Position.Y-2, 3, 3)
				ctx.Call("fillRect", invader.Position.X+3, invader.Position.Y-2, 3, 3)
			}
		}

		// Draw bullets
		for _, bullet := range state.Bullets {
			if bullet.Alive {
				if bullet.IsPlayerBullet {
					ctx.Set("fillStyle", "#00ff00")
				} else {
					ctx.Set("fillStyle", "#ff0000")
				}
				ctx.Call("fillRect", bullet.Position.X-1, bullet.Position.Y, 2, 8)
			}
		}

	case game.GameOver:
		ctx.Set("fillStyle", "#ff0000")
		ctx.Set("font", "48px monospace")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", "GAME OVER", g.width/2, g.height/2)

		ctx.Set("font", "20px monospace")
		ctx.Set("fillStyle", "#ffffff")
		ctx.Call("fillText", fmt.Sprintf("SCORE: %d", state.Score), g.width/2, g.height/2+50)
	}

	// Draw UI (score, lives, etc)
	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("font", "16px monospace")
	ctx.Set("textAlign", "left")
	ctx.Call("fillText", fmt.Sprintf("SCORE: %06d", state.Score), 10, 30)

	ctx.Set("textAlign", "center")
	ctx.Call("fillText", fmt.Sprintf("HIGH: %06d", state.HighScore), g.width/2, 30)

	ctx.Set("textAlign", "right")
	ctx.Call("fillText", fmt.Sprintf("LIVES: %d", state.Lives), g.width-10, 30)

	// Draw wave number
	if state.Mode == game.Playing {
		ctx.Set("textAlign", "center")
		ctx.Set("fillStyle", "#00ffff")
		ctx.Call("fillText", fmt.Sprintf("WAVE %d", state.Wave), g.width/2, g.height-20)
	}
}

// updateUI updates the HTML UI elements
func (g *Game) updateUI() {
	state := g.engine.GetState()
	doc := js.Global().Get("document")

	// Update score
	if scoreElem := doc.Call("getElementById", "score"); !scoreElem.IsUndefined() && !scoreElem.IsNull() {
		scoreElem.Set("textContent", fmt.Sprintf("%06d", state.Score))
	}

	// Update high score
	if highScoreElem := doc.Call("getElementById", "highScore"); !highScoreElem.IsUndefined() && !highScoreElem.IsNull() {
		highScoreElem.Set("textContent", fmt.Sprintf("%06d", state.HighScore))
	}

	// Update lives
	if livesElem := doc.Call("getElementById", "lives"); !livesElem.IsUndefined() && !livesElem.IsNull() {
		livesElem.Set("textContent", fmt.Sprintf("%d", state.Lives))
	}

	// Update level
	if levelElem := doc.Call("getElementById", "level"); !levelElem.IsUndefined() && !levelElem.IsNull() {
		levelElem.Set("textContent", fmt.Sprintf("%d", state.Wave))
	}

	// Update status message
	if statusElem := doc.Call("getElementById", "status"); !statusElem.IsUndefined() && !statusElem.IsNull() {
		var status string
		switch state.Mode {
		case game.AttractMode:
			status = "PRESS START TO PLAY"
		case game.Playing:
			status = "PLAYING"
		case game.GameOver:
			status = "GAME OVER"
		case game.HighScore:
			status = "NEW HIGH SCORE!"
		}
		statusElem.Set("textContent", status)
	}
}

// initializeGame sets up the game and starts it
func initializeGame() {
	canvas := js.Global().Get("document").Call("getElementById", "gameCanvas")
	if canvas.IsUndefined() || canvas.IsNull() {
		log.Fatal("Could not find canvas element with id 'gameCanvas'")
		return
	}

	game := NewGame(canvas)

	// Setup event listeners for camera controls (placeholder)
	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()

		switch key {
		case " ": // Spacebar to start/stop
			if game.running {
				game.Stop()
				log.Println("Game stopped")
			} else {
				game.Start()
				log.Println("Game started")
			}
		case "Enter":
			// Handle Enter key press for game
			log.Println("Enter pressed")
		case "Escape":
			game.Stop()
			log.Println("Game stopped")
		}

		return nil
	}))

	// Start the game automatically
	if game != nil {
		log.Println("Game created, calling Start()")
		game.Start()
		log.Println("Game started successfully")
	} else {
		log.Println("Failed to create game instance")
	}
}

func main() {
	fmt.Println("BOBN WASM module loaded")

	// Check document state immediately
	readyState := js.Global().Get("document").Get("readyState").String()
	fmt.Printf("Document readyState: %s\n", readyState)

	// Wait for the DOM to be ready
	ready := make(chan struct{})

	if readyState == "complete" || readyState == "interactive" {
		// DOM is already ready
		fmt.Println("DOM is already ready, initializing immediately")
		close(ready)
	} else {
		// Wait for DOM to be ready
		fmt.Println("Waiting for DOM to be ready...")
		js.Global().Get("document").Call("addEventListener", "DOMContentLoaded", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("DOMContentLoaded event fired")
			close(ready)
			return nil
		}))
	}

	// Wait for DOM ready or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ready:
		// Add a small delay to ensure canvas is fully rendered
		time.Sleep(100 * time.Millisecond)
		initializeGame()
	case <-ctx.Done():
		log.Fatal("Timeout waiting for DOM to be ready")
	}

	// Keep the program running
	select {}
}