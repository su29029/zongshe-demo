package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/wen-qu/kit-xuesou-backend/general/errors"
	"github.com/wen-qu/kit-xuesou-backend/user/model"
	"github.com/wen-qu/kit-xuesou-backend/user/service"
	"regexp"
)

type Endpoints struct {
	Login endpoint.Endpoint
	Register endpoint.Endpoint
	UpdateProfile endpoint.Endpoint
	ReadProfile endpoint.Endpoint
}

func Login(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var rsp model.LoginResponse
		req := request.(model.LoginRequest)

		if len(req.Tel) == 0 || len(req.ValidationCode) == 0 {
			return rsp, errors.BadRequest("para:001", "missing parameter")
		}
		if ok, _ := regexp.Match("^1[3-9]\\d{9}$", []byte(req.Tel)); !ok {
			return rsp, errors.BadRequest("para:002", "invalid parameter: tel")
		}
		// TODO: check validation code
		loginRsp, err := userService.InspectUser(model.InspectRequest{
			Tel:      req.Tel,
		})

		if err != nil {
			return rsp, err
		}
		if loginRsp.User != nil && len(loginRsp.User.Uid) > 0 {
			rsp.Uid = loginRsp.User.Uid
			rsp.Status = 200
			rsp.Msg = "success"
			// TODO: generate token
			return rsp, nil
		}

		_, err = Register(userService)(ctx, model.RegisterRequest{
			Tel:            req.Tel,
			ValidationCode: req.ValidationCode,
		})
		if err != nil {
			return rsp, err
		}

		return Login(userService)(ctx, req)
	}
}

func Register(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var rsp model.RegisterResponse
		req := request.(model.RegisterRequest)

		if len(req.Tel) == 0 || len(req.ValidationCode) == 0 {
			return rsp, errors.BadRequest("para:001", "missing parameter")
		}
		if ok, _ := regexp.Match("^1[3-9]\\d{9}$", []byte(req.Tel)); !ok {
			return rsp, errors.BadRequest("para:002", "invalid parameter: tel")
		}

		// TODO: check the validation code.
		//rspCheck, err := SecClient.CheckValidation(ctx, &security.CheckValidationRequest{
		//	Code: req.ValidationCode,
		//	Tel:  req.Tel,
		//})
		//if err != nil {
		//	return err
		//}
		//if rspCheck.Status == 401 {
		//	rsp.Status = 401
		//	rsp.Msg = "invalid validation code"
		//}

		regRsp, err := userService.AddUser(model.AddRequest{
			User: &model.User{
				Tel: req.Tel,
			},
		})

		if err != nil {
			return rsp, err
		}

		if regRsp.Status == 400 {
			return rsp, errors.Forbidden("user:001", "registered")
		}

		rsp.Status = 200
		rsp.Msg = "success"

		return rsp, nil
	}
}

func UpdateProfile(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var rsp model.UpdateProfileResponse
		req := request.(model.UpdateProfileRequest)

		if len(req.Uid) == 0 {
			return rsp, errors.BadRequest("para:001", "missing parameters")
		}

		switch req.InformationType {
		case 1: // general
			fmt.Println(req.Profile)
			ok, err := UpdateGeneralProfile(userService, req.Uid, req.Profile)
			if err != nil {
				return rsp, err
			}
			if !ok {
				return rsp, errors.Forbidden("user:001", "user not existed")
			}
			rsp.Type = "general"
			rsp.Status = 200
			rsp.Msg = ""
		case 2: // order
		case 3: // discount
		case 4: // likes
		case 5: // order_review
		case 6: // classes
		case 7: // collections
		}

		return rsp, nil
	}
}

func ReadProfile(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var rsp model.ReadProfileResponse
		req := request.(model.ReadProfileRequest)

		if len(req.Uid) == 0 && len(req.Tel) == 0 {
			return rsp, errors.BadRequest("para:001", "missing parameters")
		}
		fmt.Println(req.Uid, req.Tel)
		if ok, _ := regexp.Match("^user_\\d{10}$", []byte(req.Uid)); len(req.Uid) > 0 && !ok {
			return rsp, errors.BadRequest("para:002", "invalid parameter: uid")
		}
		if ok, _ := regexp.Match("^1[3-9]\\d{9}$", []byte(req.Tel)); len(req.Tel) > 0 && !ok {
			return rsp, errors.BadRequest("para:002", "invalid parameter: tel")
		}
		switch req.InformationType {
		case 1: // general
			var err error
			rsp.Profile, err = ReadGeneralProfile(userService, req.Uid, req.Tel)
			if err != nil {
				return rsp, err
			}
			rsp.Uid = req.Uid
			rsp.Type = "general"
		case 2: // order
		case 3: // discount
		case 4: // likes
		case 5: // order_review
		case 6: // classes
		case 7: // collections
		}
		return nil, nil
	}
}