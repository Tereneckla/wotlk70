package hunter

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

var ItemSetBeastLord = core.NewItemSet(core.ItemSet{
	Name: "Beast Lord Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in kill_command.go
		},
	},
})

var ItemSetDemonStalker = core.NewItemSet(core.ItemSet{
	Name: "Demon Stalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in multi_shot.go
		},
	},
})

var ItemSetRiftStalker = core.NewItemSet(core.ItemSet{
	Name: "Rift Stalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Handled in steady_shot.go
		},
	},
})

var ItemSetGronnstalker = core.NewItemSet(core.ItemSet{
	Name: "Gronnstalker's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in rotation.go
		},
		4: func(agent core.Agent) {
			// Handled in steady_shot.go
		},
	},
})

func (hunter *Hunter) talonOfAlarActive() float64 {
	return core.TernaryFloat64(hunter.TalonOfAlarAura != nil && hunter.TalonOfAlarAura.IsActive(), 40, 0)
}

func init() {
	core.NewItemEffect(30488, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		hunter.TalonOfAlarAura = hunter.RegisterAura(core.Aura{
			Label:    "Improved Shots",
			ActionID: core.ActionID{SpellID: 37507},
			Duration: time.Second * 6,
		})

		hunter.RegisterAura(core.Aura{
			Label:    "Talon of Al'ar",
			Duration: core.NeverExpires,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != hunter.ArcaneShot || !result.Landed() {
					return
				}

			},
		})
	})

	core.NewItemEffect(32336, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()
		const manaGain = 8.0
		manaMetrics := hunter.NewManaMetrics(core.ActionID{SpellID: 46939})

		hunter.RegisterAura(core.Aura{
			Label:    "Black Bow of the Betrayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}
				hunter.AddMana(sim, manaGain, manaMetrics)
			},
		})
	})

	core.NewItemEffect(32487, func(agent core.Agent) {
		hunter := agent.(HunterAgent).GetHunter()

		procAura := hunter.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32487}, stats.Stats{stats.AttackPower: 275, stats.RangedAttackPower: 275}, time.Second*8)
		const procChance = 0.15

		hunter.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != hunter.SteadyShot {
					return
				}
				if sim.RandomFloat("Ashtongue Talisman of Swiftness") > procChance {
					return
				}
				procAura.Activate(sim)
			},
		})
	})

}
