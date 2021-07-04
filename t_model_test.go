package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	a, errNew := NewSimple("source.yaml")
	require.Nil(t, errNew)

	a.Spool(os.Stdout)
}
