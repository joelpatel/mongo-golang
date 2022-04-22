package main

import (
	"fmt"
	"net/http"

	"github.com/joelpatel/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	user_controller := controllers.NewUserController()
	router.GET("/user/:id", user_controller.GetUser)
	router.POST("/user", user_controller.CreateUser)
	router.DELETE("/user/:id", user_controller.DeleteUser)
	// router.PUT("", )

	fmt.Printf("Starting server at post 9000...\n")
	http.ListenAndServe(":9000", router)
}
