package mech

// Tweaks on top of the base brain. Zero values keep the base behaviour.
const (
	// seedSalt shifts the RNG seed. Practice seeds are fixed, so a new salt is a new sample of the same strategy.
	seedSalt int64 = 404

	// healHold: practice (squad) mode only. After a fight, while no engaged foe is in view and HP is below this
	// fraction, stop looting and farming (dealing damage resets the regen timer), stand still (silent) or take
	// cover, and regenerate before the next fight. 0 = off.
	healHold = 0.0
)
