package gogoing

type Event interface {
	GetType() uint8
	GetEventID() int
	GetSess() Session
	SetSess(sess Session)
}

type DefaultEvent struct {
	Type uint8
	ID int
	//Data []byte
	Sess Session
}

func (e DefaultEvent)String() string  {
	return string(e.Type)
}

func (e DefaultEvent) GetType() uint8 {
	return e.Type
}

func (e DefaultEvent) GetEventID() int {
	return e.ID
}

func (e DefaultEvent) GetSess() Session{
	return e.Sess
}

func (e DefaultEvent) SetSess(sess Session) {
	e.Sess = sess
}

type CloseEvent struct {
	DefaultEvent
}

type ExceptEvent struct {
	DefaultEvent
}

type InternetEvent struct {
	Type uint8
	ID int
	//Data []byte
	Sess Session
	Data []byte
}

func (e *InternetEvent)String() string  {
	return string(e.Type)
}

func (e *InternetEvent) GetType() uint8 {
	return e.Type
}

func (e *InternetEvent) GetEventID() int {
	return e.ID
}

func (e *InternetEvent) GetSess() Session{
	return e.Sess
}

func (e *InternetEvent) SetSess(sess Session) {
	e.Sess = sess
}

func NewEvent(eventType uint8, id int, sess Session) Event {
	return &DefaultEvent{Type:eventType, ID:id, Sess:sess}
}