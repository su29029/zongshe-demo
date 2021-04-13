package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wen-qu/kit-xuesou-backend/general/db"
	"github.com/wen-qu/kit-xuesou-backend/user/endpoint"
	"github.com/wen-qu/kit-xuesou-backend/user/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	addr := flag.String("http", ":8090", "http listen address")
	flag.Parse()

	ctx := context.Background()

	srv := service.NewService()
	errChan := make(chan error)

	go func(){
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := endpoint.Endpoints{
		Login: endpoint.Login(srv),
		Register: endpoint.Register(srv),
		ReadProfile: endpoint.ReadProfile(srv),
		UpdateProfile: endpoint.UpdateProfile(srv),
	}

	go func(){
		log.Println("user service is running on port: ", *addr)
		db.Init()
		handler := NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*addr, handler)
	}()

	log.Fatalln(<-errChan)
}