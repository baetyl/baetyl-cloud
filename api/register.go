package api

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

func (api *API) GetBatch(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName := c.Param("batchName")
	return api.registerService.GetBatch(batchName, ns)
}

// CreateBatch create one node
func (api *API) CreateBatch(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	b := generateDefaultBatch(ns)
	batch, err := api.parseBatch(b, c)
	if err != nil {
		return nil, err
	}

	batch.Namespace = ns
	batch.Labels[common.LabelBatch] = batch.Name
	err = api.verifyBatch(batch)
	if err != nil {
		return nil, err
	}
	// create batch in database
	res, err := api.registerService.CreateBatch(batch)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) UpdateBatch(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName := c.Param("batchName")
	b, err := api.registerService.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	labels := b.Labels
	b.Labels = nil
	batch, err := api.parseBatch(b, c)
	if err != nil {
		return nil, err
	}
	if b.Labels == nil {
		b.Labels = labels
	} else {
		b.Labels[common.LabelBatch] = labels[common.LabelBatch]
	}
	return api.registerService.UpdateBatch(batch)
}

func (api *API) DeleteBatch(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName := c.Param("batchName")
	return nil, api.registerService.DeleteBatch(batchName, ns)
}

func (api *API) GenInitCmdFromBatch(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	_, err := api.registerService.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	cmd, err := api.genCmd(string(common.Batch), ns, batchName)
	if err != nil {
		return nil, err
	}
	return map[string]string{"cmd": cmd}, nil
}

func (api *API) CreateRecord(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	r := &models.Record{}
	r, err := api.parseRecord(r, c)
	if err != nil {
		return nil, err
	}
	if r.FingerprintValue == "" {
		r.FingerprintValue = common.UUIDPrune()
	}
	r.Namespace = ns
	if r.Name == "" {
		r.Name = r.FingerprintValue
	}
	if r.NodeName == "" {
		r.NodeName = strings.ToLower(r.FingerprintValue)
	}
	record := &models.Record{
		Namespace:        ns,
		BatchName:        batchName,
		Name:             r.Name,
		FingerprintValue: r.FingerprintValue,
		NodeName:         r.NodeName,
		ActiveTime:       time.Unix(common.DefaultActiveTime, 0),
	}
	return api.registerService.CreateRecord(record)
}

func (api *API) UpdateRecord(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName, recordName := c.Param("batchName"), c.Param("recordName")
	_, err := api.registerService.GetBatch(batchName, ns)
	if err != nil {
		return nil, err
	}
	record, err := api.registerService.GetRecord(batchName, recordName, ns)
	if err != nil {
		return nil, err
	}
	record, err = api.parseRecord(record, c)
	if err != nil {
		return nil, err
	}
	return api.registerService.UpdateRecord(record)
}

func (api *API) DeleteRecord(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName, recordName := c.Param("batchName"), c.Param("recordName")
	return nil, api.registerService.DeleteRecord(batchName, recordName, ns)
}

func (api *API) GetRecord(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	batchName, recordName := c.Param("batchName"), c.Param("recordName")
	return api.registerService.GetRecord(batchName, recordName, ns)
}

func (api *API) GenRecordRandom(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	param := &struct {
		Num int `json:num,omitempty`
	}{}
	err := c.LoadBody(param)
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	items, err := api.registerService.GenRecordRandom(ns, batchName, param.Num)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total": len(items),
		"items": items,
	}, nil
}

func (api *API) DownloadTemplateFile(c *common.Context) (interface{}, error) {
	batchName := c.Param("batchName")
	bytesBuffer := &bytes.Buffer{}

	writer := csv.NewWriter(bytesBuffer)
	_ = writer.Write([]string{batchName + "_registerNode1"})
	_ = writer.Write([]string{batchName + "_Node1"})
	_ = writer.Write([]string{batchName + "_1234"})

	writer.Flush()

	c.Header("Content-Disposition", "attachment;filename=template.csv")
	c.Data(http.StatusOK, "text/csv", bytesBuffer.Bytes())

	return nil, nil
}

func (api *API) GenRecordFromUpload(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	f, err := c.FormFile("file")
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
		return nil, err
	}

	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	var fvs []string
	for {
		fv, err := reader.Read()
		if fv == nil || err == io.EOF {
			break
		}
		fvs = append(fvs, fv[0])
	}
	if len(fvs) == 0 {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", "file is empty"))
		return nil, err
	}
	items, err := api.registerService.GenRecordFromUpload(ns, batchName, fvs)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": len(items),
		"items": items,
	}, nil
}

func (api *API) ListBatch(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.registerService.ListBatch(ns, params)
}

func (api *API) ListRecord(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	params := &models.Filter{}
	if err := c.Bind(&params); err != nil {
		return nil, err
	}
	params.Format()
	return api.registerService.ListRecord(batchName, ns, params)
}

func (api *API) DownloadRecords(c *common.Context) (interface{}, error) {
	ns, batchName := c.GetNamespace(), c.Param("batchName")
	return api.registerService.DownloadRecords(batchName, ns)
}

func (api *API) parseRecord(record *models.Record, c *common.Context) (*models.Record, error) {
	err := c.LoadBody(record)
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return record, err
}

func (api *API) parseBatch(batch *models.Batch, c *common.Context) (*models.Batch, error) {
	err := c.LoadBody(batch)
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return batch, err
}

func (api *API) verifyBatch(batch *models.Batch) error {
	if batch.Fingerprint.Type < common.FingerprintSN || batch.Fingerprint.Type > ((common.FingerprintMachineID<<1)-1) {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "Fingerprint.Type error"))
	}
	if batch.EnableWhitelist != common.EnableWhitelist && batch.EnableWhitelist != common.DisableWhitelist {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "EnableWhitelist must be 0 or 1"))
	}
	if batch.SecurityType == common.Token && batch.SecurityKey == "" {
		batch.SecurityKey = common.UUIDPrune()
	}
	if batch.CallbackName != "" {
		_, err := api.callbackService.Get(batch.CallbackName, batch.Namespace)
		return err
	}

	return nil
}

func generateDefaultBatch(ns string) *models.Batch {
	name := common.UUIDPrune()
	return &models.Batch{
		Name:            name,
		Namespace:       ns,
		QuotaNum:        200,
		EnableWhitelist: 1,
		SecurityType:    common.Token,
		Labels:          map[string]string{common.LabelBatch: name},
		Fingerprint: models.Fingerprint{
			Type:       common.FingerprintSN,
			SnPath:     path.Join(common.DefaultSNPath, common.DefaultSNFile),
			InputField: common.DefaultInputField,
		},
	}
}
