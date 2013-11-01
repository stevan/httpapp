package sessions

import (
	"testing"
)

func TestSession(t *testing.T) {
	s := NewSession("0xd3cafbad")

	if s.Id != "0xd3cafbad" {
		t.Errorf("Got the wrong id - got: %v expected: %v", s.Id, "0xd3cafbad")
	}

	if v, ok := s.Data["foo"]; ok {
		t.Errorf("Expected no value, got: %v", v)
	}

	s.Data["foo"] = 100

	if v, ok := s.Data["foo"]; ok {
		if v.(int) != 100 {
			t.Errorf("Got the wrong value - got: %v expected: %v", v, 100)
		}
	} else {
		t.Error("Expected to have a value")
	}

	delete(s.Data, "foo")

	if v, ok := s.Data["foo"]; ok {
		t.Errorf("Expected no value, got: %v", v)
	}
}

func TestGenerateId(t *testing.T) {
	s := NewCookieState("test-cookie")

	id := s.Generate()

	if !s.ValidateSessionId(id) {
		t.Errorf("id was Invalid %s", id)
	}

	id2 := s.Generate()

	if id == id2 {
		t.Errorf("Got duplicate session ids - got: %v expected: %v", id, id2)
	}
}
