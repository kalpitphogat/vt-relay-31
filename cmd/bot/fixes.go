package main

// Small fixes from public bots (zero values keep the 6.31 behaviour).
const (
	// fixCubeLevel: in-fight cube detours skip cubes more than 3.5 m above or below us (on roofs: cannot be picked up).
	fixCubeLevel = true
	// fixNoAirBoost: practice mode, no boost while airborne (a jet climb) except when escaping or dodging.
	fixNoAirBoost = false
	// fixFFAUntil: a death we did not cause switches to free-for-all play only before this match time (s); 0 = any time.
	fixFFAUntil = 80.0
)
