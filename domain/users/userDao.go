package users

import (
	"fmt"
	"server-monitoring/databases"
	"server-monitoring/domain/model"
	loogers "server-monitoring/utils/looger"
	"server-monitoring/utils/rest_error"

	"sync"
)

const (
	selectCountItems     = "SELECT COUNT(*) FROM user;"
	queryInsert          = "INSERT INTO user(phone_number,verify_code)VALUES(?,?);"
	queryInsertUserPass  = "INSERT INTO user(username,email,password)VALUES(?,?,?);"
	queryFindUserByPhone = "SELECT id ,  phone_number ,  COALESCE(first_name, '') as first_name, COALESCE(last_name, '') as last_name ,status FROM user WHERE phone_number = ? AND status = 1;"
	queryFindUserByEmail = "SELECT id ,  COALESCE(phone_number ,'') as phone_number,username, email ,password,COALESCE(first_name, '') as first_name, COALESCE(last_name, '') as last_name ,status,is_super_admin FROM user WHERE email = ? AND status = 1;"
	queryFindUserByID    = "SELECT id ,  COALESCE(phone_number ,'') as phone_number,username, email ,password,COALESCE(first_name, '') as first_name, COALESCE(last_name, '') as last_name ,status,is_super_admin FROM user WHERE id = ? AND status = 1;"
	querySelectUser      = " SELECT id ,  COALESCE(phone_number ,'') as phone_number,username, email ,password,COALESCE(first_name, '') as first_name, COALESCE(last_name, '') as last_name ,status,is_super_admin FROM user LIMIT ?,?;"
	LimitOffset          = 20
)

func (u *User) Register(phone string) rest_error.RestErr {
	stmt, err := databases.Client.Prepare(queryInsert)
	if err != nil {
		loogers.Error("Error while preparing insert user query", err)
		return rest_error.NewInternalServerError(fmt.Sprintf("erro Insert %s", "DataBase Error"), err)
	}
	defer stmt.Close()
	insertUser, saveErr := stmt.Exec(u.PhoneNumber, u.VerifyCode)
	if saveErr != nil {
		loogers.Error("Error while insert user by phone", saveErr)
		return rest_error.NewInternalServerError(fmt.Sprintf("erro Insert %s", "DataBase Error"), err)
	}
	lastid, err := insertUser.LastInsertId()
	if err != nil {
		loogers.Error("Error while get last id", err)
		return rest_error.NewInternalServerError(fmt.Sprintf("erro Insert user: %s", "database error"), err)
	}

	u.Id = lastid
	return nil
}

func (u *User) Get() rest_error.RestErr {
	if err := u.ValidatePhone(); err != nil {
		return err
	}
	stmt, err := databases.Client.Prepare(queryFindUserByPhone)
	if err != nil {
		loogers.Error("Error while preparing find user", err)
		return rest_error.NewInternalServerError(fmt.Sprintf("erro login user: %s", "database error"), err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.PhoneNumber)
	if err := result.Scan(&u.Id, &u.PhoneNumber, &u.FirstName, &u.LastName, &u.Status); err != nil {
		loogers.Error("error when trying get user by phone", err)
		return rest_error.NewInternalServerError("error to find user", err)
	}

	return nil
}

func (u *User) UserByEmail(email string) error {

	stmt, err := databases.Client.Prepare(queryFindUserByEmail)
	if err != nil {
		loogers.Error("Error while preparing find user", err)
		return model.ErrCode
	}
	defer stmt.Close()

	result := stmt.QueryRow(email)

	if err := result.Scan(&u.Id, &u.PhoneNumber, &u.UserName, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Status, &u.IsSuperAdmin); err != nil {
		loogers.Error("error when trying get user by email", err)
		return model.ErrNoResult
	}

	return nil
}
func (u *User) UserById() error {

	stmt, err := databases.Client.Prepare(queryFindUserByEmail)
	if err != nil {
		loogers.Error("Error while preparing find user", err)
		return model.ErrCode
	}
	defer stmt.Close()

	result := stmt.QueryRow(u.Id)

	if err := result.Scan(&u.Id, &u.PhoneNumber, &u.UserName, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Status, &u.IsSuperAdmin); err != nil {
		loogers.Error("error when trying get user by email", err)
		return model.ErrNoResult
	}

	return nil
}

func (u *User) UserByIdC(userid int64, work chan<- User, wg *sync.WaitGroup) error {

	defer wg.Done()
	var user User
	stmt, err := databases.Client.Prepare(queryFindUserByID)
	if err != nil {
		loogers.Error("Error while preparing find user", err)
		return model.ErrCode
	}
	defer stmt.Close()

	result := stmt.QueryRow(userid)

	if err := result.Scan(&user.Id, &user.PhoneNumber, &user.UserName, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Status, &user.IsSuperAdmin); err != nil {
		loogers.Error("error when trying get user by email", err)
		return model.ErrNoResult
	}
	work <- user
	return nil
}

func (u *User) CreateUserWithPassword() error {

	smt, err := databases.Client.Prepare(queryInsertUserPass)
	if err != nil {
		loogers.Error("Error while preparing find user", err)
		return model.ErrCode
	}
	defer smt.Close()
	insertUser, saveError := smt.Exec(u.UserName, u.Email, u.Password)
	if saveError != nil {
		loogers.Error("Error while insert user by phone", saveError)
		return model.ErrCode
	}
	lastid, err := insertUser.LastInsertId()
	if err != nil {
		loogers.Error("Error while get last id", err)
		return model.ErrCode
	}
	u.Id = lastid
	return nil
}

// get all products
func (u *User) FindAll(pageNumber int) (*model.Pagination, error) {
	//prepare query

	stmt, err := databases.Client.Prepare(querySelectUser)
	if err != nil {
		loogers.Error("error when traying to prrepare product", err)
		return nil, model.ErrCode
	}
	defer stmt.Close()

	var startItem int
	if pageNumber == 1 {
		startItem = 0
	} else {
		startItem = (pageNumber - 1) * LimitOffset
	}
	//run query with params if params exists
	rows, errRow := stmt.Query(startItem, LimitOffset)
	if errRow != nil {
		loogers.Error("error when trying to prepare select product", err)
		return nil, model.ErrCode
	}
	//close rows after return
	defer rows.Close()
	//create list of product
	var paginateProduct model.Pagination
	var users []User
	for rows.Next() {
		var user User
		//scan product
		if err := rows.Scan(&user.Id, &user.PhoneNumber, &user.UserName, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Status, &user.IsSuperAdmin); err != nil {
			continue
		}

		//append files

		users = append(users, user)
	}

	//load items
	count, err := u.CountItemPerPage()
	if err != nil {

	}
	paginateProduct.Items = users
	paginateProduct.CurrentPage = pageNumber //current page

	paginateProduct.Count = count                           // number of items in table
	paginateProduct.Pages = count / LimitOffset             // number if page
	paginateProduct.Next = len(users) == LimitOffset        // check Next page exist
	paginateProduct.Previous = startItem >= LimitOffset     // check Previous page exist
	paginateProduct.ItemPerPage = LimitOffset               // Item per page loaded
	paginateProduct.ShowPerRow = paginateProduct.Pages > 10 // Item per page loaded
	paginateProduct.LastPage = paginateProduct.Pages        // Item per page loaded
	// Item per page loaded
	paginateProduct.PagesList = paginateProduct.GetPages() // Item per page loaded
	return &paginateProduct, nil
}

//get cpount Items
func (*User) CountItemPerPage() (int, error) {
	stmt, err := databases.Client.Prepare(selectCountItems)
	if err != nil {
		loogers.Error("error when prepare count product", err)
		return 0, model.ErrCode
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	var count int
	if errScan := row.Scan(&count); errScan != nil {
		loogers.Error("error when prepare count product", err)
		return 0, model.ErrCode
	}
	return count, nil
}
