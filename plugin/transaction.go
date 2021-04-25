package plugin

type TransactionFactory interface {
	BeginTx() (interface{}, error)
	Commit(interface{})
	Rollback(interface{})
}
