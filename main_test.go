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

	//Test formatter
	t.Run("Formatter", Formatter)

	//Test Base Functions
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

const (
	resetRegex     string = `\x1b\[0m`
	messageRegex   string = "Testing" + resetRegex
	timestampRegex string = `\x1b\[90m\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}` + resetRegex
	stackRegex     string = `\x1b\[90mmain_test\.go:\d+` + resetRegex
)

var debugRegex string = strings.ReplaceAll(levelNameFormatted[LevelDebug], `[`, `\[`)
var infoRegex string = strings.ReplaceAll(levelNameFormatted[LevelInfo], `[`, `\[`)
var warnRegex string = strings.ReplaceAll(levelNameFormatted[LevelWarn], `[`, `\[`)
var errorRegex string = strings.ReplaceAll(levelNameFormatted[LevelError], `[`, `\[`)
var fatalRegex string = strings.ReplaceAll(levelNameFormatted[LevelFatal], `[`, `\[`)

func Formatter(t *testing.T) {
	var result string
	result = formatMessage("Hey", 0)
	if result != "Hey0" {
		t.Fail()
	}

	result = formatMessage("Hey %d", 0)
	if result != "Hey 0" {
		t.Fail()
	}
}

func BaseFunctions(t *testing.T) {
	// create a buffer
	var buf bytes.Buffer

	// overwrite the default writer
	stdout = &buf

	RunLogFunctions()

	lines := strings.Split(buf.String(), "\n")

	//Check expected Output
	var currentLine int = 0

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
	// create a buffer
	var buf bytes.Buffer

	// overwrite the default writer
	stdout = &buf

	t.Log("Running with LogLevel Debug")
	RunLogFunctions()

	lines := strings.Split(buf.String(), "\n")

	//Check expected Output
	var currentLine int = 0

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

	// create a new buffer
	buf = *bytes.NewBuffer([]byte{})

	// overwrite the writer
	stdout = &buf

	SetLogLevel(LevelInfo)
	t.Log("Running with LogLevel Info")
	RunLogFunctions()

	lines = strings.Split(buf.String(), "\n")

	//Check expected Output
	currentLine = 0

	// timestamp (YYYY/MM/DD HH:MM:SS) + message
	reg = regexp.MustCompile(timestampRegex + " Testing")
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected print log line not recieved")
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

	// create a new buffer
	buf = *bytes.NewBuffer([]byte{})

	// overwrite the writer
	stdout = &buf

	SetLogLevel(LevelWarn)
	t.Log("Running with LogLevel Warn")
	RunLogFunctions()

	lines = strings.Split(buf.String(), "\n")

	//Check expected Output
	currentLine = 0

	// timestamp (YYYY/MM/DD HH:MM:SS) + message
	reg = regexp.MustCompile(timestampRegex + " Testing")
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected print log line not recieved")
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

	// create a new buffer
	buf = *bytes.NewBuffer([]byte{})

	// overwrite the writer
	stdout = &buf

	SetLogLevel(LevelError)
	t.Log("Running with LogLevel Error")
	RunLogFunctions()

	lines = strings.Split(buf.String(), "\n")

	//Check expected Output
	currentLine = 0

	// timestamp (YYYY/MM/DD HH:MM:SS) + message
	reg = regexp.MustCompile(timestampRegex + " Testing")
	if !reg.MatchString(lines[currentLine]) {
		t.Errorf("expected print log line not recieved")
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

	//reset log Level
	SetLogLevel(LevelDebug)
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

	// create a buffer
	buf := *bytes.NewBuffer([]byte{})

	// overwrite writer
	stdout = &buf

	// Use a deferred function to recover from panic
	defer func() {
		r := recover() // Catch panic

		if r == nil {
			// This line will not be reached if Fatal works correctly
			t.Log("Expected Fatal to panic, but it did not.")
			t.Fail()
		} else if r != testingMessage {
			t.Errorf("Unexpected panic message: %v", r)
		} else {
			// Validate Logging
			lines := strings.Split(buf.String(), "\n")
			var currentLine int = 0

			// timestamp + FATAL + filename:line + message
			reg := regexp.MustCompile(timestampRegex + " " + fatalRegex + " " + stackRegex + " " + `\x1b\[1m\x1b\[31m` + testingMessage)
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

			stdout = os.Stdout
		}
	}()

	// Call the Fatal function, which should panic
	Fatal(testingMessage)

	stdout = os.Stdout
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
