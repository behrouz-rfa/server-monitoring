package adminservice

import (
	"server-monitoring/domain/model"
	"server-monitoring/domain/users"
)

var (
	AdminUserService adminUserServiceInterface = &adminUserService{}
)

type adminUserServiceInterface interface {

	Get(*users.User) error
	Create(*users.User) error
	GetAllOrder(number int) (*model.Pagination, error)
}
type adminUserService struct {
}

func (u *adminUserService) GetAllOrder(number int) (*model.Pagination, error) {
	var ue users.User
	return ue.FindAll(number)
}

func (u *adminUserService) Create(user *users.User) error {
	return user.CreateUserWithPassword()
}

func (u *adminUserService) Get(user *users.User)  error {

	return user.UserByEmail(user.Email)
}
