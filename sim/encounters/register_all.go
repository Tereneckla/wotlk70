package encounters

import (
	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/encounters/icc"
	"github.com/Tereneckla/wotlk/sim/encounters/naxxramas"
	"github.com/Tereneckla/wotlk/sim/encounters/toc"
	"github.com/Tereneckla/wotlk/sim/encounters/ulduar"
)

func init() {
	naxxramas.Register()
	ulduar.Register()
	toc.Register()
	icc.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
