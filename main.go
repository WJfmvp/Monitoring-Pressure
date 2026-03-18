package Monitoring_Pressure

import "Monitoring-Pressure/router"

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
