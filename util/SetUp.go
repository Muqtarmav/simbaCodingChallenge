package util

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

//Delete CreatedModels clears models that are created during tests
func DeleteCreatedModels(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}

	var entries []entity
	hookName := "cleanupHook"

	db.Callback().Create().After("gorm:create").Register(hookName,
		func(scope *gorm.Scope) {
			fmt.Printf("Inserted entities of %s with %s = %v\n",
				scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
			entries = append(entries, entity{table: scope.TableName(),
				keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
		})
	return func() {
		defer db.Callback().Create().Remove(hookName)
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from %s table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+"=?", entry.key).Delete("")
		}
		if !inTransaction {
			tx.Commit()
		}
	}
}
