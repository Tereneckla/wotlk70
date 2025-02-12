import {
	Consumes,
	Debuffs,
	EquipmentSpec, Explosive, Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	RaidBuffs,
	UnitReference, Spec,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';

import {
	BalanceDruid_Options as BalanceDruidOptions,
	BalanceDruid_Rotation as BalanceDruidRotation,
	BalanceDruid_Rotation_IsUsage,
	BalanceDruid_Rotation_MfUsage,
	BalanceDruid_Rotation_Type as RotationType,
	BalanceDruid_Rotation_WrathUsage,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { Player } from "../core/player";
import { APLRotation } from '../core/proto/apl.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const Phase1Talents = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '5032003125331303213305311231--2',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfStarfire,
			major2: DruidMajorGlyph.GlyphOfStarfall,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Default,
	maintainFaerieFire: true,
	useSmartCooldowns: true,
	mfUsage: BalanceDruid_Rotation_MfUsage.BeforeLunar,
	isUsage: BalanceDruid_Rotation_IsUsage.OptimizeIs,
	wrathUsage: BalanceDruid_Rotation_WrathUsage.RegularWrath,
	useStarfire: true,
	useBattleRes: false,
	playerLatency: 200,
});

export const DefaultOptions = BalanceDruidOptions.create({
	innervateTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	prepopPotion: Potions.HastePotion,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	sanctifiedRetribution: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
	wrathOfAirTotem: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
	sunderArmor: true,
	totemOfWrath: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
};

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":42554,"enchant":3820,"gems":[41285,40049]},
		{"id":40680},
		{"id":37673,"enchant":3810,"gems":[42144]},
		{"id":41610,"enchant":3859},
		{"id":39547,"enchant":3832,"gems":[42144,40026]},
		{"id":37884,"enchant":2332,"gems":[0]},
		{"id":39544,"enchant":3604,"gems":[42144,0]},
		{"id":40696,"enchant":3601,"gems":[40014,39998]},
		{"id":37854,"enchant":3719},
		{"id":44202,"enchant":3606,"gems":[39998]},
		{"id":40585},
		{"id":43253,"gems":[40026]},
		{"id":37873},
		{"id":40682},
		{"id":45085,"enchant":3834},
		{"id":40698},
		{"id":40712}
	]}`),
};

export const ROTATION_PRESET_P3_APL = {
name: 'Basic P3 APL',
rotation: SavedRotation.create({
	specRotationOptionsJson: BalanceDruidRotation.toJsonString(DefaultRotation),
	rotation: APLRotation.fromJsonString(`{
      "type": "TypeAPL",
      "prepullActions": [
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1.5s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":48461}}},"doAtValue":{"const":{"val":"-1.5s"}}}
      ],
      "priorityList": [
        {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"5"}}}},"castSpell":{"spellId":{"tag":-1,"spellId":2825}}}},
        {"action":{"castSpell":{"spellId":{"itemId":41119}}}},
        {"action":{"multidot":{"spellId":{"spellId":48463},"maxDots":1,"maxOverlap":{"const":{"val":"0ms"}}}}},
        {"action":{"castSpell":{"spellId":{"spellId":53201}}}},
        {"action":{"castSpell":{"spellId":{"spellId":65861}}}},
        {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"10s"}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8"}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
        {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"10s"}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8"}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15s"}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
        {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"sourceUnit":{},"auraId":{"spellId":48518}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8s"}}}}]}},"castSpell":{"spellId":{"spellId":48465}}}},
        {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"sourceUnit":{},"auraId":{"spellId":48517}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48517}}},"rhs":{"const":{"val":"14.8s"}}}}]}},"castSpell":{"spellId":{"spellId":48461}}}},
        {"action":{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48468}}}}},{"auraInternalCooldown":{"auraId":{"spellId":48518}}}]}},"castSpell":{"spellId":{"spellId":48468}}}},
        {"action":{"condition":{"auraInternalCooldown":{"sourceUnit":{},"auraId":{"spellId":48518}}},"castSpell":{"spellId":{"spellId":48465}}}},
        {"action":{"castSpell":{"spellId":{"spellId":48461}}}}
      ]
    }`),
}),
};