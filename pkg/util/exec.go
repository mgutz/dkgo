package util

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/google/shlex"
)

type CommandResult struct {
	Output []byte `json:"stdout"`
	RC     int    `json:"rc"`
}

func Exec(statement string) (*CommandResult, error) {
	args, err := shlex.Split(statement)
	if err != nil {
		return &CommandResult{[]byte(err.Error()), 1}, nil
	}

	return run(args[0], args[1:]...)
}

func Execf(format string, args ...any) (*CommandResult, error) {
	statement := fmt.Sprintf(format, args...)
	return Exec(statement)
}

func ExecWithStdin(statement string, stdinText string) (*CommandResult, error) {
	args, err := shlex.Split(statement)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(args[0], args[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinText)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		rc := -1
		if err, ok := err.(*exec.ExitError); ok {
			rc = err.ExitCode()
		}
		return &CommandResult{[]byte(err.Error()), rc}, nil
	}

	return &CommandResult{out, 0}, nil
}

func run(name string, args ...string) (*CommandResult, error) {
	cmd := exec.Command(name, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		rc := -1
		if err, ok := err.(*exec.ExitError); ok {
			rc = err.ExitCode()
		}
		return &CommandResult{[]byte(err.Error()), rc}, nil
	}

	return &CommandResult{out, 0}, nil
}

func Bash(statement string) (*CommandResult, error) {
	return run("bash", "-c", statement)
}

func Bashf(format string, args ...any) (*CommandResult, error) {
	statement := fmt.Sprintf(format, args...)
	return Bash(statement)
}
