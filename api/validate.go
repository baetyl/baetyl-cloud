package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

// ValidateResourceForCreating validate when resource create
func (api *API) ValidateResourceForCreating(c *common.Context) (interface{}, error) {
	resource := struct {
		Name string `json:"name,omitempty"`
	}{}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf[:]))

	err = json.Unmarshal(buf, &resource)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	if !common.ValidNonBaetyl(resource.Name) {
		return nil, common.Error(common.ErrInvalidName, common.Field("nonBaetyl", "Name"))
	}
	return nil, nil
}

// ValidateResourceForDeleting validate when resource delete
func (api *API) ValidateResourceForDeleting(c *common.Context) (interface{}, error) {
	name := c.GetNameFromParam()
	if !common.ValidNonBaetyl(name) {
		return nil, common.Error(common.ErrInvalidName, common.Field("nonBaetyl", "Name"))
	}
	return nil, nil
}
