package jsonschema

import (
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/xeipuuv/gojsonschema"
)

func RegisterKSUIDFormat() {
	gojsonschema.FormatCheckers.Add("ksuid", ksuidFormatChecker{})
}

type ksuidFormatChecker struct{}

func (f ksuidFormatChecker) IsFormat(input any) bool {
	str, ok := input.(string)
	if !ok {
		return false
	}

	_, err := ksuid.Parse(str)
	return err == nil
}
