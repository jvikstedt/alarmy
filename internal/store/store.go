package store

type Store interface {
	Project() ProjectStore
}

type ProjectStore interface {
	GetAll(projects interface{}) error
	Create(project interface{}) error
	FindByID(id int, project interface{}) error
	CountAll() (int, error)
	Clear()
}
