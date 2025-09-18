package wasm

import (
	"errors"
	"syscall/js"
	"time"
)

// JSBridge handles JavaScript interop for WASM
type JSBridge struct {
	document js.Value
	window   js.Value
	canvas   js.Value
	context  js.Value

	// Event listeners
	keydownListener  js.Func
	keyupListener    js.Func
	resizeListener   js.Func
	focusListener    js.Func
	blurListener     js.Func

	// Input state tracking
	keysPressed map[string]bool
	keysJustPressed map[string]bool

	// Animation frame callback
	animationCallback js.Func
	lastFrameTime     float64

	// Canvas properties
	canvasWidth  int
	canvasHeight int
	deviceRatio  float64
}

// NewJSBridge creates a new JavaScript bridge
func NewJSBridge() *JSBridge {
	bridge := &JSBridge{
		document:    js.Global().Get("document"),
		window:      js.Global(),
		keysPressed: make(map[string]bool),
		keysJustPressed: make(map[string]bool),
		deviceRatio: 1.0,
	}

	// Get device pixel ratio for high DPI displays
	if ratio := bridge.window.Get("devicePixelRatio"); !ratio.IsUndefined() {
		bridge.deviceRatio = ratio.Float()
	}

	return bridge
}

// GetContext returns the canvas 2D context
func (b *JSBridge) GetContext() js.Value {
	return b.context
}

// InputState represents the current input state
type InputState struct {
	LeftPressed      bool
	RightPressed     bool
	UpPressed        bool
	DownPressed      bool
	FirePressed      bool
	FireJustPressed  bool
	PauseJustPressed bool
	EnterJustPressed bool
}

// GetInputState returns the current input state
func (b *JSBridge) GetInputState() InputState {
	state := InputState{
		LeftPressed:      b.keysPressed["ArrowLeft"],
		RightPressed:     b.keysPressed["ArrowRight"],
		UpPressed:        b.keysPressed["ArrowUp"],
		DownPressed:      b.keysPressed["ArrowDown"],
		FirePressed:      b.keysPressed[" "] || b.keysPressed["Space"],
		FireJustPressed:  b.keysJustPressed[" "] || b.keysJustPressed["Space"],
		PauseJustPressed: b.keysJustPressed["Escape"] || b.keysJustPressed["p"] || b.keysJustPressed["P"],
		EnterJustPressed: b.keysJustPressed["Enter"],
	}

	// Clear just pressed keys after reading
	for key := range b.keysJustPressed {
		b.keysJustPressed[key] = false
	}

	return state
}

// Initialize sets up the JavaScript bridge with canvas and event listeners
func (b *JSBridge) Initialize(canvasID string) error {
	// Get canvas element
	b.canvas = b.document.Call("getElementById", canvasID)
	if b.canvas.IsUndefined() {
		return errors.New("Canvas element not found: " + canvasID)
	}

	// Get 2D rendering context
	b.context = b.canvas.Call("getContext", "2d")
	if b.context.IsUndefined() {
		return errors.New("Failed to get 2D context")
	}

	// Setup canvas size
	b.setupCanvas()

	// Setup event listeners
	b.setupEventListeners()

	return nil
}

// setupCanvas configures the canvas for high DPI displays
func (b *JSBridge) setupCanvas() {
	// Get actual canvas size from CSS
	rect := b.canvas.Call("getBoundingClientRect")
	cssWidth := rect.Get("width").Float()
	cssHeight := rect.Get("height").Float()

	// Set actual canvas size accounting for device pixel ratio
	b.canvasWidth = int(cssWidth * b.deviceRatio)
	b.canvasHeight = int(cssHeight * b.deviceRatio)

	// Set canvas buffer size
	b.canvas.Set("width", b.canvasWidth)
	b.canvas.Set("height", b.canvasHeight)

	// Scale the canvas back down using CSS
	b.canvas.Get("style").Set("width", cssWidth)
	b.canvas.Get("style").Set("height", cssHeight)

	// Scale the context to match device pixel ratio
	b.context.Call("scale", b.deviceRatio, b.deviceRatio)

	// Set default context properties
	b.context.Set("imageSmoothingEnabled", false)
	b.context.Set("textAlign", "left")
	b.context.Set("textBaseline", "top")
}

// setupEventListeners sets up keyboard and other event listeners
func (b *JSBridge) setupEventListeners() {
	// Keyboard event listeners
	b.keydownListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()
		code := event.Get("code").String()

		// Track just pressed only if key wasn't already pressed
		if !b.keysPressed[key] {
			b.keysJustPressed[key] = true
		}
		if !b.keysPressed[code] {
			b.keysJustPressed[code] = true
		}

		b.keysPressed[key] = true
		b.keysPressed[code] = true

		// Prevent default for game keys
		if b.isGameKey(key) || b.isGameKey(code) {
			event.Call("preventDefault")
		}

		return nil
	})

	b.keyupListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()
		code := event.Get("code").String()

		b.keysPressed[key] = false
		b.keysPressed[code] = false

		return nil
	})

	// Window resize listener
	b.resizeListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		b.setupCanvas()
		return nil
	})

	// Focus/blur listeners for pausing
	b.focusListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Clear all keys when gaining focus to prevent stuck keys
		for key := range b.keysPressed {
			b.keysPressed[key] = false
		}
		return nil
	})

	b.blurListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Clear all keys when losing focus
		for key := range b.keysPressed {
			b.keysPressed[key] = false
		}
		return nil
	})

	// Add event listeners
	b.document.Call("addEventListener", "keydown", b.keydownListener)
	b.document.Call("addEventListener", "keyup", b.keyupListener)
	b.window.Call("addEventListener", "resize", b.resizeListener)
	b.window.Call("addEventListener", "focus", b.focusListener)
	b.window.Call("addEventListener", "blur", b.blurListener)

	// Make canvas focusable and focus it
	b.canvas.Set("tabIndex", 0)
	b.canvas.Call("focus")
}

// isGameKey checks if a key is used by the game
func (b *JSBridge) isGameKey(key string) bool {
	gameKeys := map[string]bool{
		"ArrowLeft":  true,
		"ArrowRight": true,
		" ":          true,
		"Space":      true,
		"Escape":     true,
		"KeyA":       true,
		"KeyD":       true,
		"KeyP":       true,
		"Enter":      true,
	}
	return gameKeys[key]
}

// StartAnimationLoop starts the animation loop using requestAnimationFrame
func (b *JSBridge) StartAnimationLoop(callback func(float64)) {
	b.animationCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentTime := args[0].Float()

		if b.lastFrameTime == 0 {
			b.lastFrameTime = currentTime
		}

		deltaTime := (currentTime - b.lastFrameTime) / 1000.0 // Convert to seconds
		b.lastFrameTime = currentTime

		// Call the game callback
		callback(deltaTime)

		// Request next frame
		b.window.Call("requestAnimationFrame", b.animationCallback)
		return nil
	})

	// Start the loop
	b.window.Call("requestAnimationFrame", b.animationCallback)
}

// GetInput returns the current input state
func (b *JSBridge) GetInput() (left, right, fire, fireJustPressed, pauseJustPressed bool) {
	// Movement keys
	left = b.keysPressed["ArrowLeft"] || b.keysPressed["KeyA"]
	right = b.keysPressed["ArrowRight"] || b.keysPressed["KeyD"]

	// Fire key
	fire = b.keysPressed[" "] || b.keysPressed["Space"]

	// For "just pressed" detection, we need to track previous state
	// This is a simplified version - a more robust system would track frame-to-frame changes
	fireJustPressed = fire
	pauseJustPressed = b.keysPressed["Escape"] || b.keysPressed["KeyP"]

	return
}

// IsKeyPressed checks if a specific key is currently pressed
func (b *JSBridge) IsKeyPressed(key string) bool {
	return b.keysPressed[key]
}

// ClearCanvas clears the entire canvas with the specified color
func (b *JSBridge) ClearCanvas(color string) {
	b.context.Set("fillStyle", color)
	b.context.Call("fillRect", 0, 0, b.canvasWidth/int(b.deviceRatio), b.canvasHeight/int(b.deviceRatio))
}

// Drawing functions

// DrawRect draws a filled rectangle
func (b *JSBridge) DrawRect(x, y, width, height float64, color string) {
	b.context.Set("fillStyle", color)
	b.context.Call("fillRect", x, y, width, height)
}

// DrawRectOutline draws a rectangle outline
func (b *JSBridge) DrawRectOutline(x, y, width, height, lineWidth float64, color string) {
	b.context.Set("strokeStyle", color)
	b.context.Set("lineWidth", lineWidth)
	b.context.Call("strokeRect", x, y, width, height)
}

// DrawCircle draws a filled circle
func (b *JSBridge) DrawCircle(x, y, radius float64, color string) {
	b.context.Set("fillStyle", color)
	b.context.Call("beginPath")
	b.context.Call("arc", x, y, radius, 0, 2*3.14159)
	b.context.Call("fill")
}

// DrawText draws text at the specified position
func (b *JSBridge) DrawText(text string, x, y float64, font, color string) {
	b.context.Set("font", font)
	b.context.Set("fillStyle", color)
	b.context.Call("fillText", text, x, y)
}

// DrawTextCentered draws centered text
func (b *JSBridge) DrawTextCentered(text string, x, y float64, font, color string) {
	b.context.Set("font", font)
	b.context.Set("fillStyle", color)
	b.context.Set("textAlign", "center")
	b.context.Call("fillText", text, x, y)
	b.context.Set("textAlign", "left") // Reset
}

// DrawLine draws a line between two points
func (b *JSBridge) DrawLine(x1, y1, x2, y2, lineWidth float64, color string) {
	b.context.Set("strokeStyle", color)
	b.context.Set("lineWidth", lineWidth)
	b.context.Call("beginPath")
	b.context.Call("moveTo", x1, y1)
	b.context.Call("lineTo", x2, y2)
	b.context.Call("stroke")
}

// Utility functions

// GetCanvasSize returns the canvas dimensions
func (b *JSBridge) GetCanvasSize() (int, int) {
	return b.canvasWidth / int(b.deviceRatio), b.canvasHeight / int(b.deviceRatio)
}

// GetDevicePixelRatio returns the device pixel ratio
func (b *JSBridge) GetDevicePixelRatio() float64 {
	return b.deviceRatio
}

// SetCanvasSize sets the canvas size explicitly
func (b *JSBridge) SetCanvasSize(width, height int) {
	b.canvasWidth = int(float64(width) * b.deviceRatio)
	b.canvasHeight = int(float64(height) * b.deviceRatio)

	b.canvas.Set("width", b.canvasWidth)
	b.canvas.Set("height", b.canvasHeight)

	b.canvas.Get("style").Set("width", width)
	b.canvas.Get("style").Set("height", height)

	b.context.Call("scale", b.deviceRatio, b.deviceRatio)
}

// Log logs a message to the browser console
func (b *JSBridge) Log(message string) {
	b.window.Get("console").Call("log", message)
}

// LogError logs an error to the browser console
func (b *JSBridge) LogError(message string) {
	b.window.Get("console").Call("error", message)
}

// ShowAlert shows an alert dialog
func (b *JSBridge) ShowAlert(message string) {
	b.window.Call("alert", message)
}

// GetCurrentTime returns the current time in milliseconds
func (b *JSBridge) GetCurrentTime() float64 {
	return b.window.Get("performance").Call("now").Float()
}

// SetTitle sets the document title
func (b *JSBridge) SetTitle(title string) {
	b.document.Set("title", title)
}

// DOM manipulation

// GetElementByID gets an element by its ID
func (b *JSBridge) GetElementByID(id string) js.Value {
	return b.document.Call("getElementById", id)
}

// SetElementText sets the text content of an element
func (b *JSBridge) SetElementText(elementID, text string) {
	element := b.GetElementByID(elementID)
	if !element.IsUndefined() {
		element.Set("textContent", text)
	}
}

// SetElementHTML sets the HTML content of an element
func (b *JSBridge) SetElementHTML(elementID, html string) {
	element := b.GetElementByID(elementID)
	if !element.IsUndefined() {
		element.Set("innerHTML", html)
	}
}

// Audio support (placeholder for future implementation)

// PlaySound plays a sound effect (to be implemented with Web Audio API)
func (b *JSBridge) PlaySound(soundID string) {
	// Placeholder - would implement Web Audio API calls here
	b.Log("Playing sound: " + soundID)
}

// Storage support

// SetLocalStorage sets a value in localStorage
func (b *JSBridge) SetLocalStorage(key, value string) {
	storage := b.window.Get("localStorage")
	if !storage.IsUndefined() {
		storage.Call("setItem", key, value)
	}
}

// GetLocalStorage gets a value from localStorage
func (b *JSBridge) GetLocalStorage(key string) string {
	storage := b.window.Get("localStorage")
	if !storage.IsUndefined() {
		item := storage.Call("getItem", key)
		if !item.IsUndefined() && item.Type() == js.TypeString {
			return item.String()
		}
	}
	return ""
}

// Cleanup releases all resources
func (b *JSBridge) Cleanup() {
	// Remove event listeners
	if !b.keydownListener.IsUndefined() {
		b.document.Call("removeEventListener", "keydown", b.keydownListener)
		b.keydownListener.Release()
	}
	if !b.keyupListener.IsUndefined() {
		b.document.Call("removeEventListener", "keyup", b.keyupListener)
		b.keyupListener.Release()
	}
	if !b.resizeListener.IsUndefined() {
		b.window.Call("removeEventListener", "resize", b.resizeListener)
		b.resizeListener.Release()
	}
	if !b.focusListener.IsUndefined() {
		b.window.Call("removeEventListener", "focus", b.focusListener)
		b.focusListener.Release()
	}
	if !b.blurListener.IsUndefined() {
		b.window.Call("removeEventListener", "blur", b.blurListener)
		b.blurListener.Release()
	}
	if !b.animationCallback.IsUndefined() {
		b.animationCallback.Release()
	}

	// Clear key state
	b.keysPressed = make(map[string]bool)
}

// Performance monitoring

// GetFPS calculates and returns the current FPS
func (b *JSBridge) GetFPS(deltaTime float64) float64 {
	if deltaTime > 0 {
		return 1.0 / deltaTime
	}
	return 0
}

// UpdateInputState updates the input state with proper "just pressed" detection
func (b *JSBridge) UpdateInputState(state *InputState) {
	// For now, simple implementation without tracking previous state
	// TODO: Implement proper "just pressed" detection with frame tracking
	state.LeftPressed = b.keysPressed["ArrowLeft"] || b.keysPressed["KeyA"]
	state.RightPressed = b.keysPressed["ArrowRight"] || b.keysPressed["KeyD"]
	state.UpPressed = b.keysPressed["ArrowUp"] || b.keysPressed["KeyW"]
	state.DownPressed = b.keysPressed["ArrowDown"] || b.keysPressed["KeyS"]
	state.FirePressed = b.keysPressed[" "] || b.keysPressed["Space"]
	state.EnterJustPressed = b.keysPressed["Enter"]
	state.PauseJustPressed = b.keysPressed["KeyP"] || b.keysPressed["Escape"]
}

// GetTime returns the current time (for compatibility with time.Time)
func GetTime() time.Time {
	// This is a simplified version - in a real implementation,
	// you might want to sync with JavaScript's Date.now()
	return time.Now()
}