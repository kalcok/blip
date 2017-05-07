package logger

const (
	DEFAULT = iota
	ERROR = iota
	WARNING = iota
	INFO = iota
	DEBUG = iota
)

type Logger interface {
	Log(level int, msg string)
	SetLevel(level int)
	Level() int
	LevelString() string
}

type baseLogger struct {
	level int
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

func (logger *baseLogger) LevelString() (level_name string){
	switch logger.level {
	case ERROR:
		level_name = "ERROR"
	case WARNING:
		level_name = "WARNING"
	case INFO:
		level_name = "INFO"
	case DEBUG:
		level_name = "DEBUG"
	default:
		level_name = "CUSTOM"
		
	}
	return
}