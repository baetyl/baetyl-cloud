package models

import "strings"

type ListView struct {
	Total    int         `json:"total"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	Items    interface{} `json:"items,omitempty"`
}

type Filter struct {
	PageNo   int    `form:"pageNo" json:"pageNo"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Name     string `form:"name,omitempty" json:"name,omitempty"`
}

type ListOptions struct {
	LabelSelector string `form:"selector,omitempty" json:"selector,omitempty"`
	FieldSelector string `form:"fieldSelector,omitempty" json:"fieldSelector,omitempty"`
	Limit         int64  `form:"limit,omitempty" json:"limit,omitempty"`
	Continue      string `form:"continue,omitempty" json:"omitempty"`
	Filter        `json:",inline"`
}

func (f *Filter) GetLimitOffset() int {
	if f.PageNo <= 0 {
		f.PageNo = 1
	}
	return (f.PageNo - 1) * f.GetLimitNumber()
}

func (f *Filter) GetLimitNumber() int {
	return f.PageSize
}

func (f *Filter) GetFuzzyName() string {
	if f.Name == "" {
		return "%"
	} else if !strings.Contains(f.Name, "%") {
		return "%" + f.Name + "%"
	}
	return f.Name
}
