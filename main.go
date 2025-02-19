package main

import (
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	webSocketByUserMutex sync.Mutex
	connections          = make(map[*websocket.Conn]bool)
	wsUpgrade            = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	e := echo.New()

	e.Pre(middleware.HTTPSRedirect())

	e.GET("/check", healthCheck)
	e.GET("/ws", wsFunc)
	e.GET("/ws2", wsFunc2)

	certFile := "fullchain.pem"
	keyFile := "privkey.pem"

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("Certificado não encontrado: %v", err)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("Chave privada não encontrada: %v", err)
	}

	e.Logger.Fatal(e.StartTLS(":7890", certFile, keyFile))

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

func wsFunc2(c echo.Context) error {
	ws, err := wsUpgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return err
	}
	defer ws.Close()

	webSocketByUserMutex.Lock()
	connections[ws] = true
	webSocketByUserMutex.Unlock()

	log.Println("Novo cliente conectado!")

	// Ouve por mensagens do cliente e faz o broadcast para todos os conectados
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			break
		}

		// Imprime a mensagem recebida
		log.Printf("Mensagem recebida: %s\n", msg)

		// Envia a mensagem para todos os WebSockets conectados
		broadcastMessage(msg)

		// Envia uma resposta de confirmação ao cliente
		err = ws.WriteMessage(websocket.TextMessage, []byte("Mensagem recebida"))
		if err != nil {
			log.Println("Erro ao enviar mensagem:", err)
			break
		}
	}

	// Remove a conexão WebSocket da lista quando o cliente se desconectar
	webSocketByUserMutex.Lock()
	delete(connections, ws)
	webSocketByUserMutex.Unlock()

	log.Println("Cliente desconectado!")

	return nil
}

// Função para fazer o broadcast da mensagem para todos os WebSockets conectados
func broadcastMessage(msg []byte) {
	// Envia a mensagem para todos os WebSockets conectados
	webSocketByUserMutex.Lock()
	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Erro ao enviar mensagem para WebSocket:", err)
			conn.Close()
			delete(connections, conn)
		}
	}
	webSocketByUserMutex.Unlock()
}
