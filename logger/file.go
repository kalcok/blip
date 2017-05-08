package logger

import (
	"fmt"
	"os"
	"time"
)

type fileLogger struct{
	baseLogger
	logFile *os.File
	logFilePath string
}

func NewFileLogger(log_path string, level int) (logger *fileLogger){
	logger = new(fileLogger)
	logger.SetLevel(level)
	logger.SetLogFile(log_path)
	return
}

func (logger *fileLogger)Log(level int, msg string){
	now := time.Now().Format(time.RFC3339)
	msg = fmt.Sprintf("%s - %s - %s\n", logger.LevelItoa(level), now, msg)
	if level <= logger.level{
		logger.logFile.WriteString(msg)
	}
}

func (logger *fileLogger) SetLogFile(path string){
	file, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil{
		panic(fmt.Sprintf("Failed to open logging file. %s", err))
	}
	logger.logFilePath = path
	logger.logFile = file
}

func (logger *fileLogger) Close(){
	if logger.logFile != nil{
		logger.logFile.Close()
	}
}