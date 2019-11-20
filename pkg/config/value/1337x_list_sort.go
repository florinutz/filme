package value

import (
	"fmt"
	"strings"
)

// region local_types
// region criteria
type sortCriteria uint8

const (
	SortCriteriaTime sortCriteria = iota
	SortCriteriaSize
	SortCriteriaSeeders
	SortCriteriaLeechers
)

var possibleCriteriaValues = map[sortCriteria]string{
	SortCriteriaTime:     "time",
	SortCriteriaSize:     "size",
	SortCriteriaSeeders:  "seeders",
	SortCriteriaLeechers: "leechers",
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
	SortOrderAsc sortOrder = iota
	SortOrderDesc
)

var possibleOrderValues = map[sortOrder]string{
	SortOrderAsc:  "asc",
	SortOrderDesc: "desc",
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

type LeetxListSortValue struct {
	Criteria sortCriteria
	Order    sortOrder
}

func (v *LeetxListSortValue) String() string {
	return fmt.Sprintf("%s-%s", v.Criteria, v.Order)
}

// endregion
// endregion

func (v *LeetxListSortValue) Set(value string) error {
	incomingValue, err := NewLeetxListSortValue(value)
	if err != nil {
		return fmt.Errorf("can't set value: %w", err)
	}
	*v = *incomingValue
	return nil
}

func (*LeetxListSortValue) Type() string {
	return "sort-pair"
}

func GetAllLeetxListSortValues() (values []string) {
	for _, cv := range possibleCriteriaValues {
		for _, ov := range possibleOrderValues {
			values = append(values, fmt.Sprintf("%s-%s", cv, ov))
		}
	}
	return
}

func NewLeetxListSortValue(value string) (*LeetxListSortValue, error) {
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

	return &LeetxListSortValue{*criteria, *order}, nil
}
