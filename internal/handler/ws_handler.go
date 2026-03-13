package handler

import (
	"cinema-booking/backend/internal/ws"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) Handle(c echo.Context) error {
	showtimeID := c.QueryParam("showtime_id")
	if showtimeID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "showtime_id is required",
		})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &ws.Client{
		Conn:       conn,
		ShowtimeID: showtimeID,
	}

	h.hub.Register(client)
	defer h.hub.Unregister(client)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
