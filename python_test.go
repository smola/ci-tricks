package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallPython(t *testing.T) {
	require := require.New(t)

	version := os.Getenv(PythonVersionKey)
	if version == "" {
		t.Skip(fmt.Sprintf("%s not set", PythonVersionKey))
	}

	t.Log(fmt.Sprintf("Python version: %s\n", pythonVersion()))
	require.True(hasPythonVersion(version))
}
