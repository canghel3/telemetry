package log

import (
	"bytes"
	"fmt"
	"github.com/Ginger955/telemetry/level"
	"gotest.tools/v3/assert"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"
)

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

func ExampleFile() {
	File("./testdata/out.log").Info().Log([]byte("Foo"))

	toFile := File("./testdata/out.log")
	toFile.Info().Log([]byte("Bar"))
}

func ExampleStdout() {
	Stdout().Info().Log([]byte("Foo"))

	stdout := Stdout()
	stdout.Info().Log([]byte("Bar"))
}

func ExampleOutputDriver() {
	OutputDriver(os.Stdout).Info().Log([]byte("Foo"))
	OutputDriver(os.Stderr).Info().Log([]byte("Bar"))
}

func TestOutputToFile(t *testing.T) {
	_, err := os.Create(file)
	assert.NilError(t, err)

	t.Run("LEVELS", func(t *testing.T) {
		content := []byte("the quick brown fox jumps over the lazy dog")
		moreContent := []byte("พยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩, ⟨д⟩, ⟨ж⟩, ⟨з⟩, ⟨к⟩, ⟨л⟩, ⟨м⟩, ⟨н⟩, ⟨п⟩, ⟨р⟩, ⟨с⟩, ⟨т⟩, ⟨ф⟩, ⟨х⟩, ⟨ц⟩, ⟨ч⟩, ⟨ш⟩, ⟨щ⟩")

		t.Run("NO LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			File(file).NoLevel().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.None().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).NoLevel().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.None().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			File(file).Error().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Error().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).Error().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Error().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("WARN LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			File(file).Warn().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Warn().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).Warn().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Warn().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("INFO LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			File(file).Info().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Info().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).Info().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Info().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			File(file).Debug().Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(level.Debug().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).Debug().Log(moreContent)

			retrieved, err = os.ReadFile(file)
			assert.NilError(t, err)

			expected = []byte(level.Debug().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)
		})

		t.Run("CUSTOM LEVEL", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			criticalLevel := level.Custom("CRITICAL")

			File(file).Level(criticalLevel).Log(content)

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			expected := []byte(criticalLevel.Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(retrieved, expected) == true)

			File(file).Level(criticalLevel).Log(moreContent)

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

			toFile := File(file)
			tx := BeginTx()
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

			tx := BeginTx()
			tx.Append(File(file).Debug().Msg(l1))
			tx.Append(File(file).NoLevel().Msg(l2))
			tx.Append(File(file).Warn().Msg(l3))
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

			toFile := File(file)
			tx := BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))
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
			os.WriteFile(file, nil, 0600)

			meta := map[any]any{
				1:       "2",
				'a':     true,
				3.14159: imag(complex(1, 1)),
			}

			toFile := File(file)
			tx := BeginTxWithMetadata(meta)
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))
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
			os.WriteFile(file, nil, 0600)

			toFile := File(file)
			tx := BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))
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

		t.Run("DID NOT COMMIT", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := File(file)
			tx := BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})

		t.Run("ROLLBACK", func(t *testing.T) {
			os.WriteFile(file, nil, 0600)

			toFile := File(file)
			tx := BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))

			retrieved, err := os.ReadFile(file)
			assert.NilError(t, err)

			assert.Assert(t, len(retrieved) == 0)
		})
	})

	t.Run("CONFIG FILE", func(t *testing.T) {

	})

	t.Run("FAILED TO WRITE LOG", func(t *testing.T) {
		r, w, err := os.Pipe()
		assert.NilError(t, err)

		os.Stderr = w

		ed := &errorDriver{}
		OutputDriver(ed).Info().Log([]byte("schmeckermeister"))

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
				toFile := File(file)
				toFile.Info().Log(content)
			}()
		}

		wg.Wait()
	})

	t.Run("NIL CONTENT", func(t *testing.T) {
		os.WriteFile(file, nil, 0600)

		File(file).Info().Log(nil)

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

			Stdout().NoLevel().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.None().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().NoLevel().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.None().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})

		t.Run("ERROR LEVEL", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Error().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.Error().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Error().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.Error().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})

		t.Run("WARN LEVEL", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Warn().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.Warn().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Warn().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.Warn().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})

		t.Run("INFO LEVEL", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Info().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.Info().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Info().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.Info().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})

		t.Run("DEBUG LEVEL", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Debug().Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(level.Debug().Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Debug().Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(level.Debug().Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})

		t.Run("CUSTOM LEVEL", func(t *testing.T) {
			criticalLevel := level.Custom("CRITICAL")

			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Level(criticalLevel).Log(content)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected := []byte(criticalLevel.Type() + " " + string(content) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)

			read.Reset()

			initial = os.Stdout
			r, w, err = os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			Stdout().Level(criticalLevel).Log(moreContent)

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected = []byte(criticalLevel.Type() + " " + string(moreContent) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
		})
	})

	t.Run("TRANSACTIONS", func(t *testing.T) {
		l1 := []byte("i enjoy doing the same thing over and over again and expecting different results")
		l2 := []byte("freedom is not for everyone. at least that's what the sign says")
		l3 := []byte("I LIVE, I DIE, I LIVE AGAIN. WITNESS ME GOING TO VALHALLA!")

		t.Run("VERBOSE", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			toFile := Stdout()
			tx := BeginTx()
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(toFile.Warn().Msg(l3))
			tx.Log()

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected1 := []byte(level.Info().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.Custom("MAJOR").Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected1) == true)
			assert.Assert(t, bytes.Contains(read.Bytes(), expected2) == true)
			assert.Assert(t, bytes.Contains(read.Bytes(), expected3) == true)
		})

		t.Run("IN-LINE", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			tx := BeginTx()
			tx.Append(Stdout().Debug().Msg(l1))
			tx.Append(Stdout().NoLevel().Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))
			tx.Log()

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			expected1 := []byte(level.Debug().Type() + " " + string(l1) + "\n")
			expected2 := []byte(level.None().Type() + " " + string(l2) + "\n")
			expected3 := []byte(level.Warn().Type() + " " + string(l3) + "\n")
			assert.Assert(t, bytes.Contains(read.Bytes(), expected1) == true)
			assert.Assert(t, bytes.Contains(read.Bytes(), expected2) == true)
			assert.Assert(t, bytes.Contains(read.Bytes(), expected3) == true)
		})

		t.Run("WITH METADATA", func(t *testing.T) {
			t.SkipNow()
			os.WriteFile(file, nil, 0600)

			meta := map[any]any{
				1:       "2",
				'a':     true,
				3.14159: imag(complex(1, 1)),
			}

			toFile := File(file)
			tx := BeginTxWithMetadata(meta)
			tx.Append(toFile.Info().Msg(l1))
			tx.Append(toFile.Level(level.Custom("MAJOR")).Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))
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

		t.Run("DID NOT LOG", func(t *testing.T) {
			initial := os.Stdout
			r, w, err := os.Pipe()
			assert.NilError(t, err)

			os.Stdout = w

			tx := BeginTx()
			tx.Append(Stdout().Debug().Msg(l1))
			tx.Append(Stdout().NoLevel().Msg(l2))
			tx.Append(Stdout().Warn().Msg(l3))

			err = w.Close()
			assert.NilError(t, err)

			os.Stdout = initial

			var read bytes.Buffer
			_, err = io.Copy(&read, r)
			assert.NilError(t, err)

			assert.Assert(t, len(read.Bytes()) == 0)
		})
	})

	t.Run("CONFIG", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("CONCURRENT WRITING TO STDOUT", func(t *testing.T) {
		wg := sync.WaitGroup{}
		//launch 100 goroutines to write to the same file
		for i := 0; i < 100; i++ {
			wg.Add(1)
			content := []byte("some content " + strconv.Itoa(i))
			go func() {
				defer wg.Done()
				toStdout := Stdout()
				toStdout.Info().Log(content)
			}()
		}

		wg.Wait()
	})

	t.Run("NIL CONTENT", func(t *testing.T) {
		initial := os.Stdout
		r, w, err := os.Pipe()
		assert.NilError(t, err)

		os.Stdout = w

		Stdout().Info().Log(nil)

		err = w.Close()
		assert.NilError(t, err)

		os.Stdout = initial

		var read bytes.Buffer
		_, err = io.Copy(&read, r)
		assert.NilError(t, err)

		expected := []byte(level.Info().Type() + " \n")
		assert.Assert(t, bytes.Contains(read.Bytes(), expected) == true)
	})
}

func TestOutputToCustom(t *testing.T) {
	t.Run("FAILED TO WRITE LOG", func(t *testing.T) {
		initial := os.Stderr
		r, w, err := os.Pipe()
		assert.NilError(t, err)

		os.Stderr = w

		ed := &errorDriver{}
		OutputDriver(ed).Info().Log([]byte("schmeckermeister"))

		err = w.Close()
		assert.NilError(t, err)

		os.Stderr = initial

		var read bytes.Buffer
		_, err = io.Copy(&read, r)
		assert.NilError(t, err)

		assert.Assert(t, bytes.Contains(read.Bytes(), []byte("failed to write log")))
		assert.Assert(t, bytes.Contains(read.Bytes(), []byte(level.Info().Type())))
		assert.Assert(t, bytes.Contains(read.Bytes(), []byte("schmeckermeister: intentional write error")))
	})
}
