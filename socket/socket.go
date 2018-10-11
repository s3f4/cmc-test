package socket

import (
	"sync"
	"html/template"
	"net/http"
	"path/filepath"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
	"flag"
	"time"
	"dehaa.com/core/controller"
	"encoding/json"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		//tracer:  trace.Off(), this code is close command line output
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			fmt.Print(r.clients)
			delete(r.clients, client)
			close(client.send)
		case message := <-r.forward:
			for client := range r.clients {
				client.send <- message
			}
		}
	}
}

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- message
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for {
		cmcController := controller.CoinMarketCapController{}
		marketData, _ := cmcController.GlobalData()
		m, _ := json.Marshal(marketData)
		_ = c.socket.WriteMessage(websocket.TextMessage, m)
		fmt.Println([]byte(time.Now().String()))
		time.Sleep(1 * time.Second)
	}

}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.once.Do(func() {
		th.templ = template.Must(template.ParseFiles(filepath.Join("web/tpl", th.filename)))
	})
	var x interface{}

	if x == nil {
		fmt.Println("x == nil")
	}

	th.templ.Execute(w, r)
}

func StartSocket() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
