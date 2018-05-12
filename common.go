package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var Tricks []func(*Env) error

type Env struct {
	Provider string
	OS       string
	Arch     string
}

const (
	Travis   = "travis"
	Appveyor = "appveyor"

	Linux   = "linux"
	OSX     = "darwin"
	Windows = "windows"

	Amd64 = "amd64"
	X84   = "386"

	DefaultTimeout = 10 * time.Minute
)

func RunTricks() error {
	env, err := GetEnv()
	if err != nil {
		return err
	}

	for _, trick := range Tricks {
		if err := trick(env); err != nil {
			return err
		}
	}

	return nil
}

func GetEnv() (*Env, error) {
	p, err := GetProvider()
	if err != nil {
		return nil, err
	}

	return &Env{
		Provider: p,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}, nil
}

func GetProvider() (string, error) {
	if os.Getenv("TRAVIS") != "" {
		return Travis, nil
	}

	if os.Getenv("APPVEYOR") != "" {
		return Appveyor, nil
	}

	return "", fmt.Errorf("unrecognized CI provider")
}

func Run(cmd string, args ...string) error {
	fmt.Println("Run:", cmd, strings.Join(args, " "))
	c := exec.CommandContext(GetTimeoutContext(DefaultTimeout), cmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func GetTimeoutContext(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
