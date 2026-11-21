package mech

import "math"

// Tweaks on top of the base brain. Zero values keep the base behaviour.
const (
	// seedSalt shifts the RNG seed. Practice seeds are fixed, so a new salt is a new sample of the same strategy.
	seedSalt int64 = 832

	// healHold: practice (squad) mode only. After a fight, while no engaged foe is in view and HP is below this
	// fraction, stop looting and farming (dealing damage resets the regen timer), stand still (silent) or take
	// cover, and regenerate before the next fight. 0 = off.
	healHold = 0.0
)

// Ideas from kyoto-max (practice/squad mode only; false = base behaviour).
const (
	// kyEarly: before 200 s with < 10 cubes, collect cubes within 150 m and farm crates/drones down to 40% HP.
	kyEarly = false
	// kyExplore: in that early phase, with nothing to loot, walk to unvisited 100 m grid points inside the zone.
	kyExplore = false 
	// kyBoost: boost when threatened (energy > 30) and when closing on a target more than 30 m away (energy > 40).
	kyBoost = false
	// kyUpgrades: upgrade order plating, armor, nanites, servos, jets, laser, reactor, gun, ...
	kyUpgrades = false
	// kyHoldFire: below 60% HP, do not shoot a healthy target more than 120 m away (stay hidden, let regen start).
	kyHoldFire = false
	// kyListen: when idle, look toward a recent nearby sound that no visible mech explains.
	kyListen = true
)

var kyotoUpgradeOrder = []string{"plating", "armor", "nanites", "servos", "jets", "laser_damage", "laser_charge", "reactor",
	"gun_damage", "gun_rate", "lightning_efficiency", "lightning_range", "gun_mag", "gun_velocity",
	"capacitor", "gyros", "lightning_chain", "mortar_reload", "mortar_splash", "mortar_velocity", "sensors", "laser_pierce"}

func pickFrom(of *Offer, order []string) int {
	best, bestRank := 0, 1<<30
	for i, c := range of.Choices {
		rank := len(order)
		for r, id := range order {
			if id == c.ID {
				rank = r
				break
			}
		}
		if rank < bestRank {
			best, bestRank = i, rank
		}
	}
	return best
}

// kyEarlyPhase is the kyoto-max early practice phase.
func (b *Brain) kyEarlyPhase(o *Obs) bool {
	return b.pracOn && !b.notPractice && o.T < 200 && o.Self.Cubes < 10
}

type heardSound struct {
	kind    string
	bearing float64
	far     bool
	t       float64
}

var listenKinds = map[string]bool{"gun": true, "laser": true, "laser_charge": true, "lightning": true,
	"boost": true, "jet": true, "mortar": true, "step": true}

// listen remembers recent sounds that no visible mech explains.
func (b *Brain) listen(o *Obs) {
	keep := b.heard[:0]
	for _, h := range b.heard {
		if b.now-h.t < 3 {
			keep = append(keep, h)
		}
	}
	b.heard = keep
	for _, s := range o.Sounds {
		if !listenKinds[s.Kind] {
			continue
		}
		vis := false
		for i := range o.Mechs {
			m := &o.Mechs[i]
			if m.Kind == "mech" && math.Abs(angDiff(math.Atan2(m.Pos[1]-o.Self.Pos[1], m.Pos[0]-o.Self.Pos[0]), s.Bearing)) < 0.3 {
				vis = true
				break
			}
		}
		if !vis {
			b.heard = append(b.heard, heardSound{s.Kind, s.Bearing, s.Dist == "far", b.now - s.Age})
		}
	}
}

func (b *Brain) lastHeard() *heardSound {
	var best *heardSound
	for i := range b.heard {
		h := &b.heard[i]
		if h.far || b.now-h.t > 2 {
			continue
		}
		if best == nil || h.t > best.t {
			best = h
		}
	}
	return best
}

// kyExploreGoal: nearest unvisited 100 m grid point inside the current zone that we can reach in time.
func (b *Brain) kyExploreGoal(o *Obs) (float64, float64, bool) {
	me := &o.Self
	z := &o.Zone
	cz := b.curZone(z)
	r := cz.R - 30
	n := int(math.Floor(r / 100))
	x, y := me.Pos[0], me.Pos[1]
	var bx, by, bd float64
	found := false
	for i := -n; i <= n; i++ {
		for j := -n; j <= n; j++ {
			px, py := cz.X+float64(i)*100, cz.Y+float64(j)*100
			if hyp(px-cz.X, py-cz.Y) > r || !b.W.InBounds(px, py) || !b.W.Free(px, py) {
				continue
			}
			if nf := b.W.NearestFree(b.W.Cell(px, py)); nf < 0 || b.W.Comp[nf] != b.W.MainComp {
				continue
			}
			d := hyp(px-x, py-y)
			if b.now+d*1.4/me.Stats.WalkSpeed+5 > b.zoneDeadline(z, px, py, 6) {
				continue
			}
			near := false
			for _, p := range b.trail {
				if hyp(px-p[0], py-p[1]) < 70 {
					near = true
					break
				}
			}
			if near {
				continue
			}
			if !found || d < bd {
				bx, by, bd, found = px, py, d, true
			}
		}
	}
	return bx, by, found
}

// kyoto-max v34-v41 lineage on the squad (practice) parameters (0/false = base).
const (
	kvFightBoost  = 40.0   // boost-speed jinking in fights above this energy (v38: 40)
	kvDodgeBoost  = 25.0   // boost on every dodge above this energy (v39: 25)
	kvThreatBoost = 30.0   // boost whenever threatened and moving above this energy (v40: 30)
	kvBackstab    = 40.0   // target bonus for foes showing their back, rear damage in lethal release (v36: 40)
	kvSpeedUp     = true // servos/jets/reactor right after plating and armor (v41)
)

func kyotoSquad(p Params) Params {
	if kvFightBoost > 0 {
		p.FightBoost = kvFightBoost
	}
	if kvDodgeBoost > 0 {
		p.DodgeBoost = kvDodgeBoost
	}
	if kvThreatBoost > 0 {
		p.ThreatBoost = kvThreatBoost
	}
	if kvBackstab > 0 {
		p.Backstab = kvBackstab
	}
	if kvSpeedUp {
		p.SpeedUp = true
	}
	return p
}
