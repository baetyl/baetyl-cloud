package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"time"
)

//go:generate mockgen -destination=../mock/service/register.go -package=service github.com/baetyl/baetyl-cloud/v2/service RegisterService

type RegisterService interface {
	GetBatch(name, ns string) (*models.Batch, error)
	CreateBatch(batch *models.Batch) (*models.Batch, error)
	UpdateBatch(batch *models.Batch) (*models.Batch, error)
	DeleteBatch(name, ns string) error
	GetRecord(batchName, recordName, ns string) (*models.Record, error)
	GetRecordByFingerprint(batchName, ns, value string) (*models.Record, error)
	CreateRecord(record *models.Record) (*models.Record, error)
	UpdateRecord(record *models.Record) (*models.Record, error)
	DeleteRecord(batchName, recordName, ns string) error
	DownloadRecords(batchName, ns string) ([]byte, error)
	GenRecordRandom(ns, batchName string, num int) ([]string, error)
	ListBatch(ns string, page *models.Filter) (*models.ListView, error)
	ListRecord(batchName, ns string, page *models.Filter) (*models.ListView, error)
}

type registerService struct {
	cfg          *config.CloudConfig
	modelStorage plugin.ModelStorage
	dbStorage    plugin.DBStorage
}

// NewRegisterService New Register Service
func NewRegisterService(config *config.CloudConfig) (RegisterService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	ds, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}
	return &registerService{
		cfg:          config,
		modelStorage: ms.(plugin.ModelStorage),
		dbStorage:    ds.(plugin.DBStorage),
	}, nil
}

func (r *registerService) GetBatch(name, ns string) (*models.Batch, error) {
	batch, err := r.dbStorage.GetBatch(name, ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if batch == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "batch"), common.Field("name", name))
	}
	return batch, nil
}

func (r *registerService) CreateBatch(batch *models.Batch) (*models.Batch, error) {
	_, err := r.dbStorage.CreateBatch(batch)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	res, err := r.GetBatch(batch.Name, batch.Namespace)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *registerService) UpdateBatch(batch *models.Batch) (*models.Batch, error) {
	_, err := r.dbStorage.UpdateBatch(batch)
	if err != nil {
		return nil, err
	}
	res, err := r.GetBatch(batch.Name, batch.Namespace)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *registerService) DeleteBatch(name, ns string) error {
	count, err := r.dbStorage.CountRecord(name, "%", ns)
	if err != nil {
		return err
	}
	if count > 0 {
		return common.Error(common.ErrRegisterDeleteRecord, common.Field("name", name))
	}
	res, err := r.dbStorage.DeleteBatch(name, ns)
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", err))
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if affect == 0 {
		return common.Error(common.ErrResourceNotFound, common.Field("type", "batch"), common.Field("name", name))
	}
	return nil
}

func (r *registerService) GetRecord(batchName, recordName, ns string) (*models.Record, error) {
	_, err := r.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	record, err := r.dbStorage.GetRecord(batchName, recordName, ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if record == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "record"), common.Field("name", recordName))
	}
	return record, nil
}

func (r *registerService) GetRecordByFingerprint(batchName, ns, value string) (*models.Record, error) {
	return r.dbStorage.GetRecordByFingerprint(batchName, ns, value)
}

func (r *registerService) CreateRecord(record *models.Record) (*models.Record, error) {
	batch, err := r.GetBatch(record.BatchName, record.Namespace)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "batch"), common.Field("name", record.BatchName))
	}
	count, err := r.dbStorage.CountRecord(record.BatchName, "%", record.Namespace)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if count >= batch.QuotaNum {
		return nil, common.Error(common.ErrRegisterQuotaNumOut, common.Field("num", batch.QuotaNum))
	}
	_, err = r.dbStorage.CreateRecord([]models.Record{*record})
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	res, err := r.GetRecord(record.BatchName, record.Name, record.Namespace)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *registerService) UpdateRecord(record *models.Record) (*models.Record, error) {
	_, err := r.dbStorage.UpdateRecord(record)
	if err != nil {
		return nil, err
	}
	res, err := r.GetRecord(record.BatchName, record.Name, record.Namespace)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *registerService) DeleteRecord(batchName, recordName, ns string) error {
	_, err := r.GetBatch(batchName, ns)
	if err != nil {
		return err
	}
	record, err := r.GetRecord(batchName, recordName, ns)
	if err != nil {
		return err
	}
	if record.Active == common.Activated {
		return common.Error(common.ErrRegisterRecordActivated)
	}
	res, err := r.dbStorage.DeleteRecord(batchName, recordName, ns)
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", err))
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return common.Error(common.ErrResourceNotFound, common.Field("type", "record"), common.Field("name", recordName))
	}
	return nil
}

func (r *registerService) DownloadRecords(batchName, ns string) ([]byte, error) {
	_, err := r.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	count, err := r.dbStorage.CountRecord(batchName, "%", ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	records, err := r.dbStorage.ListRecord(batchName, "%", ns, 1, count)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	data := ""
	for _, record := range records {
		data += record.FingerprintValue + "\n"
	}
	return []byte(data), nil
}

func (r *registerService) GenRecordRandom(ns, batchName string, num int) ([]string, error) {
	batch, err := r.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	count, err := r.dbStorage.CountRecord(batchName, "%", ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if num < 1 || (num+count) > batch.QuotaNum {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "num error"))
	}

	data := []string{}
	records := []models.Record{}
	for i := 0; i < num; i++ {
		fv := common.UUIDPrune()
		data = append(data, fv)
		record := models.Record{
			Name:             fv,
			Namespace:        ns,
			BatchName:        batchName,
			FingerprintValue: fv,
			NodeName:         fv,
			Active:           common.Inactivated,
			ActiveTime:       time.Unix(common.DefaultActiveTime, 0),
		}
		records = append(records, record)
	}
	_, err = r.dbStorage.CreateRecord(records)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return data, nil
}

func (r *registerService) ListBatch(ns string, page *models.Filter) (*models.ListView, error) {
	batchs, err := r.dbStorage.ListBatch(ns, page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	count, err := r.dbStorage.CountBatch(ns, page.Name)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return &models.ListView{
		Total:    count,
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Items:    batchs,
	}, nil
}

func (r *registerService) ListRecord(batchName, ns string, page *models.Filter) (*models.ListView, error) {
	records, err := r.dbStorage.ListRecord(batchName, page.Name, ns, page.PageNo, page.PageSize)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	count, err := r.dbStorage.CountRecord(batchName, page.Name, ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return &models.ListView{
		Total:    count,
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Items:    records,
	}, nil
}
