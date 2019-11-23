package value

import "github.com/sirupsen/logrus"

type DebugLevelValue struct {
	logrus.Level
}

func (v *DebugLevelValue) Set(lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	v.Level = level

	return nil
}

func (*DebugLevelValue) Type() string {
	return "level"
}

func GetAllLevels() (result []string) {
	for _, level := range logrus.AllLevels {
		result = append(result, level.String())
	}

	return
}
