package log

import (
	"fmt"
	"github.com/google/uuid"
)

type Tx struct {
	id         string
	commited   bool
	attributes map[any]any
	logs       []Logger
}

func (l *Logger) BeginTx() *Tx {
	return &Tx{
		id:         uuid.New().String(),
		logs:       []Logger{},
		attributes: nil,
		commited:   false,
	}
}

func (l *Logger) BeginTxWithMetadata(metadata map[any]any) *Tx {
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

func (tx *Tx) Commit() error {
	if !tx.commited {
		for _, log := range tx.logs {
			err := log.outputDriver.Log(formatTransactionOutput(&log, tx.id))
			if err != nil {
				return err
			}
		}
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

func formatTransactionOutput(log *Logger, id string) []byte {
	return append([]byte("| TRANSACTION - "+id+" | "), log.buf...)
}
