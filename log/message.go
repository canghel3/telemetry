package log

import (
	"bytes"
	"fmt"
	"github.com/Ginger955/telemetry/level"
	"os"
	"time"
)

type Message struct {
	content  []byte
	level    level.Level
	metadata map[any]any
	output   *Output
}

func newMessage(output *Output, level level.Level) *Message {
	return &Message{
		content:  nil,
		level:    level,
		metadata: output.meta,
		output:   output,
	}
}

// Metadata sets the metadata only for this message.
func (m *Message) Metadata(meta map[any]any) *Message {
	m.metadata = meta
	return m
}

// Msg is only meant for use in log transactions.
func (m *Message) Msg(msg string) *Message {
	m.content = []byte(msg)
	return m
}

// Msgf is only meant for use in log transactions.
func (m *Message) Msgf(msg string, format ...any) *Message {
	m.content = []byte(fmt.Sprintf(msg, format...))
	return m
}

// Log logs to the corresponding output driver.
func (m *Message) Log(msg string) {
	m.content = []byte(msg)
	m.log()
}

// Logf logs to the corresponding output driver based on the given format.
func (m *Message) Logf(msg string, format ...any) {
	m.content = []byte(fmt.Sprintf(msg, format...))
	m.log()
}

func (m *Message) log() {
	var output = m.content
	if !m.output.config.Formatting.LogConfig.FormattingDisabled {
		output = m.formatLogOutput()
	}

	_, err := m.output.driver.Write(output)
	if err != nil {
		//write the error encountered during writing to os.Stderr
		//we could write to the log output driver because it implements the io.Writer,
		//but if the output driver is fatally broken, the writing failure will be lost as well
		//and debugging becomes more difficult.
		fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", m.content, err.Error())
	}
}

func (m *Message) formatLogOutput() []byte {
	//default order: TIMESTAMP LEVEL METADATA BUFFER
	var buffer bytes.Buffer

	// format timestamp
	timestampFormat := "2006-01-02 15:04:05"
	if len(m.output.config.Formatting.LogConfig.Timestamp) > 0 {
		timestampFormat = m.output.config.Formatting.LogConfig.Timestamp
	}

	buffer.WriteString(time.Now().Format(timestampFormat))
	buffer.WriteByte(' ')

	// format level
	buffer.WriteString(m.level.Type())
	buffer.WriteByte(' ')

	// format metadata
	for k, v := range m.metadata {
		fmt.Fprintf(&buffer, "%v:%v ", k, v)
	}

	//if len(m.metadata) > 0 {
	//	buffer.WriteByte(' ')
	//}

	// add content
	buffer.Write(m.content)
	buffer.WriteByte('\n')

	return buffer.Bytes()
}
