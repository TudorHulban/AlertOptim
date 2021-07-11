package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	source = "source.yaml"
	target = "target.yaml"
)

func TestSimple(t *testing.T) {
	a, errNew := NewOptim(source)
	require.Nil(t, errNew)

	f, errFile := os.Create(target)
	require.Nil(t, errFile)
	defer f.Close()

	a.Spool(f)

	require.Nil(t, sameContent(source, target))
}
