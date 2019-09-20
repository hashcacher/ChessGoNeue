package matchRequest

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"service"
)

func TestMatchMeOK(t *testing.T) {
	s := service.NewWebservice()

	req, err := http.NewRequest("POST", "/v1/matchMe",
		bytes.NewReader([]byte(`{ "clientID": "1234" }`)))

	if err != nil {
		t.Fatal(err)
	}

	rr := NewTestResponseRecorder()
	handler := http.HandlerFunc(s.matchMeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("should be waiting for a match but got response %v",
			rr.Body.String())
	}

	time.Sleep(1e6)
	if rr.Body.String() != expected {
		t.Errorf("should still be waiting for a match but got response %v",
			rr.Body.String())
	}

	// Should still be blank as no match is found
}

func TestMatchMeFindMatch(t *testing.T) {
	s := service.NewWebservice()

	// Client 1 sends matchme
	req, err := http.NewRequest("POST", "/v1/matchMe",
		bytes.NewReader([]byte(`{ "clientID": "1234" }`)))

	if err != nil {
		t.Fatal(err)
	}

	rr := NewTestResponseRecorder()
	s.matchMeHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Client 2 sends matchme
	req, err = http.NewRequest("POST", "/v1/matchMe",
		bytes.NewReader([]byte(`{ "clientID": "5678" }`)))
	if err != nil {
		t.Fatal(err)
	}

	rr2 := NewTestResponseRecorder()
	s.matchMeHandler(rr2, req)

	// Wait a millisec
	time.Sleep(1e6)

	// Client2 should have a match
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected = `{"err":null,"haveMatch":true,"areWhite":true}`
	if rr2.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Client 1 should have a match too now.
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
