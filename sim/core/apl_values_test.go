package core

import (
	"testing"
	"time"

	"github.com/Tereneckla/wotlk70/sim/core/proto"
)

func TestValueConst(t *testing.T) {
	sim := &Simulation{}
	unit := &Unit{}

	stringVal := unit.newValueConst(&proto.APLValueConst{Val: "test str"})
	if stringVal.GetString(sim) != "test str" {
		t.Fatalf("Unexpected string value %s", stringVal.GetString(sim))
	}

	intVal := unit.newValueConst(&proto.APLValueConst{Val: "10"})
	if intVal.GetInt(sim) != 10 {
		t.Fatalf("Unexpected int value %d", intVal.GetInt(sim))
	}

	floatVal := unit.newValueConst(&proto.APLValueConst{Val: "10.123"})
	if floatVal.GetFloat(sim) != 10.123 {
		t.Fatalf("Unexpected float value %f", floatVal.GetFloat(sim))
	}

	durVal := unit.newValueConst(&proto.APLValueConst{Val: "10.123s"})
	if durVal.GetDuration(sim) != time.Millisecond*10123 {
		t.Fatalf("Unexpected duration value %s", durVal.GetDuration(sim))
	}

	coercedDurVal := unit.coerceTo(floatVal, proto.APLValueType_ValueTypeDuration)
	if _, ok := coercedDurVal.(*APLValueConst); !ok {
		t.Fatalf("Failed to skip coerce wrapper for duration value")
	}
	if coercedDurVal.GetDuration(sim) != time.Millisecond*10123 {
		t.Fatalf("Unexpected coerced duration value %s", coercedDurVal.GetDuration(sim))
	}
}
