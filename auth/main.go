package main

import (
	auth "e_commerce/kitex_gen/auth/authservice"
	"log"
)

func Init() {

}

func main() {
	svr := auth.NewServer(new(AuthServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
