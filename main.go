package main

import (
	"depmod/db"
	"depmod/routes"
)

func main() {
	db.Init()
	e := routes.Init()

	e.Logger.Fatal(e.Start(":8000"))
}
