package log

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

type Tx struct {
	messages []*Message
	id       string
	commited bool
	metadata map[any]any
}

func BeginTx() *Tx {
	return &Tx{
		id:       uuid.New().String(),
		messages: []*Message{},
		metadata: nil,
		commited: false,
	}
}

func BeginTxWithMetadata(metadata map[any]any) *Tx {
	return &Tx{
		messages: []*Message{},
		id:       uuid.New().String(),
		metadata: metadata,
		commited: false,
	}
}

func (tx *Tx) Append(message *Message) {
	if !tx.commited {
		tx.messages = append(tx.messages, message)
	}
}

// Log send the existing message entries to their respective output driver.
// Any error is written to os.Stderr
func (tx *Tx) Log() {
	if !tx.commited {
		tx.commited = true
		for _, msg := range tx.messages {
			//TODO: enable/disable formatting based on config
			formattedOutput := tx.formatTransactionOutput(msg)
			_, err := msg.output.driver.Write(formattedOutput)
			if err != nil {
				//write the error encountered during logging to os.Stderr. wip: any configured file
				//we could write to the log output driver because it implements the required w io.Writer,
				//but if the output driver is fatally broken, we also lose the error messages.
				fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", formattedOutput, err.Error())
			}
		}
	}
}

// TODO: complete formatter
func (tx *Tx) formatTransactionOutput(msg *Message) []byte {
	var buffer bytes.Buffer

	// format timestamp
	timestampFormat := "2006-01-02 15:04:05"
	if len(msg.output.config.Formatting.LogConfig.Timestamp) > 0 {
		timestampFormat = msg.output.config.Formatting.LogConfig.Timestamp
	}

	buffer.WriteString(time.Now().Format(timestampFormat))
	buffer.WriteByte(' ')
	buffer.WriteString("| TRANSACTION - " + tx.id + " |")
	buffer.WriteByte(' ')

	// format level
	buffer.WriteString(msg.level.Type())
	buffer.WriteByte(' ')

	// format metadata
	for k, v := range tx.metadata {
		fmt.Fprintf(&buffer, "%v:%v ", k, v)
	}
	//
	//if len(tx.metadata) > 0 {
	//	buffer.WriteByte(' ')
	//}

	// add content
	buffer.Write(msg.content)
	buffer.WriteByte('\n')

	return buffer.Bytes()
}
