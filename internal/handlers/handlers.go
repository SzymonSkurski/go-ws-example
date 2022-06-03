package handlers

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"), //will scan given directory for *.jet
	jet.InDevelopmentMode(),             //this should be a env variable
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse define response sent back from websocket
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades conncection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	//3th argument is response header
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>connected to server</small></em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response) //will marshall and pass json to client
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			err = nil
			//ingonre
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// Home home page handler
func Home(w http.ResponseWriter, r *http.Request) {
	//here could pass data map instead of nil
	// err := renderPage(w, "home.jet", nil)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	var response WsJsonResponse
	response.Message = "hello"
	jData, err := json.Marshal(response)
	if err != nil {
		log.Panicln(err) //handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan
		switch e.Action {
		case "username":
			//get list of all usernames and send it back
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.Message = fmt.Sprintf("some message, and action was %s", e.Action)
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn) //delete current socket user
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%v</strong>: %v", e.Username, e.Message)
			broadcastToAll(response)
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renderPage render page from jet template and write response
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	fmt.Println("Rendering: `" + tmpl + "`");
	view, err := views.GetTemplate(tmpl) //read temple from html/<tmpl>.jet file
	if err != nil {
		log.Panicln(err)
		return err //handle as 404 error - not found
	}

	err = view.Execute(w, data, nil) //pass data (map) to template and write
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}
