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

const file = "xyz.log"

type customLevel struct {
	levelType string
}

func newCustomLevel(levelType string) *customLevel {
	return &customLevel{levelType: levelType}
}

func (cl *customLevel) Type() string {
	return cl.levelType
}

func TestOutputToFile(t *testing.T) {
	_, err := os.Create(file)
	assert.NilError(t, err)

	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			log.File(file).NoLevel().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.None().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).NoLevel().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.None().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			log.File(file).Error().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Error().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Error().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Error().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("WARN LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			log.File(file).Warn().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Warn().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Warn().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Warn().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("INFO LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			log.File(file).Info().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Info().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Info().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Info().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			log.File(file).Debug().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Debug().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Debug().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Debug().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("CUSTOM LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			criticalLevel := level.Custom("CRITICAL")

			log.File(file).Level(criticalLevel).Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(criticalLevel.Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			log.File(file).Level(criticalLevel).Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(criticalLevel.Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})
	})

	t.Run("TRANSACTIONS", func(t *testing.T) {
		l1 := []byte("i enjoy doing the same thing over and over again and expecting different results")
		l2 := []byte("freedom is not for everyone. at least that's what the sign says")
		l3 := []byte("I LIVE, I DIE, I LIVE AGAIN. WITNESS ME GOING TO VALHALLA!")

		t.Run("USING VARIABLES", func(t *testing.T) {
			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
		})

		t.Run("WITHOUT USING VARIABLES", func(t *testing.T) {
			//tx := log.BeginTx()
			//tx.Append(log.File(.))
		})

		t.Run("APPEND DIFFERENT OUTPUT THAN TRANSACTION WAS INITIALIZED WITH", func(t *testing.T) {

		})

	})

	t.Run("NIL CONTENT", func(t *testing.T) {
		os.WriteFile(file, nil, 0600)

		log.File(file).Info().Log(nil)

		retrieved, err := os.ReadFile(file)
		assert.NilError(t, err)

		expected := []byte(level.Info().Type() + " \n")
		assert.Assert(t, bytes.Contains(retrieved, expected) == true)
	})

	os.WriteFile(file, nil, 0600)
}

func TestOutputToStdout(t *testing.T) {

}
