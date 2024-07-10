package main

import (
	"bytes"
	"gotest.tools/v3/assert"
	"os"
	"telemetry/log"
	"testing"
)

//1. Test output to file with all severity levels and random content
//2. Test output to stdout with all severity levels and random content
//3. Test output to custom driver with all severity levels and random content

func TestOutputToFile(t *testing.T) {
	file := "xyz"
	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			log.File(file).NoLevel().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, bytes.Contains(retrieved, content))

			log.File(file).NoLevel().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, bytes.Contains(retrieved, moreContent))
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			log.File(file).Error().Write(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, bytes.Contains(retrieved, content))

			log.File(file).Error().Write(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, bytes.Contains(retrieved, moreContent))
		})

		t.Run("WARN LEVEL", func(t *testing.T) {

		})

		t.Run("INFO LEVEL", func(t *testing.T) {

		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {

		})
	})
}
