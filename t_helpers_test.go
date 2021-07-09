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
