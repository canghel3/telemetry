package main

import (
	"bytes"
	"fmt"
	"gotest.tools/v3/assert"
	"io"
	"os"
	"strconv"
	"sync"
	"telemetry/level"
	"telemetry/log"
	"testing"
)

//1. Test output to file with all severity levels and random content
//2. Test output to stdout with all severity levels and random content
//3. Test output to custom driver with all severity levels and random content

const file = "xyz.log"

type errorDriver struct {
	msg []byte
}

func (ed *errorDriver) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("intentional write error")
}

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

		t.Run("VERBOSE", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == true)
		})

		t.Run("IN-LINE", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			tx := log.BeginTx()
			tx.Append(log.File(file).Debug().Msg(l1))
			tx.Append(log.File(file).NoLevel().Msg(l2))
			tx.Append(log.File(file).Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Debug().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.None().Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == true)
		})

		t.Run("MIXED OUTPUTS", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == false)
		})

		t.Run("WITH METADATA", func(t *testing.T) {
			t.SkipNow()
			os.WriteFile(file, nil, 0600)

			meta := map[any]any{
				1:       "2",
				'a':     true,
				3.14159: imag(complex(1, 1)),
			}

			toFile := log.File(file)
			tx := log.BeginTxWithMetadata(meta)
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == false)
		})

		t.Run("LOG", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("DID NOT COMMIT", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})

		t.Run("ROLLBACK", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})

		t.Run("ALREADY COMMITED", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Commit()
			err = tx.Commit()

			assert.Error(t, err, "transaction already committed or rolled back")
		})

		t.Run("ALREADY ROLLED BACK", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Rollback()
			err := tx.Rollback()

			assert.Error(t, err, "transaction already committed or rolled back")
		})
	})

	t.Run("CONFIG FILE", func(t *testing.T) {

	})

	t.Run("FAILED TO WRITE LOG", func(t *testing.T) {
		r, w, err := os.Pipe()
		assert.NilError(t, err)

		os.Stderr = w

		ed := &errorDriver{}
		log.OutputDriver(ed).Info().Log([]byte("schmeckermeister"))

		err = w.Close()
		assert.NilError(t, err)

		var read bytes.Buffer
		_, err = io.Copy(&read, r)
		assert.NilError(t, err)

		assert.Assert(t, bytes.Contains(read.Bytes(), []byte("intentional write error")))
	})

	t.Run("CONCURRENT WRITING TO THE SAME FILE", func(t *testing.T) {
		os.WriteFile(file, nil, 0600)

		wg := sync.WaitGroup{}
		//launch 100 goroutines to write to the same file
		for i := 0; i < 100; i++ {
			wg.Add(1)
			content := []byte("some content " + strconv.Itoa(i))
			go func() {
				defer wg.Done()
				toFile := log.File(file)
				toFile.Info().Log(content)
			}()
		}

		wg.Wait()
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
	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			log.Stdout().NoLevel().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.None().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			expected = []byte(level.None().Type() + " " + string(moreContent) + "\n")

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			log.Stdout().NoLevel().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.None().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
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

		t.Run("VERBOSE", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == true)
		})

		t.Run("IN-LINE", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			tx := log.BeginTx()
			tx.Append(log.File(file).Debug().Msg(l1))
			tx.Append(log.File(file).NoLevel().Msg(l2))
			tx.Append(log.File(file).Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Debug().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.None().Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == true)
		})

		t.Run("MIXED OUTPUTS", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == false)
		})

		t.Run("WITH METADATA", func(t *testing.T) {
			t.SkipNow()
			os.WriteFile(file, nil, 0600)

			meta := map[any]any{
				1:       "2",
				'a':     true,
				3.14159: imag(complex(1, 1)),
			}

			toFile := log.File(file)
			tx := log.BeginTxWithMetadata(meta)
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))
			tx.Log()

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected1) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected2) == true)
			assert.Assert(t, bytes.Contains(retrieved, expected3) == false)
		})

		t.Run("LOG", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("DID NOT COMMIT", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})

		t.Run("ROLLBACK", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(log.Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})

		t.Run("ALREADY COMMITED", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Commit()
			err := tx.Commit()

			assert.Error(t, err, "transaction already committed or rolled back")
		})

		t.Run("ALREADY ROLLED BACK", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := log.File(file)
			tx := log.BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Rollback()
			err := tx.Rollback()

			assert.Error(t, err, "transaction already committed or rolled back")
		})
	})

	t.Run("CONFIG FILE", func(t *testing.T) {

	})

	t.Run("FAILED TO WRITE LOG", func(t *testing.T) {
		t.SkipNow()
		ed := &errorDriver{}
		log.OutputDriver(ed).Info().Log([]byte("schmeckermeister"))

		//f, err := os.OpenFile(os.Stderr.Name(), os.O_RDONLY, 0777)
		//assert.NilError(t, err)

		read, err := os.ReadFile(os.Stderr.Name())
		assert.NilError(t, err)

		assert.NilError(t, err)
		assert.Assert(t, bytes.Equal(read, []byte("failed to write log: intentional write error\n")))
	})

	t.Run("CONCURRENT WRITING TO THE SAME FILE", func(t *testing.T) {
		os.WriteFile(file, nil, 0600)

		wg := sync.WaitGroup{}
		//launch 100 goroutines to write to the same file
		for i := 0; i < 100; i++ {
			wg.Add(1)
			content := []byte("some content " + strconv.Itoa(i))
			go func() {
				defer wg.Done()
				toFile := log.File(file)
				toFile.Info().Log(content)
			}()
		}

		wg.Wait()
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
