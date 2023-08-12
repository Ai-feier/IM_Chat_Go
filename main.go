package main

import (
	"fmt"
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMYSQL()
	utils.InitRedis()

	fmt.Println("MAIN:", utils.DB)
	r := router.Router()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
