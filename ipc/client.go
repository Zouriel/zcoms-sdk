package ipc

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"time"
)

// Client talks to the core daemon's IPC socket.
type Client struct {
	socket string
}

// New returns a client for the given socket path. Use DefaultSocketPath for the
// standard location.
func New(socket string) *Client { return &Client{socket: socket} }

// NewDefault returns a client pointed at ~/.config/zcoms/daemon.sock.
func NewDefault() (*Client, error) {
	p, err := DefaultSocketPath()
	if err != nil {
		return nil, err
	}
	return &Client{socket: p}, nil
}

// Available reports whether the daemon is listening.
func (c *Client) Available() bool {
	conn, err := net.DialTimeout("unix", c.socket, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// Do sends one request and reads one response. readDeadline bounds the wait for
// the response (zero = no deadline, for blocking ops like ask).
func (c *Client) Do(req Request, readDeadline time.Time) (Response, error) {
	conn, err := net.DialTimeout("unix", c.socket, 2*time.Second)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	line, err := json.Marshal(req)
	if err != nil {
		return Response{}, err
	}
	if _, err := conn.Write(append(line, '\n')); err != nil {
		return Response{}, err
	}
	if !readDeadline.IsZero() {
		_ = conn.SetReadDeadline(readDeadline)
	}
	respLine, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil && len(respLine) == 0 {
		return Response{}, err
	}
	var resp Response
	if err := json.Unmarshal(respLine, &resp); err != nil {
		return Response{}, err
	}
	if !resp.OK {
		return resp, errors.New(resp.Error)
	}
	return resp, nil
}

// Send delivers a one-way message; To is an @username or numeric chat id.
func (c *Client) Send(to, text string) (Response, error) {
	return c.Do(Request{Op: "send", To: to, Text: text}, time.Now().Add(30*time.Second))
}

// SendFile uploads a local file (waits for the upload).
func (c *Client) SendFile(to, path, caption string) (Response, error) {
	return c.Do(Request{Op: "sendfile", To: to, Path: path, Text: caption}, time.Now().Add(31*time.Minute))
}

// Read fetches the last count history messages of a chat.
func (c *Client) Read(to string, count int, download bool) (Response, error) {
	d := time.Now().Add(60 * time.Second)
	if download {
		d = time.Now().Add(5 * time.Minute)
	}
	return c.Do(Request{Op: "read", To: to, Count: count, Download: download}, d)
}

// Unread returns unread Telegram 1:1 messages from non-allow-listed senders.
func (c *Client) Unread() ([]UnreadItem, error) {
	resp, err := c.Do(Request{Op: "unread"}, time.Now().Add(60*time.Second))
	return resp.Unread, err
}

// MarkRead marks the given Telegram messages in a chat as read.
func (c *Client) MarkRead(chatID int64, messageIDs []int64) error {
	_, err := c.Do(Request{Op: "mark_read", ChatID: chatID, MessageIDs: messageIDs}, time.Now().Add(30*time.Second))
	return err
}

// Resolve maps an @username (or numeric id) to a Telegram chat id.
func (c *Client) Resolve(to string) (int64, error) {
	resp, err := c.Do(Request{Op: "resolve", To: to}, time.Now().Add(30*time.Second))
	return resp.ChatID, err
}
