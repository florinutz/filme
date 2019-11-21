package sort

import (
	"fmt"
	"strings"
)

// region local_types
// region criteria
type sortCriteria uint8

const (
	CriteriaTime sortCriteria = iota
	CriteriaSize
	CriteriaSeeders
	CriteriaLeechers
)

var possibleCriteriaValues = map[sortCriteria]string{
	CriteriaTime:     "time",
	CriteriaSize:     "size",
	CriteriaSeeders:  "seeders",
	CriteriaLeechers: "leechers",
}

func (c sortCriteria) String() string {
	return possibleCriteriaValues[c]
}

func getSortCriteriaFromStr(criteria string) *sortCriteria {
	for c, v := range possibleCriteriaValues {
		if criteria == v {
			return &c
		}
	}
	return nil
}

// endregion

// region order
type sortOrder uint8

const (
	OrderAsc sortOrder = iota
	OrderDesc
)

var possibleOrderValues = map[sortOrder]string{
	OrderAsc:  "asc",
	OrderDesc: "desc",
}

func getSortOrderFromStr(order string) *sortOrder {
	for o, v := range possibleOrderValues {
		if order == v {
			return &o
		}
	}
	return nil
}

func (o sortOrder) String() string {
	return possibleOrderValues[o]
}

type Value struct {
	Criteria sortCriteria
	Order    sortOrder
}

func (v *Value) String() string {
	return fmt.Sprintf("%s-%s", v.Criteria, v.Order)
}

// endregion
// endregion

func (v *Value) Set(value string) error {
	incomingValue, err := NewValue(value)
	if err != nil {
		return fmt.Errorf("can't set value: %w", err)
	}
	*v = *incomingValue
	return nil
}

func (*Value) Type() string {
	return "sort-pair"
}

func GetAllValues() (values []string) {
	for _, cv := range possibleCriteriaValues {
		for _, ov := range possibleOrderValues {
			values = append(values, fmt.Sprintf("%s-%s", cv, ov))
		}
	}
	return
}

func NewValue(value string) (*Value, error) {
	err := fmt.Errorf("value '%s' is not a valid criteria-order sort pair", value)

	spl := strings.Split(value, "-")
	if len(spl) != 2 {
		return nil, err
	}

	var criteria *sortCriteria
	if criteria = getSortCriteriaFromStr(spl[0]); criteria == nil {
		return nil, fmt.Errorf("invalid sort criteria in value %s", value)
	}

	var order *sortOrder
	if order = getSortOrderFromStr(spl[1]); order == nil {
		return nil, fmt.Errorf("invalid sort order in value '%s'", value)
	}

	return &Value{*criteria, *order}, nil
}
