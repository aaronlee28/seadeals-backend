package helper

import (
	"fmt"
	"gorm.io/gorm"
)

func CommitOrRollback(tx *gorm.DB, err *error) {
	if *err != nil {
		fmt.Println("Rollback")
		tx.Rollback()
		return
	}
	fmt.Println("Commit")
	tx.Commit()
}
