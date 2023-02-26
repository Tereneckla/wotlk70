package retribution

import (
	"testing"

	_ "github.com/Tereneckla/wotlk/sim/common" // imported to get item effects included.
	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman, proto.Race_RaceDraenei, proto.Race_RaceDwarf},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: Phase1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Retribution Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Retribution Paladin SOC",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Command,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: defaultRetRotation,
					},
				},
			},
			{
				Label: "Retribution Paladin SOR",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Righteousness,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: defaultRetRotation,
					},
				},
			},
			{
				Label: "Retribution Paladin SOV 2 Target Swapping",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Vengeance,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: &proto.RetributionPaladin_Rotation{
							ConsSlack:            500,
							ExoSlack:             500,
							UseDivinePlea:        false,
							DivinePleaPercentage: 0.75,
							HolyWrathThreshold:   4,
							SovTargets:           2,
						},
					},
				},
			},
		},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeMace,
			},
			ArmorType: proto.ArmorType_ArmorTypePlate,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeLibram,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceBloodElf,
				Class:         proto.Class_ClassPaladin,
				TalentsString: StandardTalents,
				Glyphs:        StandardGlyphs,
				Equipment:     Phase1Gear,
				Consumes:      FullConsumes,
				Spec:          DefaultOptions,
				Buffs:         core.FullIndividualBuffs,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}

var StandardTalents = "050501-05-05232051203331302133231331"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfJudgement),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfConsecration),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfKings),
}

var defaultRetRotation = &proto.RetributionPaladin_Rotation{
	ConsSlack:            500,
	ExoSlack:             500,
	UseDivinePlea:        false,
	DivinePleaPercentage: 0.75,
	HolyWrathThreshold:   4,
}

var defaultRetOptions = &proto.RetributionPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Seal:      proto.PaladinSeal_Vengeance,
	Aura:      proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options:  defaultRetOptions,
		Rotation: defaultRetRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodRoastedClefthoof,
	SuperSapper:     true,
}

var Phase1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 34244
	},
	{
		"id": 32260
	},
	{
		"id": 34388
	},
	{
		"id": 34241
	},
	{
		"id": 34369
	},
	{
		"id": 34431
	},
	{
		"id": 34378
	},
	{
		"id": 34485
	},
	{
		"id": 34188
	},
	{
		"id": 34561
	},
	{
		"id": 32497
	},
	{
		"id": 34361
	},
	{
		"id": 34427
	},
	{
		"id": 32505
	},
	{
		"id": 34989
	},
	{},
	{
		"id": 37574
	}
]}`)
