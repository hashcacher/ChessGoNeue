namespace ChessGo
{
    public class GamePublic {
        public int id;
        public string opponent;
        public long timeLeft;
        public bool myTurn;
    }

    public class MyGamesResponse {
        GamePublic[] games;
    }
}
