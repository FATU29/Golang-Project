package dbFactory

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/internal/constant"
)

func GetDatabase[T any](cfg *config.Config) IDatabase[T] {
	dbType := cfg.DbType

	if dbType == constant.MySql {
		return &MySql[T]{
			cfg: cfg,
		}
	}
	return nil
}
