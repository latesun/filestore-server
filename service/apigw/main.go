package main

import (
	"filestore-server/service/apigw/route"
)

func main() {
	r := route.Router()
	panic(r.Run(":8080"))
}
