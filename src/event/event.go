package event

import "gogoing"

type Event struct {
	Type uint8
	//Data []byte
	Sess gogoing.Session
}
