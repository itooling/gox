package oth

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	manager  WebsocketManager
	upgrader websocket.Upgrader
)

func init() {
	manager = WebsocketManager{
		data: make(map[string]*WebsocketClient),
	}

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			token, err := r.Cookie("token")
			if err != nil || token == nil {
				return false
			}
			return true
		},
	}
}

type WebsocketMessage struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type WebsocketClient struct {
	ID     string
	Conn   *websocket.Conn
	Input  chan []byte
	Output chan []byte
}

type WebsocketManager struct {
	lock sync.RWMutex
	data map[string]*WebsocketClient
}

func (w *WebsocketManager) Set(k string, v *WebsocketClient) {
	w.lock.Lock()
	w.data[k] = v
	w.lock.Unlock()
	w.Handle(k)
}

func (w *WebsocketManager) Get(k string) *WebsocketClient {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if v, ok := w.data[k]; ok && v != nil {
		return v
	}
	return nil
}

func (w *WebsocketManager) Del(k string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if v, ok := w.data[k]; ok && v != nil {
		close(w.data[k].Input)
		w.data[k].Conn.Close()
	}
	delete(w.data, k)
}

func (w *WebsocketManager) Keys() []string {
	w.lock.RLock()
	defer w.lock.RUnlock()
	ids := make([]string, 0)
	for k, _ := range w.data {
		ids = append(ids, k)
	}
	return ids
}

func (w *WebsocketManager) Read(k string) []byte {
	if v := w.Get(k); v != nil {
		return <-v.Input
	}
	return nil
}

func (w *WebsocketManager) Write(k string, msg []byte) {
	if v := w.Get(k); v != nil {
		v.Output <- msg
	}
}

func (w *WebsocketManager) Handle(k string) {
	//read msg
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
		}()
		for {
			if v := w.Get(k); v != nil {
				_, msg, err := v.Conn.ReadMessage()
				if err != nil {
					w.Del(k)
					break
				}
				v.Input <- msg
			}
		}
	}()

	//write msg
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
		}()
		for {
			if v := w.Get(k); v != nil {
				msg, ok := <-v.Output
				if !ok {
					w.Del(k)
					break
				}
				err := v.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					w.Del(k)
					break
				}
			}
		}
	}()
}

func (w *WebsocketManager) Link(k1, k2 string) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
		}()
		for {
			v1, v2 := w.Get(k1), w.Get(k2)
			if v1 != nil && v2 != nil {
				if msg, ok := <-v1.Input; ok {
					v2.Output <- msg
					continue
				}
				break
			}
		}
	}()

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
		}()
		for {
			v1, v2 := w.Get(k1), w.Get(k2)
			if v1 != nil && v2 != nil {
				if msg, ok := <-v2.Input; ok {
					v1.Output <- msg
					continue
				}
				break
			}
		}
	}()
}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wm := WebsocketMessage{}
	err = ws.ReadJSON(&wm)
	if err != nil || wm.From == "" {
		ws.Close()
		return
	}

	ws.SetCloseHandler(func(code int, text string) error {
		manager.Del(wm.From)
		return nil
	})

	manager.Set(wm.From, &WebsocketClient{
		ID:     wm.From,
		Conn:   ws,
		Input:  make(chan []byte),
		Output: make(chan []byte),
	})
}

func WsRead(k string) []byte {
	return manager.Read(k)
}

func WsWrite(k string, msg []byte) {
	manager.Write(k, msg)
}

func WsKeys() []string {
	return manager.Keys()
}

func WsLink(k1, k2 string) {
	manager.Link(k1, k2)
}
