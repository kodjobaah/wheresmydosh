package db

import (
	"fmt"
	"time"

	pg "github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type User struct {
	ID          int       `sql:",pk"`
	UserID      string    `sql:"user_id,unique"`
	FirstName   string    `sql:"first_name"`
	MiddleName  string    `sql:"middle_name"`
	LastName    string    `sql:"last_name"`
	PhoneNumber string    `sql:"phone_number"`
	Email       string    `sql:"email"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
	IsActive    bool      `sql:"is_active"`
	//Cards       []*CardDetails // has many
	//Transfer []*Transfer
}

type CardDetails struct {
	ID            int `sql:",pk"`
	UserID        int
	Name          string
	Number        string
	Expiry        time.Time
	CardType      string `sql:"type"`
	SortCode      int
	AccountNumber int       `sql:"accnt_number"`
	CreatedAt     time.Time `sql:"created_at"`
	UpdatedAt     time.Time `sql:"updated_at"`
	IsActive      bool      `sql:"is_active"`
	User          User
}

type Transfer struct {
	ID            int `sql:",pk"`
	Amount        int
	SenderID      int         `sql:"from"`
	Sender        User        `pg:"fk:from"`
	ReceiverID    int         `sql:"send_to"`
	Receiver      *User       `pg:"fk:send_to"`
	CardDetailsID int         `sql:"card_info"`
	CardDetails   CardDetails `pg:"fk:card_info"`
	CreatedAt     time.Time   `sql:"created_at"`
	UpdatedAt     time.Time   `sql:"updated_at"`
	IsActive      bool        `sql:"is_active"`
}

func CreateDBTables(db *pg.DB) error {

	models := []interface{}{&User{}, &CardDetails{}, &Transfer{}}
	for _, model := range models {
		err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			fmt.Printf("%v", err)
		}

		createErr := db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		})

		if createErr != nil {
			fmt.Printf("Error creating table User %v", createErr)
			return createErr
		}
		fmt.Printf("Table User created successfully: %v", createErr)
	}

	return nil

}
