package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockSQLResult struct {
	lastId int64
	affect int64
}

func (m *mockSQLResult) LastInsertId() (int64, error) {
	return m.lastId, nil
}

func (m *mockSQLResult) RowsAffected() (int64, error) {
	return m.affect, nil
}

func genBatchTestCase() *models.Batch {
	batch := &models.Batch{
		Name:            "87ed3dba2b8f11eabc62186590da6863",
		Namespace:       "default",
		Description:     "dest test",
		QuotaNum:        200,
		EnableWhitelist: 0,
		SecurityType:    "1",
		SecurityKey:     "111",
		CallbackName:    "",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		Labels:          map[string]string{"batch": "87ed3dba2b8f11eabc62186590da6863"},
		Fingerprint:     models.Fingerprint{},
	}
	return batch
}

func genRecordTestCase() *models.Record {
	return &models.Record{
		Name:             "11eaa870fa163e812e31765d7cfc2f01",
		Namespace:        "default",
		BatchName:        "87ed3dba2b8f11eabc62186590da6863",
		FingerprintValue: "123",
		Active:           0,
		NodeName:         "node name",
		ActiveIP:         "0.0.0.0",
		ActiveTime:       time.Now(),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
}

func TestDefaultRegisterService_GetBatch(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(batch, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.GetBatch(batch.Name, batch.Namespace)
	assert.NoError(t, err)
	assert.EqualValues(t, batch, res)
}

func TestDefaultRegisterService_GetBatch_DBErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(nil, fmt.Errorf("db error")).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.GetBatch(batch.Name, batch.Namespace)
	assert.NotNil(t, err)
}

func TestDefaultRegisterService_GetBatch_NilErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.GetBatch(batch.Name, batch.Namespace)
	assert.NotNil(t, err)
}

func TestDefaultRegisterService_CreateBatch(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateBatch(batch).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.CreateBatch(batch)
	assert.NoError(t, err)
	assert.EqualValues(t, batch, res)
}

func TestDefaultRegisterService_CreateBatch_CreateError(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().CreateBatch(batch).Return(nil, common.Error(common.ErrDatabase, common.Field("error", "create error"))).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateBatch(batch)
	assert.Error(t, err, "Problem with database operation.create error")
}

func TestDefaultRegisterService_CreateBatch_GetError(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().CreateBatch(batch).Return(nil, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(nil, common.Error(common.ErrDatabase, common.Field("error", "get error"))).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateBatch(batch)
	assert.Error(t, err, "Problem with database operation.get error")
}

func TestDefaultRegisterService_UpdateBatch(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().UpdateBatch(batch).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.UpdateBatch(batch)
	assert.NoError(t, err)
	assert.EqualValues(t, batch, res)
}

func TestDefaultRegisterService_GetRecord(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.GetRecord(record.BatchName, record.Name, record.Namespace)
	assert.NoError(t, err)
	assert.EqualValues(t, record, res)
}

func TestDefaultRegisterService_GetRecord_DBErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, fmt.Errorf("db error")).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.GetRecord(record.BatchName, record.Name, record.Namespace)
	assert.NotNil(t, err)
}

func TestDefaultRegisterService_GetRecord_NilErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.GetRecord(record.BatchName, record.Name, record.Namespace)
	assert.NotNil(t, err)
}

func TestDefaultRegisterService_GetRecordByFingerprint(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	mockObject.dbStorage.EXPECT().GetRecordByFingerprint(record.Name, record.Namespace, record.FingerprintValue).Return(record, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.GetRecordByFingerprint(record.Name, record.Namespace, record.FingerprintValue)
	assert.NoError(t, err)
	assert.EqualValues(t, record, res)
}

func TestDefaultRegisterService_CreateRecord(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(1, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateRecord([]models.Record{*record}).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateRecord(record)
	assert.NoError(t, err)
}

func TestDefaultRegisterService_CreateRecord_GetErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(nil,
		common.Error(common.ErrDatabase, common.Field("error", "get batch error"))).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateRecord(record)
	assert.Error(t, err, "Problem with database operation.get batch error")
}

func TestDefaultRegisterService_CreateRecord_CountErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(0, common.Error(common.ErrDatabase, common.Field("error", "count error"))).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateRecord(record)
	assert.Error(t, err, "Problem with database operation.count error")
}

func TestDefaultRegisterService_CreateRecord_CreateErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(1, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateRecord([]models.Record{*record}).Return(nil, common.Error(common.ErrDatabase, common.Field("error", "create record error"))).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateRecord(record)
	assert.Error(t, err, "Problem with database operation.create record error")
}

func TestDefaultRegisterService_CreateRecord_GetRecordErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(1, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(nil, common.Error(common.ErrDatabase, common.Field("error", "get record error"))).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateRecord([]models.Record{*record}).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	_, err = rs.CreateRecord(record)
	assert.Error(t, err, "Problem with database operation.get record error")
}

func TestDefaultRegisterService_DeleteBatch(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	res := mockSQLResult{
		lastId: 0,
		affect: 0,
	}
	mockObject.dbStorage.EXPECT().CountRecord(batch.Name, "%", batch.Namespace).Return(0, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().DeleteBatch(batch.Name, batch.Namespace).Return(&res, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	err = rs.DeleteBatch(batch.Name, batch.Namespace)
	assert.NotNil(t, err)

	res = mockSQLResult{
		lastId: 0,
		affect: 1,
	}
	err = rs.DeleteBatch(batch.Name, batch.Namespace)
	assert.NoError(t, err)
}

func TestDefaultRegisterService_DeleteRecord(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	res := mockSQLResult{
		lastId: 0,
		affect: 0,
	}
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().DeleteRecord(record.BatchName, record.Name, record.Namespace).Return(&res, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	err = rs.DeleteRecord(record.BatchName, record.Name, record.Namespace)
	assert.NotNil(t, err)

	res = mockSQLResult{
		lastId: 0,
		affect: 1,
	}
	err = rs.DeleteRecord(record.BatchName, record.Name, record.Namespace)
	assert.NoError(t, err)
}

func TestDefaultRegisterService_DeleteRecord_GetErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, fmt.Errorf("db error")).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	err = rs.DeleteRecord(record.BatchName, record.Name, record.Namespace)
	assert.Error(t, err, common.Error(common.ErrDatabase, common.Field("error", "db error")))
}

func TestDefaultRegisterService_DeleteRecord_DeleteErr(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().DeleteRecord(record.BatchName, record.Name, record.Namespace).Return(nil, fmt.Errorf("db error")).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	record.Active = common.Activated
	err = rs.DeleteRecord(record.BatchName, record.Name, record.Namespace)
	assert.Error(t, err, common.Error(common.ErrRegisterRecordActivated))
	record.Active = common.Inactivated
	assert.Error(t, err, common.Error(common.ErrDatabase, common.Field("error", "db error")))
}

func TestDefaultRegisterService_UpdateRecord(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetRecord(record.BatchName, record.Name, record.Namespace).Return(record, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().UpdateRecord(record).Return(nil, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.UpdateRecord(record)
	assert.NoError(t, err)
	assert.EqualValues(t, record, res)
}

func TestDefaultRegisterService_DownloadRecords(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	record := genRecordTestCase()
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(1, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().ListRecord(record.BatchName, "%", record.Namespace, 1, 1).Return([]models.Record{*record}, nil).AnyTimes()

	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.DownloadRecords(record.BatchName, record.Namespace)
	assert.NoError(t, err)
	assert.EqualValues(t, []byte(record.FingerprintValue+"\n"), res)
}

func TestDefaultRegisterService_GenRecordRandom(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	record := genRecordTestCase()
	batch := genBatchTestCase()
	batch.Name = record.BatchName
	mockObject.dbStorage.EXPECT().GetBatch(record.BatchName, record.Namespace).Return(batch, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CountRecord(record.BatchName, "%", record.Namespace).Return(1, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateRecord(gomock.Any()).Return(nil, nil).AnyTimes()

	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)
	res, err := rs.GenRecordRandom(record.Namespace, record.BatchName, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestDefaultRegisterService_ListBatch(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.dbStorage.EXPECT().ListBatch(batch.Namespace, page.Name, page.PageNo, page.PageSize).Return([]models.Batch{*batch}, nil).Times(2)
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)

	mockObject.dbStorage.EXPECT().CountBatch(batch.Namespace, "%").Return(1, nil).Times(1)

	res, err := rs.ListBatch(batch.Namespace, page)
	assert.NoError(t, err)
	assert.EqualValues(t, *batch, res.Items.([]models.Batch)[0])

	// bad case
	mockObject.dbStorage.EXPECT().CountBatch(batch.Namespace, "%").Return(0, fmt.Errorf("db err")).Times(1)
	_, err = rs.ListBatch(batch.Namespace, page)
	assert.Error(t, err)
}

func TestDefaultRegisterService_ListBatch_Err(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.dbStorage.EXPECT().ListBatch(batch.Namespace, page.Name, page.PageNo, page.PageSize).Return(nil, fmt.Errorf("db err")).Times(1)
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)

	_, err = rs.ListBatch(batch.Namespace, page)
	assert.NotNil(t, err)
}

func TestDefaultRegisterService_ListRecord(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	record := genRecordTestCase()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.dbStorage.EXPECT().ListRecord(batch.Name, page.Name, batch.Namespace, page.PageNo, page.PageSize).Return([]models.Record{*record}, nil).AnyTimes()
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)

	mockObject.dbStorage.EXPECT().CountRecord(batch.Name, "%", batch.Namespace).Return(1, nil).Times(1)

	res, err := rs.ListRecord(batch.Name, batch.Namespace, page)
	assert.NoError(t, err)
	assert.EqualValues(t, *record, res.Items.([]models.Record)[0])

	//bad case
	mockObject.dbStorage.EXPECT().CountRecord(batch.Name, "%", batch.Namespace).Return(0, fmt.Errorf("db err")).Times(1)
	_, err = rs.ListRecord(batch.Name, batch.Namespace, page)
	assert.Error(t, err)
}

func TestDefaultRegisterService_ListRecord_Err(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	batch := genBatchTestCase()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.dbStorage.EXPECT().ListRecord(batch.Name, page.Name, batch.Namespace, page.PageNo, page.PageSize).Return(nil, fmt.Errorf("db err")).Times(1)
	rs, err := NewRegisterService(mockObject.conf)
	assert.NoError(t, err)

	_, err = rs.ListRecord(batch.Name, batch.Namespace, page)
	assert.NotNil(t, err)
}
