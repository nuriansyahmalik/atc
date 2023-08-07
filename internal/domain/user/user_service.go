package user

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type UserService interface {
	Create(requestFormat UserRequestFormat, userID uuid.UUID) (user Users, err error)
	Login(requestFormat LoginRequestFormat) (user Users, err error)
}

type UserServiceImpl struct {
	UserRepository UserRepository
	Config         *configs.Config
}

func ProvideUserServiceImpl(userRepository UserRepository, config *configs.Config) *UserServiceImpl {
	return &UserServiceImpl{UserRepository: userRepository, Config: config}
}

func (u *UserServiceImpl) Create(requestFormat UserRequestFormat, userID uuid.UUID) (user Users, err error) {
	user, err = user.UsersRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}
	if err != nil {
		return user, failure.BadRequest(err)

	}
	err = u.UserRepository.Create(user)
	if err != nil {
		return
	}
	return
}
func (u *UserServiceImpl) Login(requestFormat LoginRequestFormat) (user Users, err error) {
	user, err = user.LoginRequestFormat(requestFormat)
	if err != nil {
		return
	}
	user, err = u.UserRepository.ResolveByEmail(user.Email)
	if err != nil {
		return user, failure.BadRequest(err)
	}
	return
}
