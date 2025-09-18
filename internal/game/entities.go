package game

import (
	"math"
	"time"
)

// Vector2 represents a 2D vector for position and velocity
type Vector2 struct {
	X, Y float64
}

// Add returns the sum of two vectors
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Scale returns a scaled vector
func (v Vector2) Scale(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Magnitude returns the magnitude of the vector
func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Bounds represents a rectangular boundary
type Bounds struct {
	X, Y, Width, Height float64
}

// Contains checks if a point is within the bounds
func (b Bounds) Contains(x, y float64) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}

// Intersects checks if two bounds intersect (AABB collision)
func (b Bounds) Intersects(other Bounds) bool {
	return b.X < other.X+other.Width &&
		b.X+b.Width > other.X &&
		b.Y < other.Y+other.Height &&
		b.Y+b.Height > other.Y
}

// PlayerShip represents the player's ship
type PlayerShip struct {
	Position     Vector2
	Velocity     Vector2
	Bounds       Bounds
	Alive        bool
	MaxSpeed     float64
	Acceleration float64
	Friction     float64

	// Animation state
	AnimFrame    int
	AnimTimer    float64

	// Shooting state
	CanShoot     bool
	LastShotTime time.Time
	FireRate     float64 // shots per second
}

// NewPlayerShip creates a new player ship at the specified position
func NewPlayerShip(x, y float64) *PlayerShip {
	const shipWidth = 24
	const shipHeight = 16

	return &PlayerShip{
		Position:     Vector2{X: x, Y: y},
		Velocity:     Vector2{X: 0, Y: 0},
		Bounds:       Bounds{X: x - shipWidth/2, Y: y - shipHeight/2, Width: shipWidth, Height: shipHeight},
		Alive:        true,
		MaxSpeed:     200.0, // pixels per second
		Acceleration: 800.0, // pixels per second squared
		Friction:     400.0, // pixels per second squared
		CanShoot:     true,
		FireRate:     4.0, // 4 shots per second
		LastShotTime: time.Now(),
	}
}

// Update updates the player ship's position and state
func (p *PlayerShip) Update(deltaTime float64, screenWidth float64) {
	if !p.Alive {
		return
	}

	// Update position based on velocity
	p.Position = p.Position.Add(p.Velocity.Scale(deltaTime))

	// Update bounds
	p.Bounds.X = p.Position.X - p.Bounds.Width/2
	p.Bounds.Y = p.Position.Y - p.Bounds.Height/2

	// Keep player within screen bounds
	halfWidth := p.Bounds.Width / 2
	if p.Position.X < halfWidth {
		p.Position.X = halfWidth
		p.Velocity.X = 0
	} else if p.Position.X > screenWidth-halfWidth {
		p.Position.X = screenWidth - halfWidth
		p.Velocity.X = 0
	}

	// Update shooting cooldown
	if !p.CanShoot && time.Since(p.LastShotTime).Seconds() > 1.0/p.FireRate {
		p.CanShoot = true
	}

	// Update animation
	p.AnimTimer += deltaTime
	if p.AnimTimer > 0.1 { // 10 FPS animation
		p.AnimFrame = (p.AnimFrame + 1) % 2
		p.AnimTimer = 0
	}
}

// ApplyInput applies input forces to the player ship
func (p *PlayerShip) ApplyInput(left, right bool, deltaTime float64) {
	if !p.Alive {
		return
	}

	// Apply acceleration based on input
	if left && !right {
		p.Velocity.X -= p.Acceleration * deltaTime
	} else if right && !left {
		p.Velocity.X += p.Acceleration * deltaTime
	} else {
		// Apply friction when no input
		if p.Velocity.X > 0 {
			p.Velocity.X -= p.Friction * deltaTime
			if p.Velocity.X < 0 {
				p.Velocity.X = 0
			}
		} else if p.Velocity.X < 0 {
			p.Velocity.X += p.Friction * deltaTime
			if p.Velocity.X > 0 {
				p.Velocity.X = 0
			}
		}
	}

	// Clamp velocity to max speed
	if p.Velocity.X > p.MaxSpeed {
		p.Velocity.X = p.MaxSpeed
	} else if p.Velocity.X < -p.MaxSpeed {
		p.Velocity.X = -p.MaxSpeed
	}
}

// TryShoot attempts to create a bullet if shooting is allowed
func (p *PlayerShip) TryShoot() *Bullet {
	if !p.Alive || !p.CanShoot {
		return nil
	}

	p.CanShoot = false
	p.LastShotTime = time.Now()

	// Create bullet at player position, moving upward
	return NewBullet(p.Position.X, p.Position.Y-p.Bounds.Height/2, 0, -400, true)
}

// InvaderType represents different types of invaders
type InvaderType int

const (
	InvaderTypeSmall InvaderType = iota
	InvaderTypeMedium
	InvaderTypeLarge
)

// Invader represents an enemy invader
type Invader struct {
	Type      InvaderType
	Position  Vector2
	Bounds    Bounds
	Alive     bool
	Points    int
	Direction int // -1 for left, 1 for right

	// Animation state
	AnimFrame int
	AnimTimer float64

	// Shooting state (for advanced invaders)
	CanShoot     bool
	LastShotTime time.Time
	ShootChance  float64 // probability per second
}

// NewInvader creates a new invader
func NewInvader(invaderType InvaderType, x, y float64, points int) *Invader {
	var width, height float64
	var shootChance float64

	switch invaderType {
	case InvaderTypeSmall:
		width, height = 16, 16
		shootChance = 0.03 // 3% chance per second (reduced from 10%)
	case InvaderTypeMedium:
		width, height = 20, 16
		shootChance = 0.02 // 2% chance per second (reduced from 5%)
	case InvaderTypeLarge:
		width, height = 24, 16
		shootChance = 0.01 // 1% chance per second (reduced from 2%)
	}

	return &Invader{
		Type:         invaderType,
		Position:     Vector2{X: x, Y: y},
		Bounds:       Bounds{X: x - width/2, Y: y - height/2, Width: width, Height: height},
		Alive:        true,
		Points:       points,
		Direction:    1, // Initially moving right
		CanShoot:     true,
		ShootChance:  shootChance,
		LastShotTime: time.Now(),
	}
}

// Update updates the invader's animation state
func (i *Invader) Update(deltaTime float64) {
	if !i.Alive {
		return
	}

	// Update animation
	i.AnimTimer += deltaTime
	if i.AnimTimer > 0.5 { // 2 FPS animation
		i.AnimFrame = (i.AnimFrame + 1) % 2
		i.AnimTimer = 0
	}
}

// Move moves the invader by the specified offset
func (i *Invader) Move(deltaX, deltaY float64) {
	if !i.Alive {
		return
	}

	i.Position.X += deltaX
	i.Position.Y += deltaY

	// Update bounds
	i.Bounds.X = i.Position.X - i.Bounds.Width/2
	i.Bounds.Y = i.Position.Y - i.Bounds.Height/2
}

// TryShoot attempts to create a bullet if shooting conditions are met
func (i *Invader) TryShoot(deltaTime float64) *Bullet {
	if !i.Alive || !i.CanShoot {
		return nil
	}

	// Random shooting based on shoot chance
	shootProbability := i.ShootChance * deltaTime
	if math.Mod(float64(time.Now().UnixNano()/1000), 1.0) < shootProbability {
		i.LastShotTime = time.Now()
		// Create bullet moving downward
		return NewBullet(i.Position.X, i.Position.Y+i.Bounds.Height/2, 0, 200, false)
	}

	return nil
}

// Bullet represents a projectile
type Bullet struct {
	Position       Vector2
	Velocity       Vector2
	Bounds         Bounds
	Alive          bool
	IsPlayerBullet bool
	Damage         int
}

// NewBullet creates a new bullet
func NewBullet(x, y, velX, velY float64, isPlayerBullet bool) *Bullet {
	const bulletWidth = 2
	const bulletHeight = 8

	return &Bullet{
		Position:       Vector2{X: x, Y: y},
		Velocity:       Vector2{X: velX, Y: velY},
		Bounds:         Bounds{X: x - bulletWidth/2, Y: y - bulletHeight/2, Width: bulletWidth, Height: bulletHeight},
		Alive:          true,
		IsPlayerBullet: isPlayerBullet,
		Damage:         1,
	}
}

// Update updates the bullet's position
func (b *Bullet) Update(deltaTime float64, screenWidth, screenHeight float64) {
	if !b.Alive {
		return
	}

	// Update position
	b.Position = b.Position.Add(b.Velocity.Scale(deltaTime))

	// Update bounds
	b.Bounds.X = b.Position.X - b.Bounds.Width/2
	b.Bounds.Y = b.Position.Y - b.Bounds.Height/2

	// Remove bullets that go off screen
	if b.Position.Y < 0 || b.Position.Y > screenHeight ||
		b.Position.X < 0 || b.Position.X > screenWidth {
		b.Alive = false
	}
}

// UFO represents the bonus enemy UFO
type UFO struct {
	Position  Vector2
	Velocity  Vector2
	Bounds    Bounds
	Alive     bool
	Points    int
	Direction int // -1 for left, 1 for right

	// State tracking
	SpawnTime    time.Time
	MaxLifetime  time.Duration
}

// NewUFO creates a new UFO
func NewUFO(startX, y float64, direction int) *UFO {
	const ufoWidth = 32
	const ufoHeight = 16
	const ufoSpeed = 100.0 // pixels per second

	velocity := Vector2{X: ufoSpeed * float64(direction), Y: 0}
	points := []int{100, 150, 200, 300}[int(time.Now().UnixNano()/1000000)%4] // Random point value

	return &UFO{
		Position:    Vector2{X: startX, Y: y},
		Velocity:    velocity,
		Bounds:      Bounds{X: startX - ufoWidth/2, Y: y - ufoHeight/2, Width: ufoWidth, Height: ufoHeight},
		Alive:       true,
		Points:      points,
		Direction:   direction,
		SpawnTime:   time.Now(),
		MaxLifetime: 15 * time.Second, // UFO disappears after 15 seconds
	}
}

// Update updates the UFO's position and state
func (u *UFO) Update(deltaTime float64, screenWidth float64) {
	if !u.Alive {
		return
	}

	// Update position
	u.Position = u.Position.Add(u.Velocity.Scale(deltaTime))

	// Update bounds
	u.Bounds.X = u.Position.X - u.Bounds.Width/2
	u.Bounds.Y = u.Position.Y - u.Bounds.Height/2

	// Remove UFO if it goes off screen or exceeds lifetime
	if u.Position.X < -u.Bounds.Width || u.Position.X > screenWidth+u.Bounds.Width ||
		time.Since(u.SpawnTime) > u.MaxLifetime {
		u.Alive = false
	}
}

// ShouldSpawnUFO determines if a UFO should be spawned based on game state
func ShouldSpawnUFO(lastUFOTime time.Time, gameTime float64) bool {
	// Spawn UFO every 20-40 seconds randomly
	minInterval := 20.0
	maxInterval := 40.0

	timeSinceLastUFO := time.Since(lastUFOTime).Seconds()
	spawnThreshold := minInterval + (maxInterval-minInterval)*math.Mod(gameTime*0.123, 1.0)

	return timeSinceLastUFO > spawnThreshold
}