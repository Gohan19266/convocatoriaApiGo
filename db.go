package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func db() *sql.DB {
	usuario := "uysv9t18xsovpgmn"
	pass := "lqXoWZythDxqtkIJbcOp"
	host := "tcp(birvxmqrxnbzpdvm0ogz-mysql.services.clever-cloud.com)"
	nombreBaseDeDatos := "birvxmqrxnbzpdvm0ogz"
	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL")
	return db
}
