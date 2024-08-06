package log

import (
	"fmt"
	"github.com/google/uuid"
	"os"
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
	output := make([]byte, 0)

	t := "| TRANSACTION - " + tx.id + " |"

	output = append(output, t...)
	output = append(output, ' ')

	var meta2bytes = make([]byte, 0)
	for k, v := range tx.metadata {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...)
	}

	output = append(output, meta2bytes...)
	output = append(output, ' ')
	output = append(output, msg.content...)
	output = append(output, '\n')

	return output
}
