package service

import (
	"database/sql"
	"fmt"
	"github.com/wen-qu/kit-xuesou-backend/general/db"
	"github.com/wen-qu/kit-xuesou-backend/general/errors"
	"github.com/wen-qu/kit-xuesou-backend/user/model"
	"github.com/wen-qu/kit-xuesou-backend/user/util"
	"strconv"
	"time"
)

type IUserService interface {
	InspectUser(req model.InspectRequest) (model.InspectResponse, error)
	UpdateUser(req model.UpdateRequest) (model.UpdateResponse, error)
	AddUser(req model.AddRequest) (model.AddResponse, error)
	DeleteUser(req model.DeleteRequest) (model.DeleteResponse, error)
}

type userService struct {}

func NewService() IUserService {
	return userService{}
}

func (u userService)InspectUser(req model.InspectRequest) (model.InspectResponse, error) {
	var user model.User
	var rsp model.InspectResponse
	var row *sql.Row

	if len(req.Uid) > 0 {
		row = db.GetUserDB().QueryRow("select uid, username, password, tel, age, sex, email, " +
			"address, class_num, img from user_profile_table where uid = ?", req.Uid)
	} else if len(req.Tel) > 0 { // some questions about different logistics last code and this code
		row = db.GetUserDB().QueryRow("select uid, username, password, tel, age, sex, email, " +
			"address, class_num, img from user_profile_table where tel = ?", req.Tel)
	} else {
		return rsp, errors.BadRequest("para:001", "missing parameters: uid or tel")
	}

	err := row.Scan(
		&user.Uid, &user.Username, &user.Password,
		&user.Tel, &user.Age, &user.Sex,
		&user.Email, &user.Address, &user.ClassNum, &user.Img)

	if err == sql.ErrNoRows {
		return rsp, nil
	} else if err != nil {
		return rsp, errors.InternalServerError("user.InspectUser:fatal:001", err.Error())
	}

	rsp.User = &user

	return rsp, nil
}

func (u userService)UpdateUser(req model.UpdateRequest) (model.UpdateResponse, error) {
	var rsp model.UpdateResponse

	currentUser, err := u.InspectUser(model.InspectRequest{
		Uid: req.User.Uid,
		Tel: req.User.Tel,
	})
	if err != nil {
		return rsp, errors.InternalServerError("user.UpdateUser:fatal:001", err.Error())
	}
	if currentUser.User == nil {
		return rsp, errors.Forbidden("user:001", "user not existed")
	}
	fmt.Println("req.User: ", req.User)
	fmt.Println("currentUser: ", currentUser.User)
	// not existed field, set it to the currentUser.User's field, existed field, don't replace it.
	if err := util.CopyStruct(req.User, currentUser.User); err != nil {
		return rsp, err
	}
	fmt.Println("after, req.User: ",req.User)
	_, err = db.GetUserDB().Exec("update user_profile_table set username = ?, password = ?, tel = ?, " +
		"age = ?, sex = ?, email = ?, address = ?, class_num = ?, img = ? where uid = ? ",
		req.User.Username, req.User.Password,
		req.User.Tel, req.User.Age, req.User.Sex,
		req.User.Email, req.User.Address, req.User.ClassNum,
		req.User.Img, req.User.Uid)

	if err != nil {
		return rsp, errors.InternalServerError("user.UpdateUser:fatal:003", err.Error())
	}

	rsp.Status = 200
	rsp.Msg = "success"

	return rsp, nil
}

func (u userService)AddUser(req model.AddRequest) (model.AddResponse, error) {
	var rsp model.AddResponse
	if req.User == nil {
		return rsp, errors.BadRequest("para:001", "missing parameters")
	}
	user, err := u.InspectUser(model.InspectRequest{
		Tel: req.User.Tel,
	})
	if err != nil {
		return rsp, errors.InternalServerError("user.AddUser:fatal:001", err.Error())
	}
	if user.User != nil {
		rsp.Status = 400
		rsp.Msg = "registered"
		return rsp, nil
	}


	uid := "user_" + strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(uid)
	if _, err := db.GetUserDB().Exec("insert into user_profile_table " +
		"(uid, username, password, tel, email, sex, age, address, class_num, img) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", uid, req.User.Username, req.User.Password,
		req.User.Tel, req.User.Email, req.User.Sex, req.User.Age, req.User.Address,
		req.User.ClassNum, req.User.Img); err != nil {
		return rsp, errors.InternalServerError("user.AddUser:fatal:002", err.Error())
	}

	// create table [uid]_user_class_table, [uid]_user_chatting_table, [uid]_user_evaluations_table
	tableName := uid + "_user_class_table"
	if _, err := db.GetUserDB().Exec("create table `" + tableName + "` (" +
		"`uid` varchar(18) not null," +
		"`class_id` varchar(19) not null," +
		"`bought_time` varchar(20) not null," +
		"`agency_id` varchar(20) not null" +
		") engine=innodb default charset=utf8"); err != nil {
		return rsp, errors.InternalServerError("user.AddUser:fatal:003", err.Error())
	}
	tableName = uid + "_user_chatting_table"
	if _, err := db.GetUserDB().Exec("create table `" + tableName + "` (" +
		"`chat_id` varchar(18) primary key not null," +
		"`uid` varchar(18) not null," +
		"`msg_num` int not null," +
		"`agency_icon` varchar(60)," +
		"`agency_id` varchar(20) not null," +
		"`agency_name` varchar(50) not null" +
		") engine=innodb default charset=utf8"); err != nil {
		return rsp, errors.InternalServerError("user.AddUser:fatal:004", err.Error())
	}
	tableName = uid + "_user_evaluation_table"
	if _, err := db.GetUserDB().Exec("create table `" + tableName + "` (" +
		"`evaluation_id` varchar(20) primary key not null," +
		"`favicon` varchar(60)," +
		"`rating` float not null," +
		"`username` varchar(50) not null," +
		"`agency_id` varchar(20) not null," +
		"`uid` varchar(18) not null," +
		"`class_id` varchar(19) not null," +
		"`detail` varchar(10000)," +
		"`pics` varchar(700)" +
		") engine=innodb default charset=utf8"); err != nil {
		return rsp, errors.InternalServerError("user.AddUser:fatal:005", err.Error())
	}

	rsp.Status = 200
	rsp.Msg = "success"

	return rsp, nil
}

func (u userService)DeleteUser(req model.DeleteRequest) (model.DeleteResponse, error) {
	var rsp model.DeleteResponse
	if len(req.Tel) == 0 && len(req.Uid) == 0 {
		return rsp, errors.BadRequest("para:002", "missing parameters")
	}

	goalUser, err := u.InspectUser(model.InspectRequest{
		Uid: req.Tel,
		Tel: req.Uid,
	})
	if err != nil {
		return rsp, errors.InternalServerError("user.DeleteUser:fatal:001", err.Error())
	}

	if goalUser.User == nil {
		return rsp, errors.Forbidden("user:001", "user not existed")
	}

	if len(req.Tel) > 0 {
		// _, err = db.GetUserDB().Exec("delete from user where tel = ?", req.Tel)
	} else if len(req.Uid) > 0 {
		// _, err = db.GetUserDB().Exec("delete from user where uid = ?", req.Uid)
	}

	// if err != nil {
	// 	return rsp, errors.InternalServerError("user-srv.UserSrv.DeleteUser:fatal:002", err.Error())
	// }

	// TODO: delete user's chatting table, evaluations table, class table

	rsp.Status = 200
	rsp.Msg = "success"

	return model.DeleteResponse{}, nil
}
