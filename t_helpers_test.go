package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const toChop = "xxx.1"

func TestChop(t *testing.T) {
	require.Error(t, chopLastRow(toChop, ""))
	require.Nil(t, chopLastRow(toChop, "2\n"))
}

func TestPos(t *testing.T) {
	require.Equal(t, -1, startPos(""))
	require.Equal(t, 0, startPos("x  "))
	require.Equal(t, 2, startPos("  x  "))
}
