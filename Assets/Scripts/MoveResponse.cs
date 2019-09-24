namespace ChessGo
{
    public class MoveResponse {
        public string move; // e.g. a2,a4
        public bool gameOver;
        public bool whiteWins;
        public string winReason; // Checkmate / stalemate / king surround
        public int eloAdj;
    }
}
