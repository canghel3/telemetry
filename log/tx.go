package log

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

type Tx struct {
	id         string
	commited   bool
	attributes map[any]any
	logs       []Logger
}

func BeginTx() *Tx {
	return &Tx{
		id:         uuid.New().String(),
		logs:       []Logger{},
		attributes: nil,
		commited:   false,
	}
}

func BeginTxWithMetadata(metadata map[any]any) *Tx {
	return &Tx{
		logs:       []Logger{},
		id:         uuid.New().String(),
		attributes: metadata,
		commited:   false,
	}
}

func (tx *Tx) Append(log *Logger) {
	if !tx.commited {
		tx.logs = append(tx.logs, *log)
	}
}

// Commit sends the existing logs to the corresponding output driver.
// If an error occurs during writing, the error message is written to os.Stderr and the process continues.
func (tx *Tx) Commit() error {
	if !tx.commited {
		for _, log := range tx.logs {
			formattedOutput := formatTransactionOutput(*tx, log)
			_, err := log.outputDriver.Write(append(formattedOutput, '\n'))
			if err != nil {
				//write the error encountered during writing to os.Stderr
				//we could write to the log output driver because it implements the io.Writer,
				//but if the output driver is fatally broken, the writing failure will be lost as well
				//and debugging becomes more difficult.
				fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", formattedOutput, err.Error())
			}
		}
		tx.commited = true
		return nil
	}

	return fmt.Errorf("transaction already committed or rolled back")
}

// Rollback is not really required. If you don't need the transaction anymore,
// just let the garbage collector take care of it.
func (tx *Tx) Rollback() error {
	if !tx.commited {
		tx.logs = []Logger{}
		tx.commited = true
		tx.attributes = nil
	}

	return fmt.Errorf("transaction already committed or rolled back")
}

func formatTransactionOutput(tx Tx, log Logger) []byte {
	output := make([]byte, 0)

	t := "| TRANSACTION - " + tx.id + " |"

	output = append(output, t...)
	output = append(output, ' ')

	var meta2bytes = make([]byte, 0)
	meta2bytes = append(meta2bytes, "METADATA: "...)
	for k, v := range tx.attributes {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v", k, v))...)
	}

	output = append(output, meta2bytes...)
	output = append(output, ' ')
	output = append(output, log.buf...)
	return output
}
