package gogoing

type Event struct {
	Type uint8
	//Data []byte
	Sess Session
}

func (e Event)String() string  {
	return string(e.Type)
}
type DataEvent struct {
	Event
	Data []byte
}

func (e DataEvent)String() string  {
	return string(e.Type)+string(e.Data)
}