package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"sync"

	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

type MyTemplate struct {
	templates *template.Template
}

// Render renders a template document
func (t *MyTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func MyWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	connectionPool.Lock()
	connectionPool.connections[ws] = struct{}{}

	defer func(connection *websocket.Conn){
		connectionPool.Lock()
		delete(connectionPool.connections, connection)
		connectionPool.Unlock()
	}(ws)

	connectionPool.Unlock()

	for {
		// Read
		_, msgRead, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("Message Read: %s\n", msgRead)

		// Write
		msgToBeWritten := []byte("Hello, Client")
		err = sendMessageToAllPool(msgToBeWritten)
		if err != nil {
			return err
		}
		fmt.Printf("Message Written: %s\n", msgToBeWritten)
	}
}

func sendMessageToAllPool(message []byte) error {
	connectionPool.RLock()
	defer connectionPool.RUnlock()
	for connection := range connectionPool.connections {
		err := connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = t

	e.Static("/", "assets")
	e.Static("/js", "js")
	e.Static("/css", "css")
	e.Static("/", "public")

	// Render
	e.GET("/ws", MyWS)

	e.Logger.Fatal(e.Start(":8000"))
}
