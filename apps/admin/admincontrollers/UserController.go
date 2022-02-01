package admincontrollers

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
	"github.com/labstack/echo"
	"log"
	"net/http"
	Jwtconfig "server-monitoring/domain/Jwt"
	"server-monitoring/domain/settings"

	"server-monitoring/apps/admin/adminservice"
	"server-monitoring/domain/model"
	"server-monitoring/domain/users"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/passhash"
	"server-monitoring/shared/session"
	"server-monitoring/shared/view"
	"strconv"
	"time"
)

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
)

var (
	UserController userControllerInterface = &userController{}
)

type userControllerInterface interface {
	UsersGet(echo.Context) error
	LogoutGET(echo.Context) error
	LoginGET(echo.Context) error
	LoginPost(echo.Context) error
	RegisterGET(echo.Context) error
	RegisterPost(echo.Context) error
}
type userController struct {
}

func (u userController) UsersGet(c echo.Context) error {
	v := view.New(c.Request())
	v.Name = "admin/user/index"
	v.Vars["URL"] = "/admin"
	//trying to get products
	page := c.QueryParams().Get("page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		pageNumber = 1
	}

	items, err := adminservice.AdminUserService.GetAllOrder(pageNumber)
	if err != nil {

	}

	v.Vars["paginate"] = items
	v.RenderAdmin(c.Response())
	return nil
}

func (u userController) LogoutGET(c echo.Context) error {
	// Get session
	sess := session.Instance(c.Request())

	// If user is authenticated
	if sess.Values[consts.USER_ID] != nil {
		session.Empty(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(c.Request(), c.Response())
	}

	http.Redirect(c.Response(), c.Request(), "/user/login", http.StatusFound)
	return nil
}

func (u userController) RegisterGET(c echo.Context) error {
	sess := session.Instance(c.Request())

	// Display the view
	v := view.New(c.Request())
	v.Name = "front/user/register"
	v.Vars["token"] = csrfbanana.Token(c.Response(), c.Request(), sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, c.Request().Form, v.Vars)
	v.Render(c.Response())
	return nil
}

func (u userController) LoginGET(c echo.Context) error {
	// Get session
	sess := session.Instance(c.Request())

	// Display the view
	v := view.New(c.Request())
	v.Name = "front/user/login"
	v.Vars["token"] = csrfbanana.Token(c.Response(), c.Request(), sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, c.Request().Form, v.Vars)
	v.Render(c.Response())
	return nil
}

func (u userController) LoginPost(c echo.Context) error {
	// Get session
	sess := session.Instance(c.Request())

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Brute force register prevented")
		http.Redirect(c.Response(), c.Request(), "/user/login", http.StatusFound)
		return errors.New("bad request")
	}

	// Validate with required fields
	if validate, missingField := view.Validate(c.Request(), []string{"username", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(c.Request(), c.Response())
		u.LoginGET(c)
		return nil
	}

	// Form values
	username := c.FormValue("username")
	password := c.FormValue("password")

	var user settings.Setting
	user.Username = username
	// Get database result
	err := adminservice.SettingService.Login(&user)

	// Determine if user exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(c.Request(), c.Response())
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(c.Request(), c.Response())
	} else if passhash.MatchString(user.Password, password) {
		if user.Status != 1 {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(c.Request(), c.Response())
		} else {

			session.Empty(sess)
			sess.AddFlash(view.Flash{"Login successful!", view.FlashSuccess})
			sess.Values[consts.USER_ID] = user.ID.String()
			if user.IsSuperAdmin == consts.IS_SUPER_ADMIN {
				// Login successfully
				// Set custom claims
				claims := &Jwtconfig.JwtCustomClaims{
					user.Username,
					true,
					jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
					},
				}

				// Create token with claims
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				// Generate encoded token and send it as response.
				t, err := token.SignedString([]byte("secret"))
				if err != nil {
					return err
				}

				sess.Values[consts.IS_SUPER_ADMIN] = user.IsSuperAdmin == consts.IS_SUPER_ADMIN
				sess.Values[consts.TOKEN] = t
			}
			sess.Values[consts.USER_EMAIL] = user.Email
			sess.Values[consts.USER_MAME] = user.Username

			sess.Values[consts.FUll_NAME] = fmt.Sprintf("%s", user.Username)
			sess.Save(c.Request(), c.Response())
			if user.IsSuperAdmin == consts.IS_SUPER_ADMIN {
				http.Redirect(c.Response(), c.Request(), "/admin", http.StatusFound)
			} else {
				http.Redirect(c.Response(), c.Request(), "/", http.StatusFound)
			}
			return nil
		}
	} else {
		//loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(c.Request(), c.Response())
	}

	// Show the login page again
	u.LoginGET(c)
	return nil
}

func (u userController) RegisterPost(c echo.Context) error {
	// Get session
	sess := session.Instance(c.Request())

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Brute force register prevented")
		http.Redirect(c.Response(), c.Request(), "/user/register", http.StatusFound)
		return errors.New("bad request")
	}

	// Validate with required fields
	if validate, missingField := view.Validate(c.Request(), []string{"email", "username", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(c.Request(), c.Response())
		u.RegisterGET(c)
		return errors.New("bad request")
	}

	// Form values
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	mo_number := c.FormValue("mo_number")

	password_has, errp := passhash.HashString(password)

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(c.Request(), c.Response())
		http.Redirect(c.Response(), c.Request(), "/user/register", http.StatusFound)
		return errp
	}
	var user users.User
	user.Email = email

	err := adminservice.AdminUserService.Get(&user)
	if err == model.ErrNoResult {
		user.Email = email
		user.Password = password_has
		user.UserName = username
		user.PhoneNumber = mo_number
		err := adminservice.AdminUserService.Create(&user)

		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(c.Request(), c.Response())
			http.Redirect(c.Response(), c.Request(), "/user/register", http.StatusFound)
			return nil
		} else {
			sess.AddFlash(view.Flash{"Account created successfully for: " + email, view.FlashSuccess})
			sess.Save(c.Request(), c.Response())
			http.Redirect(c.Response(), c.Request(), "/user/login", http.StatusFound)
			return nil
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(c.Request(), c.Response())
	} else { // Else the user already exists
		sess.AddFlash(view.Flash{"Account already exists for: " + email, view.FlashError})
		sess.Save(c.Request(), c.Response())
	}
	return nil
}

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		//sess.Values[sessLoginAttempt] = 1
	} else {
		//	sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}
