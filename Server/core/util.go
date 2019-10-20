package core

func RemoveMatchRequest(s []MatchRequest, i int) []MatchRequest {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveGame(s []*Game, i int) []*Game {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
