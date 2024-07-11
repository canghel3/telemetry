package log

import "fmt"

type Tx struct {
	id         string
	commited   bool
	attributes map[any]any
	logs       []Log
}

func (l *Log) BeginTx() *Tx {
	return &Tx{
		id:   "", //set id to uuid v4
		logs: []Log{*l},
	}
}

func (l *Log) BeginTxWithMetadata(metadata map[any]any) *Tx {
	return &Tx{
		logs:       []Log{*l},
		id:         "",
		attributes: metadata,
		commited:   false,
	}
}

func (tx *Tx) Append(log *Log) {
	if !tx.commited {
		tx.logs = append(tx.logs, *log)
	}
}

func (tx *Tx) Commit() error {
	if !tx.commited {
		for _, log := range tx.logs {
			err := log.outputDriver.Write(formatLogOutput(&log))
			if err != nil {
				return err

			}
		}
		return nil
	}

	return fmt.Errorf("transaction already committed")
}
