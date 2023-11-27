package lobby

import (
	"errors"
	"maps"
	"net"
	"sync"
)

type UserLobbyManager struct {
	mu          *sync.Mutex
	connections map[string]net.Conn
}

func NewUserLobbyManager() *UserLobbyManager {
	return &UserLobbyManager{
		mu:          &sync.Mutex{},
		connections: make(map[string]net.Conn),
	}
}

func (l *UserLobbyManager) AddUser(name string, conn net.Conn) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.connections[name]; ok {
		return errors.New("user already exists")
	}

	l.connections[name] = conn

	return nil
}

func (l *UserLobbyManager) RemoveUser(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.connections, name)
}

func (l *UserLobbyManager) GetUser(name string) net.Conn {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.connections[name]
}

func (l *UserLobbyManager) GetUsers() map[string]net.Conn {
	l.mu.Lock()
	defer l.mu.Unlock()

	users := make(map[string]net.Conn, len(l.connections))
	maps.Copy(users, l.connections)

	return users
}

func (l *UserLobbyManager) IterateOverUsers(fn func(string, net.Conn)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for name, conn := range l.connections {
		fn(name, conn)
	}
}
