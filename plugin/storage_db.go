package plugin

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../mock/plugin/storage_db.go -package=plugin github.com/baetyl/baetyl-cloud/plugin DBStorage

// DBStorage DBStorage
type DBStorage interface {
	Transact(func(*sqlx.Tx) error) error

	// index
	CreateIndex(namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error)
	ListIndex(namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error)
	DeleteIndex(namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error)
	CreateIndexTx(tx *sqlx.Tx, namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error)
	ListIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error)
	DeleteIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error)
	RefreshIndex(namespace string, keyA, keyB common.Resource, valueA string, valueBs []string) error

	// batch
	GetBatch(name, ns string) (*models.Batch, error)
	ListBatch(ns, name string, page, size int) ([]models.Batch, error)
	CreateBatch(batch *models.Batch) (sql.Result, error)
	UpdateBatch(batch *models.Batch) (sql.Result, error)
	DeleteBatch(name, ns string) (sql.Result, error)
	CountBatch(ns, name string) (int, error)
	CountBatchByCallback(callbackName, ns string) (int, error)
	GetBatchTx(tx *sqlx.Tx, name, ns string) (*models.Batch, error)
	ListBatchTx(tx *sqlx.Tx, ns, name string, page, size int) ([]models.Batch, error)
	CreateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error)
	UpdateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error)
	DeleteBatchTx(tx *sqlx.Tx, name, ns string) (sql.Result, error)
	CountBatchTx(tx *sqlx.Tx, ns, name string) (int, error)
	CountBatchByCallbackTx(tx *sqlx.Tx, callbackName, ns string) (int, error)

	// record
	CountRecord(batchName, fingerprintValue, ns string) (int, error)
	GetRecord(batchName, recordName, ns string) (*models.Record, error)
	GetRecordByFingerprint(batchName, ns, value string) (*models.Record, error)
	ListRecord(batchName, fingerprintValue, ns string, page, size int) ([]models.Record, error)
	CreateRecord(records []models.Record) (sql.Result, error)
	UpdateRecord(record *models.Record) (sql.Result, error)
	DeleteRecord(batchName, recordName, ns string) (sql.Result, error)
	GetRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (*models.Record, error)
	CountRecordTx(tx *sqlx.Tx, batchName, fingerprintValue, ns string) (int, error)
	GetRecordByFingerprintTx(tx *sqlx.Tx, batchName, ns, value string) (*models.Record, error)
	ListRecordTx(tx *sqlx.Tx, batchName, fingerprintValue, ns string, page, size int) ([]models.Record, error)
	CreateRecordTx(tx *sqlx.Tx, records []models.Record) (sql.Result, error)
	UpdateRecordTx(tx *sqlx.Tx, record *models.Record) (sql.Result, error)
	DeleteRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (sql.Result, error)

	// task
	CreateTask(task *models.Task) (sql.Result, error)
	UpdateTask(task *models.Task) (sql.Result, error)
	GetTask(traceId string) (*models.Task, error)
	DeleteTask(traceId string) (sql.Result, error)
	CountTask(task *models.Task) (int, error)

	GetTaskTx(tx *sqlx.Tx, traceId string) (*models.Task, error)
	CreateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error)
	UpdateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error)
	DeleteTaskTx(tx *sqlx.Tx, traceId string) (sql.Result, error)

	// callback
	GetCallback(name, namespace string) (*models.Callback, error)
	CreateCallback(callback *models.Callback) (sql.Result, error)
	UpdateCallback(callback *models.Callback) (sql.Result, error)
	DeleteCallback(name, ns string) (sql.Result, error)
	GetCallbackTx(tx *sqlx.Tx, name, namespace string) (*models.Callback, error)
	CreateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error)
	UpdateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error)
	DeleteCallbackTx(tx *sqlx.Tx, name, ns string) (sql.Result, error)

	// application
	CreateApplication(app *specV1.Application) (sql.Result, error)
	UpdateApplication(app *specV1.Application, oldVersion string) (sql.Result, error)
	DeleteApplication(name, namespace, version string) (sql.Result, error)
	GetApplication(name, namespace, version string) (*specV1.Application, error)
	ListApplication(name, namespace string, pageNo, pageSize int) ([]specV1.Application, error)
	CreateApplicationWithTx(tx *sqlx.Tx, app *specV1.Application) (sql.Result, error)
	UpdateApplicationWithTx(tx *sqlx.Tx, app *specV1.Application, oldVersion string) (sql.Result, error)
	DeleteApplicationWithTx(tx *sqlx.Tx, name, namespace, version string) (sql.Result, error)
	CountApplication(tx *sqlx.Tx, name, namespace string) (int, error)
	// system config
	GetSysConfig(tp, key string) (*models.SysConfig, error)
	ListSysConfig(tp string, page, size int) ([]models.SysConfig, error)
	ListSysConfigAll(tp string) ([]models.SysConfig, error)
	CreateSysConfig(sysConfig *models.SysConfig) (sql.Result, error)
	UpdateSysConfig(sysConfig *models.SysConfig) (sql.Result, error)
	DeleteSysConfig(tp, key string) (sql.Result, error)

	Shadow
}
