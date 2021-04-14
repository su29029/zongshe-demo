package main

import (
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/wen-qu/kit-xuesou-backend/user/endpoint"
	"github.com/wen-qu/kit-xuesou-backend/user/transport"
	"net/http"
)

func NewHTTPServer(ctx context.Context, endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(authMiddleWare)

	s := r.PathPrefix("/user").Subrouter()

	s.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.Login,
		transport.DecodeLoginRequest,
		transport.Encode,
	))
	s.Methods("POST").Path("/register").Handler(httptransport.NewServer(
		endpoints.Register,
		transport.DecodeRegisterRequest,
		transport.Encode,
	))
	s.Methods("GET").Path("/readProfile").Handler(httptransport.NewServer(
		endpoints.ReadProfile,
		transport.DecodeReadProfileRequest,
		transport.Encode,
	))
	s.Methods("POST").Path("/updateProfile").Handler(httptransport.NewServer(
		endpoints.UpdateProfile,
		transport.DecodeUpdateProfileRequest,
		transport.Encode,
	))
	return r
}

func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}