using System;
namespace ChessGo
{
    [Serializable]
    public class GamePublic {
        public int id;
        public string opponent;
        public long timeLeft;
        public bool myTurn;
    }
}
