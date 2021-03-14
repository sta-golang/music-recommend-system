package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"runtime/debug"
)

type mysqlTransaction struct {
}

var onceMysqlTransaction = mysqlTransaction{}

func newMysqlTransaction() *mysqlTransaction {
	return &onceMysqlTransaction
}

func (mt *mysqlTransaction) doTransaction(txFn func(tx *sqlx.Tx) error) (err error) {

	transaction, err := client(dbMusicRecommendNameTest).Beginx()
	if err != nil {
		log.Error(err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Fatal(p)
			log.Fatalf(str.BytesToString(debug.Stack()))
			rErr := transaction.Rollback()
			if rErr != nil {
				log.Error(rErr)
			}
			err = fmt.Errorf("panic err : %v", p)
		} else if err != nil {
			log.Error(err)
			rErr := transaction.Rollback()
			if rErr != nil {
				log.Error(rErr)
			}
		} else {
			err = transaction.Commit()
			if err != nil {
				log.Error(err)
			}
			return
		}
	}()
	err = txFn(transaction)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (mt *mysqlTransaction) Transaction(txFn func(tx *sqlx.Tx) error) error {
	return mt.doTransaction(txFn)
}
