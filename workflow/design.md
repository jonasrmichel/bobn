# Camera-Based Space Invaders Web Game Design

## Overall Goal
Build a retro arcade-style Space Invaders game that uses camera-based head tracking for ship control, implemented in Go with WebAssembly for client-side processing and a simple HTTP server for hosting.

## Requirements
1. **Single-player gameplay** with head tracking controls
2. **WASM-based client** for low-latency camera processing
3. **Retro arcade theme** inspired by classic Space Invaders and Asteroids
4. **2D head movement** controls both X and Y ship position
5. **Modern browsers only** (Chrome, Firefox, Safari latest versions)
6. **Casual gaming experience** with 10s of concurrent players capability

## Architecture Overview

### Technology Stack
- **Language**: Go (both server and WASM client)
- **Client**: Go compiled to WebAssembly for camera processing and game logic
- **Server**: Simple HTTP server to serve static files
- **Graphics**: HTML5 Canvas with retro CRT effects
- **Audio**: Web Audio API with generated arcade sounds
- **Camera**: getUserMedia API with motion detection

### Project Structure
```
bobn/
├── cmd/
│   ├── server/
│   │   └── main.go           # HTTP server for WASM/static files
│   └── wasm/
│       └── main.go           # WASM game entry point
├── internal/
│   ├── game/
│   │   ├── engine.go         # Core game loop
│   │   ├── invaders.go       # Space Invaders mechanics
│   │   ├── asteroids.go      # Asteroids bonus rounds
│   │   ├── collision.go      # Physics detection
│   │   └── state.go          # Game state management
│   └── wasm/
│       ├── camera.go         # Head tracking logic
│       ├── oscilloscope.go   # Camera visualization
│       ├── renderer.go       # Canvas rendering
│       ├── audio.go          # Sound generation
│       └── storage.go        # LocalStorage for high scores
├── web/
│   ├── index.html            # Arcade cabinet UI
│   ├── arcade.css            # Retro CRT effects
│   ├── main.wasm             # Compiled game
│   └── wasm_exec.js          # Go WASM support
└── assets/
    └── sprites.go            # Embedded pixel art
```

## Key Design Decisions

### 1. Camera-Based Control System
- **2D Head Tracking**: Tracks head position in X/Y axes
- **Motion Detection**: Simple brightness-based center of mass tracking
- **Smoothing**: 70/30 blend to reduce jitter
- **Dead Zones**: 5% threshold to ignore small movements
- **Calibration**: Quick 5-point calibration (center, corners)

### 2. Ship Movement Physics
- **Horizontal Movement**: Full screen width (90% range)
- **Vertical Movement**: Limited to bottom 40% of screen (Space Invaders style)
- **Acceleration-based**: Smooth acceleration/deceleration for natural feel
- **Different Speeds**: Horizontal faster than vertical movement
- **Visual Feedback**: Ship tilts based on movement direction

### 3. Visual Design - Retro Arcade
- **CRT Effects**: Scanlines, phosphor glow, screen curvature
- **Color Palette**: Classic phosphor green, amber, cyan, magenta
- **Typography**: "Press Start 2P" pixel font
- **Arcade Cabinet**: Full cabinet UI with marquee, bezel, control panel
- **Oscilloscope**: Live waveform visualization of head tracking

### 4. Game Modes
- **Space Invaders Mode**: Classic gameplay with limited vertical movement
- **Asteroids Bonus Rounds**: Full 2D movement, vector graphics style
- **TILT System**: Warnings when tracking is lost, auto-pause

### 5. Performance Targets
- **Camera Processing**: 30 FPS capture, 10 FPS processing
- **Game Logic**: 20 Hz server tick rate
- **Rendering**: 60 FPS smooth animation
- **Latency**: <70ms perceived delay
- **WASM Size**: Target <1MB compiled

## Design Session Transcript Summary

### Initial Requirements Discussion
- User specified: 10s of players, simple head tracking (left-right movement), casual web game, modern browsers only, single server deployment
- Decided on WASM for client-side processing to reduce latency

### Architecture Evolution
1. Started with multi-session server design for multiplayer
2. Simplified to single-player after user clarification
3. Removed WebSocket complexity in favor of pure client-side game
4. Added oscilloscope visualization for camera status per user request

### Visual Theme Discussion
- User requested retro arcade theme
- Designed full arcade cabinet aesthetic with CRT effects
- Created phosphor color palette and scanline effects
- Added classic arcade UI elements (INSERT COIN, HIGH SCORES)

### Control System Evolution
1. Initially designed for X-axis only (left-right)
2. User requested Y-axis movement as well
3. Designed 2D tracking with different movement constraints per game mode
4. Added gesture detection options (nod to fire, auto-fire zones)

### Key Technical Decisions
- **No WebRTC**: Too complex for casual game, WebSocket unnecessary for single-player
- **No P2P**: All processing client-side in WASM
- **Simple Motion Detection**: Brightness-based tracking sufficient for casual play
- **LocalStorage**: For high scores and calibration data
- **No Backend Database**: Keeping it simple for MVP

## Patterns & Best Practices

### WASM Integration
- Use `js.Value` for DOM manipulation
- Minimize JS boundary crossings
- Process camera frames in Go for consistency
- Use Web Workers if needed for heavy processing

### Game Loop
- Fixed timestep for physics (50ms)
- Variable rendering with interpolation
- Separate update and render cycles
- Request Animation Frame for smooth visuals

### Camera Processing
- Downsample to 320x240 for performance
- Process every 3rd frame (10 FPS effective)
- Use moving average for smoothing
- Implement fallback to keyboard controls

### Error Handling
- Graceful degradation when camera unavailable
- Clear user messages for calibration
- Auto-recovery from tracking loss
- TILT warnings for edge cases

## Technical Debt & Future Enhancements

### MVP Scope (Current Implementation)
- Single-player only
- Basic motion detection
- Simple collision detection
- Local high scores only
- Basic sound effects

### Future Enhancements (Post-MVP)
- Multiplayer support with WebSocket
- Advanced gesture recognition
- Power-ups and special weapons
- Global leaderboards
- Mobile touch controls
- Progressive difficulty
- More bonus round types
- Facial expression detection for actions

## Success Criteria
1. Game runs at 60 FPS in modern browsers
2. Head tracking latency <100ms
3. No server required after initial load
4. Works with built-in webcams
5. Playable without calibration (reasonable defaults)
6. Retro aesthetic feels authentic
7. Game is fun and engaging for 5+ minutes