package admincontrollers

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"server-monitoring/apps/admin/adminservice"
	"server-monitoring/shared/consts"
	"server-monitoring/shared/session"
	"server-monitoring/shared/view"
)

var (
	RequestController requestControllerInterfaces = requestController{}
)

type requestController struct {
}

type requestControllerInterfaces interface {
	Index(ctx echo.Context) error
}

func (i requestController) Index(c echo.Context) error {
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
	v.Name = "admin/requests/index"
	v.Vars["URL"] = "/admin"

	key := c.QueryParam("key")
	if len(key) == 0 {
		key = ""
	}

	requests, err := adminservice.HomeServices.LoadRequestsFilter(1, key)
	//surveyList, err := adminservice.FormService.GetSuerveyList()
	if err != nil {
		fmt.Println(err)
	}
	v.Vars["requests"] = requests

	//v.Vars["SurveyCount"] = len(surveyList)
	//v.Vars["UserCount"] = adminservice.AdminUserService.Count()
	//v.Vars["UserCountSurvey"] = adminservice.AdminUserService.Count()
	//var su survays.Survay
	//adminservice.AdminUserService.GetActiveSurvey(&su)
	//v.Vars["Survey"] = su
	v.RenderAdmin(c.Response())

	return nil
}
