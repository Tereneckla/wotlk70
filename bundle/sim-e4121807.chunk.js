import{A as e,f4 as t,eR as a,f5 as n,eS as l,dY as s,P as i,h as o,f6 as d,cb as r,cc as p,dU as m,K as c,f7 as u,E as f,cj as S,cm as h,co as g,bD as I,a2 as y,ab as v,F as T,T as R,a5 as b,aD as O,bn as A,w as P,B as w,aE as H}from"./detailed_results-7b150079.chunk.js";import{k as W,b as D,c as E,y as k,i as B,B as M,I as x,O as C,T as J,W as N,X as j,Y as F,Z as G,$ as V,_ as L,t as z}from"./individual_sim_ui-7ca50b32.chunk.js";const q=W({fieldName:"customRotation",numColumns:2,values:[{actionId:e.fromSpellId(53408),value:t.JudgementOfWisdom},{actionId:e.fromSpellId(27180),value:t.HammerOfWrath},{actionId:e.fromSpellId(27173),value:t.Consecration},{actionId:e.fromSpellId(27139),value:t.HolyWrath},{actionId:e.fromSpellId(27138),value:t.Exorcism},{actionId:e.fromSpellId(32700),value:t.AvengersShield},{actionId:e.fromSpellId(53595),value:t.HammerOfTheRighteous},{actionId:e.fromSpellId(27179),value:t.HolyShield}]}),U={inputs:[D({fieldName:"hammerFirst",label:"Open with HotR",labelTooltip:"Open with Hammer of the Righteous instead of Shield of Righteousness in the standard rotation. Recommended for AoE."}),D({fieldName:"squeezeHolyWrath",label:"Squeeze Holy Wrath",labelTooltip:"Squeeze a Holy Wrath cast during sufficiently hasted GCDs (Bloodlust) in the standard rotation."}),E({fieldName:"waitSlack",label:"Max Wait Time (ms)",labelTooltip:"Maximum time in milliseconds to prioritize waiting for next Hammer/Shield to maintain 969. Affects standard and custom priority."}),D({fieldName:"useCustomPrio",label:"Use custom priority",labelTooltip:"Deviates from the standard 96969 rotation, using the priority configured below. Will still attempt to keep a filler GCD between Hammer and Shield."}),q]},_=k({fieldName:"aura",label:"Aura",values:[{name:"None",value:a.NoPaladinAura},{name:"Devotion Aura",value:a.DevotionAura},{name:"Retribution Aura",value:a.RetributionAura}]}),K=k({fieldName:"seal",label:"Seal",labelTooltip:"The seal active before encounter",values:[{name:"Vengeance",value:n.Vengeance},{name:"Command",value:n.Command}]}),Y=k({fieldName:"judgement",label:"Judgement",labelTooltip:"Judgement debuff you will use on the target during the encounter.",values:[{name:"Wisdom",value:l.JudgementOfWisdom},{name:"Light",value:l.JudgementOfLight}]}),X=B({fieldName:"useAvengingWrath",label:"Use Avenging Wrath"}),Z={name:"Baseline Example",data:s.create({talentsString:"-05005135200132311333312321-511302012003",glyphs:{major1:i.GlyphOfSealOfVengeance,major2:i.GlyphOfRighteousDefense,major3:i.GlyphOfDivinePlea,minor1:o.GlyphOfSenseUndead,minor2:o.GlyphOfLayOnHands,minor3:o.GlyphOfBlessingOfKings}})},$=d.create({hammerFirst:!1,squeezeHolyWrath:!0,waitSlack:300,useCustomPrio:!1,customRotation:r.create({spells:[p.create({spell:t.HammerOfTheRighteous}),p.create({spell:t.HolyShield}),p.create({spell:t.HammerOfWrath}),p.create({spell:t.Consecration}),p.create({spell:t.AvengersShield}),p.create({spell:t.JudgementOfWisdom}),p.create({spell:t.Exorcism})]})}),Q={name:"Default (969)",rotation:m.create({specRotationOptionsJson:d.toJsonString(d.create({})),rotation:c.fromJsonString('{\n\t\t\t"type": "TypeAPL",\n\t\t\t"prepullActions": [\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":48952}}},"doAtValue":{"const":{"val":"-3s"}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":54428}}},"doAtValue":{"const":{"val":"-1500ms"}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}\n\t\t\t],\n\t\t\t"priorityList": [\n\t\t\t\t{"action":{"autocastOtherCooldowns":{}}},\n\t\t\t\t{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"spellTimeToReady":{"spellId":{"spellId":53595}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":61411}}}},\n\t\t\t\t{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"spellTimeToReady":{"spellId":{"spellId":61411}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":53595}}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":48806}}}},\n\t\t\t\t{"action":{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":61411}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":53595}}}}},{"cmp":{"op":"OpLe","lhs":{"min":{"vals":[{"spellTimeToReady":{"spellId":{"spellId":61411}}},{"spellTimeToReady":{"spellId":{"spellId":53595}}}]}},"rhs":{"const":{"val":"350ms"}}}}]}},"wait":{"duration":{"min":{"vals":[{"spellTimeToReady":{"spellId":{"spellId":61411}}},{"spellTimeToReady":{"spellId":{"spellId":53595}}}]}}}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":48819}}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":48952}}}},\n\t\t\t\t{"action":{"castSpell":{"spellId":{"spellId":53408}}}},\n\t\t\t\t{"action":{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":61411}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":53595}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48819}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48952}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":53408}}}}}]}},"wait":{"duration":{"min":{"vals":[{"spellTimeToReady":{"spellId":{"spellId":61411}}},{"spellTimeToReady":{"spellId":{"spellId":53595}}},{"spellTimeToReady":{"spellId":{"spellId":48819}}},{"spellTimeToReady":{"spellId":{"spellId":48952}}},{"spellTimeToReady":{"spellId":{"spellId":53408}}}]}}}}}\n\t\t\t]\n\t\t}')})},ee=u.create({aura:a.RetributionAura,judgement:l.JudgementOfWisdom}),te=f.create({flask:S.FlaskOfStoneblood,food:h.FoodDragonfinFilet,defaultPotion:g.IndestructiblePotion,prepopPotion:g.IndestructiblePotion}),ae={name:"Preraid Preset",tooltip:M,enableWhen:e=>!0,gear:I.fromJsonString('{"items": [\n\t\t{"id":42549,"enchant":3818,"gems":[41396,49110]},\n\t\t{"id":40679},\n\t\t{"id":37635,"enchant":3852,"gems":[40015]},\n\t\t{"id":44188,"enchant":3605},\n\t\t{"id":39638,"enchant":1953,"gems":[36767,40089]},\n\t\t{"id":37682,"enchant":3850,"gems":[0]},\n\t\t{"id":39639,"enchant":3860,"gems":[36767,0]},\n\t\t{"id":37379,"enchant":3601,"gems":[40022,40008]},\n\t\t{"id":37292,"enchant":3822,"gems":[40089]},\n\t\t{"id":44243,"enchant":3606},\n\t\t{"id":37186},\n\t\t{"id":37257},\n\t\t{"id":44063,"gems":[36767,40015]},\n\t\t{"id":37220},\n\t\t{"id":37179,"enchant":2673},\n\t\t{"id":43085,"enchant":3849},\n\t\t{"id":40707}\n\t]}')},ne={name:"P1 Preset",tooltip:M,enableWhen:e=>!0,gear:I.fromJsonString('{"items": [\n\t\t{"id":40581,"enchant":3818,"gems":[41380,36767]},\n\t\t{"id":40387},\n\t\t{"id":40584,"enchant":3852,"gems":[40008]},\n\t\t{"id":40410,"enchant":3605},\n\t\t{"id":40579,"enchant":3832,"gems":[36767,40022]},\n\t\t{"id":39764,"enchant":3850,"gems":[0]},\n\t\t{"id":40580,"enchant":3860,"gems":[40008,0]},\n\t\t{"id":39759,"enchant":3601,"gems":[40008,40008]},\n\t\t{"id":40589,"enchant":3822},\n\t\t{"id":39717,"enchant":3606,"gems":[40089]},\n\t\t{"id":40718},\n\t\t{"id":40107},\n\t\t{"id":44063,"gems":[36767,40089]},\n\t\t{"id":37220},\n\t\t{"id":40345,"enchant":3788},\n\t\t{"id":40400,"enchant":3849},\n\t\t{"id":40707}\n\t]}')},le={name:"P2 Preset",tooltip:M,enableWhen:e=>!0,gear:I.fromJsonString('{\n      "items": [\n        {"id":46175,"enchant":3818,"gems":[41380,40088]},\n        {"id":45485,"gems":[40088]},\n        {"id":46177,"enchant":3852,"gems":[40034]},\n        {"id":45496,"enchant":3605,"gems":[40034]},\n        {"id":46039,"enchant":3832,"gems":[36767,36767]},\n        {"id":45111,"enchant":3850,"gems":[0]},\n        {"id":45487,"enchant":3860,"gems":[40008,40008,0]},\n        {"id":45825,"enchant":3601,"gems":[40008]},\n        {"id":45594,"enchant":3822,"gems":[40034,45880,40088]},\n        {"id":45988,"enchant":3606,"gems":[40008,40008]},\n        {"id":45471,"gems":[40088]},\n        {"id":45326},\n        {"id":45158},\n        {"id":46021},\n        {"id":45947,"enchant":3788,"gems":[40088]},\n        {"id":45587,"enchant":3849,"gems":[36767]},\n        {"id":45145}\n      ]\n    }')};class se extends x{constructor(e,t){super(e,t,{cssClass:"protection-paladin-sim-ui",cssScheme:"paladin",knownIssues:[],epStats:[y.StatStamina,y.StatStrength,y.StatAgility,y.StatAttackPower,y.StatMeleeHit,y.StatSpellHit,y.StatMeleeCrit,y.StatExpertise,y.StatMeleeHaste,y.StatArmorPenetration,y.StatSpellPower,y.StatArmor,y.StatBonusArmor,y.StatDefense,y.StatBlock,y.StatBlockValue,y.StatDodge,y.StatParry,y.StatResilience,y.StatNatureResistance,y.StatShadowResistance,y.StatFrostResistance],epPseudoStats:[v.PseudoStatMainHandDps],epReferenceStat:y.StatSpellPower,displayStats:[y.StatHealth,y.StatArmor,y.StatBonusArmor,y.StatStamina,y.StatStrength,y.StatAgility,y.StatAttackPower,y.StatMeleeHit,y.StatMeleeCrit,y.StatMeleeHaste,y.StatExpertise,y.StatArmorPenetration,y.StatSpellPower,y.StatSpellHit,y.StatDefense,y.StatBlock,y.StatBlockValue,y.StatDodge,y.StatParry,y.StatResilience,y.StatNatureResistance,y.StatShadowResistance,y.StatFrostResistance],modifyDisplayStats:e=>{let t=new T;return R.freezeAllAndDo((()=>{e.getMajorGlyphs().includes(i.GlyphOfSealOfVengeance)&&e.getSpecOptions().seal==n.Vengeance&&(t=t.addStat(y.StatExpertise,10*b))})),{talents:t}},defaults:{gear:le.gear,epWeights:T.fromMap({[y.StatArmor]:.07,[y.StatBonusArmor]:.06,[y.StatStamina]:1.14,[y.StatStrength]:1,[y.StatAgility]:.62,[y.StatAttackPower]:.26,[y.StatExpertise]:.69,[y.StatMeleeHit]:.79,[y.StatMeleeCrit]:.3,[y.StatMeleeHaste]:.17,[y.StatArmorPenetration]:.04,[y.StatSpellPower]:.13,[y.StatBlock]:.52,[y.StatBlockValue]:.28,[y.StatDodge]:.46,[y.StatParry]:.61,[y.StatDefense]:.54},{[v.PseudoStatMainHandDps]:3.33}),consumes:te,rotation:$,talents:Z.data,specOptions:ee,raidBuffs:O.create({giftOfTheWild:A.TristateEffectImproved,powerWordFortitude:A.TristateEffectImproved,strengthOfEarthTotem:A.TristateEffectImproved,arcaneBrilliance:!0,unleashedRage:!0,leaderOfThePack:A.TristateEffectRegular,icyTalons:!0,totemOfWrath:!0,swiftRetribution:!0,moonkinAura:A.TristateEffectRegular,sanctifiedRetribution:!0,manaSpringTotem:A.TristateEffectRegular,bloodlust:!0,thorns:A.TristateEffectImproved,devotionAura:A.TristateEffectImproved,shadowProtection:!0}),partyBuffs:P.create({}),individualBuffs:w.create({blessingOfKings:!0,blessingOfSanctuary:!0,blessingOfWisdom:A.TristateEffectImproved,blessingOfMight:A.TristateEffectImproved}),debuffs:H.create({judgementOfWisdom:!0,judgementOfLight:!0,misery:!0,faerieFire:A.TristateEffectImproved,ebonPlaguebringer:!0,totemOfWrath:!0,shadowMastery:!0,bloodFrenzy:!0,mangle:!0,exposeArmor:!0,sunderArmor:!0,vindication:!0,thunderClap:A.TristateEffectImproved,insectSwarm:!0})},playerIconInputs:[],rotationInputs:U,includeBuffDebuffInputs:[C],excludeBuffDebuffInputs:[],otherInputs:{inputs:[J,N,j,F,G,V,L,_,X,Y,K,z]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[Z],rotations:[Q],gear:[ae,ne,le]}})}}export{$ as D,Z as G,se as P,ee as a,te as b,ne as c,le as d};
//# sourceMappingURL=sim-e4121807.chunk.js.map
