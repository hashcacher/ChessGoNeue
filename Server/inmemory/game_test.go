package inmemory_test

import (
	"reflect"
	"testing"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/inmemory"
)

// func TestStoreOK(t *testing.T) {

// 	type output struct {
// 		id  int
// 		err error
// 	}

// 	tests := []struct {
// 		gamesMapInjection   map[int]core.Game
// 		in                  core.Game
// 		expectOut           output
// 		expectGamesMapAfter map[int]core.Game
// 	}{
// 		{
// 			gamesMapInjection: make(map[int]core.Game),
// 			in:                core.Game{BlackUser: 1, WhiteUser: 2},
// 			expectOut:         output{id: 1, err: nil},
// 			expectGamesMapAfter: map[int]core.Game{
// 				1: core.Game{ID: 1, BlackUser: 1, WhiteUser: 2},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		games := inmemory.NewGames(tt.gamesMapInjection)
// 		gotID, gotErr := games.Store(tt.in)
// 		// Check error
// 		if gotErr != nil {
// 			if tt.expectOut.err != nil {
// 				if gotErr.Error() != tt.expectOut.err.Error() {
// 					t.Fatalf("got error: %v, expected: %v", gotErr, tt.expectOut.err)
// 				}
// 			} else {
// 				t.Fatalf("unexpected error: %v", gotErr)
// 			}
// 		}
// 		// Check ID
// 		if gotID != tt.expectOut.id {
// 			t.Fatalf("got id: %v, expected: %v", gotID, tt.expectOut.id)
// 		}
// 		// Check state after
// 		if !reflect.DeepEqual(tt.gamesMapInjection, tt.expectGamesMapAfter) {
// 			t.Fatalf("got: %v, expected: %v", tt.gamesMapInjection, tt.expectGamesMapAfter)
// 		}

// 	}
// }

// Successful create user call
func TestStoreOK(t *testing.T) {
	gamesMap := make(map[int]core.Game)
	games := inmemory.NewGames(gamesMap)

	expectGamesMap := map[int]core.Game{
		1: core.Game{ID: 1, BlackUser: 1, WhiteUser: 2},
	}
	expectID := 1

	gotID, gotErr := games.Store(core.Game{BlackUser: 1, WhiteUser: 2})

	if gotID != expectID {
		t.Fatalf("expected output id of %v, got %v", expectID, gotID)
	}

	if gotErr != nil {
		t.Fatalf("unexpected error: %v", gotErr)
	}
	if !reflect.DeepEqual(gamesMap, expectGamesMap) {
		t.Fatalf("got: %v, expected: %v", gamesMap, expectGamesMap)
	}
}

func TestStoreMultipleOK(t *testing.T) {
	gamesMap := make(map[int]core.Game)
	games := inmemory.NewGames(gamesMap)

	expectGamesMap := map[int]core.Game{
		1: core.Game{ID: 1, BlackUser: 1, WhiteUser: 2},
		2: core.Game{ID: 2, BlackUser: 2, WhiteUser: 3},
	}

	// First store
	expectID := 1
	gotID, gotErr := games.Store(core.Game{BlackUser: 1, WhiteUser: 2})
	if gotID != expectID {
		t.Fatalf("expected output id of %v, got %v", expectID, gotID)
	}
	if gotErr != nil {
		t.Fatalf("unexpected error: %v", gotErr)
	}
	// Second store
	expectID = 2
	gotID, gotErr = games.Store(core.Game{BlackUser: 2, WhiteUser: 3})
	if gotID != expectID {
		t.Fatalf("expected output id of %v, got %v", expectID, gotID)
	}
	if gotErr != nil {
		t.Fatalf("unexpected error: %v", gotErr)
	}
	// Check state
	if !reflect.DeepEqual(gamesMap, expectGamesMap) {
		t.Fatalf("got: %v, expected: %v", gamesMap, expectGamesMap)
	}
}
