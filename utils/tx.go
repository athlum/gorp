package utils

import (
	"github.com/juju/errors"
	"github.com/athlum/gorp"
)

type TxHandler func(tx gorp.SqlExecutor) error

func (h TxHandler) Do(tx *gorp.Transaction) error {
	if h == nil {
		return nil
	}

	if err := h(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Annotatef(err, "TxHandler.do.rollback")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Annotatef(err, "TxHandler.do.commit")
	}

	return nil
}
