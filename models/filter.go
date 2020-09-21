package models

import "strings"

type ListView struct {
	Total    int         `json:"total,omitempty"`
	PageNo   int         `json:"pageNo,omitempty"`
	PageSize int         `json:"pageSize,omitempty"`
	Items    interface{} `json:"items,omitempty"`
}

type Filter struct {
	PageNo   int    `form:"pageNo,omitempty"`
	PageSize int    `form:"pageSize,omitempty"`
	Name     string `form:"name,omitempty"`
}

func (f *Filter) GetLimitNumber() int {
	if f.PageNo <= 0 {
		f.PageNo = 1
	}
	return f.PageNo
}

func (f *Filter) GetLimitOffset() int {
	return f.PageSize
}

func (f *Filter) GetFuzzyName() string {
	if f.Name == "" {
		f.Name = "%"
	} else if !strings.Contains(f.Name, "%") {
		f.Name = "%" + f.Name + "%"
	}
	return f.Name
}
