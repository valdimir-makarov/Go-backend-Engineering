package repository

import (
	"database/sql"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Repository interface {
	AddClient(conn *websocket.Conn, username string) (string, error)
	RemoveClient(conn *websocket.Conn) (string, error)
	GetUsername(conn *websocket.Conn) string
	BroadcastMessage(message []byte) (bool, error)
}
type DbPostgress struct {
	Db *sql.DB
}
type webSocketRepository struct {
	Clientconnection map[*websocket.Conn]string
	mu               sync.Mutex
}

func NewWebSocketRepo() Repository {
	return &webSocketRepository{

		Clientconnection: make(map[*websocket.Conn]string),
	}

}

// BroadcastMessage implements Repository.

// RemoveClient implements Repository.

//now lets create an abstruction of websocketRepository

// AddClient implements Repository.
func (r *webSocketRepository) AddClient(conn *websocket.Conn, username string) (string, error) {

	// handle go routines
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clientconnection[conn] = username
	log.Println("clients Added")
	return "yes  i have added a Connection SuccessFully", nil
}

// BroadcastMessage implements Repository.
func (r *webSocketRepository) BroadcastMessage(message []byte) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var brErr error
	success := true

	for it := range r.Clientconnection {

		err := it.WriteMessage(websocket.TextMessage, message)

		if err != nil {

			log.Fatalf("error while message Broad Cast")
			it.Close()

			delete(r.Clientconnection, it)
			brErr = err

			success = false
		}

	}
	return success, brErr

}
func (w *webSocketRepository) RemoveClient(conn *websocket.Conn) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	username := w.Clientconnection[conn]
	delete(w.Clientconnection, conn)
	log.Printf("Client removed: %s\n", username)
	return "Left ", nil
}

func (r *webSocketRepository) GetUsername(conn *websocket.Conn) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Clientconnection[conn]
}

// // RemoveClient implements Repository.
// func (d *DbPostgress) RemoveClient(conn *websocket.Conn) (string, error) {
// 	panic("unimplemented")
// }

// func NewRepo(repo *DbPostgress) Repository {

// 	db, err := sql.Open("Postgress", pslInfo)
// 	if err != nil {
// 		log.Fatalf("the Db connection error")
// 	}
// 	return &DbPostgress{
// 		Db: db,
// 	}

// }
