package admincontrollers

import (
	"github.com/labstack/echo"
	"net/http"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/session"
	"server-monitoring/shared/view"
)

var (
	IndexController indexControllerInterfaces = indexController{}
)

type indexController struct {
}

type indexControllerInterfaces interface {
	Index(ctx echo.Context) error
}

const (

)

func (i indexController) Index(c echo.Context) error {
	//sess := session.Instance(c.Request())
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
	v.Name = "admin/index/index"
	v.Vars["URL"] = "/admin"
	//response,err :=adminservice.HomeServices.GetSsh()
	//GetUserList,err :=adminservice.HomeServices.GetUserList()

	//surveyList, err := adminservice.FormService.GetSuerveyList()
	//if err != nil {
	//
	//}
	//
	//v.Vars["surveyList"] = surveyList
	//v.Vars["SurveyCount"] = len(surveyList)
	//v.Vars["UserCount"] = adminservice.AdminUserService.Count()
	//v.Vars["UserCountSurvey"] = adminservice.AdminUserService.Count()
	//var su survays.Survay
	//adminservice.AdminUserService.GetActiveSurvey(&su)
	//v.Vars["Survey"] = su
	v.RenderAdmin(c.Response())

	return nil
}
