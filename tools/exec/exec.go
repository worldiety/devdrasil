package tools

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"worldiety.net/scm/git/prj-ops-controller.git/api"
	"sync"
)

type LogLevel int

const INFO LogLevel = 0
const ERROR LogLevel = 1

type Env struct {
	consoleLoggerInfo  *log.Logger
	consoleLoggerError *log.Logger
	//Dir is the working directory
	Dir string

	Variables map[string]string
}

func NewEnv() *Env {
	env := &Env{}
	env.consoleLoggerInfo = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.consoleLoggerError = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.Variables = make(map[string]string)
	return env
}

func NewLogCapturedEnv() (env *Env, buf *strings.Builder) {
	env = &Env{}
	buf = &strings.Builder{}
	dispatcher := &WriterDispatcher{}
	dispatcher.Writer = append(dispatcher.Writer, os.Stdout, buf)
	env.consoleLoggerInfo = log.New(dispatcher, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.consoleLoggerError = log.New(dispatcher, "ERROR ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.Variables = make(map[string]string)
	return env, buf
}

func NewLogViewCapturedEnv() (env *Env, buf *WriterDispatcherLogFileView) {
	env = &Env{}
	dispatcher := &WriterDispatcherLogFileView{}
	dispatcher.Delegate = os.Stdout
	dispatcher.LogFileView = &api.LogFileView{}
	env.consoleLoggerInfo = log.New(dispatcher, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.consoleLoggerError = log.New(dispatcher, "ERROR ", log.Ldate|log.Ltime|log.Lmicroseconds)
	env.Variables = make(map[string]string)
	return env, dispatcher
}

func (e *Env) Logf(level LogLevel, format string, args ...interface{}) {
	switch level {
	case INFO:
		e.consoleLoggerInfo.Printf(format, args...)
	case ERROR:
		e.consoleLoggerError.Printf(format, args...)
	}
}

func (e *Env) Log(level LogLevel, text string) {
	switch level {
	case INFO:
		e.consoleLoggerInfo.Println(text)
	case ERROR:
		e.consoleLoggerError.Println(text)
	}
}

func ExecDump(env *Env, out io.Writer, cmdName string, cmdArgs ...string) error {

	for key, value := range env.Variables {
		env.Logf(INFO, "SET %s=%s", key, value)
	}

	tmp := cmdName
	for _, arg := range cmdArgs {
		tmp += " " + arg
	}
	env.Logf(INFO, tmp)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = env.Dir
	cmd.Env = make([]string, 0)
	for _, kv := range os.Environ() {
		cmd.Env = append(cmd.Env, kv)
	}

	for key, value := range env.Variables {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	cmdStdReader, err := cmd.StdoutPipe()
	if err != nil {
		env.Logf(ERROR, "Error creating StdoutPipe for %s [%s]: %s", cmdName, cmdArgs, err)
		return err
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		env.Logf(ERROR, "Error creating StderrPipe for %s [%s]: %s", cmdName, cmdArgs, err)
		return err
	}

	scannerErr := bufio.NewScanner(cmdErrReader)
	go func() {
		for scannerErr.Scan() {
			txt := scannerErr.Text()
			env.Log(ERROR, txt)
		}
	}()

	err = cmd.Start()
	if err != nil {
		env.Logf(ERROR, "Error starting Cmd: %s", err)
		return err
	}

	b, e := ioutil.ReadAll(cmdStdReader)
	if e != nil {
		env.Logf(ERROR, "Error consuming out for %s [%s]: %s", cmdName, cmdArgs, err)
		return err
	}
	out.Write(b)

	err = cmd.Wait()
	if err != nil {
		env.Logf(ERROR, "Error waiting for Cmd: %s", err)
		return err
	}

	return nil
}

func Exec(env *Env, out *[]string, cmdName string, cmdArgs ...string) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for key, value := range env.Variables {
		env.Logf(INFO, "SET %s=%s", key, value)
	}

	tmp := cmdName
	for _, arg := range cmdArgs {
		tmp += " " + arg
	}
	env.Logf(INFO, tmp)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = env.Dir
	cmd.Env = make([]string, 0)
	for _, kv := range os.Environ() {
		cmd.Env = append(cmd.Env, kv)
	}

	for key, value := range env.Variables {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	cmdStdReader, err := cmd.StdoutPipe()
	if err != nil {
		env.Logf(ERROR, "Error creating StdoutPipe for %s [%s]: %s", cmdName, cmdArgs, err)
		return err
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		env.Logf(ERROR, "Error creating StderrPipe for %s [%s]: %s", cmdName, cmdArgs, err)
		return err
	}

	scannerStd := bufio.NewScanner(cmdStdReader)
	go func() {
		for scannerStd.Scan() {
			txt := scannerStd.Text()
			if out != nil {
				*out = append(*out, txt)
			}
			env.Log(INFO, txt)
		}
		wg.Done()
	}()

	scannerErr := bufio.NewScanner(cmdErrReader)
	go func() {
		for scannerErr.Scan() {
			txt := scannerErr.Text()
			if out != nil {
				*out = append(*out, txt)
			}
			env.Log(ERROR, txt)
		}
	}()

	err = cmd.Start()
	if err != nil {
		env.Logf(ERROR, "Error starting Cmd: %s", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		env.Logf(ERROR, "Error waiting for Cmd: %s", err)
		return err
	}

	wg.Wait()
	return nil
}

type WriterDispatcher struct {
	Writer []io.Writer
}

func (w *WriterDispatcher) Write(p []byte) (n int, err error) {
	for _, writer := range w.Writer {
		n, err = writer.Write(p)
		if err != nil {
			return
		}
	}
	return 0, nil
}

type WriterDispatcherLogFileView struct {
	Delegate    io.Writer
	LogFileView *api.LogFileView
	LineNumbers int
}

func (w *WriterDispatcherLogFileView) Write(p []byte) (n int, err error) {
	w.LogFileView.AddLine(&api.NumberedLine{Line: w.LineNumbers, Text: string(p)})
	w.LineNumbers++
	return w.Delegate.Write(p)
}

type Lines struct {
	Strings []string
}
