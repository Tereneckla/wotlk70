import { Consumes } from '../core/proto/common.js';
import { WeaponImbue } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';

import {
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Options as FeralDruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
	FeralDruid_Rotation_AplType,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const OmenTalents = {
	name: 'Omen',
	data: SavedTalents.create({
		talentsString: '-5032021323220100531202303104-20350001',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfRip,
			major2: DruidMajorGlyph.GlyphOfShred,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const BerserkTalents = {
	name: 'Berserk',
	data: SavedTalents.create({
		talentsString: '-503202132322010053120230310511-1043',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfRip,
			major2: DruidMajorGlyph.GlyphOfShred,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,

	bearWeaveType: FeralDruid_Rotation_BearweaveType.None,
	minCombosForRip: 5,
	minCombosForBite: 5,

	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 4.0,
	berserkBiteThresh: 25.0,
	berserkFfThresh: 15.0,
	powerbear: false,
	minRoarOffset: 12.0,
	ripLeeway: 3.0,
	maintainFaerieFire: true,
	hotUptime: 0.0,
	snekWeave: false,
	flowerWeave: false,
	raidTargets: 30,
	maxFfDelay: 0.1,
	prePopOoc: true,
});

export const DefaultOptions = FeralDruidOptions.create({
	latencyMs: 100,
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodGrilledMudfish,
	defaultPotion: Potions.PotionOfSpeed,
	weaponMain: WeaponImbue.ImbueAdamantiteWeightStone,
});

export const ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Legacy Default',
	rotation: SavedRotation.create({
		specRotationOptionsJson: FeralDruidRotation.toJsonString(DefaultRotation),
	}),
}

export const PreRaid_PRESET = {
	name: 'PreRaid',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":42550,"enchant":3817,"gems":[41398,39996]},
		{"id":40678},
		{"id":37139,"enchant":3808,"gems":[39996]},
		{"id":37840,"enchant":3605},
		{"id":37219,"enchant":3832},
		{"id":44203,"enchant":3845,"gems":[0]},
		{"id":37409,"enchant":3604,"gems":[0]},
		{"id":40694,"gems":[49110,39996]},
		{"id":37644,"enchant":3823},
		{"id":44297,"enchant":3606},
		{"id":37642},
		{"id":37624},
		{"id":40684},
		{"id":37166},
		{"id":37883,"enchant":3827},
		{},
		{"id":40713}
  ]}`),
};
