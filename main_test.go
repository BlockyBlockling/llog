package llog

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	//Running Main Tests
	t.Log("Running Main Tests:")

	//Test stdout
	t.Run("Stdout Test", StdOutTest)

	t.Run("Base Functions", BaseFunctions)

	//Test Levels
	SetLogLevel(LevelDebug)
	t.Run("GetLevel", GetLevels)
	t.Run("Level", LogLevels)
}

// Bad writer to test edge case of os.Stdout failing
type badWriter struct{}

func (badWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write failed")
}

func StdOutTest(t *testing.T) {
	// create a buffer
	var buf bytes.Buffer

	// overwrite the default write	oldstdout := stdout
	stdout = &buf

	// use your printing function
	printStdout("test")

	//check output
	if buf.String() != "test" {
		t.Log("printStdout failed to print correctly")
		t.Fail()
	}

	//Destroy pipe to cover panic
	stdout = badWriter{}

	// Use a deferred function to recover from panic
	defer func() {
		r := recover() // Catch panic

		if r == nil {
			// This line will not be reached if Fatal works correctly
			t.Log("Expected Fatal to panic, but it did not.")
			t.Fail()
		} else if r != "Failed to print to Stdout" {
			t.Errorf("Unexpected panic message: %v", r)
		}

		//reset stdout
		stdout = os.Stdout
	}()

	printStdout("test")

	//reset stdout
	stdout = os.Stdout

}

func BaseFunctions(t *testing.T) {
	// create a buffer
	var buf bytes.Buffer

	// overwrite the default writer
	stdout = &buf
	buf.Reset()

	RunLogFunctions()

	lines := strings.Split(buf.String(), "\n")

	//Check expected Output
	var currentLine int = 0

	//TODO: Automate colorcodes on regex with the colocodes on the liberary
	//like strings.ReplaceAll(levelNameFormatted[LevelDebug], "[", `\[`)
	resetRegex := `\x1b\[0m`
	messageRegex := "Testing" + resetRegex
	timestampRegex := `\x1b\[90m\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}` + resetRegex
	stackRegex := `\x1b\[90mmain_test\.go:\d+` + resetRegex
	debugRegex := `\x1b\[1m\x1b\[34mDEBU` + resetRegex
	infoRegex := `\x1b\[92mINFO` + resetRegex
	warnRegex := `\x1b\[33mWARN` + resetRegex
	errorRegex := `\x1b\[31mERR` + resetRegex

	// timestamp (YYYY/MM/DD HH:MM:SS) + message
	reg := regexp.MustCompile(timestampRegex + " Testing")
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected print log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	// timestamp + DEBU + message
	reg = regexp.MustCompile(timestampRegex + " " + debugRegex + " " + messageRegex)
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected debug log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	// timestamp + DEBU + filename:line + message
	reg = regexp.MustCompile(timestampRegex + " " + debugRegex + " " + stackRegex + " " + messageRegex)
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected debugstack log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	//timespamp + INFO + message
	reg = regexp.MustCompile(timestampRegex + " " + infoRegex + " " + messageRegex)
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected info log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	// timestamp + WARN + filename:line+ message
	reg = regexp.MustCompile(timestampRegex + " " + warnRegex + " " + stackRegex + " " + `\x1b\[33m` + messageRegex)
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected warn log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	// timestamp + ERR + filename:line + message
	reg = regexp.MustCompile(timestampRegex + " " + errorRegex + " " + stackRegex + " " + `\x1b\[31m` + messageRegex)
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected error log line not recieved")
		t.Log("Regex: " + reg.String())
		t.Logf("%q\n", lines[currentLine])
		clean := reg.ReplaceAllString(lines[currentLine], "")
		t.Errorf("Not matching: %q", clean)
	}
	currentLine++

	if lines[currentLine] == "" {
		//Empty Line (as expected)
		currentLine++
	}

	if currentLine != len(lines) {
		t.Errorf("More lines Printed than analysed. Missing %d line/s.", len(lines)-currentLine)
		t.Log(buf.String())
	}

	//reset stdout
	stdout = os.Stdout

	//Test Fatal
	t.Run("Fatal Functions", TestFatal)

	//Test Err != Nil functions
	t.Run("Err!=Nil Functions", TestNil)
}

func GetLevels(t *testing.T) {
	level, err := GetLevelByName("InvalidLevel")
	if level != 0 || err == nil {
		t.Log("Invalid output on Invalid Level")
		t.Fail()
	}
	level, err = GetLevelByName("Debug")
	if level != LevelDebug || err != nil {
		t.Log("Failed to get debug level by name")
		t.Fail()
	}
	level, err = GetLevelByName("Info")
	if level != LevelInfo || err != nil {
		t.Log("Failed to get info level by name")
		t.Fail()
	}
	level, err = GetLevelByName("Warn")
	if level != LevelWarn || err != nil {
		t.Log("Failed to get warn level by name")
		t.Fail()
	}
	level, err = GetLevelByName("Error")
	if level != LevelError || err != nil {
		t.Log("Failed to get error level by name")
		t.Fail()
	}
	level, err = GetLevelByName("Fatal")
	if level != LevelFatal || err != nil {
		t.Log("Failed to get fatal level by name")
		t.Fail()
	}
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
	Print("Testing")
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

type testerror struct{}

func (m *testerror) Error() string {
	return "testing error"
}

func GetError(isNil bool) error {
	if isNil {
		return nil
	} else {
		return &testerror{}
	}
}

func TestNil(t *testing.T) {
	//ErrNil
	var errNil error
	errNotNil := errors.New("error test")
	ErrNil(errNil)
	ErrNil(errNotNil)

	//FatalNil
	t.Run("FataNil", TestFatalNil)
}

func TestFatalNil(t *testing.T) {
	var errNil error
	errNotNil := errors.New("error test")

	// Call the Fatal function, which should panic
	result := FatalNil(errNil)
	if result {
		//Should be false
		t.Log("Was true should be false")
		t.Fail()
	}

	// Use a deferred function to recover from panic
	defer func() {
		r := recover() // Catch panic

		if r == nil {
			// This line will not be reached if Fatal works correctly
			t.Log("Expected Fatal to panic, but it did not.")
			t.Fail()
		} else if r != "error test" {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	// Call the Fatal function, which should panic
	result = FatalNil(errNotNil)
	if !result {
		//Should be true
		t.Log("Was false should be true")
		t.Fail()
	}
}
