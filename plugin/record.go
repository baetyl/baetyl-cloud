package plugin

import (
	"database/sql"
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../mock/plugin/record.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Record

// Record interface of Record
type Record interface {
	CountRecord(batchName, fingerprintValue, ns string) (int, error)
	GetRecord(batchName, recordName, ns string) (*models.Record, error)
	GetRecordByFingerprint(batchName, ns, value string) (*models.Record, error)
	ListRecord(batchName, ns string, filter *models.Filter) ([]models.Record, error)
	CreateRecord(records []models.Record) (sql.Result, error)
	UpdateRecord(record *models.Record) (sql.Result, error)
	DeleteRecord(batchName, recordName, ns string) (sql.Result, error)
	GetRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (*models.Record, error)
	CountRecordTx(tx *sqlx.Tx, batchName, fingerprintValue, ns string) (int, error)
	GetRecordByFingerprintTx(tx *sqlx.Tx, batchName, ns, value string) (*models.Record, error)
	ListRecordTx(tx *sqlx.Tx, batchName, ns string, filter *models.Filter) ([]models.Record, error)
	CreateRecordTx(tx *sqlx.Tx, records []models.Record) (sql.Result, error)
	UpdateRecordTx(tx *sqlx.Tx, record *models.Record) (sql.Result, error)
	DeleteRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (sql.Result, error)
	io.Closer
}
