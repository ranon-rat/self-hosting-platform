package executionlogs

type ExecutionLogsRepoDB interface {
	Create(log *NewLog) error
	// creo que esto deberia de funcionar?
	Get(oldId, projectID int) ([]Logs, error)
	DeleteOldMessages(idProject int) error
}
