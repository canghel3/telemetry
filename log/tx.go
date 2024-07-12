package log

import (
	"fmt"
	"github.com/google/uuid"
)

type Tx struct {
	id         string
	commited   bool
	attributes map[any]any
	logs       Logger
}

func (l *Logger) BeginTx() *Tx {
	return &Tx{
		id:   uuid.New().String(),
		logs: *l,
	}
}

func (l *Logger) BeginTxWithMetadata(metadata map[any]any) *Tx {
	return &Tx{
		logs:       *l,
		id:         uuid.New().String(),
		attributes: metadata,
		commited:   false,
	}
}

func (tx *Tx) Append(log *Logger) {
	if !tx.commited {
		tx.logs = *log
	}
}

func (tx *Tx) Commit() error {
	if !tx.commited {
		return tx.logs.outputDriver.Write(tx.logs.buf)
	}

	return fmt.Errorf("transaction already committed")
}
