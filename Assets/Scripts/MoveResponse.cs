namespace ChessGo
{
    public class MoveResponse {
        public string move; // e.g. a2,a4
        public long whiteLeft;
        public long blackLeft;
        public string blackTurnStarted;
        public string whiteTurnStarted;
        public bool gameOver;
        public bool whiteWins;
        public string winReason; // Checkmate / stalemate / king surround
        public int eloAdj;
    }
}
