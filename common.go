package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var Tricks []func(*Env) error

type Env struct {
	Provider string
	OS       string
	Arch     string
	User     *user.User
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

	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot get current user: %s", err)
	}

	return &Env{
		Provider: p,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		User:     u,
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
	return RunWithEnv(os.Environ(), cmd, args...)
}

func RunWithEnv(env []string, cmd string, args ...string) error {
	fmt.Println("Run:", cmd, strings.Join(args, " "))
	ctx, cancel := GetTimeoutContext(DefaultTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, cmd, args...)

	// XXX: Workaround for https://github.com/golang/go/issues/23019
	// Otherwise, go test will hang forever on Go 1.10.
	// Workaround found on https://go-review.googlesource.com/c/go/+/42271/3/misc/android/go_android_exec.go#36
	c.Stderr = struct{ io.Writer }{os.Stderr}
	c.Stdout = struct{ io.Writer }{os.Stdout}

	c.Env = env
	return c.Run()
}

func GetTimeoutContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func WaitForTCP(addr string) error {
	for i := 0; i < 30; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("could not connect to: %s", addr)
}
