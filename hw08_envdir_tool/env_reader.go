package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ExtractValueFromString(str string) string {
	return strings.TrimRight(strings.ReplaceAll(str, "\x00", "\n"), " \n\t")
}

func ExtractValueFromFile(path string) (EnvValue, error) {
	fp, err := os.Open(path)
	if err != nil {
		return EnvValue{}, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	if s.Scan() {
		return EnvValue{ExtractValueFromString(s.Text()), false}, nil
	}
	if err := s.Err(); err != nil {
		return EnvValue{}, err
	}
	return EnvValue{"", true}, nil
}

func IsCorrectEnvVarName(str string) bool {
	// Based on https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap08.html
	//
	// Environment variable names used by the utilities
	// in the Shell and Utilities volume of POSIX.1-2024 consist
	// solely of uppercase letters, digits, and
	// the <underscore> ('_') from the characters defined in
	// Portable Character Set and do not begin with a digit.
	isCorrect, err := regexp.MatchString(`^[A-Za-z_]+[A-Za-z0-9_]*$`, str)
	if err != nil {
		panic(err)
	}
	return isCorrect
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !IsCorrectEnvVarName(name) {
			// Ignore
			continue
		}
		envValue, err := ExtractValueFromFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		env[name] = envValue
	}
	return env, nil
}
