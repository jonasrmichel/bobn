package game

import (
	"math"
)

// CollisionSystem handles all collision detection and response
type CollisionSystem struct {
	// Spatial partitioning could be added here for optimization
	// For now, we'll use simple brute force collision detection
}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{}
}

// CheckAABBCollision performs Axis-Aligned Bounding Box collision detection
func CheckAABBCollision(a, b Bounds) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// CheckPointInBounds checks if a point is within bounds
func CheckPointInBounds(x, y float64, bounds Bounds) bool {
	return x >= bounds.X && x < bounds.X+bounds.Width &&
		y >= bounds.Y && y < bounds.Y+bounds.Height
}

// CheckCircleCollision performs circle-to-circle collision detection
func CheckCircleCollision(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < (r1 + r2)
}

// CheckCircleRectCollision checks collision between a circle and rectangle
func CheckCircleRectCollision(circleX, circleY, radius float64, rect Bounds) bool {
	// Find the closest point on the rectangle to the circle center
	closestX := math.Max(rect.X, math.Min(circleX, rect.X+rect.Width))
	closestY := math.Max(rect.Y, math.Min(circleY, rect.Y+rect.Height))

	// Calculate the distance between the circle center and the closest point
	dx := circleX - closestX
	dy := circleY - closestY

	// Check if the distance is less than the circle radius
	return (dx*dx + dy*dy) < (radius * radius)
}

// CollisionResult represents the result of a collision check
type CollisionResult struct {
	Collided      bool
	PenetrationX  float64
	PenetrationY  float64
	ContactPointX float64
	ContactPointY float64
}

// CheckAABBCollisionWithDetails performs AABB collision with detailed results
func CheckAABBCollisionWithDetails(a, b Bounds) CollisionResult {
	result := CollisionResult{}

	// Check if collision occurs
	if !CheckAABBCollision(a, b) {
		return result
	}

	result.Collided = true

	// Calculate overlap amounts
	overlapX1 := (a.X + a.Width) - b.X
	overlapX2 := (b.X + b.Width) - a.X
	overlapY1 := (a.Y + a.Height) - b.Y
	overlapY2 := (b.Y + b.Height) - a.Y

	// Choose minimum overlap for separation
	if overlapX1 < overlapX2 {
		result.PenetrationX = overlapX1
	} else {
		result.PenetrationX = -overlapX2
	}

	if overlapY1 < overlapY2 {
		result.PenetrationY = overlapY1
	} else {
		result.PenetrationY = -overlapY2
	}

	// Calculate contact point (center of overlapping area)
	overlapLeft := math.Max(a.X, b.X)
	overlapRight := math.Min(a.X+a.Width, b.X+b.Width)
	overlapTop := math.Max(a.Y, b.Y)
	overlapBottom := math.Min(a.Y+a.Height, b.Y+b.Height)

	result.ContactPointX = (overlapLeft + overlapRight) / 2
	result.ContactPointY = (overlapTop + overlapBottom) / 2

	return result
}

// CheckBulletInvaderCollision checks collision between a bullet and invader with pixel-perfect detection
func CheckBulletInvaderCollision(bullet *Bullet, invader *Invader) bool {
	if !bullet.Alive || !invader.Alive {
		return false
	}

	// First do AABB check for early rejection
	if !CheckAABBCollision(bullet.Bounds, invader.Bounds) {
		return false
	}

	// For Space Invaders, AABB is usually sufficient since bullets are small
	// and invaders are relatively large. More detailed collision detection
	// could be added here if needed (e.g., pixel-perfect collision).
	return true
}

// CheckBulletPlayerCollision checks collision between a bullet and player
func CheckBulletPlayerCollision(bullet *Bullet, player *PlayerShip) bool {
	if !bullet.Alive || !player.Alive || bullet.IsPlayerBullet {
		return false
	}

	return CheckAABBCollision(bullet.Bounds, player.Bounds)
}

// CheckBulletUFOCollision checks collision between a bullet and UFO
func CheckBulletUFOCollision(bullet *Bullet, ufo *UFO) bool {
	if !bullet.Alive || !ufo.Alive || !bullet.IsPlayerBullet {
		return false
	}

	return CheckAABBCollision(bullet.Bounds, ufo.Bounds)
}

// CheckPlayerInvaderCollision checks direct collision between player and invader
func CheckPlayerInvaderCollision(player *PlayerShip, invader *Invader) bool {
	if !player.Alive || !invader.Alive {
		return false
	}

	return CheckAABBCollision(player.Bounds, invader.Bounds)
}

// CheckBoundaryCollision checks if an entity is within screen boundaries
func CheckBoundaryCollision(bounds Bounds, screenWidth, screenHeight float64) (left, right, top, bottom bool) {
	left = bounds.X < 0
	right = bounds.X+bounds.Width > screenWidth
	top = bounds.Y < 0
	bottom = bounds.Y+bounds.Height > screenHeight
	return
}

// KeepInBounds constrains bounds within screen boundaries
func KeepInBounds(bounds *Bounds, screenWidth, screenHeight float64) {
	if bounds.X < 0 {
		bounds.X = 0
	}
	if bounds.X+bounds.Width > screenWidth {
		bounds.X = screenWidth - bounds.Width
	}
	if bounds.Y < 0 {
		bounds.Y = 0
	}
	if bounds.Y+bounds.Height > screenHeight {
		bounds.Y = screenHeight - bounds.Height
	}
}

// LineIntersection checks if two line segments intersect
func LineIntersection(x1, y1, x2, y2, x3, y3, x4, y4 float64) (bool, float64, float64) {
	// Calculate the denominator
	denom := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(denom) < 1e-10 {
		return false, 0, 0 // Lines are parallel
	}

	// Calculate intersection parameters
	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / denom
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / denom

	// Check if intersection occurs within both line segments
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		// Calculate intersection point
		intersectX := x1 + t*(x2-x1)
		intersectY := y1 + t*(y2-y1)
		return true, intersectX, intersectY
	}

	return false, 0, 0
}

// CheckBulletBarrierCollision checks collision between bullet and barrier
func CheckBulletBarrierCollision(bullet *Bullet, barriers [][]bool, barrierBlockSize float64) (bool, int, int) {
	if !bullet.Alive || len(barriers) == 0 {
		return false, -1, -1
	}

	// Calculate which barrier blocks the bullet overlaps
	bulletLeft := int(bullet.Bounds.X / barrierBlockSize)
	bulletRight := int((bullet.Bounds.X + bullet.Bounds.Width) / barrierBlockSize)
	bulletTop := int((bullet.Bounds.Y - float64(len(barriers[0]))*barrierBlockSize) / barrierBlockSize)
	bulletBottom := int((bullet.Bounds.Y + bullet.Bounds.Height - float64(len(barriers[0]))*barrierBlockSize) / barrierBlockSize)

	// Clamp to barrier array bounds
	if bulletLeft < 0 {
		bulletLeft = 0
	}
	if bulletRight >= len(barriers) {
		bulletRight = len(barriers) - 1
	}
	if bulletTop < 0 {
		bulletTop = 0
	}
	if bulletBottom >= len(barriers[0]) {
		bulletBottom = len(barriers[0]) - 1
	}

	// Check for collision with barrier blocks
	for x := bulletLeft; x <= bulletRight; x++ {
		for y := bulletTop; y <= bulletBottom; y++ {
			if x >= 0 && x < len(barriers) && y >= 0 && y < len(barriers[0]) && barriers[x][y] {
				return true, x, y
			}
		}
	}

	return false, -1, -1
}

// DestroyBarrierBlock destroys a barrier block and surrounding blocks for impact effect
func DestroyBarrierBlock(barriers [][]bool, x, y int, radius int) {
	if x < 0 || x >= len(barriers) || y < 0 || y >= len(barriers[0]) {
		return
	}

	// Destroy blocks in a circular pattern
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			newX := x + dx
			newY := y + dy

			if newX >= 0 && newX < len(barriers) && newY >= 0 && newY < len(barriers[0]) {
				// Use circular destruction pattern
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance <= float64(radius) {
					barriers[newX][newY] = false
				}
			}
		}
	}
}

// GetClosestPoint returns the closest point on a rectangle to a given point
func GetClosestPoint(pointX, pointY float64, rect Bounds) (float64, float64) {
	closestX := math.Max(rect.X, math.Min(pointX, rect.X+rect.Width))
	closestY := math.Max(rect.Y, math.Min(pointY, rect.Y+rect.Height))
	return closestX, closestY
}

// GetDistance calculates the distance between two points
func GetDistance(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return math.Sqrt(dx*dx + dy*dy)
}

// GetDistanceSquared calculates the squared distance between two points (faster than GetDistance)
func GetDistanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

// NormalizeVector normalizes a vector to unit length
func NormalizeVector(x, y float64) (float64, float64) {
	length := math.Sqrt(x*x + y*y)
	if length == 0 {
		return 0, 0
	}
	return x / length, y / length
}

// ReflectVector reflects a vector off a surface with the given normal
func ReflectVector(vecX, vecY, normalX, normalY float64) (float64, float64) {
	// Reflection formula: r = d - 2(d·n)n
	// where d is the incident vector, n is the normal, r is the reflected vector
	dotProduct := vecX*normalX + vecY*normalY
	reflectedX := vecX - 2*dotProduct*normalX
	reflectedY := vecY - 2*dotProduct*normalY
	return reflectedX, reflectedY
}

// SeparateEntities separates two overlapping entities
func SeparateEntities(bounds1, bounds2 *Bounds, mass1, mass2 float64) {
	if !CheckAABBCollision(*bounds1, *bounds2) {
		return
	}

	// Calculate overlap
	overlapX := math.Min(bounds1.X+bounds1.Width, bounds2.X+bounds2.Width) - math.Max(bounds1.X, bounds2.X)
	overlapY := math.Min(bounds1.Y+bounds1.Height, bounds2.Y+bounds2.Height) - math.Max(bounds1.Y, bounds2.Y)

	// Determine separation direction (choose minimum overlap)
	if overlapX < overlapY {
		// Separate horizontally
		totalMass := mass1 + mass2
		separation1 := overlapX * (mass2 / totalMass)
		separation2 := overlapX * (mass1 / totalMass)

		if bounds1.X < bounds2.X {
			bounds1.X -= separation1
			bounds2.X += separation2
		} else {
			bounds1.X += separation1
			bounds2.X -= separation2
		}
	} else {
		// Separate vertically
		totalMass := mass1 + mass2
		separation1 := overlapY * (mass2 / totalMass)
		separation2 := overlapY * (mass1 / totalMass)

		if bounds1.Y < bounds2.Y {
			bounds1.Y -= separation1
			bounds2.Y += separation2
		} else {
			bounds1.Y += separation1
			bounds2.Y -= separation2
		}
	}
}