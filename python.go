package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Tricks = append(Tricks, InstallPython)
}

const PythonVersionKey = "PYTHON_VERSION"

func InstallPython(env *Env) error {
	version := os.Getenv(PythonVersionKey)
	if version == "" {
		return nil
	}

	fmt.Println("Installing Python", version)
	if hasPythonVersion(version) {
		fmt.Printf("Python %s already installed", version)
		return nil
	}

	var err error
	switch env.Provider {
	case Appveyor:
		err = installPythonAppveyor(env, version)
	case Travis:
		err = installPythonTravis(env, version)
	default:
		err = fmt.Errorf("unsupported provider: %s", env.Provider)
	}

	if err == nil {
		err = Run("python", "--version")
	}

	if err != nil {
		return fmt.Errorf("error installing Python: %s", err)
	}

	return nil
}

func installPythonAppveyor(env *Env, version string) error {
	if env.OS != Windows {
		return fmt.Errorf("Python is supported only on Windows")
	}

	// Based on https://packaging.python.org/guides/supporting-windows-using-appveyor/
	path := `C:\Python`
	path += strings.Replace(version, ".", "", -1)
	if env.Arch == Amd64 {
		path += "-x64"
	}

	if err := Run("setx", "PYTHON", path, "/m"); err != nil {
		return err
	}

	if err := Run("setx", "PATH",
		fmt.Sprintf("%s;%s", path, os.Getenv("PATH")),
		"/m"); err != nil {
		return err
	}

	return nil
}

func installPythonTravis(env *Env, version string) error {
	switch env.OS {
	case Linux:
		return installPythonTravisLinux(env, version)
	case OSX:
		return installPythonTravisOSX(env, version)
	default:
		return fmt.Errorf("Python is only supported on Linux and macOS")
	}
}

func installPythonTravisLinux(env *Env, version string) error {
	if err := Run("sudo", "service", "postgresql", "stop"); err != nil {
		return err
	}

	if err := Run("sudo", "service", "postgresql", "start", version); err != nil {
		return err
	}

	return nil
}

var pythonMacOSVersions = map[string]string{
	"2.7": "https://www.python.org/ftp/python/2.7.15/python-2.7.15-macosx10.9.pkg",
	"3.6": "https://www.python.org/ftp/python/3.6.5/python-3.6.5-macosx10.9.pkg",
}

func installPythonTravisOSX(env *Env, version string) error {
	url, ok := pythonMacOSVersions[version]
	if !ok {
		return fmt.Errorf("Python %s not supported", version)
	}

	tmp := os.TempDir()
	path := filepath.Join(tmp, "python.pkg")
	if err := Run("wget", "-qO", path, url); err != nil {
		return err
	}

	if err := Run("sudo", "installer", "-pkg", path, "-target", "/"); err != nil {
		return err
	}

	return nil
}

func hasPythonVersion(version string) bool {
	actualVersion := pythonVersion()
	if actualVersion == "" {
		return false
	}

	actualVersion = strings.Replace(actualVersion, "Python ", "", 1)
	return strings.HasPrefix(actualVersion, version+".")
}

func pythonVersion() string {
	cmd := exec.Command("python", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
