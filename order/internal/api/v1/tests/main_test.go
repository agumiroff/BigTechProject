package tests

import (
	"os"
	"testing"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func TestMain(m *testing.M) {
	// Initialize logger for all tests
	logger.SetNopLogger()

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}
