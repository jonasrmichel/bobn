# Space Invaders - Camera Edition

A modern take on the classic Space Invaders game, featuring camera-based controls built with Go and WebAssembly.

## Project Structure

```
.
├── cmd/
│   ├── server/          # HTTP server for serving the web application
│   └── wasm/            # WebAssembly game implementation
├── internal/
│   ├── game/            # Core game logic (shared between server and WASM)
│   └── wasm/            # WASM-specific utilities
├── web/                 # Web assets (HTML, CSS, JS, WASM)
├── assets/              # Game assets (sprites, sounds, etc.)
├── scripts/             # Build and deployment scripts
└── Makefile             # Build automation
```

## Quick Start

### Prerequisites
- Go 1.19+ installed
- Modern web browser with WASM support

### Build and Run

1. **Build everything:**
   ```bash
   make all
   ```

2. **Run the development server:**
   ```bash
   make run-server
   ```

3. **Open your browser:**
   Navigate to `http://localhost:8080`

### Development Commands

```bash
# Build server only
make server

# Build WASM only
make wasm

# Build web assets (WASM + wasm_exec.js)
make web

# Run development server with auto-rebuild
make dev

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Show all available commands
make help
```

## Game Features

### Current Status
- [x] Basic project structure
- [x] WASM build pipeline
- [x] HTTP server for web delivery
- [x] Basic game loop and canvas rendering
- [ ] Camera input integration
- [ ] Player movement detection
- [ ] Enemy AI and movement
- [ ] Collision detection
- [ ] Scoring system
- [ ] Sound effects

### Planned Features
- **Camera Controls**: Use device camera to detect player movement
- **Classic Gameplay**: Faithful recreation of Space Invaders mechanics
- **Modern Enhancements**: Improved graphics and sound
- **Responsive Design**: Works on desktop and mobile devices

## Architecture

### Server (cmd/server/)
- HTTP server serving static web assets
- Health check endpoint
- Graceful shutdown handling

### WASM Game (cmd/wasm/)
- Game loop implementation
- Canvas rendering
- Input handling
- Game state management

### Internal Packages
- `internal/game/`: Core game logic shared between components
- `internal/wasm/`: WASM-specific utilities and helpers

## Controls

- **SPACEBAR**: Start/Stop game
- **ESCAPE**: Stop game
- **Camera controls**: Coming soon!

## Development

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality
```bash
# Format code
make fmt

# Vet code
make vet

# Run linter (requires golangci-lint)
make lint
```

### Installing Development Tools
```bash
make install-tools
```

## Browser Support

- Chrome/Chromium 57+
- Firefox 52+
- Safari 11+
- Edge 16+

## Contributing

1. Ensure Go 1.19+ is installed
2. Run `make deps` to download dependencies
3. Run `make fmt` before committing
4. Ensure all tests pass with `make test`

## License

This project is open source and available under the [MIT License](LICENSE).