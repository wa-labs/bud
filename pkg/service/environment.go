package service

import (
	"os"
	"strconv"
)

// EnvString looks for an environment var that is a string and if it doesn't
// exist it returns the fallback variable passed to it.
func EnvString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

// EnvInt looks for an environment var that is a int and if it doesn't
// exist it returns the fallback variable passed to it.
func EnvInt(env string, fallback int) int {
	e := os.Getenv(env)
	i, err := strconv.Atoi(e)
	if e == "" || err != nil {
		return fallback
	}
	return i
}
