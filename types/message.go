package types

import qk "github.com/quickfixgo/quickfix"

type Message struct {
	SessionID qk.SessionID
	Message   *qk.Message
}
