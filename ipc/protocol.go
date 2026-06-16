// Package ipc is the shared protocol between the zcoms core daemon (which owns
// the single Telegram session) and the component processes (bridge/triage/
// errands). It speaks newline-delimited JSON over the daemon's Unix socket: one
// Request line in, one Response line out (the `subscribe` op streams Events).
package ipc

import (
	"os"
	"path/filepath"
)

const socketName = "daemon.sock"

// DefaultSocketPath returns ~/.config/zcoms/daemon.sock (the core daemon's IPC
// socket). Components dial this to reuse the daemon's Telegram session.
func DefaultSocketPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "zcoms", socketName), nil
}

// Request is one command sent to the daemon. Op selects which fields are used.
type Request struct {
	Op       string `json:"op"`                 // send|sendfile|read|ask|unread|mark_read|resolve|errand_*
	To       string `json:"to,omitempty"`       // @username or numeric chat id / errand target
	Text     string `json:"text,omitempty"`     // message body / question / caption
	Path     string `json:"path,omitempty"`     // local file path (sendfile)
	Count    int    `json:"count,omitempty"`    // history messages (read)
	Download bool   `json:"download,omitempty"` // download media in a read

	// mark_read
	ChatID     int64   `json:"chat_id,omitempty"`
	MessageIDs []int64 `json:"message_ids,omitempty"`

	// Errand ops.
	Brief     string `json:"brief,omitempty"`
	Deliver   bool   `json:"deliver,omitempty"`
	AutoStart bool   `json:"auto_start,omitempty"`
	ID        string `json:"id,omitempty"`
}

// Message is one history message returned by the "read" op (mirrors the fields
// `zc tg chat` prints).
type Message struct {
	MessageID int64  `json:"message_id"`
	ChatID    int64  `json:"chat_id"`
	Date      int64  `json:"date"`
	Outgoing  bool   `json:"outgoing"`
	Sender    string `json:"sender"`
	Kind      string `json:"kind"`
	Text      string `json:"text"`
	File      string `json:"file,omitempty"`
}

// UnreadItem is one unread 1:1 message from a non-allow-listed sender, returned
// by the "unread" op (Telegram only — components merge WhatsApp via the sidecar).
type UnreadItem struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
	When   int64  `json:"when"` // unix seconds
	ChatID int64  `json:"chat_id"`
	MsgID  int64  `json:"msg_id"`
}

// Response is the daemon's reply to a Request.
type Response struct {
	OK        bool         `json:"ok"`
	MessageID int64        `json:"message_id,omitempty"`
	ChatID    int64        `json:"chat_id,omitempty"`
	Reply     string       `json:"reply,omitempty"`
	Label     string       `json:"label,omitempty"`
	Messages  []Message    `json:"messages,omitempty"`
	Unread    []UnreadItem `json:"unread,omitempty"`
	Error     string       `json:"error,omitempty"`
}
