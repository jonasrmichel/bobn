package wasm

import (
	"log"
	"math"
	"syscall/js"
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
	currentFrame  []uint8  // Store current frame for ASCII generation
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

	// Store frame data for ASCII generation
	frameSize := c.width * c.height * 4
	if len(c.currentFrame) != frameSize {
		c.currentFrame = make([]uint8, frameSize)
	}

	// Copy frame data
	for i := 0; i < frameSize; i++ {
		c.currentFrame[i] = uint8(data.Index(i).Int())
	}

	// Simple brightness-based motion detection
	var sumX, sumY, totalBrightness float64
	pixelCount := 0

	// Sample every 8th pixel for better performance
	for y := 0; y < c.height; y += 8 {
		for x := 0; x < c.width; x += 8 {
			idx := (y*c.width + x) * 4

			// Get pixel brightness
			r := float64(c.currentFrame[idx])
			g := float64(c.currentFrame[idx+1])
			b := float64(c.currentFrame[idx+2])
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

// updateOscilloscope updates the oscilloscope visualization with ASCII art
func (c *CameraController) updateOscilloscope(x, y float64) {
	if !c.oscilloscope.Truthy() {
		return
	}

	ctx := c.oscilloscope.Call("getContext", "2d")
	width := c.oscilloscope.Get("width").Float()
	height := c.oscilloscope.Get("height").Float()

	// Clear canvas
	ctx.Set("fillStyle", "#000000")
	ctx.Call("fillRect", 0, 0, width, height)

	// Set font for ASCII art
	ctx.Set("font", "10px monospace")
	ctx.Set("fillStyle", "#00ff00")
	ctx.Set("textAlign", "left")
	ctx.Set("textBaseline", "top")

	if c.tracking {
		// Generate ASCII art representation of camera view
		asciiArt := c.generateASCIIArt()

		// Draw each line of ASCII art
		lineHeight := 12.0
		startY := 10.0

		for i, line := range asciiArt {
			ctx.Call("fillText", line, 10, startY+float64(i)*lineHeight)
		}

		// Draw position indicator at bottom
		ctx.Set("fillStyle", "#ffff00")
		posText := "POS: ["
		for i := 0; i < 10; i++ {
			if float64(i-5)/5.0 < x && x <= float64(i-4)/5.0 {
				posText += "■"
			} else {
				posText += "·"
			}
		}
		posText += "]"
		ctx.Call("fillText", posText, 10, height-20)
	} else {
		// Show "NO SIGNAL" when camera not active
		ctx.Set("font", "16px monospace")
		ctx.Set("fillStyle", "#00ff00")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", "NO SIGNAL", width/2, height/2)
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

// generateASCIIArt generates ASCII art representation of the camera view
func (c *CameraController) generateASCIIArt() []string {
	if len(c.currentFrame) == 0 {
		return []string{"NO DATA"}
	}

	// ASCII characters for different brightness levels
	asciiChars := " .:-=+*#%@"

	// Create smaller ASCII grid (sample the image)
	asciiWidth := 30
	asciiHeight := 12

	lines := make([]string, asciiHeight)

	// Calculate sampling intervals
	xStep := c.width / asciiWidth
	yStep := c.height / asciiHeight

	for ay := 0; ay < asciiHeight; ay++ {
		line := ""
		for ax := 0; ax < asciiWidth; ax++ {
			// Sample a region of pixels
			x := ax * xStep
			y := ay * yStep

			// Average brightness in the region
			var totalBrightness float64
			sampleCount := 0

			// Sample a few pixels in the region
			for dy := 0; dy < yStep && y+dy < c.height; dy += 2 {
				for dx := 0; dx < xStep && x+dx < c.width; dx += 2 {
					idx := ((y+dy)*c.width + (x+dx)) * 4
					if idx < len(c.currentFrame)-3 {
						r := float64(c.currentFrame[idx])
						g := float64(c.currentFrame[idx+1])
						b := float64(c.currentFrame[idx+2])
						brightness := (r + g + b) / 3.0

						// Edge detection: look for significant brightness changes
						isEdge := false
						if x+dx+2 < c.width {
							nextIdx := ((y+dy)*c.width + (x+dx+2)) * 4
							if nextIdx < len(c.currentFrame)-3 {
								nextR := float64(c.currentFrame[nextIdx])
								nextG := float64(c.currentFrame[nextIdx+1])
								nextB := float64(c.currentFrame[nextIdx+2])
								nextBrightness := (nextR + nextG + nextB) / 3.0
								if math.Abs(brightness-nextBrightness) > 50 {
									isEdge = true
								}
							}
						}

						if isEdge {
							brightness = 255 // Highlight edges
						}

						totalBrightness += brightness
						sampleCount++
					}
				}
			}

			// Map brightness to ASCII character
			if sampleCount > 0 {
				avgBrightness := totalBrightness / float64(sampleCount)
				charIndex := int(avgBrightness * float64(len(asciiChars)-1) / 255.0)
				if charIndex >= len(asciiChars) {
					charIndex = len(asciiChars) - 1
				}
				line += string(asciiChars[charIndex])
			} else {
				line += " "
			}
		}
		lines[ay] = line
	}

	return lines
}

// StartCalibration starts the calibration process
func (c *CameraController) StartCalibration() {
	// Simple auto-calibration based on current position
	c.centerX = c.smoothedX
	c.centerY = c.smoothedY
	c.calibrated = true
	log.Println("Camera calibrated")
}