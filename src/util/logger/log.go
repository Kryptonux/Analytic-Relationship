package logger

import (
	"fmt"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func LogMessage(colorStart, colorEnd, title string, message string) {
	now := time.Now().Format("15:04:05")

	maxNowLength := 8
	maxTitleLength := 7
	now = padString(now, maxNowLength)
	title = padString(title, maxTitleLength)

	fmt.Printf("[%s] %s[%s]%s %s\n", now, colorStart, title, colorEnd, message)
}

func Crash(message string, err error) {
	LogMessage(colorRed, colorReset, "ERROR", message)
	if err != nil {
		fmt.Printf("%s%s%s\n", colorRed, err.Error(), colorReset)
	}
	panic(err)
}

func Warning(message string) {
	LogMessage(colorYellow, colorReset, "WARNING", message)
}

func Info(message string) {
	LogMessage("", colorReset, "INFO", message)
}

func padString(input string, maxLength int) string {
	if len(input) >= maxLength {
		return input
	}
	padding := strings.Repeat(" ", maxLength-len(input))
	return input + padding
}
