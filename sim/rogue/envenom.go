package rogue

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
)

func (rogue *Rogue) registerEnvenom() {
	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 32684},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.deadlyPoisonProcChanceBonus += 0.15
			rogue.UpdateInstantPoisonPPM(0.75)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.deadlyPoisonProcChanceBonus -= 0.15
			rogue.UpdateInstantPoisonPPM(0.0)
		},
	})
	cost := 35.0
	if rogue.HasSetBonus(ItemSetAssassination, 4) {
		cost -= 10
	}
	chanceToRetainStacks := []float64{0, 0.33, 0.66, 1}[rogue.Talents.MasterPoisoner]
	deathMantleDamage := core.TernaryFloat64(rogue.HasSetBonus(Tier5, 2), 40, 0)
	rogue.Envenom = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 32684},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags() | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          cost,
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
				rogue.applyDeathmantle(sim, spell, cast)
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0 && rogue.DeadlyPoison.Dot(target).IsActive()
		},

		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			// - the aura is active even if the attack fails to land
			// - the aura is applied before the hit effect
			// See: https://github.com/where-fore/rogue-wotlk/issues/32
			rogue.EnvenomAura.Duration = time.Second * time.Duration(1+comboPoints)
			rogue.EnvenomAura.Activate(sim)

			dp := rogue.DeadlyPoison.Dot(target)
			// - 215 base is scaled by consumed doses (<= comboPoints)
			// - apRatio is independent of consumed doses (== comboPoints)
			consumed := core.MinInt32(dp.GetStacks(), comboPoints)
			baseDamage := 147*float64(consumed) + 0.09*float64(comboPoints)*spell.MeleeAttackPower() + deathMantleDamage*float64(comboPoints)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
				if !sim.Proc(chanceToRetainStacks, "Master Poisoner") {
					if newStacks := dp.GetStacks() - comboPoints; newStacks > 0 {
						dp.SetStacks(sim, newStacks)
					} else {
						dp.Cancel(sim)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (rogue *Rogue) EnvenomDuration(comboPoints int32) time.Duration {
	return time.Second * (1 + time.Duration(comboPoints))
}
