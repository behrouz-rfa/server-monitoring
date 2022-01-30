package admincontrollers

import (
	"github.com/labstack/echo"
	"net/http"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/session"
	"server-monitoring/shared/view"
)

var (
	NetworkController networkControllerInterfaces = networkController{}
)

type networkController struct {
}



type networkControllerInterfaces interface {
	Index(ctx echo.Context) error
}

func (n networkController) Index(c echo.Context) error {
	sess := session.Instance(c.Request())

	if sess.Values[consts.IS_SUPER_ADMIN] == nil {
		session.Empty(sess)
		return c.Redirect(http.StatusFound, "/user/login")
		//v := view.New(c.Request())
		//v.Name = "front/user/login"
		//v.Vars["first_name"] = sess.Values["first_name"]
		//v.Render(c.Response())
	}
	v := view.New(c.Request())
	v.Name = "admin/network/index"
	v.Vars["URL"] = "/admin/network"

	////_, err = adminservice.HomeServices.NetInfo()
	//_, err = adminservice.HomeServices.HostInfo()
	//_, err = adminservice.HomeServices.DiskInfo()
	//GetUserList,err :=adminservice.HomeServices.GetUserList()

	//surveyList, err := adminservice.FormService.GetSuerveyList()

	//v.Vars["SurveyCount"] = len(surveyList)
	//v.Vars["UserCount"] = adminservice.AdminUserService.Count()
	//v.Vars["UserCountSurvey"] = adminservice.AdminUserService.Count()
	//var su survays.Survay
	//adminservice.AdminUserService.GetActiveSurvey(&su)
	//v.Vars["Survey"] = su
	v.RenderAdmin(c.Response())
	return nil
}