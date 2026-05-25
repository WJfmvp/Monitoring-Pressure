package main

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/router"
)

func main() {
	r := router.SetupRouter()
	db.InitDB()
	r.Run(":8080")
}
