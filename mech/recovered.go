package mech

// Practice (squad) targeting settings taken from the recovered controller. Zero values keep the 6.17 behaviour.
const (
	rcStick     = 60.0   // target stickiness bonus (6.17: 200, recovered: 60)
	rcLowBonus  = 100.0   // bonus for targets under 30% HP (6.17: 200, recovered: 100)
	rcNoFinish  = true // recovered: no close-in on a weak foe (6.17: Finish 20 m)
	rcLethalK   = 0.8   // share of partial-laser damage that must cover the foe's HP (6.17: 0.7, recovered: 0.8)
)

func recoveredSquad(p Params) Params {
	if rcStick > 0 {
		p.Stick = rcStick
	}
	if rcLowBonus > 0 {
		p.LowBonus = rcLowBonus
	}
	if rcNoFinish {
		p.Finish = 0
	}
	if rcLethalK > 0 {
		p.LethalK = rcLethalK
	}
	return p
}
