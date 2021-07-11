package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	toChop     = "xxx.1"
	empty      = "xxx.2"
	line       = "xxx.3"
	lineWSpace = "xxx.4"
)

// func TestChop(t *testing.T) {
// 	require.Error(t, chopLastRow(toChop, ""))
// 	require.Nil(t, chopLastRow(toChop, "2\n"))
// }

func TestPos(t *testing.T) {
	require.Equal(t, -1, startPos(""))
	require.Equal(t, 0, startPos("x  "))
	require.Equal(t, 2, startPos("  x  "))
}

func TestSame(t *testing.T) {
	require.Nil(t, sameContent(toChop, toChop, os.Stdin))
	require.Error(t, sameContent(toChop, empty, os.Stdin))
	require.Error(t, sameContent(line, lineWSpace, os.Stdin))
}
