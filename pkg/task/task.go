package task

type Task interface {
	Run() error
	Name() string
}
