package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inPath := "testdata/input.txt"

	t.Run("simple case", func(t *testing.T) {
		outPath := "simple.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, 0)

		require.NoError(t, err)
	})

	t.Run("negative limit", func(t *testing.T) {
		outPath := "negative-limit.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, -1)

		require.ErrorAs(t, err, &ErrInvalidBoundaries)
	})

	t.Run("negative offset", func(t *testing.T) {
		outPath := "negative-offset.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, -1, 0)

		require.ErrorAs(t, err, &ErrInvalidBoundaries)

		_, err = os.Stat(outPath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("too big offset", func(t *testing.T) {
		outPath := "big-offset.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 10000000, 0)
		require.ErrorAs(t, err, &ErrInvalidBoundaries)

		_, err = os.Stat(outPath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("irregular input file", func(t *testing.T) {
		inPath := "testdata"
		outPath := "irrgular.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, 0)
		require.ErrorAs(t, err, &ErrInvalidFileType)

		_, err = os.Stat(outPath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("input doesn't exist", func(t *testing.T) {
		inPath := "foobarxxx_123"
		outPath := "outpath-foobar.txt"
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, 0)
		require.ErrorAs(t, err, &ErrInputOutput{})

		_, err = os.Stat(outPath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("irregular output file", func(t *testing.T) {
		outPath := "testdata"

		err := Copy(inPath, outPath, 0, 0)
		require.ErrorAs(t, err, &ErrInputOutput{})
	})

	t.Run("empty input or output", func(t *testing.T) {
		err := Copy("", "outpath", 0, 0)
		require.ErrorAs(t, err, &ErrInvalidFileType)

		err = Copy("inpath", "", 0, 0)
		require.ErrorAs(t, err, &ErrInvalidFileType)
	})

	t.Run("the same path", func(t *testing.T) {
		outPath := inPath
		before, err := os.ReadFile(outPath)
		require.NoError(t, err)

		err = Copy(inPath, outPath, 0, 0)
		require.ErrorAs(t, err, &ErrInvalidFileType)

		after, err := os.ReadFile(outPath)
		require.NoError(t, err)

		require.Equal(t, before, after)
	})
}
