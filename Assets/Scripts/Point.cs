using System;
namespace ChessGo
{
    public struct Point
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



        public override bool Equals(System.Object obj)
        {
            // If parameter is null return false.
            if (obj == null)
            {
                return false;
            }

            // If parameter cannot be cast to Point return false.
            Point p = (Point)obj;
            if ((System.Object)p == null)
            {
                return false;
            }

            // Return true if the fields match:
            return (this.row == p.row) && (this.col == p.col);
        }

        //public override bool Equals(Point p2)
        //{
        //    return this.row == p2.row && this.col == p2.col;
        //}


    }
}

