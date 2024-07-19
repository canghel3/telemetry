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

// Log send the existing log entries to their respective output driver.
// Any error is written to os.Stderr
// TODO: add configurable output and error destination;
func (tx *Tx) Log() {
	if !tx.commited {
		tx.commited = true
		for _, log := range tx.logs {
			formattedOutput := formatTransactionOutput(*tx, log)
			_, err := log.outputDriver.Write(append(formattedOutput, '\n'))
			if err != nil {
				//write the error encountered during logging to os.Stderr. wip: any configured file
				//we could write to the log output driver because it implements the required w io.Writer,
				//but if the output driver is fatally broken, we also lose the error logs.
				fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", formattedOutput, err.Error())
			}
		}
	}
}

func formatTransactionOutput(tx Tx, log Logger) []byte {
	output := make([]byte, 0)

	t := "| TRANSACTION - " + tx.id + " |"

	output = append(output, t...)
	output = append(output, ' ')

	var meta2bytes = make([]byte, 0)
	//meta2bytes = append(meta2bytes, "METADATA: "...)
	for k, v := range tx.attributes {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...)
	}

	output = append(output, meta2bytes...)
	output = append(output, ' ')
	output = append(output, log.buf...)
	return output
}
