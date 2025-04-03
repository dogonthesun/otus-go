package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func UpdateEnviron(env []string, update Environment) []string {
	renv := make(map[string]string, len(env))
	for _, e := range env {
		s := strings.SplitN(e, "=", 2)
		renv[s[0]] = s[1]
	}
	for k, v := range update {
		switch {
		case v.NeedRemove:
			delete(renv, k)
		default:
			renv[k] = v.Value
		}
	}
	updatedEnv := make([]string, 0, len(renv))
	for k, v := range renv {
		updatedEnv = append(updatedEnv, fmt.Sprintf("%s=%s", k, v))
	}
	return updatedEnv
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// #nosec G204
	c := exec.Command(cmd[0], cmd[1:]...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Env = UpdateEnviron(os.Environ(), env)

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); !ok {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			return FailExitCode
		}
		return exitError.ExitCode()
	}

	return 0
}
