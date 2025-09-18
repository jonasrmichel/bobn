package main

import (
	"context"
	"fmt"
	"log"
	"math"
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
	g.running = true
	g.gameLoop()
}

// Stop ends the game loop
func (g *Game) Stop() {
	g.running = false
}

// gameLoop runs the main game loop
func (g *Game) gameLoop() {
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

		g.update(deltaTime)
		g.render()

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

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

		// Merge camera input with keyboard input
		leftPressed := input.LeftPressed || (g.camera.IsEnabled() && g.cameraX < -0.3)
		rightPressed := input.RightPressed || (g.camera.IsEnabled() && g.cameraX > 0.3)

		// Auto-fire when camera is enabled and player is moving
		firePressed := input.FirePressed
		if g.camera.IsEnabled() && math.Abs(g.cameraY) < 0.3 {
			// Fire when head is centered vertically
			firePressed = true
		}

		// Process input and update game state
		g.engine.ProcessInput(
			leftPressed,
			rightPressed,
			firePressed,
			input.FireJustPressed,
			input.PauseJustPressed || input.EnterJustPressed,
		)
		g.engine.Update(fixedTimeStep / 1000.0) // Convert to seconds

		g.accumulator -= fixedTimeStep
	}

	// Update UI elements in HTML
	g.updateUI()
}

// render handles drawing the game
func (g *Game) render() {
	// Render the game using the renderer
	g.renderer.RenderGame(g.engine.GetState())
}

// updateUI updates the HTML UI elements
func (g *Game) updateUI() {
	state := g.engine.GetState()
	doc := js.Global().Get("document")

	// Update score
	if scoreElem := doc.Call("getElementById", "score"); !scoreElem.IsUndefined() {
		scoreElem.Set("textContent", fmt.Sprintf("%06d", state.Score))
	}

	// Update high score
	if highScoreElem := doc.Call("getElementById", "highScore"); !highScoreElem.IsUndefined() {
		highScoreElem.Set("textContent", fmt.Sprintf("%06d", state.HighScore))
	}

	// Update lives
	if livesElem := doc.Call("getElementById", "lives"); !livesElem.IsUndefined() {
		livesElem.Set("textContent", fmt.Sprintf("%d", state.Lives))
	}

	// Update level
	if levelElem := doc.Call("getElementById", "level"); !levelElem.IsUndefined() {
		levelElem.Set("textContent", fmt.Sprintf("%d", state.Wave))
	}

	// Update status message
	if statusElem := doc.Call("getElementById", "status"); !statusElem.IsUndefined() {
		var status string
		switch state.Mode {
		case game.AttractMode:
			status = "INSERT COIN TO PLAY"
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
	log.Println("Initializing BOBN WASM game")

	canvas := js.Global().Get("document").Call("getElementById", "gameCanvas")
	if canvas.IsUndefined() {
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
		game.Start()
		log.Println("Game started successfully")
	} else {
		log.Println("Failed to create game instance")
	}
}

func main() {
	fmt.Println("BOBN WASM module loaded")

	// Wait for the DOM to be ready
	ready := make(chan struct{})

	js.Global().Get("document").Call("addEventListener", "DOMContentLoaded", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(ready)
		return nil
	}))

	// If DOM is already loaded, start immediately
	if js.Global().Get("document").Get("readyState").String() == "complete" {
		close(ready)
	}

	// Wait for DOM ready or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ready:
		initializeGame()
	case <-ctx.Done():
		log.Fatal("Timeout waiting for DOM to be ready")
	}

	// Keep the program running
	select {}
}