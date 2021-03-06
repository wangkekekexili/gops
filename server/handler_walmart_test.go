package server

import (
	"testing"

	"github.com/wangkekekexili/gops/model"
)

func TestWalmartHandler_extractNameAndCondition(t *testing.T) {
	w := &WalmartHandler{}
	tests := []struct {
		fullname  string
		ok        bool
		name      string
		condition string
	}{
		{
			fullname: "PlayStation 4 Slim 500GB Uncharted 4 Bundle",
			ok:       false,
		},
		{
			fullname:  "Overwatch Game of the Year Edition (PlayStation 4)",
			ok:        true,
			name:      "Overwatch Game of the Year Edition",
			condition: model.ProductConditionNew,
		},
		{
			fullname:  "Minecraft (PS4)",
			ok:        true,
			name:      "Minecraft",
			condition: model.ProductConditionNew,
		},
		{
			fullname:  "Watch Dogs (PS4) - Pre-Owned",
			ok:        true,
			name:      "Watch Dogs",
			condition: model.ProductConditionPreowned,
		},
		{
			fullname:  "Tom Clancy's The Division - Pre-Owned (PS4)",
			ok:        true,
			name:      "Tom Clancy's The Division",
			condition: model.ProductConditionPreowned,
		},
	}
	for _, test := range tests {
		nameGot, conditionGot, okGot := w.extractNameAndCondition(test.fullname)
		if test.ok && !okGot {
			t.Fatalf("expected to parse %v successfully", test.fullname)
		}
		if !test.ok {
			continue
		}
		if test.name != nameGot {
			t.Fatalf("expected to get name %v; got %v", test.name, nameGot)
		}
		if test.condition != conditionGot {
			t.Fatalf("expected to get condition %v; got %v", test.condition, conditionGot)
		}
	}
}
