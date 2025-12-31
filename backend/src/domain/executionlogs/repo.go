package executionlogs

type ExecutionLogsRepoDB interface {
	Create(log *NewLog) error
	// creo que esto deberia de funcionar?
	Get(oldId int) ([]Logs, error)
}
