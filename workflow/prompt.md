# Implementation Checklist: Camera-Based Space Invaders

## Phase 1: Project Setup & Core Structure

### Go Module & Dependencies
- [ ] Initialize Go module: `go mod init github.com/jonasrmichel/bobn`
- [ ] Add minimal dependencies (no external game libraries for MVP)
- [ ] Set up build scripts for both server and WASM

### Directory Structure
- [ ] Create `cmd/server/main.go` - HTTP server entry point
- [ ] Create `cmd/wasm/main.go` - WASM game entry point
- [ ] Create `internal/game/` directory for game logic
- [ ] Create `internal/wasm/` directory for browser-specific code
- [ ] Create `web/` directory for static assets
- [ ] Create `assets/` directory for embedded sprites

## Phase 2: Basic HTTP Server

### Server Implementation (`cmd/server/main.go`)
- [ ] Create simple HTTP server on port 8080
- [ ] Serve static files from `web/` directory
- [ ] Add CORS headers for WASM
- [ ] Set correct MIME type for .wasm files
- [ ] Add basic logging

### Build Configuration
- [ ] Create `Makefile` or build script
- [ ] Add WASM build target: `GOOS=js GOARCH=wasm go build`
- [ ] Copy `wasm_exec.js` from Go distribution
- [ ] Add server build target

## Phase 3: WASM Foundation

### WASM Entry Point (`cmd/wasm/main.go`)
- [ ] Set up main function with channel to prevent exit
- [ ] Initialize JavaScript bindings
- [ ] Set up panic handler for debugging
- [ ] Create game instance

### JavaScript Bridge (`internal/wasm/bridge.go`)
- [ ] Create DOM manipulation helpers
- [ ] Canvas context acquisition
- [ ] RequestAnimationFrame binding
- [ ] Event listener setup (keyboard fallback)

## Phase 4: HTML/CSS Retro UI

### HTML Structure (`web/index.html`)
- [ ] Create arcade cabinet layout
- [ ] Add canvas element for game
- [ ] Add oscilloscope canvas for camera viz
- [ ] Insert coin / start buttons
- [ ] Score and lives display
- [ ] High scores panel

### CSS Styling (`web/arcade.css`)
- [ ] Implement CRT effects (scanlines, glow)
- [ ] Create phosphor color variables
- [ ] Arcade cabinet styling
- [ ] Neon text effects
- [ ] Responsive layout for mobile

## Phase 5: Game Core Logic

### Game State (`internal/game/state.go`)
- [ ] Define GameState struct
- [ ] Player ship position and lives
- [ ] Invader grid and positions
- [ ] Bullets array
- [ ] Score and wave number
- [ ] Game mode enum (menu, playing, game over)

### Game Engine (`internal/game/engine.go`)
- [ ] Create game loop with fixed timestep
- [ ] Update function (physics, collisions)
- [ ] Input processing
- [ ] State transitions

### Entities (`internal/game/entities.go`)
- [ ] PlayerShip struct with movement physics
- [ ] Invader struct with movement patterns
- [ ] Bullet struct with velocity
- [ ] UFO bonus enemy

### Collision Detection (`internal/game/collision.go`)
- [ ] AABB collision for bullets vs invaders
- [ ] Player vs enemy bullets
- [ ] Invader reached bottom detection

## Phase 6: Camera Integration

### Camera Module (`internal/wasm/camera.go`)
- [ ] Request camera permission
- [ ] Set up video element
- [ ] Capture frames to canvas
- [ ] Handle permission denied

### Head Tracking (`internal/wasm/tracker.go`)
- [ ] Implement simple motion detection
- [ ] Calculate brightness center of mass
- [ ] Smooth position values
- [ ] Normalize to -1 to 1 range

### Calibration (`internal/wasm/calibration.go`)
- [ ] 5-point calibration UI
- [ ] Store calibration data in LocalStorage
- [ ] Auto-recalibrate on tracking loss

### Oscilloscope Viz (`internal/wasm/oscilloscope.go`)
- [ ] Draw grid background
- [ ] Plot head position waveform
- [ ] Show tracking status
- [ ] Add phosphor decay effect

## Phase 7: Rendering System

### Canvas Renderer (`internal/wasm/renderer.go`)
- [ ] Get 2D context from canvas
- [ ] Clear and draw game frame
- [ ] Sprite rendering functions
- [ ] Text rendering for score/UI

### Sprite System (`assets/sprites.go`)
- [ ] Define pixel art as byte arrays
- [ ] Player ship sprite (13x8 pixels)
- [ ] Invader sprites (3 types, 11x8 pixels)
- [ ] Explosion animation frames
- [ ] Embed sprites in binary

### Visual Effects
- [ ] Ship tilt based on movement
- [ ] Muzzle flash on shooting
- [ ] Explosion particles
- [ ] Screen shake on player hit

## Phase 8: Game Mechanics

### Space Invaders Logic (`internal/game/invaders.go`)
- [ ] Invader grid formation
- [ ] Left-right-down movement pattern
- [ ] Speed increases as numbers decrease
- [ ] Random shooting from bottom row
- [ ] Point values per invader type

### Player Controls (`internal/game/controls.go`)
- [ ] Map head X position to ship X
- [ ] Map head Y position to ship Y (limited range)
- [ ] Auto-fire based on position
- [ ] Keyboard fallback (arrow keys + space)

### Scoring System
- [ ] Points per invader type
- [ ] UFO bonus points
- [ ] Extra life at 10,000 points
- [ ] High score tracking in LocalStorage

## Phase 9: Audio System

### Sound Generator (`internal/wasm/audio.go`)
- [ ] Create Web Audio context
- [ ] Generate retro sound effects
- [ ] Laser sound (descending square wave)
- [ ] Explosion (white noise burst)
- [ ] Invader movement bass notes
- [ ] UFO warble sound

### Music System
- [ ] Four-note bass loop
- [ ] Increase tempo as invaders descend
- [ ] Victory/defeat jingles

## Phase 10: Polish & Game Feel

### TILT System
- [ ] Detect tracking loss
- [ ] Show TILT warning
- [ ] Pause game and prompt recalibration
- [ ] Auto-resume when tracking restored

### Bonus Round (`internal/game/asteroids.go`)
- [ ] Every 3 waves switch to Asteroids mode
- [ ] Vector-style graphics
- [ ] Full 2D movement
- [ ] Different scoring rules

### Attract Mode
- [ ] Demo play when idle
- [ ] Cycle through high scores
- [ ] Show instructions
- [ ] INSERT COIN blinking

## Phase 11: Testing & Optimization

### Performance
- [ ] Profile WASM performance
- [ ] Optimize camera processing (skip frames)
- [ ] Reduce allocations in game loop
- [ ] Batch canvas operations

### Browser Testing
- [ ] Test in Chrome, Firefox, Safari
- [ ] Verify camera permissions flow
- [ ] Check WASM loading
- [ ] Test keyboard fallback

### Build Optimization
- [ ] Minify WASM with wasm-opt
- [ ] Gzip compression
- [ ] Cache headers for static assets

## Phase 12: Documentation

### Code Documentation
- [ ] Add package comments
- [ ] Document public APIs
- [ ] TODO comments for tech debt
- [ ] README with build instructions

### User Documentation
- [ ] How to play instructions
- [ ] Camera setup guide
- [ ] Troubleshooting section
- [ ] Browser requirements

## Technical Debt (TODO Comments)
- [ ] TODO: (tech debt) Add WebSocket for multiplayer
- [ ] TODO: (tech debt) Implement power-ups system
- [ ] TODO: (tech debt) Add particle effects system
- [ ] TODO: (tech debt) Global leaderboard integration
- [ ] TODO: (tech debt) Mobile touch controls
- [ ] TODO: (tech debt) Advanced face recognition
- [ ] TODO: (tech debt) Gamepad support

## Success Validation
- [ ] Game loads in <3 seconds
- [ ] Runs at 60 FPS consistently
- [ ] Camera tracking feels responsive
- [ ] Retro aesthetic is convincing
- [ ] Game is fun for 5+ minutes
- [ ] No server needed after load
- [ ] Works with built-in webcams