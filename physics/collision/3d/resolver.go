package rgcol3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ResolveMTV applies the contact minimum translation vector to a SpatialCollider.
func ResolveMTV(active SpatialCollider, hit Contact) {
	if !hit.Hit || hit.Penetration <= 0 {
		return
	}
	active.SetPosition(rl.Vector3Add(active.GetPosition(), rl.Vector3Scale(hit.Normal, hit.Penetration)))
}

// ResolveMultiMTV resolves overlap by repeatedly applying Minimum Translation Vectors.
// Helps handle cases where pushing out of one object pushes into another.
// Returns true if all overlaps were resolved within maxIter, false if still stuck.
func ResolveMultiMTV(active SpatialCollider, others []Collider, maxIter int) bool {
	for range maxIter {
		hitThisIter := false

		for _, other := range others {
			if active == other {
				continue
			}

			contact := active.Collide(other)
			if contact.Hit && contact.Penetration > epsilon {
				hitThisIter = true
				ResolveMTV(active, contact)
			}
		}

		if !hitThisIter {
			return true // fully resolved
		}
	}
	return false // still overlapping something
}
