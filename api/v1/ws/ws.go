package ws

import (
	"github.com/gofiber/contrib/v3/websocket"
)

type Ws struct{}

var WsSend = new(Ws)

func (*Ws) SendMessage(c *websocket.Conn) {
	// log.Println(c.Locals("allowed"))  // true
	// log.Println(c.Params("id"))       // 123
	// log.Println(c.Query("v"))         // 1.0
	// log.Println(c.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			// log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			// log.Println("write:", err)
			break
		}
	}
}
