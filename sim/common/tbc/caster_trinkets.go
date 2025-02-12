package tbc

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

func init() {

	// Offensive trinkets. Keep these in order by item ID.
	core.NewSimpleStatOffensiveTrinketEffect(23046, stats.Stats{stats.SpellPower: 130}, time.Second*20, time.Minute*2)  // Restrained Essence of Sapphiron
	core.NewSimpleStatOffensiveTrinketEffect(24126, stats.Stats{stats.SpellPower: 150}, time.Second*20, time.Minute*5)  // Living Ruby Serpent
	core.NewSimpleStatOffensiveTrinketEffect(25634, stats.Stats{stats.SpellPower: 113}, time.Second*20, time.Minute*2)  // Oshu'gun Relic
	core.NewSimpleStatOffensiveTrinketEffect(28949, stats.Stats{stats.SpellPower: 120}, time.Second*15, time.Second*90) // Vengeance of the Illidari
	core.NewSimpleStatOffensiveTrinketEffect(28223, stats.Stats{stats.SpellPower: 167}, time.Second*20, time.Minute*2)  // Arcanist's Stone
	core.NewSimpleStatOffensiveTrinketEffect(29132, stats.Stats{stats.SpellPower: 150}, time.Second*15, time.Second*90) // Scryer's Bloodgem
	core.NewSimpleStatOffensiveTrinketEffect(29179, stats.Stats{stats.SpellPower: 150}, time.Second*15, time.Second*90) // Xiri's Gift
	core.NewSimpleStatOffensiveTrinketEffect(29370, stats.Stats{stats.SpellPower: 158}, time.Second*20, time.Minute*2)  // Icon of the Silver Crescent
	core.NewSimpleStatOffensiveTrinketEffect(29376, stats.Stats{stats.SpellPower: 158}, time.Second*20, time.Minute*2)  // Essence of the Marytr
	core.NewSimpleStatOffensiveTrinketEffect(30340, stats.Stats{stats.SpellPower: 125}, time.Second*15, time.Second*90) // Living Ruby Serpent
	core.NewSimpleStatOffensiveTrinketEffect(32483, stats.Stats{stats.SpellHaste: 175}, time.Second*20, time.Minute*2)  // Skull of Gul'dan
	core.NewSimpleStatOffensiveTrinketEffect(33829, stats.Stats{stats.SpellPower: 211}, time.Second*20, time.Minute*2)  // Hex Shrunken Head
	core.NewSimpleStatOffensiveTrinketEffect(34429, stats.Stats{stats.SpellPower: 320}, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	// Even though these item effects are handled elsewhere, add them so they are
	// detected for automatic testing.
	for _, itemID := range core.AlchStoneItemIDs {
		core.NewItemEffect(itemID, func(core.Agent) {})
	}

	core.NewItemEffect(21625, func(agent core.Agent) { // Scarab Brooch
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 21625}

		shieldSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 26470},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Scarab Brooch Shield",
					Duration: time.Second * 30,
				},
			},
		})

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Persistent Shield",
			ActionID: core.ActionID{SpellID: 26467},
			Callback: core.CallbackOnHealDealt,
			Duration: time.Second * 30,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				shieldSpell.Shield(result.Target).Apply(sim, result.Damage*0.15)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	core.NewItemEffect(30447, func(agent core.Agent) {
		paladin := agent.GetCharacter()
		procAura := paladin.NewTemporaryStatsAura("Tome of Fiery Redemption Proc", core.ActionID{ItemID: 30447}, stats.Stats{stats.SpellPower: 290}, time.Second*15)

		icd := core.Cooldown{
			Timer:    paladin.NewTimer(),
			Duration: time.Second * 45,
		}
		paladin.RegisterAura(core.Aura{
			Label:    "Tome of Fiery Redemption",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

				if sim.RandomFloat("Tome of Fiery Redemption") > 0.15 {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(31856, func(agent core.Agent) {
		character := agent.GetCharacter()

		var apBonusPerStack stats.Stats
		apAura := character.RegisterAura(core.Aura{
			Label:     "DMC Crusade AP",
			ActionID:  core.ActionID{ItemID: 31856, Tag: 1},
			Duration:  time.Second * 10,
			MaxStacks: 20,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				apBonusPerStack = character.ApplyStatDependencies(stats.Stats{stats.AttackPower: 6, stats.RangedAttackPower: 6})
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatsDynamic(sim, apBonusPerStack.Multiply(float64(newStacks-oldStacks)))
			},
		})

		var spBonusPerStack stats.Stats
		spAura := character.RegisterAura(core.Aura{
			Label:     "DMC Crusade SP",
			ActionID:  core.ActionID{ItemID: 31856, Tag: 2},
			Duration:  time.Second * 10,
			MaxStacks: 10,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				spBonusPerStack = character.ApplyStatDependencies(stats.Stats{stats.SpellPower: 8})
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatsDynamic(sim, spBonusPerStack.Multiply(float64(newStacks-oldStacks)))
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "DMC Crusade",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					apAura.Activate(sim)
					apAura.AddStack(sim)
					apAura.Refresh(sim)
				} else if spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					if !result.Landed() {
						return
					}
					spAura.Activate(sim)
					spAura.AddStack(sim)
					spAura.Refresh(sim)
				}
			},
		})
	})
}
