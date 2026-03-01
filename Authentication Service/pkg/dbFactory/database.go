package dbFactory

import "Authentication_Service/internal/config"

type IDatabase[T any] interface {
	Connect(orm IOrm[T]) (T, error)
}

type MySql[T any] struct {
	Cfg *config.Config
}

func (m MySql[T]) Connect(orm IOrm[T]) (T, error) {
	db, err := orm.Connect(m.Cfg)
	if err != nil {
		var zero T
		return zero, err
	}
	return db, nil
}
