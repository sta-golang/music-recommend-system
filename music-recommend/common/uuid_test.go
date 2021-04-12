package common

import (
	"testing"

	"github.com/sta-golang/go-lib-utils/log"
)

func TestUUID(t *testing.T) {
	log.ConsoleLogger.Debug(UUID())
}
