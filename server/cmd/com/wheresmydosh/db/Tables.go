package db

import (
	"fmt"
	"time"

	pg "github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type User struct {
	TableName   struct{}  `sql:"users,alias:user_mapping"`
	ID          int       `sql:"item_id,pk"`
	UserID      string    `sql:"user_id,unique"`
	FirstName   string    `sql:"first_name"`
	MiddleName  string    `sql:"middle_name"`
	LastName    string    `sql:"last_name"`
	PhoneNumber string    `sql:"phone_number"`
	Email       string    `sql:"email"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
	IsActive    bool      `sql:"is_active"`
}

func CreateTableUser(db *pg.DB) error {
	/*opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	*/
	err := db.DropTable((*User)(nil), &orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})
	fmt.Printf("%v", err)

	createErr := db.CreateTable(&User{}, nil)
	if createErr == nil {
		fmt.Printf("Error creating table User %v", createErr)
		return createErr
	}
	fmt.Printf("Table User created successfully: %v", createErr)

	return nil

}
