package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sync"

	"github.com/barcek2281/MyEcho/internal/app/model"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var Upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ControllerWS struct {
	log       *log.Logger
	session   sessions.Store
	storage   *storage.Storage
	mutex     *sync.RWMutex
	clients   map[*websocket.Conn]string
	allow     map[string]map[string]bool
	broadcast chan *Message
}

type Message struct {
	Type  string   `json:"type"`
	Users []string `json:"users,omitempty"`
	From  string   `json:"from,omitempty"`
	To    string   `json:"to,omitempty"`
	Msg   string   `json:"message,omitempty"`
}

func NewControllerWS(logger *logrus.Logger, session sessions.Store, storage *storage.Storage) *ControllerWS {
	allow, err := storage.Allow().GetAllow()
	if err != nil {
		logger.Errorf("Error with allow %v", err)
		allow = make(map[string]map[string]bool)
	}
	logger.Infof("Allow: %v", allow)
	return &ControllerWS{
		mutex:     &sync.RWMutex{},
		clients:   make(map[*websocket.Conn]string),
		log:       logger,
		session:   session,
		storage:   storage,
		broadcast: make(chan *Message),
		allow:     allow,
	}

}

func (c *ControllerWS) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.log.Errorf("error with webscoket conn: %v", err)
			return
		}

		session, err := c.session.Get(r, sessionName)
		if err != nil {
			c.log.Errorf("error with sessions: %v", err)
			return
		}

		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			utils.Error(w, r, 404, errSessionTimeOut)
			c.log.Errorf("session: %v", err)
			return
		}

		m, err := c.storage.User().FindById(user_id)
		if err != nil {
			c.log.Errorf("user doesnt exist: %v", err)
			return
		}
		c.log.Infof("user with this IP, connected: %v", conn.RemoteAddr().String())

		c.mutex.Lock()
		c.clients[conn] = m.Login
		c.mutex.Unlock()

		c.broadcast <- &Message{Type: "users"}
		go c.readFromClient(conn)
	}
}

func (c *ControllerWS) ChatsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/chats.html")
		if err != nil {
			c.log.Errorf("Error with parsing file chatas: %v", err)
			return
		}

		session, err := c.session.Get(r, sessionName)
		if err != nil {
			c.log.Errorf("error with sessions: %v", err)
			return
		}

		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			utils.Error(w, r, 404, errSessionTimeOut)
			c.log.Errorf("session: %v", err)
			return
		}

		m, err := c.storage.User().FindById(user_id)
		if err != nil {
			c.log.Errorf("user doesnt exist: %v", err)
			return
		}

		err = t.Execute(w, m)
		if err != nil {
			c.log.Errorf("Error with execution: %v", err)
			return
		}
		c.log.Info("Handle /chats GET")
	}
}

func (c *ControllerWS) readFromClient(conn *websocket.Conn) {
	defer conn.Close()
	for {
		msg := new(Message)
		err := conn.ReadJSON(msg)
		if err != nil {
			c.log.Errorf("Error with reading from client ws, %v", err)
			break
		}

		var Message model.Messages
		Message.Receiver = msg.To
		Message.Sender = msg.From
		Message.Message = msg.Msg
		if msg.Type == "message" {
			if !(c.allow[msg.From][msg.To] == true && c.allow[msg.To][msg.From] == true) {
				msg.Type = "error"
				c.log.Infof("Message: %+v", msg)
				c.broadcast <- msg
				continue
			}
			err = c.storage.Msg().CreateMessage(&Message)
			if err != nil {
				c.log.Fatalf("Error to store data: %v", err)
			}
		} else if msg.Type == "invite" {
			mp, ok := c.allow[msg.From]
			if !ok || mp == nil {
				c.allow[msg.From] = make(map[string]bool)
			}
			c.allow[msg.From][msg.To] = true
			var allow model.Allow
			allow.EmailFirst = msg.From
			allow.EmailSecond = msg.To
			err := c.storage.Allow().AddAllow(allow)
			if err != nil {
				c.log.Errorf("error to add: %v", err)
			}
			
		} else if msg.Type == "accept" {
			mp, ok := c.allow[msg.From]
			if !ok || mp == nil {
				c.allow[msg.From] = make(map[string]bool)
			}
			c.allow[msg.From][msg.To] = true
			var allow model.Allow
			allow.EmailFirst = msg.From
			allow.EmailSecond = msg.To
			err := c.storage.Allow().AddAllow(allow)
			if err != nil {
				c.log.Errorf("error to accept: %v", err)
			}
			
		} else if msg.Type == "block" {
			mp, ok := c.allow[msg.From]
			if !ok || mp == nil {
				c.allow[msg.From] = make(map[string]bool)
			}
			c.allow[msg.From][msg.To] = false
			var allow model.Allow
			allow.EmailFirst = msg.From
			allow.EmailSecond = msg.To
			err := c.storage.Allow().RemoveAllow(allow)
			if err != nil {
				c.log.Errorf("error to remove: %v", err)
			}
		} else if msg.Type == "unblock" {
			mp, ok := c.allow[msg.From]
			if !ok || mp == nil {
				c.allow[msg.From] = make(map[string]bool)
			}
			c.allow[msg.From][msg.To] = true
			var allow model.Allow
			allow.EmailFirst = msg.From
			allow.EmailSecond = msg.To
			err := c.storage.Allow().AddAllow(allow)
			if err != nil {
				c.log.Errorf("error to unblock: %v", err)
			}
		}
		c.log.Infof("Message: %+v", msg)
		c.broadcast <- msg
	}

	c.mutex.Lock()
	delete(c.clients, conn)
	c.mutex.Unlock()
}

func (c *ControllerWS) WriteToClients() {
	for {
		//c.mutex.RLock()
		msg := <-c.broadcast
		if msg.Type == "message" {
			for client := range c.clients {
				go func() {
					if err := client.WriteJSON(msg); err != nil {
						c.log.Warnf("Error with sending message: %v", err)
					}
				}()
			}
		} else if msg.Type == "users" {
			var userList []string
			for _, username := range c.clients {
				if username != "" {
					userList = append(userList, username)
				}
			}
			userMessage, _ := json.Marshal(Message{Type: "users", Users: userList})

			for conn := range c.clients {
				err := conn.WriteMessage(websocket.TextMessage, userMessage)
				if err != nil {
					c.log.Errorf("error with something: %v", err)
				}
			}
		} else if msg.Type == "history" {
			messages, err := c.storage.Msg().GetMsg(msg.To, msg.From, 5)
			if err != nil {
				c.log.Errorf("error with getting messages: %v", err)
			}

			for conn, login := range c.clients {
				if login == msg.From {
					for _, message := range messages {
						historyMessage, _ := json.Marshal(Message{Type: "history", From: message.Sender, To: message.Receiver, Msg: message.Message})
						err := conn.WriteMessage(websocket.TextMessage, historyMessage)
						if err != nil {
							c.log.Errorf("error with something: %v", err)
						}
					}
				}
			}
		} else if msg.Type == "invite" {
			for conn, login := range c.clients {
				if login == msg.To {
					notification, _ := json.Marshal(Message{Type: "invite", From: msg.From, To: msg.To})
					err := conn.WriteMessage(websocket.TextMessage, notification)
					if err != nil {
						c.log.Errorf("error with something: %v", err)
					}
				}
			}
		} else if msg.Type == "error" {
			for conn, login := range c.clients {
				if login == msg.From {
					notification, _ := json.Marshal(Message{Type: "error", From: msg.From, To: msg.To})
					err := conn.WriteMessage(websocket.TextMessage, notification)
					if err != nil {
						c.log.Errorf("error with something: %v", err)
					}
				}
			}
		} else if msg.Type == "accept" {
			for conn, login := range c.clients {
				if login == msg.To {
					notification, _ := json.Marshal(Message{Type: "accept", From: msg.From, To: msg.To})
					err := conn.WriteMessage(websocket.TextMessage, notification)
					if err != nil {
						c.log.Errorf("error with something: %v", err)
					}
				}
			}
		} else if msg.Type == "block" {
			for conn, login := range c.clients {
				if login == msg.To {
					notification, _ := json.Marshal(Message{Type: "block", From: msg.From, To: msg.To})
					err := conn.WriteMessage(websocket.TextMessage, notification)
					if err != nil {
						c.log.Errorf("error with something: %v", err)
					}
				}
			}
		}else if msg.Type == "unblock" {
			for conn, login := range c.clients {
				if login == msg.To {
					notification, _ := json.Marshal(Message{Type: "unblock", From: msg.From, To: msg.To})
					err := conn.WriteMessage(websocket.TextMessage, notification)
					if err != nil {
						c.log.Errorf("error with something: %v", err)
					}
				}
			}
		}
	}
	//c.mutex.RUnlock()
}
