package models

import "strings"

type ListView struct {
	Total    int         `json:"total"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	Items    interface{} `json:"items,omitempty"`
}

type Filter struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name,omitempty"`
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
		f.Name = "%"
	} else if !strings.Contains(f.Name, "%") {
		f.Name = "%" + f.Name + "%"
	}
	return f.Name
}
