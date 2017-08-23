package eventlog

import (
	"github.com/Sirupsen/logrus"
	"fmt"
)

type HookStdout struct {
}

func (buf *HookStdout) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (buf *HookStdout) Fire(e *logrus.Entry) error {
	fmt.Printf("[%s] %s\n", e.Level, e.Message)
	return nil
}