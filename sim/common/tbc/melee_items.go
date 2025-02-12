package tbc

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

func init() {
	core.NewSimpleStatItemEffect(28484, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30) // Bulwark of Kings
	core.NewSimpleStatItemEffect(28485, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30) // Bulwark of Ancient Kings

	// Proc effects. Keep these in order by item ID.

	core.NewItemEffect(871, func(agent core.Agent) {
		character := agent.GetCharacter()
		procMask := character.GetProcMaskForItem(871)

		ppmm := character.AutoAttacks.NewPPMManager(3.0, procMask)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond,
		}

		flurryAxeSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 18797},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    procMask,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: character.AutoAttacks.MHConfig.DamageMultiplier,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: character.AutoAttacks.MHConfig.ThreatMultiplier,

			ApplyEffects: character.AutoAttacks.MHConfig.ApplyEffects,
		})

		character.RegisterAura(core.Aura{
			Label:    "Flurry Axe",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Flurry Axe") {
					return
				}
				icd.Use(sim)
				flurryAxeSpell.Cast(sim, result.Target)
			},
		})

	})

	core.NewItemEffect(9423, func(agent core.Agent) {
		character := agent.GetCharacter()
		procMask := character.GetProcMaskForItem(9423)
		ppmm := character.AutoAttacks.NewPPMManager(2.0, procMask)

		procAura := character.NewTemporaryStatsAura("Jackhammer", core.ActionID{ItemID: 9423}, stats.Stats{stats.MeleeHaste: 300, stats.SpellHaste: 300}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Jackhammer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Jackhammer") {
					return
				}
				procAura.Activate(sim)

			},
		})
	})

	core.NewItemEffect(9449, func(agent core.Agent) {
		character := agent.GetCharacter()

		// Assumes that the user will swap pummelers to have the buff for the whole fight.
		character.AddStat(stats.MeleeHaste, 500)
	})

	core.NewItemEffect(19019, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(19019)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
			},
		})

		makeDebuffAura := func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		}

		numHits := core.MinInt32(5, character.Env.GetNumTargets())
		debuffAuras := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
		for i, target := range character.Env.Encounter.TargetUnits {
			debuffAuras[i] = makeDebuffAura(target)
		}

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
					if result.Landed() {
						debuffAuras[target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Thunderfury",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Thunderfury") {
					singleTargetSpell.Cast(sim, result.Target)
					bounceSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	core.NewItemEffect(24114, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BonusDamage += 5
	})

	core.NewItemEffect(29297, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 0.03
		procAura := character.NewTemporaryStatsAura("Band of the Eternal Defender Proc", core.ActionID{ItemID: 29297}, stats.Stats{stats.Armor: 800}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Band of the Eternal Defender",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.SpellSchool != core.SpellSchoolPhysical {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Band of the Eternal Defender") < procChance {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(29962, func(agent core.Agent) {
		character := agent.GetCharacter()
		procMask := character.GetProcMaskForItem(29962)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procAuraMH := character.NewTemporaryStatsAura("Heartrazor Proc", core.ActionID{ItemID: 29962, Tag: 1}, stats.Stats{stats.AttackPower: 270, stats.RangedAttackPower: 270}, time.Second*10)
		procAuraOH := character.NewTemporaryStatsAura("Heartrazor Proc", core.ActionID{ItemID: 29962, Tag: 2}, stats.Stats{stats.AttackPower: 270, stats.RangedAttackPower: 270}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Heartrazor",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Heartrazor") {
					return
				}
				if spell.IsMH() {
					procAuraMH.Activate(sim)
				} else {
					procAuraOH.Activate(sim)
				}

			},
		})
	})

	core.NewItemEffect(29996, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(29996)
		pppm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		actionID := core.ActionID{ItemID: 29996}

		var resourceMetricsRage *core.ResourceMetrics
		var resourceMetricsEnergy *core.ResourceMetrics
		if character.HasRageBar() {
			resourceMetricsRage = character.NewRageMetrics(actionID)
		}
		if character.HasEnergyBar() {
			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rod of the Sun King",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if pppm.Proc(sim, spell.ProcMask, "Rod of the Sun King") {
					switch spell.Unit.GetCurrentPowerBar() {
					case core.RageBar:
						spell.Unit.AddRage(sim, 5, resourceMetricsRage)
					case core.EnergyBar:
						spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
					}
				}
			},
		})
	})

	core.NewItemEffect(30090, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 3.7 / 60.0
		procAura := character.NewTemporaryStatsAura("World Breaker Proc", core.ActionID{ItemID: 30090}, stats.Stats{stats.MeleeCrit: 900}, time.Second*4)

		character.RegisterAura(core.Aura{
			Label:    "World Breaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					procAura.Deactivate(sim)
					return
				}
				if sim.RandomFloat("World Breaker") > procChance {
					procAura.Deactivate(sim)
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(31332, func(agent core.Agent) {
		character := agent.GetCharacter()
		procMask := character.GetProcMaskForItem(31332)

		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond,
		}

		blinkStrikeSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 38308},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    procMask,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: character.AutoAttacks.MHConfig.DamageMultiplier,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: character.AutoAttacks.MHConfig.ThreatMultiplier,

			ApplyEffects: character.AutoAttacks.MHConfig.ApplyEffects,
		})

		character.RegisterAura(core.Aura{
			Label:    "Blinkstrike",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(procMask) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "Blinkstrike") {
					return
				}
				icd.Use(sim)
				blinkStrikeSpell.Cast(sim, result.Target)
			},
		})

	})

	core.NewItemEffect(31193, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(31193)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 24585},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(48, 54) + spell.SpellPower()
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Blade of Unquenched Thirst Trigger",
			ActionID:   core.ActionID{ItemID: 31193},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.02,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(32262, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(32262)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 40291},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeMagicHitAndCrit)
			},
		})

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Siphon Essence",
			ActionID: core.ActionID{SpellID: 40291},
			Duration: time.Second * 6,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, result.Target)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Syphon of the Nathrezim",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Syphon Of The Nathrezim") {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(32375, func(agent core.Agent) {
		character := agent.GetCharacter()

		const procChance = 0.02
		procAura := character.NewTemporaryStatsAura("Bulwark Of Azzinoth Proc", core.ActionID{ItemID: 32375}, stats.Stats{stats.Armor: 2000}, time.Second*10)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Bulwark Of Azzinoth",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.SpellSchool == core.SpellSchoolPhysical && sim.RandomFloat("Bulwark of Azzinoth") < procChance {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(34473, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Commendation of Kael'Thas Proc", core.ActionID{ItemID: 34473}, stats.Stats{stats.Dodge: 152}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 30,
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Commendation of Kael'Thas",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if aura.Unit.CurrentHealthPercent() >= 0.35 {
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

	core.NewItemEffect(12590, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(12590)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		effectAura := character.NewTemporaryStatsAura("Felstriker Proc", core.ActionID{SpellID: 16551}, stats.Stats{stats.MeleeCrit: 100 * core.CritRatingPerCritChance}, time.Second*3)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Felstriker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Felstriker") {
					effectAura.Activate(sim)
				}
			},
		})
	})
}
