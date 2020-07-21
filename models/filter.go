package models

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

func (f *Filter) Format() {
	if f.Name == "" {
		f.Name = "%"
	} else {
		f.Name = "%" + f.Name + "%"
	}
	if f.PageNo <= 0 {
		f.PageNo = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
}

type AmisData struct {
	Count int         `json:"count,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
}
type AmisListView struct {
	Status string   `json:"status,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Data   AmisData `json:"data,omitempty"`
}
