package tests

import (
	"fmt"
	"testing"
)

func failOnError(t *testing.T, err error, msg string, args ...interface{}) {
	if err != nil {
		messageWithArgs := msg
		if len(args) > 0 {
			messageWithArgs = fmt.Sprintf(msg, args...)
		}
		t.Fatalf("%s: %s", messageWithArgs, err)
	}
}
