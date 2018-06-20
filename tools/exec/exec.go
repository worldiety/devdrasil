package tools

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"fmt"
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

//executes and blocks until the program exists. This will by-pass any configured logger. All stdout and errout is captured into the returned buffer
func (env *Env) ExecCapture(cmdName string, cmdArgs ...string) ([]byte, error) {

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = env.Dir
	cmd.Env = make([]string, 0)
	for _, kv := range os.Environ() {
		cmd.Env = append(cmd.Env, kv)
	}

	for key, value := range env.Variables {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	buf, err := cmd.CombinedOutput()

	if err != nil {
		env.Logf(ERROR, "Error failed to execute cmd: %s", err)
		return buf, err
	}

	return buf, nil
}

//just like ExecPipe but interprets everything as a line
func (env *Env) ExecLines(cmdName string, cmdArgs ...string) ([]string, error) {
	buf, err := env.ExecCapture(cmdName, cmdArgs...)
	lines := strings.Split(string(buf), "\n")
	return lines, err
}

//the default execution
func (env *Env) Exec(onNewStdOut func(string), onNewErrOut func(string), cmdName string, cmdArgs ...string) error {
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
		return fmt.Errorf()
		err
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

//creates a multi-line string which is bash compatible
func (env *Env) CreateBashCommand(cmdName string, cmdArgs ...string) string {
	sb := &strings.Builder{}
	sb.WriteString("#!/bin/bash")
	sb.WriteString("cd ")
	sb.WriteString(env.Dir)
	sb.WriteString("\n")

	for key, value := range env.Variables {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteString("\n")
	}

	sb.WriteString(cmdName)
	for _, arg := range cmdArgs {
		sb.WriteString(" ")
		sb.WriteString(arg)
	}
	return sb.String()
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
	LineNumbers int
}

func (w *WriterDispatcherLogFileView) Write(p []byte) (n int, err error) {
	//w.LogFileView.AddLine(&api.NumberedLine{Line: w.LineNumbers, Text: string(p)})
	w.LineNumbers++
	return w.Delegate.Write(p)
}

type Lines struct {
	Strings []string
}
