package domain

type TaskRepostory interface {
	Insert(t *Task) (*Task, error)
	FindById(id TaskId) (*Task, error)
	RemoveById(id TaskId) error
	FindAll() []*Task
	FindAllCompleted() []*Task
	FindAllUncompleted() []*Task
}
