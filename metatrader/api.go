package metatrader

// Rebuild json access methods for all structs in file
// easyjson -all <file>.go

import (
	"net/http"

	_ "engine/docs" // docs generated by swag-cli

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware" // echo-swagger middleware
	echoSwagger "github.com/swaggo/echo-swagger"
)

// WebSockets
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// StateEntry used to
type StateEntry struct {
	Page       string `json:"page" example:"my-test-page"`
	Started    string `json:"started" example:"2020-12-20 23:10:01"`
	UpdateFreq string `json:"updateFreq" example:"minute"`
}

// StateData used to export state information through /api/state
type StateData struct {
	Online   int          `json:"online" example:"1"`
	Accounts []StateEntry `json:"accounts"`
}

// Run API server
func (f *Factory) startAPIServer(addr string) {
	// HTTP server to serve JSON data
	e := echo.New()
	e.GET("/api/stats", f.StatsAPIHandler)
	e.HEAD("/api/stats", f.StatsAPIHandler)
	e.GET("/api/rest/:page", f.RestAPIHandler)
	e.GET("/api/wss/:page", f.WssAPIHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler) // including images etc
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	f.log.Fatal(e.Start(addr).Error())
}

// StatsAPIHandler is a handler for server state api
// @Summary Provide actual list of connected accounts
// @Produce json
// @Success 200 {object} StateData
// @failure 500 {string} Server internal error
// @Router /api/stats [get]
func (f *Factory) StatsAPIHandler(c echo.Context) error {
	st := f.exportState()
	return c.JSON(http.StatusOK, st)
}

// RestAPIHandler is serving REST API calls
// @Summary Provide actual data on connected account
// @Produce json
// @Param page path string true "Account Page name"
// @Success 200 {object} Account
// @failure 404 {string} Page not found
// @failure 500 {string} Server internal error
// @Router /rest/{page} [get]
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
// @Summary Provide actual data on connected account via WebSocket connection
// @Produce json
// @Param page path string true "Account Page name"
// @Success 200 {object} Account
// @failure 404 {string} Page not found
// @failure 500 {string} Server internal error
// @Router /wss/{page} [get]
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
	acc.SendUpdateToViewer(ws)
	return nil
}
