// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package eventsource

import (
	"sync"
)

// Manager manages the eventsource Messengers
type Manager struct {
	mutex sync.Mutex

	messengers map[int64]*Messenger
	connection chan struct{}
}

var manager *Manager

func init() {
	manager = &Manager{
		messengers: make(map[int64]*Messenger),
		connection: make(chan struct{}, 1),
	}
}

// GetManager returns a Manager and initializes one as singleton if there's none yet
func GetManager() *Manager {
	return manager
}

// Register message channel
func (m *Manager) Register(uid int64) <-chan *Event {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	if !ok {
		messenger = NewMessenger(uid)
		m.messengers[uid] = messenger
	}
	select {
	case m.connection <- struct{}{}:
	default:
	}
	m.mutex.Unlock()
	return messenger.Register()
}

// Unregister message channel
func (m *Manager) Unregister(uid int64, channel <-chan *Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	messenger, ok := m.messengers[uid]
	if !ok {
		return
	}
	if messenger.Unregister(channel) {
		delete(m.messengers, uid)
	}
}

// UnregisterAll message channels
func (m *Manager) UnregisterAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, messenger := range m.messengers {
		messenger.UnregisterAll()
	}
	m.messengers = map[int64]*Messenger{}
}

// BroadcastMessage sends a message to all connected users
func (m *Manager) BroadcastMessage(message *Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, messenger := range m.messengers {
		messenger.SendMessage(message)
	}
}
func (m *Manager) SendMessage(uid int64, message *Event) {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	m.mutex.Unlock()
	if ok {
		messenger.SendMessage(message)
	}
}

// SendMessageBlocking sends a message to a particular user
func (m *Manager) SendMessageBlocking(uid int64, message *Event) {
	m.mutex.Lock()
	messenger, ok := m.messengers[uid]
	m.mutex.Unlock()
	if ok {
		messenger.SendMessageBlocking(message)
	}
}
