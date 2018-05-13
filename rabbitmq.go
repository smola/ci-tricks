package main

import (
	"fmt"
	"os"
	"time"
)

func init() {
	Tricks = append(Tricks, InstallRabbit)
}

const RabbitVersionKey = "RABBITMQ_VERSION"

func InstallRabbit(env *Env) error {
	version := os.Getenv(RabbitVersionKey)
	if version == "" {
		return nil
	}

	fmt.Println("Installing RabbitMQ", version)

	if version != "any" {
		return fmt.Errorf("Setting specific RabbitMQ version is not supported")
	}

	var err error
	switch env.Provider {
	case Appveyor:
		err = installRabbitAppveyor(env, version)
	case Travis:
		err = installRabbitTravis(env, version)
	default:
		err = fmt.Errorf("unsupported provider: %s", env.Provider)
	}

	if err == nil {
		err = WaitForTCP("127.0.0.1:5672")
	}

	if err != nil {
		return fmt.Errorf("error starting RabbitMQ: %s", err)
	}

	return nil
}

func installRabbitAppveyor(env *Env, version string) error {
	if env.OS != Windows {
		return fmt.Errorf("RabbitMQ is supported only on Windows")
	}

	if env.Arch != Amd64 {
		return fmt.Errorf("RabbitMQ is supported only on amd64")
	}

	path := fmt.Sprintf(
		`C:\Program Files\erl9.2\bin;C:\ProgramData\chocolatey\bin;%s`,
		os.Getenv("PATH"))
	environ := os.Environ()
	environ = append(environ, fmt.Sprintf("PATH=%s", path))

	if err := RunWithEnv(environ,
		"choco", "install", "rabbitmq", "--ignoredependencies", "-y",
	); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	return nil
}

func installRabbitTravis(env *Env, version string) error {
	switch env.OS {
	case Linux:
		return installRabbitTravisLinux(env, version)
	case OSX:
		return installRabbitTravisOSX(env, version)
	default:
		return fmt.Errorf("RabbitMQ is only supported on Linux and macOS")
	}
}

func installRabbitTravisLinux(env *Env, version string) error {
	return Run("sudo", "service", "rabbitmq-server", "start")
}

func installRabbitTravisOSX(env *Env, version string) error {
	// High Sierra does not come with /usr/local/sbin by default,
	// which makes linking some Homebrew packages fail.
	if err := Run("sudo", "mkdir", "-p", "/usr/local/sbin"); err != nil {
		return err
	}

	if err := Run("sudo", "chown", "-R",
		fmt.Sprintf("%s:admin", env.User.Username),
		"/usr/local/sbin"); err != nil {
		return err
	}

	if err := Run("brew", "install", "rabbitmq"); err != nil {
		return err
	}

	return Run("brew", "services", "start", "rabbitmq")
}
