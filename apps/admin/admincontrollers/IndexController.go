package admincontrollers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"server-monitoring/apps/admin/adminservice"
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
	Ws(ctx echo.Context) error
	MemoryInfo(ctx echo.Context) error
	DiskInfo(ctx echo.Context) error
	CpuInfo(ctx echo.Context) error
	CpuMemory(ctx echo.Context) error
	Block(ctx echo.Context) error
}

var (
	upgrader = websocket.Upgrader{}
)

func (i indexController) Block(ctx echo.Context) error {

	name := ctx.QueryParam("ip")
	if mssg, err := adminservice.IpTableService.BlockIP(name); err != nil {
		log.Fatalln(err)
		return err
	} else {
		return ctx.JSON(200, mssg)
	}

}
func (i indexController) Ws(ctx echo.Context) error {

	hub := adminservice.NewHub()
	go hub.Run()
	adminservice.ServeWs(hub, ctx.Response(), ctx.Request())
	//ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	//if err != nil {
	//	return err
	//}
	//defer ws.Close()
	//
	//for {
	//	// Write
	//
	//	//adminservice.HomeServices.NetInfo()
	//
	//	go func() {
	//		info, _ := net.IOCounters(true)
	//		for index, v := range info {
	//			fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	//			err := ws.WriteMessage(websocket.TextMessage, passhash.Serialize(v))
	//
	//			if err != nil {
	//				ctx.Logger().Error(err)
	//			}
	//		}
	//	}()
	//
	//
	//
	//	go func() {
	//		_, msg, err := ws.ReadMessage()
	//		if err != nil {
	//			ctx.Logger().Error(err)
	//		}
	//		fmt.Printf("%s\n", msg)
	//	}()
	//	// Read
	//
	//
	//}
	return nil
}
func (i indexController) CpuMemory(ctx echo.Context) error {
	value, err := adminservice.HomeServices.CpuInfo()
	memory, err := adminservice.HomeServices.MemoryInfo()
	if err != nil {
		return ctx.JSON(404, err)
	}
	items := make(map[string]interface{})
	items["cpu"] = value
	items["memory"] = memory.UsedPercent

	return ctx.JSON(200, items)
}
func (i indexController) CpuInfo(ctx echo.Context) error {
	value, err := adminservice.HomeServices.CpuInfo()
	if err != nil {
		return ctx.JSON(404, err)
	}
	return ctx.JSON(200, value)
}
func (i indexController) MemoryInfo(ctx echo.Context) error {
	memory, err := adminservice.HomeServices.MemoryInfo()
	if err != nil {
		return ctx.JSON(404, err)
	}
	return ctx.JSON(200, memory)
}
func (i indexController) DiskInfo(ctx echo.Context) error {
	disk, err := adminservice.HomeServices.DiskInfo()
	if err != nil {
		return ctx.JSON(404, err)
	}
	return ctx.JSON(200, disk)
}
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
	requests, err := adminservice.HomeServices.LoadRequests(1)
	//
	disk, err := adminservice.HomeServices.DiskInfo()
	host, err := adminservice.HomeServices.HostInfo()
	////_, err = adminservice.HomeServices.NetInfo()
	//_, err = adminservice.HomeServices.HostInfo()
	//_, err = adminservice.HomeServices.DiskInfo()
	//GetUserList,err :=adminservice.HomeServices.GetUserList()

	//surveyList, err := adminservice.FormService.GetSuerveyList()
	if err != nil {
		fmt.Println(err)
	}
	//
	v.Vars["requests"] = requests
	v.Vars["disk"] = disk
	v.Vars["host"] = host
	//v.Vars["SurveyCount"] = len(surveyList)
	//v.Vars["UserCount"] = adminservice.AdminUserService.Count()
	//v.Vars["UserCountSurvey"] = adminservice.AdminUserService.Count()
	//var su survays.Survay
	//adminservice.AdminUserService.GetActiveSurvey(&su)
	//v.Vars["Survey"] = su
	v.RenderAdmin(c.Response())

	return nil
}
