package logger

import (
	"strings"
	"fmt"
	"time"
	"os"
)

const (
	CURRENT = -1
	DEFAULT = iota
	ERROR = iota
	WARNING = iota
	INFO = iota
	DEBUG = iota
)

type LogEntry struct {
	level int
	msg string
}

type Logger interface {
	Log(level int, msg string)
	log()
	Close()
	Level() int
	SetLevel(level int)
	LevelItoa(level int) string
	LevelAtoi(level string) int
}

type baseLogger struct {
	level int
	log_chan chan LogEntry
}

func (logger *baseLogger) Log(level int, msg string){
	if level <= logger.level{
		// Non-blocking sending to channel
		go func() {
			select {
			case logger.log_chan <- LogEntry{level:level, msg:msg}:
			case <-time.After(time.Second):
				// XXX: Maybe find better way to handle failed logging attempts
				fmt.Fprintln(os.Stderr,
					"Logging timed out. Make sure your logger is properly initialized.")
			}
		}()
	}
}

func (logger *baseLogger) log(){
	panic("Internal log function must be implemented in specialized Logger modules.")
}

func (logger *baseLogger) init(level int){
	logger.SetLevel(level)
	logger.log_chan = make(chan LogEntry, 100)
}

func (logger *baseLogger) Close(){
	// Use this method in subclasses of baseLogger to cleanup logger resources
}

func (logger *baseLogger) Level() int{
	return logger.level
}

func (logger *baseLogger) SetLevel(level int){
	if level == DEFAULT{
		level = INFO
	}
	logger.level = level
}

func (logger *baseLogger) LevelItoa(level int) (level_name string){
	if level == CURRENT{level = logger.level}
	switch level {
	case ERROR:
		level_name = "ERROR"
	case WARNING:
		level_name = "WARNING"
	case INFO, DEFAULT:
		level_name = "INFO"
	case DEBUG:
		level_name = "DEBUG"
	default:
		level_name = "CUSTOM"
		
	}
	return
}

func (logger *baseLogger) LevelAtoi(level string) (level_const int){
	level = strings.ToLower(level)
	switch level {
	case "error":
		level_const = ERROR
	case "warning":
		level_const = WARNING
	case "info", "default":
		level_const = INFO
	case "debug":
		level_const = DEBUG
	case "current":
		level_const = logger.level
	default:
		panic(fmt.Sprintf("Can't convert level string to constant. Unknown level '%s'",level))
	}
	return
}