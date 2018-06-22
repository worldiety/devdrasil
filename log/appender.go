package log

import (
	"strconv"
	"time"
	"encoding/json"
	"io"
	"fmt"
	"os"
	"strings"
)

//An Appender adds the given log fields to something, e.g. serializes it into json and pushes it into a remote socket
type Appender interface {
	//Appends the fields, returning potential errors
	Append(fields Fields) error
}

//The WriteAppender is a specialization which writes into a settable writer
type WriteAppender interface {
	Appender
	SetWriter(writer io.Writer)
}

//A TextWriter just writes values using %v formatting directive
type TextWriter struct {
	Writer io.Writer
}

func (f *TextWriter) SetWriter(writer io.Writer) {
	f.Writer = writer
}

//simply writes %v values and escapes " with \"
func (f *TextWriter) Append(fields Fields) error {
	if val, ok := fields[FieldTime].(int64); ok {
		f.Writer.Write([]byte(time.Unix(val, 0).Format("2006-01-02T15:04:05-0700") + " "))
	}
	if val, ok := fields[FieldLevel].(Level); ok {
		str := ""
		switch val {
		case Emergency:
			str = "Y"
		case Alert:
			str = "A"
		case Critical:
			str = "C"
		case Error:
			str = "E"
		case Warning:
			str = "W"
		case Notice:
			str = "N"
		case Informational:
			str = "I"
		case Debug:
			str = "D"
		default:
			strconv.Itoa(int(val))
		}
		f.Writer.Write([]byte(str + " "))
	}

	for key, value := range fields {
		switch key {
		case FieldLevel:
			continue
		case FieldTime:
			continue
		case FieldCallerLine:
			continue
		case FieldCallerFile:
			continue
		default:
			objStr := fmt.Sprintf("%v", value)
			tmp := fmt.Sprintf("%s=\"%s\" ", key, strings.Replace(objStr, "\"", "\\\"", -1))
			_, err := f.Writer.Write([]byte(tmp))
			if err != nil {
				return err
			}
		}

	}

	if val, ok := fields[FieldCallerFile].(string); ok {
		f.Writer.Write([]byte(val))
	}

	if val, ok := fields[FieldCallerLine].(int); ok {
		f.Writer.Write([]byte(":" + strconv.Itoa(val)))
	}

	_, err := f.Writer.Write([]byte("\n"));
	return err
}

//A JSONWriter marshals the given fields into an io.Writer
type JSONWriter struct {
	Writer io.Writer
}

func (f *JSONWriter) SetWriter(writer io.Writer) {
	f.Writer = writer
}

//simply formats known fields and such into a common json output
func (f *JSONWriter) Append(fields Fields) error {
	out := make(map[string]interface{})
	for key, value := range fields {
		switch key {
		case FieldLevel:
			if level, ok := value.(Level); ok {
				str := ""
				switch level {
				case Emergency:
					str = "emergency"
				case Alert:
					str = "alert"
				case Critical:
					str = "critical"
				case Error:
					str = "error"
				case Warning:
					str = "warning"
				case Notice:
					str = "notice"
				case Informational:
					str = "informational"
				case Debug:
					str = "debug"
				default:
					strconv.Itoa(int(level))
				}
				out["level"] = str
			}
		case FieldTime:
			if unix, ok := value.(int64); ok {
				//nice listing in https://programming.guide/go/format-parse-string-time-date-example.html
				out["time"] = time.Unix(unix, 0).Format("2006-01-02T15:04:05-0700")
			}
		case FieldMessage:
			out["msg"] = value
		default:
			out[key] = value
		}

	}

	tmp, err := json.Marshal(out)
	if err != nil {
		return err
	}
	if f.Writer == nil {
		return fmt.Errorf("no writer")
	}
	_, err = f.Writer.Write(tmp)
	return err
}

//A simple writer which the given fields into an io.Writer
type FileWriteAppender struct {
	delegate WriteAppender
	file     *os.File
}

func (w *FileWriteAppender) Append(fields Fields) error {
	return w.delegate.Append(fields)
}

func (w *FileWriteAppender) Close() error {
	return w.file.Close()
}

//opens a new file writer with append option
func NewFileWriteAppender(fname string, delegate WriteAppender) (*FileWriteAppender, error) {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &FileWriteAppender{delegate, file}, nil
}

//a simple console appender
type ConsoleWriteAppender struct {
	delegate WriteAppender
}

func (w *ConsoleWriteAppender) Append(fields Fields) error {
	return w.delegate.Append(fields)
}

//opens a new console writer with append option
func NewConsoleWriteAppender(delegate WriteAppender) *ConsoleWriteAppender {
	delegate.SetWriter(os.Stdout)
	return &ConsoleWriteAppender{delegate}
}
