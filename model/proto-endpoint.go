package model

type Profile struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Tel      string `json:"tel"`
	Sex      int32  `json:"sex"`
	Age      int32  `json:"age"`
	Address  string `json:"address"`
	ClassNum int32  `json:"classNum"`
	Img      string `json:"img"`
}

type LoginRequest struct {
	Tel            string `json:"tel"`
	ValidationCode string `json:"validationCode"`
}

type LoginResponse struct {
	Status int32  `json:"status"`
	Uid    string `json:"uid"`
	Msg    string `json:"msg"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	Tel            string `json:"tel"`
	ValidationCode string `json:"validationCode"`
}

type RegisterResponse struct {
	Status int32  `json:"status"`
	Uid    string `json:"uid"`
	Msg    string `json:"msg"`
}

type ReadProfileRequest struct {
	Uid             string `json:"uid"`
	Tel             string `json:"tel"`
	InformationType int32  `json:"informationType"`
}

type ReadProfileResponse struct {
	Uid     string   `json:"uid"`
	Type    string   `json:"type"`
	Profile *Profile `json:"profile"`
	Msg     string   `json:"msg"`
	Status  int32    `json:"status"`
}

type UpdateProfileRequest struct {
	Uid             string   `json:"uid"`
	Profile         *Profile `json:"profile"`
	InformationType int32    `json:"informationType"`
}

type UpdateProfileResponse struct {
	Msg    string `json:"msg"`
	Status int32  `json:"status"`
	Type   string `json:"type"`
}