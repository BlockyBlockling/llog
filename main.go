package llog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelDebugWithStack
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPrint
)

var levelName = map[Level]string{
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
	LevelPrint: "Print",
}

var levelNameFormatted = map[Level]string{
	LevelDebug: bold + string(Blue) + "DEBU" + reset,
	LevelInfo:  string(LightGreen) + "INFO" + reset,
	LevelWarn:  string(Yellow) + "WARN" + reset,
	LevelError: string(Red) + "ERR" + reset,
	LevelFatal: bold + string(Red) + "FATAL" + reset,
}

type Color string

const (
	reset = "\033[0m"
	bold  = "\033[1m"

	Black        Color = "\033[30m"
	Red          Color = "\033[31m"
	Green        Color = "\033[32m"
	Yellow       Color = "\033[33m"
	Blue         Color = "\033[34m"
	Magenta      Color = "\033[35m"
	Cyan         Color = "\033[36m"
	LightGray    Color = "\033[37m"
	DarkGray     Color = "\033[90m"
	LightRed     Color = "\033[91m"
	LightGreen   Color = "\033[92m"
	LightYellow  Color = "\033[93m"
	LightBlue    Color = "\033[94m"
	LightMagenta Color = "\033[95m"
	LightCyan    Color = "\033[96m"
	White        Color = "\033[97m"
)

var stdout io.Writer = os.Stdout
var currentLevel Level

//TODO: Write README.md
//TODO: Improve Logging
//TODO: Improve/Centralize Logger to prevent code iterations
//TODO: Add Multiline indented Logging via Custom Function NextLine() to be implemented into the Loggers

func SetLogLevel(level Level) {
	currentLevel = level
}

func GetLevelByName(name string) (Level, error) {
	switch name {
	case levelName[LevelDebug]:
		return LevelDebug, nil
	case levelName[LevelInfo]:
		return LevelInfo, nil
	case levelName[LevelWarn]:
		return LevelWarn, nil
	case levelName[LevelError]:
		return LevelError, nil
	case levelName[LevelFatal]:
		return LevelFatal, nil
	default:
		return 0, errors.New("Level not found")
	}
}

func Print(msg any, a ...any) {
	printStdout(formatLogLevel(LevelPrint, msg, a...))
}

func Debug(msg any, a ...any) {
	if showLevel(LevelDebug) {
		printStdout(formatLogLevel(LevelDebug, msg, a...))
	}
}

func DebugWithStack(msg any, a ...any) {
	if showLevel(LevelDebug) {
		printStdout(formatLogLevel(LevelDebugWithStack, msg, a...))
	}
}

func Info(msg any, a ...any) {
	if showLevel(LevelInfo) {
		printStdout(formatLogLevel(LevelInfo, msg, a...))
	}
}

func Warn(msg any, a ...any) {
	if showLevel(LevelWarn) {
		printStdout(formatLogLevel(LevelWarn, msg, a...))
	}
}

func Error(msg any, a ...any) {
	if showLevel(LevelError) {
		printStdout(formatLogLevel(LevelError, msg, a...))
	}
}

// Recieve an Error with a possible Nil value. It will only log if err != nil
// TODO: Add an optional attribute to add custom messages to the error
func ErrNil(err error) (errNotNil bool) {
	if err != nil {
		if showLevel(LevelError) {
			printStdout(formatLogLevel(LevelError, err.Error()))
		}
		return true
	}

	return false
}

func Fatal(msg any, a ...any) {
	if showLevel(LevelFatal) {
		format := fmt.Sprint(msg)
		printStdout(formatLogLevel(LevelFatal, format, a...))

		//Exit
		os.Exit(2) //using the same exit code as panic
	}
}

func FatalNil(err error) (errNotNil bool) {
	if showLevel(LevelFatal) && err != nil {
		printStdout(formatLogLevel(LevelFatal, err.Error()))

		//Exit
		os.Exit(2) //using the same exit code as panic
	}
	return false
}

func PrintNoNewLine(level Level, msg any, a ...any) {
	if showLevel(level) {
		printStdout("\r", strings.TrimSuffix(formatLogLevel(level, msg, a...), "\n"))
	}
}

func ReplaceLine(level Level, msg any, a ...any) {
	if showLevel(level) {
		printStdout("\r", strings.TrimSuffix(formatLogLevel(level, msg, a...), "\n"))
	}
}

func formatLogLevel(level Level, msg any, a ...any) string {
	stackLocatorIndex := 3
	message := formatMessage(msg, a...)
	switch level {
	case LevelDebug:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelDebug],
			" ",
			message,
			reset,
			"\n",
		)
	case LevelDebugWithStack:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelDebug],
			" ",
			stackLoc(stackLocatorIndex),
			" ",
			message,
			reset,
			"\n",
		)
	case LevelInfo:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelInfo],
			" ",
			message,
			reset,
			"\n",
		)
	case LevelWarn:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelWarn],
			" ",
			stackLoc(stackLocatorIndex),
			" ",
			Yellow,
			message,
			reset,
			"\n",
		)
	case LevelError:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelError],
			" ",
			stackLoc(stackLocatorIndex),
			" ",
			Red,
			message,
			reset,
			"\n",
		)
	case LevelFatal:
		return fmt.Sprint(
			timestamp(),
			" ",
			levelNameFormatted[LevelFatal],
			" ",
			stackLoc(3),
			" ",
			bold,
			Red,
			message,
			reset,
			"\n",
		)
	}

	//Default print
	return fmt.Sprint(message, reset, "\n")
}

func showLevel(level Level) bool {
	return currentLevel <= level

}

func formatMessage(msg any, a ...any) string {
	// Regular expression to match fmt directives
	re := regexp.MustCompile(`%[\d.]*[sdvfeg]`)

	msgString, ok := msg.(string)
	if ok {
		if re.MatchString(msgString) {
			return fmt.Sprintf(msgString, a...)
		}
	}
	a = append([]any{msg}, a...)
	return fmt.Sprint(a...)
}

// TODO: Add an argument adding spaces between components
func printStdout(components ...any) {
	//Printing to Stdout
	_, err := fmt.Fprint(stdout, components...)
	if err != nil {
		panic("Failed to print to Stdout")
	}
}

func timestamp() string {
	return string(DarkGray) + time.Now().Format("2006/01/02 15:04:05") + reset
}

func stackLoc(skip int) string {
	cwd, _ := os.Getwd()
	cwd += "/"
	_, file, line, _ := runtime.Caller(skip)
	fileLocal := strings.TrimPrefix(file, cwd)

	return string(DarkGray) + fileLocal + ":" + strconv.Itoa(line) + reset
}
