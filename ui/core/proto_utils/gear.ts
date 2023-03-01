import { EquipmentSpec, ItemSwap } from '../proto/common.js';
import { GemColor } from '../proto/common.js';
import { ItemSlot } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { SimDatabase } from '../proto/common.js';
import { SimItem } from '../proto/common.js';
import { SimEnchant } from '../proto/common.js';
import { SimGem } from '../proto/common.js';
import { WeaponType } from '../proto/common.js';
import { arrayEquals, equalsOrBothNull } from '../utils.js';
import { distinct, getEnumValues } from '../utils.js';
import { isBluntWeaponType, isSharpWeaponType } from '../proto_utils/utils.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

import { isMetaGemActive } from './gems.js';
import { gemMatchesSocket } from './gems.js';
import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';
import { Stats } from './stats.js';

type InternalGear = Record<ItemSlot, EquippedItem | null>;

abstract class BaseGear {
	protected readonly gear: InternalGear;

	constructor(gear: Partial<InternalGear>) {
		this.getItemSlots().forEach(slot => {
			if (!gear[slot as ItemSlot])
				gear[slot as ItemSlot] = null;
		});
		this.gear = gear as InternalGear;
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear[slot];
	}

	asArray(): Array<EquippedItem | null> {
		return Object.values(this.gear);
	}

	removeUniqueGems(gear: InternalGear, newItem: EquippedItem) {
			// If the new item has unique gems, remove matching.
			newItem.gems
				.filter(gem => gem?.unique)
				.forEach(gem => {
					this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
						gear[slot] = gear[slot]?.removeGemsWithId(gem!.id) || null;
					});
				});
	}

	removeUniqueItems(gear: InternalGear, newItem: EquippedItem){
		if (newItem.item.unique) {
			this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
				if (gear[slot]?.item.id == newItem.item.id) {
					gear[slot] = null;
				}
			});
		}
	}

	validateWeaponCombo(gear: InternalGear, newSlot: ItemSlot, canDualWield2H: boolean) {
		// Check for valid weapon combos.
		if (!validWeaponCombo(gear[ItemSlot.ItemSlotMainHand]?.item, gear[ItemSlot.ItemSlotOffHand]?.item, canDualWield2H)) {
			if (newSlot == ItemSlot.ItemSlotOffHand) {
				gear[ItemSlot.ItemSlotMainHand] = null;
			} else {
				gear[ItemSlot.ItemSlotOffHand] = null;
			}
		}
	}

	abstract toDatabase(): SimDatabase
	abstract getItemSlots(): ItemSlot[]

	protected static itemToDB(item: Item): SimItem {
		return SimItem.fromJson(Item.toJson(item), { ignoreUnknownFields: true });
	}

	protected static enchantToDB(enchant: Enchant): SimEnchant {
		return SimEnchant.fromJson(Enchant.toJson(enchant), { ignoreUnknownFields: true });
	}

	protected static gemToDB(gem: Gem): SimGem {
		return SimGem.fromJson(Gem.toJson(gem), { ignoreUnknownFields: true });
	}
}

/**
 * Represents a full gear set, including items/enchants/gems for every slot.
 *
 * This is an immutable type.
 */
export class Gear extends BaseGear {

	constructor(gear: Partial<InternalGear>) {
		super(gear);
	}

	getItemSlots(): ItemSlot[] {
		return getEnumValues(ItemSlot);
	}

	equals(other: Gear): boolean {
		return this.asArray().every((thisItem, slot) => equalsOrBothNull(thisItem, other.getEquippedItem(slot), (a, b) => a.equals(b)));
	}

	/**
	 * Returns a new Gear set with the item equipped.
	 *
	 * Checks for validity and removes/exchanges items/gems as needed.
	 */
	withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null, canDualWield2H: boolean): Gear {
		// Create a new identical set of gear
		const newInternalGear = this.asMap();

		if (newItem) {
			this.removeUniqueGems(newInternalGear, newItem);
			this.removeUniqueItems(newInternalGear, newItem);
		}

		// Actually assign the new item.
		newInternalGear[newSlot] = newItem;

		this.validateWeaponCombo(newInternalGear, newSlot, canDualWield2H);

		return new Gear(newInternalGear);
	}

	getTrinkets(): Array<EquippedItem | null> {
		return [
			this.getEquippedItem(ItemSlot.ItemSlotTrinket1),
			this.getEquippedItem(ItemSlot.ItemSlotTrinket2),
		];
	}

	hasTrinket(itemId: number): boolean {
		return this.getTrinkets().map(t => t?.item.id).includes(itemId);
	}

	asMap(): InternalGear {
		const newInternalGear: Partial<InternalGear> = {};
		getEnumValues(ItemSlot).map(slot => Number(slot) as ItemSlot).forEach(slot => {
			newInternalGear[slot] = this.getEquippedItem(slot);
		});
		return newInternalGear as InternalGear;
	}

	asSpec(): EquipmentSpec {
		return EquipmentSpec.create({
			items: this.asArray().map(ei => ei ? ei.asSpec() : ItemSpec.create()),
		});
	}

	getAllGems(): Array<Gem> {
		return this.asArray()
			.map(ei => ei == null ? [] : ei.curGems())
			.flat();
	}

	getNonMetaGems(): Array<Gem> {
		return this.getAllGems().filter(gem => gem.color != GemColor.GemColorMeta);
	}

	statsFromGems(): Stats {
		let stats = new Stats();

		// Stats from just the gems.
		const gems = this.getAllGems();
		for (let i = 0; i < gems.length; i++) {
			stats = stats.add(new Stats(gems[i].stats));
		}

		// Stats from socket bonuses.
		const items = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		for (let i = 0; i < items.length; i++) {
			stats = stats.add(items[i].socketBonusStats());
		}

		return stats;
	}

	getGemsOfColor(color: GemColor): Array<Gem> {
		return this.getAllGems().filter(gem => gem.color == color);
	}

	getJCGems(): Array<Gem> {
		return this.getAllGems().filter(gem => gem.requiredProfession == Profession.Jewelcrafting);
	}

	getMetaGem(): Gem | null {
		return this.getGemsOfColor(GemColor.GemColorMeta)[0] || null;
	}

	gemColorCounts(): ({ red: number, yellow: number, blue: number }) {
		const gems = this.getAllGems();
		return {
			red: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorRed)).length,
			yellow: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorYellow)).length,
			blue: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorBlue)).length,
		};
	}

	// Returns true if this gear set has a meta gem AND the other gems meet the meta's conditions.
	hasActiveMetaGem(): boolean {
		const metaGem = this.getMetaGem();
		if (!metaGem) {
			return false;
		}

		const gemColorCounts = this.gemColorCounts();

		const gems = this.getAllGems();
		return isMetaGemActive(
			metaGem,
			gemColorCounts.red, gemColorCounts.yellow, gemColorCounts.blue);
	}

	hasInactiveMetaGem(): boolean {
		return this.getMetaGem() != null && !this.hasActiveMetaGem();
	}

	withoutMetaGem(): Gear {
		const headItem = this.getEquippedItem(ItemSlot.ItemSlotHead);
		const metaGem = this.getMetaGem();
		if (headItem && metaGem) {
			return this.withEquippedItem(ItemSlot.ItemSlotHead, headItem.removeGemsWithId(metaGem.id), true);
		} else {
			return this;
		}
	}

	hasBluntMHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotMainHand);
		return weapon != null && isBluntWeaponType(weapon.item.weaponType);
	}
	hasSharpMHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotMainHand);
		return weapon != null && isSharpWeaponType(weapon.item.weaponType);
	}
	hasBluntOHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotOffHand);
		return weapon != null && isBluntWeaponType(weapon.item.weaponType);
	}
	hasSharpOHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotOffHand);
		return weapon != null && isSharpWeaponType(weapon.item.weaponType);
	}

	getProfessionRequirements(): Array<Profession> {
		return distinct((this.asArray().filter(ei => ei != null) as Array<EquippedItem>)
			.map(ei => ei.getProfessionRequirements())
			.flat());
	}
	getFailedProfessionRequirements(professions: Array<Profession>): Array<Item | Gem | Enchant> {
		return (this.asArray().filter(ei => ei != null) as Array<EquippedItem>)
			.map(ei => ei.getFailedProfessionRequirements(professions))
			.flat();
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => Gear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => Gear.enchantToDB(ei.enchant!))),
			gems: distinct(equippedItems.map(ei => ei.curGems().map(gem => Gear.gemToDB(gem))).flat()),
		});
	}
}

/**
 * Represents a item swap gear set, including items/enchants/gems.
 *
 * This is an immutable type.
 */
export class ItemSwapGear extends BaseGear {

	constructor() {
		super({});
	}

	getItemSlots(): ItemSlot[] {
		return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged];
	}

	equipItem(slot: ItemSlot, equippedItem: EquippedItem | null, canDualWield2H: boolean) {
		if (equippedItem) {
			this.removeUniqueGems(this.gear, equippedItem);
			this.removeUniqueItems(this.gear, equippedItem);
		}
		
		this.gear[slot] = equippedItem;
		this.validateWeaponCombo(this.gear, slot, canDualWield2H);
	}

	toProto(): ItemSwap {
		return ItemSwap.create({
			mhItem: this.gear[ItemSlot.ItemSlotMainHand]?.asSpec(),
			ohItem: this.gear[ItemSlot.ItemSlotOffHand]?.asSpec(),
			rangedItem: this.gear[ItemSlot.ItemSlotRanged]?.asSpec(),
		})
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => ItemSwapGear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => ItemSwapGear.enchantToDB(ei.enchant!))),
			gems: distinct(equippedItems.map(ei => ei.curGems().map(gem => ItemSwapGear.gemToDB(gem))).flat()),
		});
	}
}
