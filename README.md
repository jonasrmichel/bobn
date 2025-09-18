#
```
    ██████╗  ██████╗ ██████╗ ███╗   ██╗
    ██╔══██╗██╔═══██╗██╔══██╗████╗  ██║
    ██████╔╝██║   ██║██████╔╝██╔██╗ ██║
    ██╔══██╗██║   ██║██╔══██╗██║╚██╗██║
    ██████╔╝╚██████╔╝██████╔╝██║ ╚████║
    ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝
     S P A C E   D E F E N D E R   ■ ▲ ●
```

```
       .  *  .   ▲    .  *   .    *    .
    *    .   /   \   .    *   .   *   .
  .   *    /  O O  \    .   *   ■   .   *
    .    /___________\  *   .  / \  .
  *   .    |  ■■■  |    .   * |___| *  .
    .  *   |_______|  .    *    .   *
  .    *    .    ●    *   .    ▲    .  *
```

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![WASM](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=webassembly&logoColor=white)](https://webassembly.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

**A retro arcade space shooter controlled by your head movements!**

[Play Now](#quick-start) • [Features](#features) • [Controls](#controls) • [Development](#development)

</div>

---

## 🚀 About BOBN

```
    ┌─────────────────────────────────┐
    │  ▲                              │
    │ /_\  INCOMING INVADERS!         │
    │  |                              │
    │      ■ ■ ■ ■ ■ ■ ■ ■ ■ ■        │
    │      ■ ■ ■ ■ ■ ■ ■ ■ ■ ■        │
    │      ■ ■ ■ ■ ■ ■ ■ ■ ■ ■        │
    │                                 │
    │            ▲                    │
    │         ──═▼═──                 │
    └─────────────────────────────────┘
```

BOBN is a browser-based Space Invaders game written in Go and compiled to WebAssembly. What makes it special? You control your ship using head movements tracked by your webcam!

### 🎮 Key Features

```
     ∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙
    ∙  ┌───┐  FEATURES  ┌───┐  ┌───┐   ∙
    ∙  │ ▲ │            │ ■ │  │ ● │   ∙
    ∙  └───┘            └───┘  └───┘   ∙
     ∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙
```

- 🎥 **Camera-Based Control** - Move your head left/right to control your ship
- 🕹️ **Classic Arcade Gameplay** - Authentic Space Invaders experience
- 🌐 **Pure Browser-Based** - No downloads or plugins required
- ⚡ **Go + WebAssembly** - Fast, efficient performance
- 🎨 **Retro CRT Effects** - Scanlines, phosphor glow, and arcade aesthetics
- 📊 **Live Head Tracking Display** - ASCII art visualization of camera feed
- 🏆 **High Score Tracking** - Compete for the top score

---

## 🎯 Quick Start

```
    ╔═══════════════════════════════╗
    ║  > INITIALIZING SYSTEM...     ║
    ║  > LOADING WASM MODULE...     ║
    ║  > READY PLAYER ONE           ║
    ╚═══════════════════════════════╝
```

### Prerequisites

- Go 1.21 or higher
- Modern web browser with WebAssembly support
- Webcam (for head tracking control)

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/jonasrmichel/bobn.git
cd bobn

# Build and run
make run

# Open in browser
# Navigate to http://localhost:8080
```

---

## 🎮 How to Play

```
      CONTROLS           MOVEMENT
    ┌──────────┐       ┌─────────┐
    │ [SPACE]  │       │  ←   →  │
    │  FIRE    │       │  HEAD   │
    └──────────┘       └─────────┘

    ┌──────────┐       ┌─────────┐
    │ [ENTER]  │       │  [ESC]  │
    │  START   │       │  PAUSE  │
    └──────────┘       └─────────┘
```

### Game Modes

1. **Attract Mode** - Press ENTER or click START to begin
2. **Playing** - Destroy all invaders before they reach you
3. **Game Over** - Your final score is displayed

### Scoring

```
    ╔══════════════════════╗
    ║  ENEMY POINT VALUES  ║
    ╠══════════════════════╣
    ║  ■ = 10 POINTS       ║
    ║  ▲ = 20 POINTS       ║
    ║  ● = 30 POINTS       ║
    ║  ◆ = ???  MYSTERY    ║
    ╚══════════════════════╝
```

---

## 🛠️ Development

```
     ∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙
    ∙    < CODE STRUCTURE >     ∙
    ∙      /    |    \          ∙
    ∙    cmd  internal  web     ∙
     ∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙∙
```

### Project Structure

```
bobn/
├── cmd/
│   ├── server/          # HTTP server
│   └── wasm/            # WASM entry point
├── internal/
│   ├── game/            # Game engine & logic
│   │   ├── engine.go    # Core game loop
│   │   ├── entities.go  # Game objects
│   │   └── state.go     # Game state management
│   └── wasm/            # WASM-specific code
│       ├── bridge.go    # JS interop
│       ├── camera.go    # Head tracking
│       └── renderer.go  # Canvas rendering
├── web/
│   ├── index.html       # Game UI
│   ├── arcade.css       # Retro styling
│   └── wasm_exec.js     # Go WASM support
└── Makefile            # Build automation
```

### Build Commands

```bash
# Build everything
make all

# Build WASM only
make wasm

# Build server only
make server

# Run development server
make run

# Clean build artifacts
make clean
```

### Development Tips

```
    ┌─────────────────────────────┐
    │  💡 PRO TIPS:               │
    │                             │
    │  • Use Chrome DevTools      │
    │  • Check console for logs   │
    │  • Test with good lighting  │
    │  • Adjust camera threshold  │
    └─────────────────────────────┘
```

---

## 🎨 Customization

### Adjusting Head Tracking Sensitivity

Edit `internal/wasm/camera.go`:

```go
// Increase/decrease sensitivity
gameX := -((c.smoothedX - 0.5) * 4.0) // Change multiplier

// Adjust detection threshold
if brightness > 80 { // Lower = more sensitive
```

### Modifying Game Difficulty

Edit `internal/game/engine.go`:

```go
// Invader speed
baseInvaderSpeed: 1.0  // Increase for harder

// Invader drop distance
invaderDropDistance: 20.0  // Pixels per drop
```

---

## 🚁 Architecture

```
    ┌──────────────────────────────────┐
    │          BROWSER                 │
    │  ┌────────────────────────────┐  │
    │  │      HTML/CSS UI          │  │
    │  ├────────────────────────────┤  │
    │  │    WASM Game Engine       │  │
    │  │  ┌──────────┬───────────┐ │  │
    │  │  │ Camera   │  Game     │ │  │
    │  │  │ Tracking │  Logic    │ │  │
    │  │  └──────────┴───────────┘ │  │
    │  ├────────────────────────────┤  │
    │  │    Canvas Rendering       │  │
    │  └────────────────────────────┘  │
    └──────────────────────────────────┘
```

### Key Components

- **Game Engine**: Fixed timestep loop at 20Hz
- **Rendering**: 60 FPS canvas updates
- **Camera Processing**: 30 FPS head tracking
- **Input System**: Keyboard and camera hybrid control

---

## 🎮 Game Features

```
     WAVES              ENEMIES
    ┌─────┐           ┌─────────┐
    │  1  │ ───────> │ ■ ■ ■ ■ │
    │  2  │ ───────> │ ■ ■ ■ ■ │
    │  3  │ ───────> │ ▲ ▲ ▲ ▲ │
    │  4+ │ ───────> │ ● ● ● ● │
    └─────┘           └─────────┘
```

- Progressive difficulty increase
- Wave-based gameplay
- UFO bonus targets
- Life system
- Score tracking

---

## 🐛 Troubleshooting

```
    ╔════════════════════════════════╗
    ║     COMMON ISSUES & FIXES      ║
    ╚════════════════════════════════╝
```

| Issue | Solution |
|-------|----------|
| Camera not working | Allow camera permissions in browser |
| Black game screen | Refresh page, check console for errors |
| Laggy controls | Ensure good lighting for camera |
| Ship not responding | Press ENTER to start game first |

---

## 📜 License

```
    ┌─────────────────────────┐
    │   MIT License - 2024    │
    │   Free to use & modify  │
    └─────────────────────────┘
```

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

```
      ▲
     /_\    CONTRIBUTORS WELCOME!
    [___]
     | |    Fork → Code → PR
     |_|
```

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 🎯 Roadmap

```
    [ ] Multiplayer support
    [ ] Power-ups & bonuses
    [ ] More enemy types
    [ ] Sound effects
    [ ] Mobile touch controls
    [ ] Leaderboard system
```

---

## 👾 Credits

```
     ╔══════════════════════════╗
     ║  CREATED BY              ║
     ║  Jonas Michel            ║
     ║                          ║
     ║  INSPIRED BY             ║
     ║  Space Invaders (1978)   ║
     ║  By Tomohiro Nishikado   ║
     ╚══════════════════════════╝
```

---

<div align="center">

```
    GAME OVER

    ▲ ■ ● ▲ ■ ● ▲ ■ ● ▲ ■ ●

    INSERT COIN TO CONTINUE

    [ PRESS START ]
```

**Made with ❤️ and Go**

[⬆ Back to Top](#)

</div>