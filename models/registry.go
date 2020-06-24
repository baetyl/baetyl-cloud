package models

import (
	"fmt"
	"reflect"
	"time"

	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/jinzhu/copier"
)

// Registry Registry
type Registry struct {
	Name              string    `json:"name,omitempty" validate:"omitempty,resourceName,nonBaetyl"`
	Namespace         string    `json:"namespace,omitempty"`
	Address           string    `json:"address"`
	Username          string    `json:"username"`
	Password          string    `json:"password,omitempty"`
	CreationTimestamp time.Time `json:"createTime,omitempty"`
	UpdateTimestamp   time.Time `json:"updateTime,omitempty"`
	Description       string    `json:"description"`
	Version           string    `json:"version,omitempty"`
}

// RegistryList Registry List
type RegistryList struct {
	Total       int          `json:"total"`
	ListOptions *ListOptions `json:"listOptions"`
	Items       []Registry   `json:"items"`
}

func (r *Registry) Equal(target *Registry) bool {
	return reflect.DeepEqual(r.Address, target.Address) &&
		reflect.DeepEqual(r.Username, target.Username) &&
		reflect.DeepEqual(r.Password, target.Password) &&
		reflect.DeepEqual(r.Description, target.Description)
}

func (r *Registry) ToSecret() *specV1.Secret {
	res := &specV1.Secret{
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	err := copier.Copy(res, r)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data = map[string][]byte{
		"password": []byte(r.Password),
		"username": []byte(r.Username),
		"address":  []byte(r.Address),
	}
	return res
}
