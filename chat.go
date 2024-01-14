package main

import (
	"errors"
	"slices"
	"time"
)

type User struct {
	id       int64
	username string
}

type ChatRoom struct {
	id           int64
	title        string
	participants []User
	joinCode     *string
}

type Message struct {
	id      int64
	author  User
	content string
	date    time.Time
	readBy  []int64
}

type RoomFactory struct {
	owner *User
}

var chats []ChatRoom
var messagesByRooms map[int64][]Message

func (f *RoomFactory) PrivateRoom(user *User) (*ChatRoom, error) {
	room := ChatRoom{
		1,
		user.username,
		[]User{*f.owner, *user},
		nil,
	}
	chats = append(chats, room)
	return &room, nil
}

func (f *RoomFactory) EmptyGroupRoom(title string, joinCode string) (*ChatRoom, error) {
	room := ChatRoom{
		1,
		title,
		[]User{*f.owner},
		&joinCode,
	}
	chats = append(chats, room)
	return &room, nil
}

func (f *RoomFactory) GroupRoom(title string, participants *[]User, joinCode string) (*ChatRoom, error) {
	room := ChatRoom{
		1,
		title,
		append([]User{*f.owner}, *participants...),
		&joinCode,
	}
	chats = append(chats, room)
	return &room, nil
}

func Join(user *User, code string) (*ChatRoom, error) {
	idx := slices.IndexFunc(chats, func(cr ChatRoom) bool {
		return *cr.joinCode == code
	})
	if idx == -1 {
		return nil, errors.New("No chat available with provided code")
	}
	chats[idx].participants = append(chats[idx].participants, *user)
    return &chats[idx], nil
}

func Leave(room *ChatRoom) error {
    idx := slices.IndexFunc(chats, func(r ChatRoom) bool { return r.id == room.id })
	if idx == -1 {
		return errors.New("Chat no longer exists")
	}
    slices.Delete(chats, idx, idx+1)
    return nil
}

func Send(message string, author *User, room *ChatRoom) (*Message, error) {
    msg := Message{
        1,
        *author,
        message,
        time.Time{},
        []int64{},
    }
    messagesByRooms[room.id] = append(messagesByRooms[room.id], msg)
    return &msg, nil
}

func Read(msg *Message, user *User, room *ChatRoom) error {
    msgs := messagesByRooms[room.id]
    idx := slices.IndexFunc(msgs, func(m Message) bool { return m.id == msg.id })
	if idx == -1 {
		return errors.New("Chat no longer exists")
	}
    if slices.Contains(msg.readBy, user.id) {
        return nil
    }
    msg.readBy = append(msg.readBy, user.id)
    return nil
}

