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
	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			lvl := level.LevelToText[level.NoLevel]
			log.File(file).NoLevel().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(lvl + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).NoLevel().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(lvl + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			lvl := level.LevelToText[level.LevelError]
			log.File(file).Error().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(lvl + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Error().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(lvl + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("WARN LEVEL", func(t *testing.T) {
			lvl := level.LevelToText[level.LevelWarn]
			log.File(file).Warn().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(lvl + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Warn().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(lvl + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("INFO LEVEL", func(t *testing.T) {
			lvl := level.LevelToText[level.LevelInfo]
			log.File(file).Info().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(lvl + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Info().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(lvl + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {
			lvl := level.LevelToText[level.LevelDebug]
			log.File(file).Debug().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(lvl + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Debug().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(lvl + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})
	})

	os.WriteFile(file, nil, 0600)
}

func TestOutputToStdout(t *testing.T) {

}
