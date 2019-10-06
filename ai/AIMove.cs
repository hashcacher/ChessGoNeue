using ChessGo;
using System;

namespace CardGameServer
{
    struct AIMove : IComparable<AIMove>
    {
        public Point[] move;
        public short score;
        public AIMove(Point[] move, short score)
        {
            this.move = move;
            this.score = score;
        }
        
        public int CompareTo(AIMove move2){
            if (score < move2.score)
                return -1;
            if (score == move2.score)
                return 0;
            if (score > move2.score)
                return 1;
            return 0;
        }
    }
}
