package wasm

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/jonasrmichel/bobn/internal/game"
)

// Renderer handles all game rendering to the canvas
type Renderer struct {
	bridge      *JSBridge
	ctx         js.Value
	pixelSize   int
	screenWidth int
	screenHeight int
}

// NewRenderer creates a new renderer
func NewRenderer(bridge *JSBridge, screenWidth, screenHeight int) *Renderer {
	return &Renderer{
		bridge:       bridge,
		ctx:          bridge.GetContext(),
		pixelSize:    2,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// SetContext sets the rendering context
func (r *Renderer) SetContext(ctx js.Value) {
	r.ctx = ctx
}

// Clear clears the canvas
func (r *Renderer) Clear() {
	r.ctx.Call("clearRect", 0, 0, r.screenWidth, r.screenHeight)

	// Draw starfield background
	r.ctx.Set("fillStyle", "#000000")
	r.ctx.Call("fillRect", 0, 0, r.screenWidth, r.screenHeight)

	// Draw stars
	r.drawStarfield()
}

// drawStarfield draws a simple starfield background
func (r *Renderer) drawStarfield() {
	r.ctx.Set("fillStyle", "#ffffff")
	// Static stars for now
	stars := [][]int{
		{100, 50}, {200, 80}, {300, 120}, {400, 30}, {500, 90},
		{150, 200}, {250, 180}, {350, 220}, {450, 150}, {550, 210},
		{120, 350}, {220, 380}, {320, 320}, {420, 330}, {520, 390},
	}
	for _, star := range stars {
		r.ctx.Set("globalAlpha", 0.5)
		r.ctx.Call("fillRect", star[0], star[1], 1, 1)
	}
	r.ctx.Set("globalAlpha", 1.0)
}

// RenderGame renders the entire game state
func (r *Renderer) RenderGame(state *game.GameState) {
	r.Clear()

	switch state.Mode {
	case game.AttractMode:
		r.renderAttractMode(state)
	case game.Playing:
		r.renderPlayingMode(state)
	case game.GameOver:
		r.renderGameOverMode(state)
	case game.HighScore:
		r.renderHighScoreMode(state)
	}

	// Always render UI elements
	r.renderUI(state)
}

// renderAttractMode renders the attract mode screen
func (r *Renderer) renderAttractMode(state *game.GameState) {
	// Title
	r.drawText("BOBN", r.screenWidth/2, 150, 48, "#00ff00", "center")
	r.drawText("SPACE INVADERS", r.screenWidth/2, 200, 24, "#00ffff", "center")

	// Instructions
	r.drawText("USE ARROW KEYS TO MOVE", r.screenWidth/2, 300, 16, "#ffff00", "center")
	r.drawText("PRESS SPACE TO FIRE", r.screenWidth/2, 330, 16, "#ffff00", "center")

	// Blinking insert coin
	if int(js.Global().Get("Date").New().Call("getTime").Float()/500)%2 == 0 {
		r.drawText("PRESS ENTER TO START", r.screenWidth/2, 400, 20, "#ff00ff", "center")
	}

	// High score
	r.drawText(fmt.Sprintf("HIGH SCORE: %06d", state.HighScore), r.screenWidth/2, 450, 16, "#ffffff", "center")
}

// renderPlayingMode renders the main game
func (r *Renderer) renderPlayingMode(state *game.GameState) {
	// Render player
	if state.Player != nil {
		r.renderPlayer(state.Player)
	}

	// Render invaders
	for _, invader := range state.Invaders {
		r.renderInvader(invader)
	}

	// Render bullets
	for _, bullet := range state.Bullets {
		r.renderBullet(bullet)
	}

	// Render UFO
	if state.UFO != nil && state.UFO.Alive {
		r.renderUFO(state.UFO)
	}

	// Barriers not implemented yet - TODO: Add barriers later
}

// renderGameOverMode renders the game over screen
func (r *Renderer) renderGameOverMode(state *game.GameState) {
	r.drawText("GAME OVER", r.screenWidth/2, r.screenHeight/2-50, 48, "#ff0000", "center")
	r.drawText(fmt.Sprintf("FINAL SCORE: %06d", state.Score), r.screenWidth/2, r.screenHeight/2+20, 24, "#ffffff", "center")

	if state.Score > state.HighScore {
		r.drawText("NEW HIGH SCORE!", r.screenWidth/2, r.screenHeight/2+60, 20, "#ffff00", "center")
	}

	if int(js.Global().Get("Date").New().Call("getTime").Float()/500)%2 == 0 {
		r.drawText("PRESS ENTER TO CONTINUE", r.screenWidth/2, r.screenHeight/2+120, 16, "#00ff00", "center")
	}
}

// renderHighScoreMode renders the high score entry screen
func (r *Renderer) renderHighScoreMode(state *game.GameState) {
	r.drawText("NEW HIGH SCORE!", r.screenWidth/2, r.screenHeight/2-50, 36, "#ffff00", "center")
	r.drawText(fmt.Sprintf("SCORE: %06d", state.Score), r.screenWidth/2, r.screenHeight/2, 24, "#ffffff", "center")
	r.drawText("PRESS ENTER TO CONTINUE", r.screenWidth/2, r.screenHeight/2+80, 16, "#00ff00", "center")
}

// renderUI renders the UI elements (score, lives, etc.)
func (r *Renderer) renderUI(state *game.GameState) {
	// Score
	r.drawText(fmt.Sprintf("SCORE: %06d", state.Score), 10, 30, 16, "#ffffff", "left")

	// High Score
	r.drawText(fmt.Sprintf("HIGH: %06d", state.HighScore), r.screenWidth/2, 30, 16, "#ffff00", "center")

	// Lives
	r.drawText("LIVES:", r.screenWidth-150, 30, 16, "#ffffff", "left")
	for i := 0; i < state.Lives; i++ {
		r.renderMiniShip(r.screenWidth-90+i*25, 25)
	}

	// Wave
	if state.Mode == game.Playing {
		r.drawText(fmt.Sprintf("WAVE %d", state.Wave), r.screenWidth/2, r.screenHeight-20, 16, "#00ffff", "center")
	}
}

// renderPlayer renders the player ship
func (r *Renderer) renderPlayer(player *game.PlayerShip) {
	if !player.Alive {
		return
	}

	// Draw ship body (triangle shape)
	r.ctx.Set("fillStyle", "#00ff00")
	r.ctx.Call("beginPath")
	r.ctx.Call("moveTo", player.Position.X, player.Position.Y)
	r.ctx.Call("lineTo", player.Position.X-15, player.Position.Y+20)
	r.ctx.Call("lineTo", player.Position.X+15, player.Position.Y+20)
	r.ctx.Call("closePath")
	r.ctx.Call("fill")

	// Draw cockpit
	r.ctx.Set("fillStyle", "#00ffff")
	r.ctx.Call("beginPath")
	r.ctx.Call("arc", player.Position.X, player.Position.Y+5, 4, 0, math.Pi*2)
	r.ctx.Call("fill")
}

// renderMiniShip renders a small ship for lives display
func (r *Renderer) renderMiniShip(x, y int) {
	r.ctx.Set("fillStyle", "#00ff00")
	r.ctx.Call("beginPath")
	r.ctx.Call("moveTo", x, y)
	r.ctx.Call("lineTo", x-8, y+10)
	r.ctx.Call("lineTo", x+8, y+10)
	r.ctx.Call("closePath")
	r.ctx.Call("fill")
}

// renderInvader renders an invader
func (r *Renderer) renderInvader(invader *game.Invader) {
	if !invader.Alive {
		return
	}

	color := "#ffffff"
	switch invader.Type {
	case game.InvaderTypeSmall:
		color = "#ff00ff"
	case game.InvaderTypeMedium:
		color = "#ffff00"
	case game.InvaderTypeLarge:
		color = "#00ffff"
	}

	// Simple invader shape
	r.ctx.Set("fillStyle", color)

	// Body
	r.ctx.Call("fillRect", invader.Position.X-10, invader.Position.Y-5, 20, 10)

	// Arms (animate)
	armOffset := 0
	if invader.AnimFrame > 0 {
		armOffset = 3
	}
	r.ctx.Call("fillRect", invader.Position.X-15, invader.Position.Y, 5, 5+armOffset)
	r.ctx.Call("fillRect", invader.Position.X+10, invader.Position.Y, 5, 5+armOffset)

	// Eyes
	r.ctx.Set("fillStyle", "#000000")
	r.ctx.Call("fillRect", invader.Position.X-6, invader.Position.Y-2, 3, 3)
	r.ctx.Call("fillRect", invader.Position.X+3, invader.Position.Y-2, 3, 3)
}

// renderBullet renders a bullet
func (r *Renderer) renderBullet(bullet *game.Bullet) {
	if !bullet.Alive {
		return
	}

	if bullet.IsPlayerBullet {
		// Player bullet - vertical line
		r.ctx.Set("strokeStyle", "#00ff00")
		r.ctx.Set("lineWidth", 2)
		r.ctx.Call("beginPath")
		r.ctx.Call("moveTo", bullet.Position.X, bullet.Position.Y)
		r.ctx.Call("lineTo", bullet.Position.X, bullet.Position.Y+8)
		r.ctx.Call("stroke")
	} else {
		// Enemy bullet - zigzag
		r.ctx.Set("strokeStyle", "#ff0000")
		r.ctx.Set("lineWidth", 2)
		r.ctx.Call("beginPath")
		r.ctx.Call("moveTo", bullet.Position.X-2, bullet.Position.Y)
		r.ctx.Call("lineTo", bullet.Position.X+2, bullet.Position.Y+3)
		r.ctx.Call("lineTo", bullet.Position.X-2, bullet.Position.Y+6)
		r.ctx.Call("lineTo", bullet.Position.X+2, bullet.Position.Y+9)
		r.ctx.Call("stroke")
	}
}

// renderUFO renders the UFO
func (r *Renderer) renderUFO(ufo *game.UFO) {
	if !ufo.Alive {
		return
	}

	// UFO body
	r.ctx.Set("fillStyle", "#ff00ff")
	r.ctx.Call("beginPath")
	r.ctx.Call("ellipse", ufo.Position.X, ufo.Position.Y, 20, 8, 0, 0, math.Pi*2)
	r.ctx.Call("fill")

	// Dome
	r.ctx.Set("fillStyle", "#ffff00")
	r.ctx.Call("beginPath")
	r.ctx.Call("arc", ufo.Position.X, ufo.Position.Y-5, 8, math.Pi, 0)
	r.ctx.Call("fill")

	// Lights
	r.ctx.Set("fillStyle", "#ffffff")
	for i := -15; i <= 15; i += 10 {
		if int(js.Global().Get("Date").New().Call("getTime").Float()/200)%2 == 0 {
			r.ctx.Call("beginPath")
			r.ctx.Call("arc", ufo.Position.X+float64(i), ufo.Position.Y, 2, 0, math.Pi*2)
			r.ctx.Call("fill")
		}
	}
}


// drawText renders text to the canvas
func (r *Renderer) drawText(text string, x, y int, size int, color, align string) {
	r.ctx.Set("font", fmt.Sprintf("%dpx 'Press Start 2P', monospace", size))
	r.ctx.Set("fillStyle", color)
	r.ctx.Set("textAlign", align)
	r.ctx.Set("textBaseline", "middle")
	r.ctx.Call("fillText", text, x, y)
}

// RenderExplosion renders an explosion effect
func (r *Renderer) RenderExplosion(x, y float64, frame int) {
	if frame >= 10 {
		return
	}

	// Expanding circle of particles
	r.ctx.Set("fillStyle", "#ff0000")
	radius := float64(frame * 3)
	alpha := 1.0 - float64(frame)/10.0
	r.ctx.Set("globalAlpha", alpha)

	// Draw particles
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		px := x + math.Cos(angle)*radius
		py := y + math.Sin(angle)*radius
		r.ctx.Call("beginPath")
		r.ctx.Call("arc", px, py, 3, 0, math.Pi*2)
		r.ctx.Call("fill")
	}

	r.ctx.Set("globalAlpha", 1.0)
}