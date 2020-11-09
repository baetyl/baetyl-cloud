package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

var cache = map[string]string{}

func (d *DB) CreateIndex(namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error) {
	return d.CreateIndexTx(nil, namespace, keyA, keyB, valueA, valueB)
}

func (d *DB) ListIndex(namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error) {
	return d.ListIndexTx(nil, namespace, keyA, byKeyB, valueB)
}

func (d *DB) DeleteIndex(namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error) {
	return d.DeleteIndexTx(nil, namespace, keyA, byKeyB, valueB)
}

func (d *DB) CreateIndexTx(tx *sqlx.Tx, namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error) {
	selectSQL := fmt.Sprintf(`INSERT INTO %s (namespace, %s, %s) VALUES (?, ?, ?)`, getTable(keyA, keyB), keyA, keyB)
	return d.Exec(tx, selectSQL, namespace, valueA, valueB)
}

func (d *DB) ListIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error) {
	selectSQL := fmt.Sprintf(`SELECT %s FROM %s WHERE namespace = ? and %s = ?`, keyA, getTable(keyA, byKeyB), byKeyB)
	var res []string
	if err := d.Query(tx, selectSQL, &res, namespace, valueB); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DB) DeleteIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error) {
	selectSQL := fmt.Sprintf(`DELETE FROM %s WHERE namespace = ? and %s = ?`, getTable(keyA, byKeyB), byKeyB)
	return d.Exec(tx, selectSQL, namespace, valueB)
}

func (d *DB) RefreshIndex(namespace string, keyA, keyB common.Resource, valueA string, valueBs []string) error {
	return d.Transact(func(tx *sqlx.Tx) error {
		if _, err := d.DeleteIndexTx(tx, namespace, keyB, keyA, valueA); err != nil {
			return err
		}
		for _, b := range valueBs {
			if _, err := d.CreateIndexTx(tx, namespace, keyA, keyB, valueA, b); err != nil {
				return err
			}
		}
		return nil
	})
}

func getTable(keyA, keyB common.Resource) string {
	keyAB := string(keyA) + "_" + string(keyB)
	if v, ok := cache[keyAB]; ok {
		return v
	}
	keyBA := string(keyB) + "_" + string(keyA)
	if v, ok := cache[keyBA]; ok {
		return v
	}
	var res string
	if keyA < keyB {
		res = "baetyl_index_" + keyAB
	} else {
		res = "baetyl_index_" + keyBA
	}
	cache[keyAB] = res
	cache[keyBA] = res
	return res
}
