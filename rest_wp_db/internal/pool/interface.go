package pool

type Pool interface {
	Start()
	CancelTask(id string) error
	Stop()
}
