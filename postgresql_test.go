package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestInstallPostgres(t *testing.T) {
	require := require.New(t)

	env, err := GetEnv()
	require.NoError(err)

	version := os.Getenv(PostgresVersionKey)
	if version == "" {
		t.Skip(fmt.Sprintf("%s not set", PostgresVersionKey))
	}

	var passStr string
	if env.Provider == Appveyor {
		passStr = ":Password12!"
	}

	url := fmt.Sprintf("postgres://postgres%s@127.0.0.1:5432?sslmode=disable", passStr)

	db, err := sql.Open("postgres", url)
	require.NoError(err)

	rows, err := db.Query("SELECT version();")
	require.NoError(err)

	var actualVersion string
	require.True(rows.Next())
	require.NoError(rows.Scan(&actualVersion))
	t.Logf("Actual version: %s", actualVersion)
	prefix := fmt.Sprintf("PostgreSQL %s.", version)
	require.True(strings.HasPrefix(actualVersion, prefix))
	require.NoError(db.Close())
}
