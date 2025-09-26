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

	//Test Fatal
	t.Run("Fatal Functions", TestFatal)
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

	//TODO: Add test for ErrNil
}

func RunLogFunctions() {
	Debug("Testing")
	DebugWithStack("Testing")
	Info("Testing")
	Warn("Testing")
	Error("Testing")
}

func TestFatal(t *testing.T) {
	testingMessage := "Testing Fatal log"

	// Use a deferred function to recover from panic
	defer func() {
		r := recover() // Catch panic

		if r == nil {
			// This line will not be reached if Fatal works correctly
			t.Log("Expected Fatal to panic, but it did not.")
			t.Fail()
		} else if r != testingMessage {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	// Call the Fatal function, which should panic
	Fatal(testingMessage)
}
