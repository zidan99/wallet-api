package db

import "gorm.io/gorm"

func Use(dbs ...*gorm.DB) *gorm.DB {
	for _, v := range dbs {
		if v != nil {
			return v
		}
	}

	return nil
}
