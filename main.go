package llog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelName = map[Level]string{
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
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

var currentLevel Level

//TODO: Improve Logging

func SetLogLevel(level Level) {
	currentLevel = level
}

func Debug(msg any, a ...any) {
	if currentLevel > LevelDebug {
		return
	}
	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelDebug],
		" ",
		message,
		"\n",
	)
}

func DebugWithStack(msg any, a ...any) {
	if currentLevel > LevelDebug {
		return
	}

	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelDebug],
		" ",
		stackLoc(2),
		" ",
		message,
		"\n",
	)
}

func Info(msg any, a ...any) {
	if currentLevel > LevelInfo {
		return
	}
	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelInfo],
		" ",
		message,
		"\n",
	)
}

func Warn(msg any, a ...any) {
	if currentLevel > LevelWarn {
		return
	}
	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelWarn],
		" ",
		stackLoc(2),
		" ",
		Yellow,
		message,
		reset,
		"\n",
	)
}

func Error(msg any, a ...any) {
	if currentLevel > LevelError {
		return
	}
	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelError],
		" ",
		stackLoc(2),
		" ",
		Red,
		message,
		reset,
		"\n",
	)
}

// Recieve an Error with a possible Nil value. It will only log if err != nil
// TODO: Add an optional attribute to add custom messages to the error
func ErrNil(err error) (errNotNil bool) {
	if currentLevel > LevelError {
		return
	}
	if err != nil {

		printStdout(
			timestamp(),
			" ",
			levelNameFormatted[LevelError],
			" ",
			stackLoc(2),
			" ",
			Red,
			err.Error(),
			reset,
			"\n",
		)

		return true
	}
	return false
}

func Fatal(msg any, a ...any) {
	if currentLevel > LevelError {
		return
	}
	format := fmt.Sprint(msg)
	message := fmt.Sprintf(format, a...)
	printStdout(
		timestamp(),
		" ",
		levelNameFormatted[LevelFatal],
		" ",
		stackLoc(2),
		" ",
		bold,
		Red,
		message,
		reset,
		"\n",
	)

	//Exit
	os.Exit(1)
}

func FatalNil(err error) (errNotNil bool) {
	if currentLevel > LevelError {
		return
	}
	if err != nil {
		printStdout(
			timestamp(),
			" ",
			levelNameFormatted[LevelFatal],
			" ",
			stackLoc(2),
			" ",
			bold,
			Red,
			err.Error(),
			reset,
			"\n",
		)

		//Exit
		os.Exit(1)
	}
	return false
}

// TODO: Add an argument adding spaces between components
func printStdout(components ...any) {
	//Printing to Stdout
	_, err := fmt.Fprint(os.Stdout, components...)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func timestamp() string {
	return string(DarkGray) + time.Now().Format("2006/01/02 15:04:05")
}

func stackLoc(skip int) string {
	cwd, _ := os.Getwd()
	cwd += "/"
	_, file, line, _ := runtime.Caller(skip)
	fileLocal := strings.TrimPrefix(file, cwd)

	return string(DarkGray) + fileLocal + ":" + strconv.Itoa(line)
}
