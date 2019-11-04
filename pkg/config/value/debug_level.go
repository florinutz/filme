package value

import "github.com/sirupsen/logrus"

type DebugLevelValue struct {
	level logrus.Level
}

func (v *DebugLevelValue) String() string {
	return v.level.String()
}

func (v *DebugLevelValue) Set(lvl string) (err error) {
	v.level, err = logrus.ParseLevel(lvl)
	return
}

func (*DebugLevelValue) Type() string {
	return "debugLevel"
}

func GetAllLevels() (result []string) {
	for _, level := range logrus.AllLevels {
		result = append(result, level.String())
	}

	return
}
