package console

import (
	"fmt"

	"github.com/fatih/color"
)

const prefix = "💬"

func Success(msg string) {
	color.Set(color.FgGreen)
	fmt.Printf("%s [OK ] %s\n", prefix, msg)
	color.Unset()
}

func SuccessF(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	Success(msg)
}

func Info(msg string) {
	color.Set(color.FgBlue)
	fmt.Printf("%s [INF] %s\n", prefix, msg)
	color.Unset()
}

func InfoF(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	Info(msg)
}

func Error(err error) {
	color.Set(color.FgRed)
	defer color.Unset()

	fmt.Printf("%s [ERR] %v\n", prefix, err)
}
