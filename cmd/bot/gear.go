package main

// Ideas from the MIT-licensed automata bot (practice mode only; zero values keep the 6.25 behaviour).
const (
	// gearNoFirst: until we hold a laser (and before gearLateS), do not start fights with house bots that have not hit us
	// in the last 6 s, carry no laser and are farther than 45 m: arm up first.
	gearNoFirst = true
	gearLateS   = 120.0
	// gearLaserSeek: how far (m) we travel for a laser pickup in practice mode (6.25: 240).
	gearLaserSeek = 300.0
)

// More laser-first ideas from automata (practice mode only; 0 = off).
const (
	// gearFightLaser: during a fight, without a laser, go take a known laser pickup within this many metres.
	gearFightLaser = 70.0
	// gearEmptySeek: while the second weapon slot is empty, look for any wanted weapon within this many metres.
	gearEmptySeek = 0.0
)

// gearFightGrab: during a fight, while the second weapon slot is empty, go take any wanted weapon pickup within
// this many metres (automata gear_grab_r). 0 = off.
const gearFightGrab = 40.0

// fixCubeLevel: in-fight cube detours skip cubes more than 3.5 m above or below us (cannot be picked up).
const fixCubeLevel = false
