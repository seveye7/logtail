package main

import (
	"fmt"
	"os"
)

// Fatal(v ...interface{})
// Fatalf(format string, v ...interface{})
// Fatalln(v ...interface{})
// Panic(v ...interface{})
// Panicf(format string, v ...interface{})
// Panicln(v ...interface{})
// Print(v ...interface{})
// Printf(format string, v ...interface{})
// Println(v ...interface{})

func parseArgs(args ...any) []any {
	var result []any
	for i, arg := range args {
		result = append(result, fmt.Sprintf("v[%v]", i), arg)
	}
	return result
}

type LogSwaper struct{}

func (l *LogSwaper) Fatal(v ...interface{}) {
	textLogger.Error("", parseArgs(v)...)
	os.Exit(1)
}

func (l *LogSwaper) Fatalf(format string, v ...interface{}) {
	textLogger.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *LogSwaper) Fatalln(v ...interface{}) {
	textLogger.Error("", parseArgs(v)...)
	os.Exit(1)
}

func (l *LogSwaper) Panic(v ...interface{}) {
	textLogger.Error("", parseArgs(v)...)
	panic(fmt.Sprint(v...))
}

func (l *LogSwaper) Panicf(format string, v ...interface{}) {
	textLogger.Error(fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func (l *LogSwaper) Panicln(v ...interface{}) {
	textLogger.Error("", parseArgs(v)...)
	panic(fmt.Sprint(v...))
}

func (l *LogSwaper) Print(v ...interface{}) {
	textLogger.Info("", parseArgs(v)...)
}

func (l *LogSwaper) Printf(format string, v ...interface{}) {
	textLogger.Info(fmt.Sprintf(format, v...))
}

func (l *LogSwaper) Println(v ...interface{}) {
	textLogger.Info("", parseArgs(v)...)
}
