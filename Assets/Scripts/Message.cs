using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace ChessGo
{
    class Message
    {
        public int message;
        public String[] parameters;

        public Message(int message)
        {
            this.message = message;
            this.parameters = new String[0];
        }
        public Message(int message, String[] parameters)
        {
            this.message = message;
            this.parameters = parameters;
        }
    }
}
