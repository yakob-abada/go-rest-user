package mapper

import "github.com/yakob-abada/go-rest-user/pkg/model"

func FromUserRequest(userReq *model.UserRequest) *model.User {
	return &model.User{
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		Nickname:  userReq.Nickname,
		Email:     userReq.Email,
		Password:  userReq.Password,
		Country:   userReq.Country,
	}
}
