package mapper

import "github.com/yakob-abada/go-rest-user/pkg/model"

func UserResponseMapper(user model.User) model.UserResponse {
	return model.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
	}
}

func UsersResponseMapper(users []model.User) []model.UserResponse {
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponseMapper(user)
	}

	return userResponses
}
