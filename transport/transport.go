package transport

import (
	"context"
	"encoding/json"
	"github.com/wen-qu/kit-xuesou-backend/general/errors"
	"github.com/wen-qu/kit-xuesou-backend/user/model"
	"github.com/wen-qu/kit-xuesou-backend/user/util"
	"net/http"
)

func DecodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.BadRequest("para:001", "missing or invalid parameters")
	}
	return req, nil
}

func DecodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req model.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.BadRequest("para:001", "missing or invalid parameters")
	}
	return req, nil
}

func DecodeReadProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req model.ReadProfileRequest
	urlQueryString := r.URL.Query().Encode()
	if err := util.DecodeGetParameters(urlQueryString, &req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeUpdateProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req model.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.BadRequest("para:001", "missing or invalid parameters")
	}
	return req, nil
}



func Encode(ctx context.Context, w http.ResponseWriter, rsp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(rsp)
}