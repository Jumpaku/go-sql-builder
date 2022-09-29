package main

import "go-sql-builder/sql"

func main() {
	sql.SelectAll().From(sql.Table("A"))
	"Select * from A"
	println("Hello")
}
