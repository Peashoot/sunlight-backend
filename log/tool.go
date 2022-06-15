package log

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

func Info(v ...interface{}) {
	log.Println(getOutput("Info", v...))
}

func Infof(format string, args ...interface{}) {
	log.Printf(getOutputf("Info", format), args...)
}

func Debug(v ...interface{}) {
	log.Println(getOutput("Debug", v...))
}

func Debugf(format string, args ...interface{}) {
	log.Printf(getOutputf("Debug", format), args...)
}

func Error(v ...interface{}) {
	log.Println(getOutput("Error", v...))
}

func Errorf(format string, args ...interface{}) {
	log.Printf(getOutputf("Error", format), args...)
}

func Warn(v ...interface{}) {
	log.Println(getOutput("Warn", v...))
}

func Warnf(format string, args ...interface{}) {
	log.Printf(getOutputf("Warn", format), args...)
}

func Fatal(v ...interface{}) {
	log.Fatalln(getOutput("Fatal", v...))
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(getOutputf("Fatal", format), args...)
}

func getOutput(level string, v ...interface{}) []interface{} {
	tmp := []interface{}{"- " + level + " -", fmt.Sprintf("[%d]", unix.Gettid())}
	return append(tmp, v...)
}

func getOutputf(level, format string) string {
	return fmt.Sprintf("- %s - [%d] ", level, unix.Gettid()) + format
}
