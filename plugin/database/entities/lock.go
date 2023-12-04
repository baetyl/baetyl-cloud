// Package entities 数据库存储基本结构与方法
package entities

type Lock struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Version string `db:"version"`
	Expire  int64  `db:"expire"`
}
