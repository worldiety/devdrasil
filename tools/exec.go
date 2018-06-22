package tools

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"sync"
	"github.com/worldiety/devdrasil/log"
	"bytes"
	"io"
)

type LogLevel int

const INFO LogLevel = 0
const ERROR LogLevel = 1

//this field contains the executed command as a bash script for debugging purposes in case of errors
const FieldBashCMD = "bashcmd"

type Env struct {
	Logger *log.Logger
	//Dir is the working directory
	Dir string

	Variables map[string]string
}

//creates a default environment
func NewEnv() *Env {
	env := &Env{}
	env.Logger = log.Default
	env.Variables = make(map[string]string)
	return env
}

//executes and blocks until the program exists. This will by-pass any configured logger.
// All stdout and errout is captured into the returned buffer
// This is useful if you don't need any progress/want to be as fast as possible/just process piped binary data
func (env *Env) ExecBytes(cmdName string, cmdArgs ...string) ([]byte, error) {
	env.Logger.Info(log.New(env.dumpCmd(cmdName, cmdArgs...)))

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
		env.Logger.Error(log.New("cmd failed").SetError(err).Put(FieldBashCMD, env.CreateBashCommand(cmdName, cmdArgs...)))
		return buf, err
	}

	return buf, nil
}

//just like Exec and logs everything properly but returns an array of lines of the emitted stdout and errout
func (env *Env) ExecLines(cmdName string, cmdArgs ...string) ([]string, error) {
	buf := &bytes.Buffer{}
	err := env.Exec(func(std string) {
		buf.WriteString(std)
	}, func(e string) {
		buf.WriteString(e)
	}, cmdName, cmdArgs...)
	lines := strings.Split(string(buf.Bytes()), "\n")
	//purge empty line from split
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, err
}

//the default execution
func (env *Env) Exec(onNewStdOut func(string), onNewErrOut func(string), cmdName string, cmdArgs ...string) error {
	env.Logger.Info(log.New(env.dumpCmd(cmdName, cmdArgs...)))
	wg := sync.WaitGroup{}
	wg.Add(1)

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
		env.Logger.Error(log.New("stdoutpipe failed").SetError(err).Put(FieldBashCMD, env.CreateBashCommand(cmdName, cmdArgs...)))
		return err
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		env.Logger.Error(log.New("stderrpipe failed").SetError(err).Put(FieldBashCMD, env.CreateBashCommand(cmdName, cmdArgs...)))
		return err
	}

	scannerStd := bufio.NewReader(cmdStdReader)
	go func() {
		for {
			txt, err := scannerStd.ReadString('\n')
			if err != nil {
				break
			}
			env.Logger.Info(log.New(txt))
			onNewStdOut(txt)

		}

		wg.Done()
	}()

	scannerErr := bufio.NewReader(cmdErrReader)
	go func() {
		for {
			txt, err := scannerErr.ReadString('\n')
			if err != nil {
				break
			}
			env.Logger.Info(log.New(txt))
			onNewStdOut(txt)

		}
	}()

	err = cmd.Start()
	if err != nil {
		env.Logger.Error(log.New("cmd failed").SetError(err).Put(FieldBashCMD, env.CreateBashCommand(cmdName, cmdArgs...)))
		return err
	}

	err = cmd.Wait()
	if err != nil {
		env.Logger.Error(log.New("cmd failed").SetError(err).Put(FieldBashCMD, env.CreateBashCommand(cmdName, cmdArgs...)))
		return err
	}

	wg.Wait()
	return nil
}

func (env *Env) dumpCmd(cmdName string, cmdArgs ...string) string {
	sb := &strings.Builder{}
	sb.WriteString(cmdName)
	for _, arg := range cmdArgs {
		sb.WriteString(" ")
		sb.WriteString(arg)
	}
	return sb.String()
}

//creates a multi-line string which is bash compatible
func (env *Env) CreateBashCommand(cmdName string, cmdArgs ...string) string {
	sb := &strings.Builder{}
	sb.WriteString("#\n")
	sb.WriteString("#!/bin/bash\n")
	sb.WriteString("# this is not the script which has been executed, but an approximation\n")
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
	sb.WriteString("\n")
	return sb.String()
}

// strip '\n' or read until EOF, return error if read error
func readline(reader io.Reader) (line []byte, err error) {
	line = make([]byte, 0, 100)
	for {
		b := make([]byte, 1)
		n, er := reader.Read(b)
		if n > 0 {
			c := b[0]
			if c == '\n' { // end of line
				break
			}
			line = append(line, c)
		}
		if er != nil {
			err = er
			return
		}
	}
	return
}
