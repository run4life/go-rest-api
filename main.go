package main

import (
	"fmt"
	"github.com/glorinli/go-jwt-simple-auth/app"
	"github.com/glorinli/go-jwt-simple-auth/controllers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func init() {
	log.SetPrefix("simple-auth")
}

func main() {
	// 新建路由器
	router := mux.NewRouter()

	// 注册jwt认证的中间件
	router.Use(app.JwtAuthentication)

	// 注册路由
	router.HandleFunc("/api/account", controllers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/account/accesstoken", controllers.Login).Methods(http.MethodGet)
	router.HandleFunc("/api/account/me", controllers.Me).Methods(http.MethodGet)

	// 获取端口号
	port := os.Getenv("golang-jwt-simple-auth-port")
	if port == "" {
		port = "8001"
	}

	fmt.Println("Port is:", port)

	// 开始服务
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print("Fail to start server", err)
	}
}
