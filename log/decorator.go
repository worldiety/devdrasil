package log

import (
	"runtime"
	"time"
)

//A Decorator adds more fields
type Decorator interface {
	//inserts/updates/removes fields. It is guaranteed that @FieldLevel has been set
	Decorate(fields Fields) error
}

type CallerDecorator struct {
}

func (d *CallerDecorator) Decorate(fields Fields) error {
	file, line := GetCaller(5)
	fields[FieldCallerFile] = file
	fields[FieldCallerLine] = line
	return nil
}

type TimeDecorator struct {
}

func (d *TimeDecorator) Decorate(fields Fields) error {
	fields.SetTime(time.Now().Unix())
	return nil
}

// Returns file and line of the current caller. A reasonable offset is
//    1 = (call of runtime.Callers) in this method
//    2 = the line which called "GetCaller(2)"
//	  3 = the line which called the method which called "GetCaller(3)"
//	  ...
func GetCaller(offset int) (file string, line int) {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(offset, fpcs)
	if n == 0 {
		return "n/a", -1
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a", -1
	}

	return fun.FileLine(fpcs[0] - 1)
}
