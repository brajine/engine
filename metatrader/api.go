package metatrader

// Rebuild json access methods for all structs in file
// easyjson -all <file>.go

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/websocket"
)

// WebSockets
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Run API server
func (f *Factory) startAPIServer(addr string) {
	// HTTP server to serve JSON data
	e := echo.New()
	e.GET("/api/stats", f.StatsAPIHandler)
	e.GET("/api/rest/:page", f.RestAPIHandler)
	e.GET("/api/wss/:page", f.WssAPIHandler)
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	f.log.Fatal(e.Start(addr).Error())
}

// StatsAPIHandler is a handler for server state api
func (f *Factory) StatsAPIHandler(c echo.Context) error {
	st := f.exportState()
	return c.JSON(http.StatusOK, st)
}

// RestAPIHandler is serving REST API calls
func (f *Factory) RestAPIHandler(c echo.Context) error {
	page := c.Param("page")

	// Check if page exists
	acc := f.PageExist(page)
	if acc == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, acc)
}

// WssAPIHandler is serving WebSocket connections
func (f *Factory) WssAPIHandler(c echo.Context) error {
	page := c.Param("page")

	// Check if page exists
	acc := f.PageExist(page)
	if acc == nil {
		return c.NoContent(http.StatusNotFound)
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Add new connection to page viewers pool, and send him update
	acc.AddViewer(ws)
	return nil
}
