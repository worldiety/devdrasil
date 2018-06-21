//just a simple structured logger
package log

import (
	"fmt"
)

type Level int

const Emergency Level = 0
const Alert Level = 1
const Critical Level = 2
const Error Level = 3
const Warning Level = 4
const Notice Level = 5
const Informational Level = 6
const Debug Level = 7

const FieldLevel = "level"
const FieldTime = "time"
const FieldMessage = "msg"
const FieldCallerFile = "file"
const FieldCallerLine = "line"
const FieldError = "error"

var Default = NewConsoleLogger()

//A very lean structured logger design
type Logger struct {
	//the level to log (inclusive), ordered by the definition of syslog. Required to optimize performance. Entries out of log level range are neither decorated nor appended.
	Level Level

	//Decorators to enrich log fields
	Decorators []Decorator

	//Appends which captures log fields
	Appenders []Appender
}

type Fields map[string]interface{}

func New(msg string) Fields {
	return Fields{}.SetMessage(msg)
}

//sets the level field
func (e Fields) SetLevel(level Level) Fields {
	e[FieldLevel] = level
	return e
}

//sets the time field
func (e Fields) SetTime(unix int64) Fields {
	e[FieldTime] = unix
	return e
}

func (e Fields) SetError(err error) Fields {
	e[FieldError] = err
	return e
}

//sets the message field
func (e Fields) SetMessage(msg string) Fields {
	e[FieldMessage] = msg
	return e
}

//a simple put builder/like pattern
func (e Fields) Put(key string, value interface{}) Fields {
	e[key] = value
	return e
}

//creates a new info log
func (l *Logger) Info(fields Fields) {
	l.add(Informational, fields)
}

//creates a new error log
func (l *Logger) Error(fields Fields) {
	l.add(Error, fields)
}

//creates a new warn log
func (l *Logger) Warn(fields Fields) {
	l.add(Warning, fields)
}

//adds an entry
func (l *Logger) Add(level Level, fields Fields) {
	l.add(level, fields)
}

//adds a new log entry, if level >= Logger.Level
func (l *Logger) add(level Level, fields Fields) {
	if l.Level >= level {
		for _, dec := range l.Decorators {
			fields.SetLevel(level)
			err := dec.Decorate(fields)
			if err != nil {
				fmt.Println("failed to decorate:", err)
			}
		}

		for _, appender := range l.Appenders {
			err := appender.Append(fields)
			if err != nil {
				fmt.Println("failed to append:", err)
			}
		}
	}
}

//A simple default configuration, ready to use
func NewConsoleLogger() *Logger {
	return &Logger{Level: Informational, Appenders: []Appender{NewConsoleWriteAppender(&TextWriter{})}, Decorators: []Decorator{&TimeDecorator{}, &CallerDecorator{}}}
}
