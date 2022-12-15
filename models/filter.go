package models

import "strings"

type Filter struct {
	PageNo   int    `form:"pageNo" json:"pageNo,omitempty"`
	PageSize int    `form:"pageSize" json:"pageSize,omitempty"`
	Name     string `form:"name,omitempty" json:"name,omitempty"`
}

type ListOptions struct {
	LabelSelector string `form:"selector,omitempty" json:"selector,omitempty"`
	NodeSelector  string `form:"nodeSelector,omitempty" json:"nodeSelector,omitempty"`
	FieldSelector string `form:"fieldSelector,omitempty" json:"fieldSelector,omitempty"`
	KeywordType   string `form:"keywordType,omitempty" json:"keywordType,omitempty"`
	Keyword       string `form:"keyword,omitempty" json:"keyword,omitempty"`
	Limit         int64  `form:"limit,omitempty" json:"limit,omitempty"`
	Continue      string `form:"continue,omitempty" json:"continue,omitempty"`
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

func (l *ListOptions) GetFuzzyKeyword() string {
	if l.Keyword == "" {
		return "%"
	} else if !strings.Contains(l.Keyword, "%") {
		return "%" + l.Keyword + "%"
	}
	return l.Keyword
}

func GetPagingParam(listOptions *ListOptions, resLen int) (start, end int) {
	start = 0
	end = resLen
	if listOptions.GetLimitNumber() > 0 {
		start = listOptions.GetLimitOffset()
		end = listOptions.GetLimitOffset() + listOptions.GetLimitNumber()
		if end > resLen {
			end = resLen
		}
		if start > resLen {
			start = 0
			end = 0
		}
	}
	return start, end
}
