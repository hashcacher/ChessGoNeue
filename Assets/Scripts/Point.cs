using System;
namespace ChessGo
{
    public class Point : IEquatable<Point>
    {
        public int row;
        public int col;

        public Point(int x, int y)
        {
            this.row = x;
            this.col = y;
        }

        public override string ToString()
        {
            return row + "," + col;
        }


        public override int GetHashCode() {
            return row*100 + col;
        }

        public bool Equals(Point p)
        {
            if (p == null) {
                return false;
            }

            // Return true if the fields match:
            return (this.row == p.row) && (this.col == p.col);
        }
    }
}

