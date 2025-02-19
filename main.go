package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var webSocketByUserMutex sync.Mutex
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	e := echo.New()

	e.GET("/check", healthCheck)
	e.GET("/ws", wsFunc)

	//e.Logger.Fatal(e.StartTLS(":7890", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.StartTLS(
		":7890",
		"/etc/letsencrypt/live/mariadobairro.com.br/fullchain.pem",
		"/etc/letsencrypt/live/mariadobairro.com.br/privkey.pem"))

}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Server Available")
}

func wsFunc(c echo.Context) error {
	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return err
	}
	defer ws.Close()

	log.Println("Novo cliente conectado!")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			break
		}

		log.Printf("Mensagem recebida: %s\n", msg)

		err = ws.WriteMessage(websocket.TextMessage, []byte("Mensagem recebida"))
		if err != nil {
			log.Println("Erro ao enviar mensagem:", err)
			break
		}
	}

	return nil
}
