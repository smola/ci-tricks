package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestInstallRabbit(t *testing.T) {
	version := os.Getenv(RabbitVersionKey)
	if version == "" {
		t.Skip(fmt.Sprintf("%s not set", RabbitVersionKey))
	}

	require := require.New(t)

	c, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	require.NoError(err)
	require.NoError(c.Close())
}
