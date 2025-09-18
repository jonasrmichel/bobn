package wasm

import (
	"log"
	"math"
	"syscall/js"
	"time"
)

// CameraController handles camera input for head tracking
type CameraController struct {
	video         js.Value
	canvas        js.Value
	ctx           js.Value
	enabled       bool
	tracking      bool

	// Calibration
	centerX       float64
	rangeX        float64
	centerY       float64
	rangeY        float64
	calibrated    bool

	// Current position
	currentX      float64
	currentY      float64
	smoothedX     float64
	smoothedY     float64

	// Motion detection
	prevFrame     []uint8
	width         int
	height        int

	// Callbacks
	onPosition    func(x, y float64)
	oscilloscope  js.Value
}

// NewCameraController creates a new camera controller
func NewCameraController() *CameraController {
	return &CameraController{
		width:  320,
		height: 240,
		centerX: 0.5,
		rangeX: 0.3,
		centerY: 0.5,
		rangeY: 0.2,
	}
}

// Initialize sets up the camera
func (c *CameraController) Initialize() error {
	doc := js.Global().Get("document")

	// Create hidden video element
	c.video = doc.Call("createElement", "video")
	c.video.Set("width", c.width)
	c.video.Set("height", c.height)
	c.video.Set("autoplay", true)
	c.video.Set("style", "display: none")
	doc.Get("body").Call("appendChild", c.video)

	// Create processing canvas
	c.canvas = doc.Call("createElement", "canvas")
	c.canvas.Set("width", c.width)
	c.canvas.Set("height", c.height)
	c.canvas.Set("style", "display: none")
	doc.Get("body").Call("appendChild", c.canvas)
	c.ctx = c.canvas.Call("getContext", "2d")

	// Get oscilloscope canvas for visualization
	c.oscilloscope = doc.Call("getElementById", "oscilloscope")

	// Request camera access
	navigator := js.Global().Get("navigator")
	mediaDevices := navigator.Get("mediaDevices")

	if !mediaDevices.Truthy() {
		log.Println("MediaDevices API not supported")
		return nil
	}

	// Set up constraints
	constraints := map[string]interface{}{
		"video": map[string]interface{}{
			"width":  c.width,
			"height": c.height,
		},
		"audio": false,
	}

	// Get user media
	promise := mediaDevices.Call("getUserMedia", constraints)

	// Handle promise
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		stream := args[0]
		c.video.Set("srcObject", stream)
		c.enabled = true
		c.tracking = true
		log.Println("Camera initialized successfully")

		// Start processing loop
		c.startProcessing()
		return nil
	}))

	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("Failed to get camera access: %v", args[0])
		c.enabled = false
		return nil
	}))

	return nil
}

// startProcessing starts the frame processing loop
func (c *CameraController) startProcessing() {
	if !c.enabled {
		return
	}

	// Process frames at 30 FPS
	js.Global().Call("setInterval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if c.tracking && c.enabled {
			c.processFrame()
		}
		return nil
	}), 33) // ~30 FPS
}

// processFrame processes a single camera frame
func (c *CameraController) processFrame() {
	// Draw video frame to canvas
	c.ctx.Call("drawImage", c.video, 0, 0, c.width, c.height)

	// Get image data
	imageData := c.ctx.Call("getImageData", 0, 0, c.width, c.height)
	data := imageData.Get("data")

	// Simple brightness-based motion detection
	var sumX, sumY, totalBrightness float64
	pixelCount := 0

	// Sample every 8th pixel for better performance
	for y := 0; y < c.height; y += 8 {
		for x := 0; x < c.width; x += 8 {
			idx := (y*c.width + x) * 4

			// Get pixel brightness
			r := data.Index(idx).Float()
			g := data.Index(idx + 1).Float()
			b := data.Index(idx + 2).Float()
			brightness := (r + g + b) / 3.0

			// Only count bright pixels (likely face/head)
			if brightness > 80 { // Lower threshold for better detection
				sumX += float64(x) * brightness
				sumY += float64(y) * brightness
				totalBrightness += brightness
				pixelCount++
			}
		}
	}

	if totalBrightness > 0 {
		// Calculate center of mass
		centerX := sumX / totalBrightness / float64(c.width)
		centerY := sumY / totalBrightness / float64(c.height)

		// Less smoothing for more responsive control
		c.smoothedX = c.smoothedX*0.3 + centerX*0.7
		c.smoothedY = c.smoothedY*0.3 + centerY*0.7

		// Convert to game coordinates (-1 to 1)
		// Invert X because camera is mirrored
		gameX := -((c.smoothedX - 0.5) * 4.0) // Increased sensitivity
		gameY := (c.smoothedY - 0.5) * 4.0

		// Smaller dead zone for more responsive control
		if math.Abs(gameX) < 0.05 {
			gameX = 0
		}
		if math.Abs(gameY) < 0.05 {
			gameY = 0
		}

		// Store current position
		c.currentX = gameX
		c.currentY = gameY

		// Update oscilloscope visualization
		c.updateOscilloscope(gameX, gameY)

		// Call position callback if set
		if c.onPosition != nil {
			c.onPosition(gameX, gameY)
		}
	}
}

// updateOscilloscope updates the oscilloscope visualization
func (c *CameraController) updateOscilloscope(x, y float64) {
	if !c.oscilloscope.Truthy() {
		return
	}

	ctx := c.oscilloscope.Call("getContext", "2d")
	width := c.oscilloscope.Get("width").Float()
	height := c.oscilloscope.Get("height").Float()

	// Clear canvas
	ctx.Call("clearRect", 0, 0, width, height)

	// Draw grid
	ctx.Set("strokeStyle", "#003300")
	ctx.Set("lineWidth", 1)

	// Vertical lines
	for i := 0.0; i < width; i += width / 10 {
		ctx.Call("beginPath")
		ctx.Call("moveTo", i, 0)
		ctx.Call("lineTo", i, height)
		ctx.Call("stroke")
	}

	// Horizontal lines
	for i := 0.0; i < height; i += height / 6 {
		ctx.Call("beginPath")
		ctx.Call("moveTo", 0, i)
		ctx.Call("lineTo", width, i)
		ctx.Call("stroke")
	}

	if c.tracking {
		// Draw X position wave
		ctx.Set("strokeStyle", "#00ff00")
		ctx.Set("lineWidth", 2)
		ctx.Call("beginPath")

		centerY := height / 4
		for i := 0.0; i < width; i++ {
			waveY := centerY + x*20*math.Sin(i*0.1+float64(time.Now().UnixMilli())*0.001)
			if i == 0 {
				ctx.Call("moveTo", i, waveY)
			} else {
				ctx.Call("lineTo", i, waveY)
			}
		}
		ctx.Call("stroke")

		// Draw Y position wave
		ctx.Set("strokeStyle", "#00ffff")
		ctx.Call("beginPath")

		centerY = height * 3 / 4
		for i := 0.0; i < width; i++ {
			waveY := centerY + y*20*math.Sin(i*0.1+float64(time.Now().UnixMilli())*0.0015)
			if i == 0 {
				ctx.Call("moveTo", i, waveY)
			} else {
				ctx.Call("lineTo", i, waveY)
			}
		}
		ctx.Call("stroke")

		// Draw position indicator
		ctx.Set("fillStyle", "#ffff00")
		ctx.Call("beginPath")
		ctx.Call("arc", width/2+x*width/4, height/2+y*height/4, 3, 0, math.Pi*2)
		ctx.Call("fill")
	} else {
		// Draw flat line when no tracking
		ctx.Set("strokeStyle", "#00ff00")
		ctx.Set("lineWidth", 1)
		ctx.Call("beginPath")
		ctx.Call("moveTo", 0, height/2)
		ctx.Call("lineTo", width, height/2)
		ctx.Call("stroke")
	}
}

// GetPosition returns the current head position
func (c *CameraController) GetPosition() (float64, float64) {
	if !c.enabled || !c.tracking {
		return 0, 0
	}
	return c.currentX, c.currentY
}

// SetPositionCallback sets the callback for position updates
func (c *CameraController) SetPositionCallback(callback func(x, y float64)) {
	c.onPosition = callback
}

// IsEnabled returns whether camera is enabled
func (c *CameraController) IsEnabled() bool {
	return c.enabled
}

// StartCalibration starts the calibration process
func (c *CameraController) StartCalibration() {
	// Simple auto-calibration based on current position
	c.centerX = c.smoothedX
	c.centerY = c.smoothedY
	c.calibrated = true
	log.Println("Camera calibrated")
}