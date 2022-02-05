package apicontrollers

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"server-monitoring/apps/admin/adminservice"
)

var (
	IndexController indexControllerInterfaces = indexController{}
)

type indexController struct {
}

type indexControllerInterfaces interface {
	MemoryInfo(ctx echo.Context) error
	DiskInfo(ctx echo.Context) error
	CpuInfo(ctx echo.Context) error
	CpuMemory(ctx echo.Context) error
}

var (
	upgrader = websocket.Upgrader{}
)

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
