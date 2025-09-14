package llog

import (
	"testing"
)

func TestMain(t *testing.T) {
	//Running Main Tests
	t.Log("Running Main Tests:")

	t.Run("Base Functions", BaseFunctions)

	//Test Levels
	SetLogLevel(LevelDebug)
	t.Run("Base Functions", LogLevels)

}

func BaseFunctions(t *testing.T) {
	RunLogFunctions()
}

func LogLevels(t *testing.T) {
	SetLogLevel(LevelDebug)
	t.Log("Running with LogLevel Debug")
	RunLogFunctions()
	SetLogLevel(LevelInfo)
	t.Log("Running with LogLevel Info")
	RunLogFunctions()
	SetLogLevel(LevelWarn)
	t.Log("Running with LogLevel Warn")
	RunLogFunctions()
	SetLogLevel(LevelError)
	t.Log("Running with LogLevel Error")
	RunLogFunctions()
}

func RunLogFunctions() {
	Debug("Testing")
	DebugWithStack("Testing")
	Info("Testing")
	Warn("Testing")
	Error("Testing")
}
