package main

import (
	"context"
	"fmt"
	"log"
	"syscall/js"
	"time"
)

// Game represents the main game state
type Game struct {
	canvas   js.Value
	ctx      js.Value
	width    int
	height   int
	running  bool
	lastTime float64
}

// NewGame creates a new game instance
func NewGame(canvas js.Value) *Game {
	ctx := canvas.Call("getContext", "2d")
	width := canvas.Get("width").Int()
	height := canvas.Get("height").Int()

	return &Game{
		canvas: canvas,
		ctx:    ctx,
		width:  width,
		height: height,
	}
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

// update handles game logic updates
func (g *Game) update(deltaTime float64) {
	// TODO: Implement game logic updates
	// This will include:
	// - Camera input processing
	// - Player movement detection
	// - Enemy movement
	// - Collision detection
	// - Score updates
}

// render handles drawing the game
func (g *Game) render() {
	// Clear canvas
	g.ctx.Set("fillStyle", "#000000")
	g.ctx.Call("fillRect", 0, 0, g.width, g.height)

	// TODO: Implement game rendering
	// This will include:
	// - Drawing player
	// - Drawing enemies
	// - Drawing projectiles
	// - Drawing UI elements (score, lives, etc.)

	// Placeholder: Draw a simple message
	g.ctx.Set("fillStyle", "#00FF00")
	g.ctx.Set("font", "24px Arial")
	g.ctx.Set("textAlign", "center")
	g.ctx.Call("fillText", "BOBN - Camera Ready!", g.width/2, g.height/2)
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
		case "Escape":
			game.Stop()
			log.Println("Game stopped")
		}

		return nil
	}))

	// Start the game
	game.Start()
	log.Println("Game started successfully")
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