package main

// Ideas from the MIT-licensed automata bot (practice mode only; zero values keep the 6.25 behaviour).
const (
	// gearNoFirst: until we hold a laser (and before gearLateS), do not start fights with house bots that have not hit us
	// in the last 6 s, carry no laser and are farther than 45 m: arm up first.
	gearNoFirst = true
	gearLateS   = 120.0
	// gearLaserSeek: how far (m) we travel for a laser pickup in practice mode (6.25: 240).
	gearLaserSeek = 0.0
)
