package router

import (
	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"os"
	"server-monitoring/apps/admin/admincontrollers"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/session"
)

const CSRF_TOKEN_HEADER = "X-XSRF-TOKEN"
const CSRF_KEY = "csrf"

func NewRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	en := casbin.NewEnforcer("config/casbin_auth_model.conf", "config"+string(os.PathSeparator)+"policy.csv")

	enforcer := Enforcer{enforcer: en}

	admin := e.Group("/admin")
	admin.Use(enforcer.Enforce)
	admin.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:" + "csrf",
	}))

	admin.GET("", admincontrollers.IndexController.Index)
	admin.GET("/", admincontrollers.IndexController.Index)
	admin.GET("/memoryinfo", admincontrollers.IndexController.MemoryInfo)
	admin.GET("/cpuinfo", admincontrollers.IndexController.CpuInfo)
	admin.GET("/mem-cpu", admincontrollers.IndexController.CpuMemory)
	admin.GET("/network", admincontrollers.NetworkController.Index)
	admin.GET("/setting", admincontrollers.SettingController.Index)
	admin.POST("/setting", admincontrollers.SettingController.PostSetting)
	admin.POST("/setting/capture-start", admincontrollers.SettingController.StartCapturing)
	admin.POST("/setting/capture-stop", admincontrollers.SettingController.StopCapturing)
	admin.GET("/requests", admincontrollers.RequestController.Index)
	admin.GET("/block", admincontrollers.IndexController.Block)

	e.GET("/ws", admincontrollers.IndexController.Ws)

	e.GET("/user/login", admincontrollers.UserController.LoginGET)
	e.POST("/user/login", admincontrollers.UserController.LoginPost)
	e.GET("/user/register", admincontrollers.UserController.RegisterGET)
	e.GET("/user/logout", admincontrollers.UserController.LogoutGET)
	e.POST("/user/register", admincontrollers.UserController.RegisterPost)

	e.Static("/assets", "assets")

	return e
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() *echo.Echo {
	return NewRouter()
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() *echo.Echo {
	return NewRouter()

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

type Enforcer struct {
	enforcer *casbin.Enforcer
}

func (e *Enforcer) Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//user, _, _ := c.Request().BasicAuth()
		sess := session.Instance(c.Request())

		go func() {
			//clientIP := utils.FromRequest(c.Request())

			if sess.Values[consts.USER_ID] != nil {
				//userVisi.ID = sess.Values[consts.USER_ID]
			}
			//userVisi.Url = c.Request().URL.Path
			//userVisi.Ip = clientIP
			//adminservice.VisitedService.Insert(&userVisi)
		}()

		if sess.Values[consts.IS_SUPER_ADMIN] != nil {
			method := c.Request().Method
			//path := c.Request().URL.Path
			//c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+sess.Values[consts.SUPER_ADMIN].(string))
			//c.Response().WriteHeader(http.StatusOK)
			result := e.enforcer.Enforce("admin", "/bar", method)

			if result {
				return next(c)
			}

		}

		return echo.ErrForbidden
	}
}
