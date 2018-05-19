package logger

import (
	"os"
	"strconv"

	"github.com/fatih/color"
)

// Info log to stdout in colors
func Info(msg ...string) {
	c := color.New(color.FgCyan)
	c.Print("Info: ")
	c.Add(color.Reset)
	c.Println(cleanUpArgs(msg))
	c.Add(color.Reset)
}

// Out no colors
func Out(msg ...string) {
	c := color.New(color.Reset)
	c.Println(cleanUpArgs(msg))
	c.Add(color.Reset)
}

func cleanUpArgs(args []string) string {
	var ret string
	for i := range args {
		ret = ret + args[i] + " "
	}

	return ret
}

// Fatal log to stderr
func Fatal(msg ...string) {
	c := color.New(color.FgRed)
	c.Print("Error: ")
	c.Add(color.Reset)
	c.Println(cleanUpArgs(msg))
	c.Add(color.Reset)
	os.Exit(1)
}

// Warn log to stdout
func Warn(msg ...string) {
	c := color.New(color.FgYellow)
	c.Print("** Warn: ")
	c.Println(cleanUpArgs(msg))
	c.Add(color.Reset)
}

// Debug shown only when DEBUG=true
func Debug(msg ...string) {
	debug, err := strconv.ParseBool(string(os.Getenv("DEBUG")))
	if err != nil {
		return
	} else if debug != true {
		return
	}

	c := color.New(color.FgMagenta)
	c.Print("** Debug: ")
	c.Add(color.Reset)
	c.Println(cleanUpArgs(msg))
	c.Add(color.Reset)
}
