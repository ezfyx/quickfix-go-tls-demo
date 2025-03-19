package types

import qk "github.com/quickfixgo/quickfix"

type Session struct {
	SessionId qk.SessionID
	Status    int
}
