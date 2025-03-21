package repository

import (
	"database/sql"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Repository interface {
	AddClient(conn *websocket.Conn) (string, error)
	RemoveClient(conn *websocket.Conn) (string, error)

	BroadcastMessage(message []byte) (bool, error)
}
type DbPostgress struct {
	Db *sql.DB
}
type webSocketRepository struct {
	Clientconnection map[*websocket.Conn]bool
	mu               sync.Mutex
}

func NewWebSocketRepo() Repository {
	return &webSocketRepository{

		Clientconnection: make(map[*websocket.Conn]bool),
	}

}

// BroadcastMessage implements Repository.
func (r *webSocketRepository) BroadcastMessage(message []byte) (bool, error) {
	panic("unimplemented")
}

// RemoveClient implements Repository.
func (r *webSocketRepository) RemoveClient(conn *websocket.Conn) (string, error) {
	panic("unimplemented")
}

//now lets create an abstruction of websocketRepository

// AddClient implements Repository.
func (r *webSocketRepository) AddClient(conn *websocket.Conn) (string, error) {

	// handle go routines
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clientconnection[conn] = true
	log.Println("clients Added")
	return "yes  i have added a Connection SuccessFully", nil
}

// BroadcastMessage implements Repository.
// func (d *DbPostgress) BroadcastMessage(message []byte) (bool, error) {
// 	panic("unimplemented")
// }

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
