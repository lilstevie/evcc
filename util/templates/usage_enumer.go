// Code generated by "enumer -type Usage -trimprefix Usage -transform=lower"; DO NOT EDIT.

package templates

import (
	"fmt"
	"strings"
)

const _UsageName = "gridpvbatterycharge"

var _UsageIndex = [...]uint8{0, 4, 6, 13, 19}

const _UsageLowerName = "gridpvbatterycharge"

func (i Usage) String() string {
	if i < 0 || i >= Usage(len(_UsageIndex)-1) {
		return fmt.Sprintf("Usage(%d)", i)
	}
	return _UsageName[_UsageIndex[i]:_UsageIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _UsageNoOp() {
	var x [1]struct{}
	_ = x[UsageGrid-(0)]
	_ = x[UsagePV-(1)]
	_ = x[UsageBattery-(2)]
	_ = x[UsageCharge-(3)]
}

var _UsageValues = []Usage{UsageGrid, UsagePV, UsageBattery, UsageCharge}

var _UsageNameToValueMap = map[string]Usage{
	_UsageName[0:4]:        UsageGrid,
	_UsageLowerName[0:4]:   UsageGrid,
	_UsageName[4:6]:        UsagePV,
	_UsageLowerName[4:6]:   UsagePV,
	_UsageName[6:13]:       UsageBattery,
	_UsageLowerName[6:13]:  UsageBattery,
	_UsageName[13:19]:      UsageCharge,
	_UsageLowerName[13:19]: UsageCharge,
}

var _UsageNames = []string{
	_UsageName[0:4],
	_UsageName[4:6],
	_UsageName[6:13],
	_UsageName[13:19],
}

// UsageString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func UsageString(s string) (Usage, error) {
	if val, ok := _UsageNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _UsageNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Usage values", s)
}

// UsageValues returns all values of the enum
func UsageValues() []Usage {
	return _UsageValues
}

// UsageStrings returns a slice of all String values of the enum
func UsageStrings() []string {
	strs := make([]string, len(_UsageNames))
	copy(strs, _UsageNames)
	return strs
}

// IsAUsage returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Usage) IsAUsage() bool {
	for _, v := range _UsageValues {
		if i == v {
			return true
		}
	}
	return false
}
