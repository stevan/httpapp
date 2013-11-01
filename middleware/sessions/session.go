package sessions

// --------------------------------------
// Session Object
// --------------------------------------

func NewSession(id string) *Session {
	opts := make(map[string]interface{})
	data := make(map[string]interface{})
	return &Session{id, data, opts}
}

type Session struct {
	Id      string
	Data    map[string]interface{}
	Options map[string]interface{}
}

func (s *Session) IsExpired() bool {
	if expire, ok := s.Options["expire"]; ok {
		return expire.(bool)
	}
	return false
}

func (s *Session) Expire() {
	for k, _ := range s.Data {
		delete(s.Data, k)
	}
	s.Options["expire"] = true
}
