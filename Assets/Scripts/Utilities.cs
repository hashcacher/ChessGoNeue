using System;
using System.Collections.Generic;
using UnityEngine;



namespace ChessGo
{
    //Static methods for computing Go surrounds, valid moves, etc.
    public class Utilities
    {
        private const int MAXROW = 11;
        private const int MAXCOL = 11;

        public Utilities()
        {

        }

        public string GetServerHost() {
            if (Application.isEditor) {
                return "localhost:8080";
            } else {
                return "https://chessgo.xyz";
            }
        }

        // Returns in a list all orthogonally adjacent but not out of bound points
        public static HashSet<Point> GetAdjacentPoints(Point p)
        {
            HashSet<Point> neighbors = new HashSet<Point>();

            //Up
            if (p.row < MAXROW) { neighbors.Add(new Point(p.row + 1, p.col)); }
            //Down
            if (p.row > 0) { neighbors.Add(new Point(p.row - 1, p.col)); }
            //Left
            if (p.col < MAXCOL) { neighbors.Add(new Point(p.row, p.col + 1)); }
            //Right
            if (p.col > 0) { neighbors.Add(new Point(p.row, p.col - 1)); }

            return neighbors;
        }

        public static HashSet<Point> GetNeighbors(HashSet<Point> group, char[,] board)
        {
            HashSet<Point> neighbors = new HashSet<Point>();

            foreach (Point p in group)
            {
                HashSet<Point> localNeighbors = GetAdjacentPoints(p);
                foreach (Point n in localNeighbors)
                {
                    if (!group.Contains(n))
                    {
                        neighbors.Add(n);
                    }
                }
            }
            return neighbors;
        }

        //Returns a list of all pieces that are directly connected 
        public static HashSet<Point> GetGroup(Point p, char[,] board)
        {
            HashSet<Point> group = new HashSet<Point>();
            group.Add(p);
            int i = 0;
            int prevSize = 1;
            do
            {
                prevSize = group.Count;
                HashSet<Point> neighbors = GetNeighbors(group, board);
                i++;

                foreach (Point n in neighbors)
                {
                    if ((IsWhiteAt(p, board) && IsWhiteAt(n, board)) || (IsBlackAt(p, board) && IsBlackAt(n, board)))
                    {
                        group.Add(n);
                    }
                }
                //Just in case of infinite loop!
                if (i == MAXCOL * MAXROW - 1)
                    Debug.LogError("Some shit went wrong in GetGroup");
            } while (prevSize != group.Count && i < MAXCOL * MAXROW);

            return group;
        }

        public static HashSet<Point> GetLiberties(Point p, char[,] board)
        {
            HashSet<Point> group = GetGroup(p, board);
            HashSet<Point> liberties = GetNeighbors(group, board);

            return liberties;
        }

        //Checks to see if group that piece at point p belongs to is totally surrounded - none of its liberties are occupied by an empty space
        public static bool IsGroupDead(Point p, char[,] board)
        {
            HashSet<Point> liberties = GetLiberties(p, board);
            foreach (Point l in liberties)
            {
                if (IsEmptyAt(l, board)) { return false; }
            }
            return true;
        }

	public static int CountLiberties(Point p, char[,] board)
	{
	    HashSet<Point> liberties = GetLiberties(p, board);
	    int count = 0;
	    foreach (Point l in liberties)
	    {
		if (IsEmptyAt(l, board)) { count++; }
	    }
	    return count;
	}

        // Checks to see if white king is under check
        public static bool IsWhiteKingChecked(char[,] board)
        {
            //Find the black king
            Point kingPos = new Point();
	    bool kingFound = false;
            for (int i = 0; i <= MAXROW && !kingFound; i++)
            {
                for (int j = 0; j <= MAXCOL && !kingFound; j++)
                {
                    if (board[i, j] == 'k')
                    {
                        //White King Found
                        kingPos = new Point(i, j);
                        kingFound = true;
                    }
                }
            }

	    if (!kingFound) {
		//The king was not found, error
		Debug.LogError("WHITE KING IS DEAD.");
	    }
            //Check to see if there is a legal move for white onto the white king.

            for (int i = 0; i <= MAXROW; i++)
            {
                for (int j = 0; j <= MAXCOL; j++)
                {
                    Point p = new Point(i, j);
                    if (IsValidMove(p, kingPos, board))
                    {
			//White king is under attack!
			Debug.LogError("WHITE KING IS CHECKED.");
                        return true;
                    }
                }
            }

            return (CountLiberties(kingPos, board) == 1);
        }

        public static bool IsBlackKingChecked(char[,] board)
        {
            //Find the black king
	    Point kingPos =  new Point();
	    bool kingFound = false;
	    for (int i = 0; i <= MAXROW && !kingFound; i++)
            {
		for (int j = 0; j <= MAXCOL && !kingFound; j++)
                {
                    if (board[i, j] == 'K')
                    {
                        //Black King Found
			kingPos = new Point(i, j);
			kingFound = true;
                    }
                }
            }

	    if (!kingFound) {
		//The king was not found, error
		Debug.LogError("WHITE KING IS DEAD.");
	    }
            //Check to see if there is a legal move for white onto the black king.

            for (int i = 0; i <= MAXROW; i++)
            {
                for (int j = 0; j <= MAXCOL; j++)
                {
                    Point p = new Point(i, j);
                    if (IsValidMove(p, kingPos, board))
                    {
			//Black king is under attack!
			Debug.LogError("BLACK KING IS CHECKED.");
                        return true;
                    }
                }
            }
	    //TODO: Toggle condition flags
	    return (CountLiberties(kingPos, board) == 1);
	}

	public static bool IsBlackCheckmated(char[,] board)
	{

	    if (IsBlackKingChecked (board)) {

		Point kingPos = new Point();
		for (int i = 0; i <= MAXROW; i++)
		{
		    for (int j = 0; j <= MAXCOL; j++)
		    {
			if (board[i, j] == 'K')
			{
			    //Black King Found
			    kingPos = new Point(i, j);
			    break;
			}
		    }
		}

		//First: Attempt to move the king out of harm's way
		for(int i = -1; i <= 1; i++)
		{
		    for(int j = -1; j<=1; j++)
		    {
			//Try to move the king somewhere
			char[,] attemptBoard = (char[,])board.Clone();
			int newRow = kingPos.row+i;
			int newCol = kingPos.col+j;
			Point newPos = new Point(newRow, newCol);
			if( newRow >= 0 && newRow <= MAXROW && newCol >= 0 && newCol <= MAXCOL && !IsBlackAt(newPos, board) )
			{
			    //Check king escape paths
			    attemptBoard[kingPos.row, kingPos.col] = '\0';
			    attemptBoard[newRow, newCol] = 'K';
			    //Try escape path, if it get king out of check, then black is not checkmated
			    if( !IsBlackKingChecked(attemptBoard) ) { return false; }
			}
		    }
		}
		//King cannot escape by moving, find attackers
		//Locating white attackers

		List<Point> attackers = new List<Point>();
		for (int i = 0; i <= MAXROW; i++)
		{
		    for (int j = 0; j <= MAXCOL; j++)
		    {
			Point curPos = new Point(i, j);
			if(IsWhiteAt(curPos, board))
			{
			    if(IsValidMove(curPos, kingPos, board))
			    {
				//curPos contains a piece that can attack the king
				attackers.Add(curPos);
			    }
			}
		    }
		}

		//White attacker(s) located, try to stop them from attacking by killing them
		for (int i = 0; i <= MAXROW; i++)
		{
		    for (int j = 0; j <= MAXCOL; j++)
		    {
			Point curPos = new Point(i, j);
			if(IsBlackAt(curPos, board))
			{
			    foreach(Point a in attackers)
			    {
				if(IsValidMove(curPos, a, board))
				{
				    //Check to see if defense will work
				    char[,] attemptBoard = (char[,])board.Clone();
				    attemptBoard[a.row, a.col] = board[curPos.row, curPos.col];
				    attemptBoard[curPos.row, curPos.col] = '\0';
				    //Try defense, if it get king out of check, then black is not checkmated
				    if( !IsBlackKingChecked(attemptBoard) ) { return false; }
				}
			    }
			}
		    }
		}

		// Attacker(s) not killable; try blocking with go stone
		// Determine possible defenses through brute force; may be inefficient
		// Try only as a last resort.

		foreach(Point a in attackers)
		{
		    //Try dropping a go stone in the path of the attacker, then see if it resolves the check
		    //Intersection defense only works against rooks, bishops, and queens.
		    if( board[a.row, a.col] == 'b' || board[a.row, a.col] == 'r' || board[a.row, a.col] == 'q' )
		    {
			//Try to defend all 
			List<Point> intercept = GetValidDestinations(a, board);
			foreach(Point i in intercept)
			{
			    char[,] attemptBoard = (char[,])board.Clone();
			    if( attemptBoard[i.row, i.col] == '\0' )
			    {
				//Try a go piece defense
				attemptBoard[i.row, i.col] = 'S';
				if( !IsBlackKingChecked(attemptBoard) ) { return false; }
			    }
			}
		    }
		}


		//All escape attempts failed, black king is checkmated
		Debug.LogError("BLACK KING IS CHECKMATED. GAME OVER.");
		return true;
	    } else {
		//Black king was not checked
		return false;
	    }
	}


	public static bool IsWhiteCheckmated(char[,] board)
	{
	    if (IsWhiteKingChecked (board)) {

		Point kingPos = new Point();
		bool kingFound = false;
		for (int i = 0; i <= MAXROW && !kingFound; i++)
		{
		    for (int j = 0; j <= MAXCOL && !kingFound; j++)
		    {
			if (board[i, j] == 'k')
			{
			    //White King Found
			    kingPos = new Point(i, j);
			    kingFound = true;
			}
		    }
		}

		//First: Attempt to move the king out of harm's way
		for(int i = -1; i <= 1; i++)
		{
		    for(int j = -1; j<=1; j++)
		    {
			//Try to move the king somewhere
			char[,] attemptBoard = (char[,])board.Clone();
			int newRow = kingPos.row+i;
			int newCol = kingPos.col+j;
			Point newPos = new Point(newRow, newCol);
			if( newRow >= 0 && newRow <= MAXROW && newCol >= 0 && newCol <= MAXCOL && !IsWhiteAt(newPos, board) )
			{
			    //Check king escape paths
			    attemptBoard[kingPos.row, kingPos.col] = '\0';
			    attemptBoard[newRow, newCol] = 'k';
			    //Try escape path, if it get king out of check, then black is not checkmated
			    if( !IsWhiteKingChecked(attemptBoard) ) { return false; }
			}
		    }
		}
		//King cannot escape by moving, find attackers
		//Locating white attackers
		List<Point> attackers = new List<Point>();
		for (int i = 0; i <= MAXROW; i++)
		{
		    for (int j = 0; j <= MAXCOL; j++)
		    {
			Point curPos = new Point(i, j);
			if(IsBlackAt(curPos, board))
			{
			    if(IsValidMove(curPos, kingPos, board))
			    {
				//curPos contains a piece that can attack the king
				attackers.Add(curPos);
			    }
			}
		    }
		}

		//White attacker(s) located, try to stop them from attacking by killing them
		for (int i = 0; i <= MAXROW; i++)
		{
		    for (int j = 0; j <= MAXCOL; j++)
		    {
			Point curPos = new Point(i, j);
			if(IsWhiteAt(curPos, board))
			{
			    foreach(Point a in attackers)
			    {
				if(IsValidMove(curPos, a, board))
				{
				    //Check to see if defense will work
				    char[,] attemptBoard = (char[,])board.Clone();
				    attemptBoard[a.row, a.col] = board[curPos.row, curPos.col];
				    attemptBoard[curPos.row, curPos.col] = '\0';
				    //Try defense, if it get king out of check, then black is not checkmated
				    if( !IsWhiteKingChecked(attemptBoard) ) { return false; }
				}
			    }
			}
		    }
		}

		// Attacker(s) not killable; try blocking with go stone
		// Determine possible defenses through brute force; may be inefficient
		// Try only as a last resort.
		foreach(Point a in attackers)
		{
		    //Try dropping a go stone in the path of the attacker, then see if it resolves the check
		    //Intersection defense only works against rooks, bishops, and queens.
		    if( board[a.row, a.col] == 'B' || board[a.row, a.col] == 'R' || board[a.row, a.col] == 'Q' )
		    {
			//Try to defend all 
			List<Point> intercept = GetValidDestinations(a, board);
			foreach(Point i in intercept)
			{
			    char[,] attemptBoard = (char[,])board.Clone();
			    if( attemptBoard[i.row, i.col] == '\0' )
			    {
				//Try a go piece defense
				attemptBoard[i.row, i.col] = 's';
				if( !IsWhiteKingChecked(attemptBoard) ) { return false; }
			    }
			}
		    }
		}
		//All escape attempts failed, white king is checkmated
		//TODO: Toggle condition flags to end game?
		Debug.LogError("WHITE KING IS CHECKMATED. GAME OVER.");
		return true;
	    } else {
		//white king was not checked, no checkmate
		return false;
	    }
	}

        public static char GetCharForPiece(GameObject o)
        {
            switch (o.name.Trim())
            {
                case "WhitePawn":
                    return 'p';
                case "BlackPawn":
                    return 'P';
                case "WhiteRook":
                    return 'r';
                case "BlackRook":
                    return 'R';
                case "WhiteKnight":
                    return 'n';
                case "BlackKnight 1":
                    return 'N';
                case "WhiteBishop":
                    return 'b';
                case "BlackBishop":
                    return 'B';
                case "WhiteQueen":
                    return 'q';
                case "BlackQueen":
                    return 'Q';
                case "WhiteKing":
                    return 'k';
                case "BlackKing":
                    return 'K';
                case "WhiteStone":
                    return 's';
                case "BlackStone":
                    return 'S';
            }
            Debug.LogError("GetCharForPiece didn't find the char for piece " + o.name);
            return '\0';
        }

        public static bool IsWhiteAt(Point p, char[,] board)
        {
            char piece = board[p.row, p.col];
            return (piece == 'p' || piece == 'r' || piece == 'n' || piece == 'b' || piece == 'q' || piece == 'k' || piece == 's');
        }

        public static bool IsBlackAt(Point p, char[,] board)
        {
            char piece = board[p.row, p.col];
            return (piece == 'P' || piece == 'R' || piece == 'N' || piece == 'B' || piece == 'Q' || piece == 'K' || piece == 'S');
        }

        public static bool IsEmptyAt(Point p, char[,] board)
        {
            return board[p.row, p.col] == '\0';
        }

        //True if p1 can move to p2
        public static bool IsValidMove(Point p1, Point p2, char[,] board)
        {
            char piece = board[p1.row, p1.col];

            if (piece == '\0' || piece == 's' || piece == 'S')
            {
                //these pieces cannot move
                return false;
            }
            else if (piece == 'p')
            {
                //White Pawn
                //white pawns can move DOWN (1 y lower) or attack diagonally

                if (p1.col == MAXCOL-1 && p2.col == MAXCOL-3 && p1.row == p2.row)
                {
                    //Double move possible? Check for space
                    return (board[p1.row, p1.col - 1] == '\0' && board[p1.row, p1.col - 2] == '\0');
                }
                else if ((p1.col - 1 == p2.col) && (p1.row == p2.row))
                {
                    //Moving ahead is legal if space is empty
                    return IsEmptyAt(p2, board);
                }
                else if ((p1.col - 1 == p2.col) && (p1.row == p2.row + 1 || p1.row == p2.row - 1))
                {
                    //Check diagonal attack
                    return IsBlackAt(p2, board);
                }
            }
            else if (piece == 'P')
            {
                //Black Pawn
                //white pawns can move up (1 y higher) or attack diagonally

                if (p1.col == 1 && p2.col == 3 && p1.row == p2.row)
                {
                    //Double move possible? Check for space
                    return (board[p1.row, p1.col + 1] == '\0' && board[p1.row, p1.col + 2] == '\0');
                }
                else if ((p1.col + 1 == p2.col) && (p1.row == p2.row))
                {
                    //Moving ahead is legal if space is empty
                    return IsEmptyAt(p2, board);
                }
                else if ((p1.col + 1 == p2.col) && (p1.row == p2.row + 1 || p1.row == p2.row - 1))
                {
                    //Check diagonal attack
                    return IsWhiteAt(p2, board);
                }
            }
            else if (piece == 'r')
            {
                //White Rook
                //Check orthogonality
                if (p1.row == p2.row)
                {
                    //moving horizontally
                    //check for pieces in the way
                    if (p1.col < p2.col)
                    {
                        //moving to the right
                        for (int i = p1.col + 1; i < p2.col; i++)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving to the left
                        for (int i = p1.col - 1; i > p2.col; i--)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                else if (p1.col == p2.col)
                {
                    //moving vertically
                    //check for pieces in the way
                    if (p1.row < p2.row)
                    {
                        //moving up
                        for (int i = p1.row + 1; i < p2.row; i++)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving down
                        for (int i = p1.row - 1; i > p2.row; i--)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'R')
            {
                //Black Rook
                //Check orthogonality
                if (p1.row == p2.row)
                {
                    //moving horizontally
                    //check for pieces in the way
                    if (p1.col < p2.col)
                    {
                        //moving to the right
                        for (int i = p1.col + 1; i < p2.col; i++)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving to the left
                        for (int i = p1.col - 1; i > p2.col; i--)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                else if (p1.col == p2.col)
                {
                    //moving vertically
                    //check for pieces in the way
                    if (p1.row < p2.row)
                    {
                        //moving up
                        for (int i = p1.row + 1; i < p2.row; i++)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving down
                        for (int i = p1.row - 1; i > p2.row; i--)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'n')
            {
                //White Knight
                //Knights can move in L pattern - two in one direction, one in the other.
                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                if ((xdist == 2 && ydist == 1) || (xdist == 1 && ydist == 2))
                {
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'N')
            {
                //Black Knight
                //Knights can move in L pattern - two in one direction, one in the other.
                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                if ((xdist == 2 && ydist == 1) || (xdist == 1 && ydist == 2))
                {
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'b')
            {
                //White Bishop
                //Bishops move diagonally
                //First, check for placement on diagonals

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                int xmult, ymult;

                //Setting search direction

                if (p2.col > p1.col)
                {
                    ymult = 1;
                }
                else
                {
                    ymult = -1;
                }

                if (p2.row > p1.row)
                {
                    xmult = 1;
                }
                else
                {
                    xmult = -1;
                }

                if (xdist == ydist)
                {
                    //Destination is on diagonal - check for pieces in between
                    for (int i = 1; i < xdist; i++)
                    {
                        if (board[p1.row + i * xmult, p1.col + i * ymult] != '\0') return false;
                    }
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'B')
            {
                //Black Bishop
                //Bishops move diagonally
                //First, check for placement on diagonals

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                int xmult, ymult;

                //Setting search direction

                if (p2.col > p1.col)
                {
                    ymult = 1;
                }
                else
                {
                    ymult = -1;
                }

                if (p2.row > p1.row)
                {
                    xmult = 1;
                }
                else
                {
                    xmult = -1;
                }

                if (xdist == ydist)
                {
                    //Destination is on diagonal - check for pieces in between
                    for (int i = 1; i < xdist; i++)
                    {
                        if (board[p1.row + i * xmult, p1.col + i * ymult] != '\0') return false;
                    }
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                else { return false; }
            }
            else if (piece == 'q')
            {
                //White Queen
                //Queens get rook and bishop moves
                if (p1.row == p2.row)
                {
                    //moving horizontally
                    //check for pieces in the way
                    if (p1.col < p2.col)
                    {
                        //moving to the right
                        for (int i = p1.col + 1; i < p2.col; i++)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving to the left
                        for (int i = p1.col - 1; i > p2.col; i--)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                else if (p1.col == p2.col)
                {
                    //moving vertically
                    //check for pieces in the way
                    if (p1.row < p2.row)
                    {
                        //moving up
                        for (int i = p1.row + 1; i < p2.row; i++)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving down
                        for (int i = p1.row - 1; i > p2.row; i--)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
                //Next, check for placement on diagonals

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                int xmult, ymult;

                //Setting search direction

                if (p2.col > p1.col)
                {
                    ymult = 1;
                }
                else
                {
                    ymult = -1;
                }

                if (p2.row > p1.row)
                {
                    xmult = 1;
                }
                else
                {
                    xmult = -1;
                }

                if (xdist == ydist)
                {
                    //Destination is on diagonal - check for pieces in between
                    for (int i = 1; i < xdist; i++)
                    {
                        if (board[p1.row + i * xmult, p1.col + i * ymult] != '\0') return false;
                    }
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }

                //Options exhausted
                return false;
            }
            else if (piece == 'Q')
            {
                //Black Queen
                //Queens get rook and bishop moves
                if (p1.row == p2.row)
                {
                    //moving horizontally
                    //check for pieces in the way
                    if (p1.col < p2.col)
                    {
                        //moving to the right
                        for (int i = p1.col + 1; i < p2.col; i++)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving to the left
                        for (int i = p1.col - 1; i > p2.col; i--)
                        {
                            if (board[p1.row, i] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                else if (p1.col == p2.col)
                {
                    //moving vertically
                    //check for pieces in the way
                    if (p1.row < p2.row)
                    {
                        //moving up
                        for (int i = p1.row + 1; i < p2.row; i++)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    else
                    {
                        //moving down
                        for (int i = p1.row - 1; i > p2.row; i--)
                        {
                            if (board[i, p1.col] != '\0') { return false; }
                        }
                    }
                    //Nothing in the way, OK if destination is empty or enemy
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
                //Next, check for placement on diagonals

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                int xmult, ymult;

                //Setting search direction

                if (p2.col > p1.col)
                {
                    ymult = 1;
                }
                else
                {
                    ymult = -1;
                }

                if (p2.row > p1.row)
                {
                    xmult = 1;
                }
                else
                {
                    xmult = -1;
                }

                if (xdist == ydist)
                {
                    //Destination is on diagonal - check for pieces in between
                    for (int i = 1; i < xdist; i++)
                    {
                        if (board[p1.row + i * xmult, p1.col + i * ymult] != '\0') return false;
                    }
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }

                //Options exhausted
                return false;
            }
            else if (piece == 'k')
            {
                //White King

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                if (xdist > 1 || ydist > 1)
                {
                    return false;
                }
                else
                {
                    return (IsEmptyAt(p2, board) || IsBlackAt(p2, board));
                }
            }
            else if (piece == 'K')
            {
                //Black King

                int xdist = Mathf.Abs(p1.col - p2.col);
                int ydist = Mathf.Abs(p1.row - p2.row);
                if (xdist > 1 || ydist > 1)
                {
                    return false;
                }
                else
                {
                    return (IsEmptyAt(p2, board) || IsWhiteAt(p2, board));
                }
            }
            return false;
        }

        public static List<Point> GetValidDestinations(Point p1, char[,] board)
        {
            List<Point> moves = new List<Point>();
            char piece = board[p1.row, p1.col];
            if (piece == '\0' || piece == 's' || piece == 'S') { return moves; }
            else
            {
                //Otherwise : Test all tiles
                for (int i = 0; i <= MAXROW; i++)
                {
                    for (int j = 0; j <= MAXCOL; j++)
                    {
                        Point p2 = new Point(i, j);
                        if (IsValidMove(p1, p2, board)) { moves.Add(p2); }
                    }
                }
                return moves;
            }
        }



        public static IEnumerator<GameObject> FadeOut(GameObject o, float duration)
        {
            float t = 0f;

            Renderer r = o.GetComponent<Renderer>();
            Color newColor = r.material.color;
            newColor.a = 0;

            while (t < 1)
            {
                r.material.color = Color.Lerp(r.material.color, newColor, t);
                t += Time.deltaTime / duration;
                yield return null;
            }
            o.SetActive(false);
        }

        public static IEnumerator<GameObject> FadeIn(GameObject o, float duration)
        {
            float t = 0f;

            Renderer r = o.GetComponent<Renderer>();
            Color newColor = r.material.color;
            newColor.a = 1.0f;

            o.SetActive(true);
            while (t < 1)
            {
                r.material.color = Color.Lerp(r.material.color, newColor, t);
                t += Time.deltaTime / duration;
                yield return null;
            }

        }

        // Moves a on object o smoothly
        public static IEnumerator<GameObject> SmoothMove(Transform o, Transform end, float seconds)
        {
            float t = 0.0f;
            Vector3 startpos = o.transform.position;
            Quaternion startrot = o.transform.rotation;
            while (t <= 1.0f)
            {
                t += Time.deltaTime / seconds;
                o.position = Vector3.Lerp(startpos, end.position, Mathf.SmoothStep(0.0f, 1.0f, t));
                o.rotation = Quaternion.Lerp(startrot, end.rotation, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }
        }

        // Moves a on object o smoothly
        public static IEnumerator<GameObject> SmoothMove(Transform o, Vector3 endpos, float seconds)
        {
            float t = 0.0f;
            Vector3 startpos = o.transform.position;
            while (t <= 1.0f)
            {
                t += Time.deltaTime / seconds;
                o.position = Vector3.Lerp(startpos, endpos, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }
        }

        internal static IEnumerator<GameObject> SmoothMoveUI(RectTransform rt, Vector3 startPos, Vector3 endPos, float time)
        {
            float elapsed = 0;
            while (elapsed < time)
            {
                elapsed += Time.deltaTime;
                rt.anchoredPosition3D = Vector3.Lerp(startPos, endPos, Mathf.SmoothStep(0, 1, elapsed/time));
                yield return null;
            }
        }
    }
}

