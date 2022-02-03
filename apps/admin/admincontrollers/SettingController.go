package admincontrollers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"runtime"
	"server-monitoring/apps/admin/adminservice"
	"server-monitoring/domain/settings"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/passhash"
	"server-monitoring/shared/session"
	"server-monitoring/shared/view"
	"server-monitoring/utils"
	"strconv"
)

var (
	SettingController settingControllerInterfaces = settingController{}
)

type settingController struct {
}

type settingControllerInterfaces interface {
	Index(ctx echo.Context) error
	PostSetting(ctx echo.Context) error
	StartCapturing(ctx echo.Context) error
	StopCapturing(ctx echo.Context) error
}

func (i settingController) StartCapturing(c echo.Context) error {
	sess := session.Instance(c.Request())
	sess.Values[consts.STARTCAPTURING] = true
	//services.StartNew <- true
	sess.Save(c.Request(), c.Response())
	sess.AddFlash(view.Flash{Message: "Capture started", Class: view.FlashSuccess})
	sess.Save(c.Request(), c.Response())
	i.Index(c)
	return nil
}
func (i settingController) StopCapturing(c echo.Context) error {
	sess := session.Instance(c.Request())
	sess.Values[consts.STARTCAPTURING] = false
	sess.Save(c.Request(), c.Response())
	//services.Quit <- true
	sess.AddFlash(view.Flash{Message: "Capture stopped", Class: view.FlashSuccess})
	sess.Save(c.Request(), c.Response())
	i.Index(c)
	return nil
}
func (n settingController) Index(c echo.Context) error {
	sess := session.Instance(c.Request())

	if sess.Values[consts.IS_SUPER_ADMIN] == nil {
		session.Empty(sess)
		return c.Redirect(http.StatusFound, "/user/login")
		//v := view.New(c.Request())
		//v.Name = "front/user/login"
		//v.Vars["first_name"] = sess.Values["first_name"]
		//v.Render(c.Response())
	}
	var setting settings.Setting
	if err := adminservice.SettingService.Get(&setting); err != nil {
		fmt.Println(err)
	}
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Println(err)
	}

	v := view.New(c.Request())
	token := c.Get("csrf")
	v.Vars["csrf_token"] = token
	v.Name = "admin/setting/index"
	v.Vars["URL"] = "/admin/setting"
	v.Vars["setting"] = setting
	v.Vars["interfaces"] = interfaces
	v.Vars["settingid"] = setting.ID.String()
	if sess.Values[consts.STARTCAPTURING] == nil {
		v.Vars["isCapture"] = false
	} else {
		v.Vars["isCapture"] = sess.Values[consts.STARTCAPTURING].(bool)
	}

	v.RenderAdmin(c.Response())
	return nil
}
func (s settingController) PostSetting(c echo.Context) error {
	sess := session.Instance(c.Request())

	if validate, missingField := view.Validate(c.Request(), []string{"website"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(c.Request(), c.Response())
		s.Index(c)
		return nil
	}
	interfaces, err := net.Interfaces()

	website := c.FormValue("website")
	password := c.FormValue("password")
	username := c.FormValue("username")
	passwordRe := c.FormValue("password-re")
	email := c.FormValue("email")
	message := c.FormValue("message")
	meta := c.FormValue("meta")
	status, _ := strconv.Atoi(c.FormValue("customRadio"))
	keyword := c.FormValue("keyword")
	tel := c.FormValue("tel")
	phone := c.FormValue("phone")
	iface, _ := strconv.Atoi(c.FormValue("interface"))

	filter := c.FormValue("filter")
	languageId, _ := strconv.Atoi(c.FormValue("language_id"))

	if passwordRe != password {
		sess.AddFlash(view.Flash{Message: "Password not match", Class: view.FlashError})
		sess.Save(c.Request(), c.Response())
		s.Index(c)
		return nil
	}
	var setting settings.Setting
	if err := adminservice.SettingService.Get(&setting); err == nil {
		fmt.Println(err)
	}
	setting.SiteName = website
	setting.Email = email
	setting.Username = username
	setting.Meta = meta
	setting.Tel = tel
	setting.Phone = phone
	setting.Message = message
	if len(password) > 5 {
		setting.Password, _ = passhash.HashString(password)
	}
	setting.Status = 0
	setting.Filter = filter
	if runtime.GOOS == "windows" {
		for _, stat := range interfaces {
			if stat.Index == iface {
				for _, addr := range stat.Addrs {
					if utils.IsIpV4(addr.Addr) {
						setting.Interface = addr.Addr
						break
					}
				}
			}
		}

	} else {
		for _, stat := range interfaces {
			if stat.Index == iface {
				setting.Interface = stat.Name
				break
			}
		}
	}

	setting.Status = int8(status)

	setting.LanguageId = languageId
	setting.Keyword = keyword

	err = adminservice.SettingService.Create(&setting)

	if err != nil {
		sess.AddFlash(view.Flash{Message: "failed to save", Class: view.FlashError})
		sess.Save(c.Request(), c.Response())
		s.Index(c)
		return nil
	}
	sess.AddFlash(view.Flash{"Saved Successfulty", view.FlashSuccess})
	sess.Values[consts.LanguageId] = languageId
	sess.Values[consts.WebsiteName] = website
	sess.Save(c.Request(), c.Response())
	http.Redirect(c.Response(), c.Request(), "/admin/setting", http.StatusFound)

	return nil
}
