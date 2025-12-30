package executionlogs

type ExecutionLogsRepoDB interface {
	Create(log *NewLog) error
	// creo que esto deberia de funcionar?
	Get(firstID int) ([]Logs, error)
}
