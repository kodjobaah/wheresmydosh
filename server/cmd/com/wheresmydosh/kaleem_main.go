package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/kaleempeeroo/wheresmydosh/server/cmd/com/wheresmydosh/db"
)

func main() {

	dba := pg.Connect(&pg.Options{
		User:     "root",
		Password: "password",
		Database: "postgres",
		Addr:     "wheresmydosh-cluster.cluster-cnuxmvkomgbc.us-east-2.rds.amazonaws.com:5432",
	})
	defer dba.Close()
	var n int
	_, err := dba.QueryOne(pg.Scan(&n), "SELECT 1")
	fmt.Println(err)
	fmt.Println(n)

	db.CreateTableUser(dba)
}
