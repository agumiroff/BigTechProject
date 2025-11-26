package integration

import (
	"testing"
)

func TestVariablesExist(t *testing.T) {
	if testEnv == nil {
		t.Log("testEnv is nil at test start - this is expected before BeforeSuite")
	}
	if suiteCtx == nil {
		t.Log("suiteCtx is nil at test start - this is expected before BeforeSuite")
	}
}
