package shaman

import (
	"fmt"
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
)

func (shaman *Shaman) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_TotemRemainingTime:
		return shaman.newValueTotemRemainingTime(rot, config.GetTotemRemainingTime())
	default:
		return nil
	}
}

type APLValueTotemRemainingTime struct {
	core.DefaultAPLValueImpl
	shaman    *Shaman
	totemType proto.ShamanTotems_TotemType
}

func (shaman *Shaman) newValueTotemRemainingTime(rot *core.APLRotation, config *proto.APLValueTotemRemainingTime) core.APLValue {
	if config.TotemType == proto.ShamanTotems_TypeUnknown {
		rot.ValidationWarning("Totem Type required.")
		return nil
	}
	return &APLValueTotemRemainingTime{
		shaman:    shaman,
		totemType: config.TotemType,
	}
}
func (value *APLValueTotemRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueTotemRemainingTime) GetDuration(sim *core.Simulation) time.Duration {
	if value.totemType == proto.ShamanTotems_Earth {
		return core.MaxDuration(0, value.shaman.NextTotemDrops[EarthTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Air {
		return core.MaxDuration(0, value.shaman.NextTotemDrops[AirTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Fire {
		return core.MaxDuration(0, value.shaman.NextTotemDrops[FireTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Water {
		return core.MaxDuration(0, value.shaman.NextTotemDrops[WaterTotem]-sim.CurrentTime)
	} else {
		return 0
	}
}
func (value *APLValueTotemRemainingTime) String() string {
	return fmt.Sprintf("Totem Remaining Time(%s)", value.totemType.String())
}
