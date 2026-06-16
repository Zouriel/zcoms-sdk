# zcoms-sdk

Shared pure-Go toolkit for [zcoms](https://github.com/Zouriel/zcoms) components
(bridge / triage / errands). No cgo/TDLib — components are plain IPC clients of
the core daemon.

- `ipc` — the daemon IPC protocol + client (send/sendfile/read/unread/mark_read/resolve).
- `agent` — config readers (settings/agents/allowlist/locations), agent-backend
  selection, session listing, and the claude/codex runner.
- `whatsapp` — client for the Baileys sidecar socket.
