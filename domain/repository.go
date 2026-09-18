package domain

type TaskRepostory interface {
	Insert(t *Task) (*Task, error)
	FindById(id TaskId) (*Task, error)
	RemoveById(id TaskId) error
	FindAll() ([]*Task, error)
	FindAllCompleted() ([]*Task, error)
	FindAllUncompleted() ([]*Task, error)
}
