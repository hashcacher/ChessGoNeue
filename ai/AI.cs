using ChessGo;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Diagnostics;

namespace CardGameServer
{
    public class AI 
    {
        Game game;
        Player player;
        char[,] board;
        Random rand; 

        public AI(Game game, Player player) { 
            this.game = game;
            this.board = game.board;
            this.player = player;
            rand = new Random();
        }

        public void MakeMove()
        {
            Utilities.PrintBoard(game.board);
            var watch = new Stopwatch(); 
            watch.Start();
            var result = Minimax(watch, 4, short.MinValue, short.MaxValue, game.board, player.amBlack);

            var move = result.move;

            //Console.WriteLine("Best score : " + result.score + " Took " + watch.ElapsedMilliseconds + " move:"  + String.Join(",", result.move));

            if(move.Length == 2)
                game.MakeChessMove(player, move[0].row, move[0].col, move[1].row, move[1].col, move.Select((m) => m.ToString()).ToArray());
            else
                game.MakeGoMove(player, move[0].row, move[0].col, move.Select((m) => m.ToString()).ToArray());
            //Utilities.PrintBoard(game.board);

        }


        private Point[] RandomMove()
        {
            var pieces = Utilities.GetMyChessPieces(player.amBlack, board);
            for(int i = 0; i < pieces.Count; i++)
            {
                var piece = pieces[rand.Next(pieces.Count)];
                var dests = Utilities.GetValidDestinations(piece, board);
                if(dests.Count > 0)
                {
                    var dest = dests[rand.Next(dests.Count)];
                    return new Point[] {piece, dest};
                }
            }
            return null;
        }

        private char[,] SimulateMove(Point[] p, char[,] board, bool black)
        {
            HashSet<Point> deadPoints = null;
            if(p.Length == 2)
            {
                board[p[1].row, p[1].col] = board[p[0].row, p[0].col];  
                board[p[0].row, p[0].col] = '\0';  
                deadPoints = Utilities.CheckSurrounded(p[1], board);
            }
            else
            {
                board[p[0].row, p[0].col] = black ? 'S' : 's';  
                deadPoints = Utilities.CheckSurrounded(p[0], board);
            }

            foreach (Point dp in deadPoints)
                board [dp.row, dp.col] = '\0';

            return board;
        }

        private bool IsGameOver(char[,] board)
        {
            bool white = false, black = false;
            for(int r=0;r < board.GetLength(0);r++)
            {
                for(int c=0;c < board.GetLength(1);c++)
                {
                    if(board[r,c] == 'k')
                        white = true;
                    if(board[r,c] == 'K')
                        black = true;
                }
            }
            return !white || !black;
        }

        private AIMove Minimax(Stopwatch watch, int depth, short alpha, short beta, char[,] board, bool maximizing)
        {
            if(depth == 0 || IsGameOver(board))
                return new AIMove(new Point[2], ScoreBoard(board));

            // Try all Chess moves
            var pieces = Utilities.GetMyChessPieces(maximizing, board);
            var moves = new List<AIMove>();
            foreach(var piece in pieces)
            {
                var dests = Utilities.GetValidDestinations(piece, board);
                foreach(var dest in dests)
                   moves.Add(new AIMove(new Point[] {piece, dest}, maximizing ? short.MinValue : short.MaxValue));
            }

            // Try go moves
            for(int r=0;r < board.GetLength(0);r++)
            {
                for(int c=0;c < board.GetLength(1);c++)
                {
                    if (board [r, c] == '\0') {
                        moves.Add(new AIMove(new Point[] { new Point (r, c) }, maximizing ? short.MinValue : short.MaxValue));
                    }
                }
            }

            
            int n = moves.Count;
            var movesCopy = new List<AIMove>();
            for(int i = 0; i < n; i++)
            {
                var idx = rand.Next(moves.Count-1);
                Console.WriteLine(idx);
                var move = moves[idx];
                var newBoard = SimulateMove (move.move, board.Clone () as char[,], maximizing);
                move.score = Minimax(watch, depth-1, alpha, beta, newBoard, !maximizing).score;
                moves[i] = move;

                if(maximizing)
                    alpha = Math.Max(alpha, move.score);
                else
                    beta = Math.Min(beta, move.score);

                if(beta <= alpha) 
                    break;

                moves.RemoveAt(idx);
                movesCopy.Add(move);
            }

            if (movesCopy.Count == 0) return new AIMove(null, maximizing ? short.MinValue : short.MaxValue);
            return maximizing ? movesCopy.Max() : movesCopy.Min();
        }

        private short ScoreBoard(char[,] board)
        {
            short total = 0;
            for(int r=0;r < board.GetLength(0);r++)
            {
                for(int c=0;c < board.GetLength(1);c++)
                {
                    char piece = board[r,c];
                    short score = ScorePiece(piece);
                    total -= score;
                }
            }
            return total;
        }

        private short ScorePiece(char piece)
        {
            switch(piece)
            {
                case '\0':
                    return 0;
                case 'p':
                    return 10;
                case 'P':
                    return -10;
                case 'r':
                    return 50;
                case 'R':
                    return -50;
                case 'n':
                    return 30;
                case 'N':
                    return -30;
                case 'b':
                    return 30;
                case 'B':
                    return -30;
                case 'q':
                    return 90;
                case 'Q':
                    return -90;
                case 'k':
                    return 900;
                case 'K':
                    return -900;
                case 's':
                    return 0;
                case 'S':
                    return 0;
                default:
                    return 0;
            }

        }

        public void Destroy()
        {
            this.game = null;
            this.player = null;
            this.board = null;
            this.rand = null;
        }
    }
}
