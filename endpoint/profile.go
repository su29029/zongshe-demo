package endpoint

import (
	"github.com/wen-qu/kit-xuesou-backend/user/model"
	"github.com/wen-qu/kit-xuesou-backend/user/service"
)

func UpdateGeneralProfile(userService service.IUserService, uid string, profile *model.Profile) (bool, error) {

	currProfile, err := userService.InspectUser(model.InspectRequest{Uid: uid})

	if err != nil {
		return false, err
	}

	if currProfile.User == nil {
		return false, nil
	}

	if _, err := userService.UpdateUser(model.UpdateRequest{
		User: &model.User{
			Uid:      uid,
			Username: profile.Username,
			Tel:      profile.Tel,
			Email:    profile.Email,
			Sex:      profile.Sex,
			Age:      profile.Age,
			Address:  profile.Address,
			Img:      profile.Img,
		},
	}); err != nil {
		return false, err
	}

	return true, nil
}

func ReadGeneralProfile(userService service.IUserService, uid string, tel string) (*model.Profile, error) {
	var res = new(model.Profile)

	profile, err := userService.InspectUser(model.InspectRequest{
		Uid: uid,
		Tel: tel,
	})

	if err != nil {
		return nil, err
	}

	if profile.User == nil {
		return nil, nil
	}

	res.Username = profile.User.Username
	res.Tel = profile.User.Tel
	res.Img = profile.User.Img
	res.Address = profile.User.Address
	res.Age = profile.User.Age
	res.Email = profile.User.Email
	res.Sex = profile.User.Sex
	res.ClassNum = profile.User.ClassNum

	return res, nil
}
