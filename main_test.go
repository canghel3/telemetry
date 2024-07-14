package main

import (
	"bytes"
	"gotest.tools/v3/assert"
	"os"
	"telemetry/level"
	"telemetry/log"
	"testing"
)

//1. Test output to file with all severity levels and random content
//2. Test output to stdout with all severity levels and random content
//3. Test output to custom driver with all severity levels and random content

func TestOutputToFile(t *testing.T) {
	file := "xyz.txt"
	_, err := os.Create(file)
	assert.NilError(t, err)

	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			log.File(file).NoLevel().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.NewNoLevel().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).NoLevel().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.NewNoLevel().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			log.File(file).Error().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.NewErrorLevel().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Error().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.NewErrorLevel().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("WARN LEVEL", func(t *testing.T) {
			log.File(file).Warn().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.NewWarnLevel().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Warn().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.NewWarnLevel().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("INFO LEVEL", func(t *testing.T) {
			log.File(file).Info().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.NewInfoLevel().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Info().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.NewInfoLevel().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {
			log.File(file).Debug().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.NewDebugLevel().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Debug().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.NewDebugLevel().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})
	})

	os.WriteFile(file, nil, 0600)
}

func TestOutputToStdout(t *testing.T) {

}
