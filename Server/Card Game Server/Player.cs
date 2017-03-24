using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace CardGameServer
{
    public class Player
    {
        int playerNumber;

        public Handler handler;

        private bool connected;
        public bool amBlack;

        public Player(int playerNum, Handler handle)
        {
            connected = false;
            playerNumber = playerNum;
            handler = handle;
        }

        public int GetPlayerNumber()
        {
            return playerNumber;
        }
        public void HasConnected()
        {
            connected = true;
        }
        public bool IsConnected()
        {
            return connected;
        }
        public void SendMessage(int message, params object[] parameters)
        {
            handler.SendMessage(message, parameters);
        }
    }
}
