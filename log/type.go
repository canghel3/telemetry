package log

import (
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
		metadata: nil,
		output:   output,
	}
}

func (m *Message) Metadata(meta map[any]any) *Message {
	m.metadata = meta
	return m
}

func (m *Message) Msg(msg string) *Message {
	m.content = []byte(msg)
	return m
}

func (m *Message) Msgf(msg string, format ...any) *Message {
	m.content = []byte(fmt.Sprintf(msg, format...))
	return m
}

func (m *Message) Log(msg string) {
	m.content = []byte(msg)

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
		fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", msg, err.Error())
	}
}

func (m *Message) Logf(msg string, format ...any) {
	m.content = []byte(fmt.Sprintf(msg, format...))
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
		fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", msg, err.Error())
	}
}

func (m *Message) formatLogOutput() []byte {
	//default order: TIMESTAMP LEVEL METADATA BUFFER
	var timestamp string
	if len(m.output.config.Formatting.LogConfig.Timestamp) > 0 {
		timestamp = time.Now().Format(m.output.config.Formatting.LogConfig.Timestamp)
	} else {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp)...)
	out = append(out, byte(' '))
	out = append(out, []byte(m.level.Type())...)
	out = append(out, byte(' '))

	var meta2bytes = make([]byte, 0)
	for k, v := range m.metadata {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...) //careful
	}
	//very careful (whitespace)
	out = append(out, meta2bytes...)
	out = append(out, m.content...)
	out = append(out, '\n')

	return out
}
