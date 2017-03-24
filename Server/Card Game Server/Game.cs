using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace CardGameServer
{
    public class Game
    {
        Player player1;
        Player player2;

        char[,] board;

        public int gameState;

        int currentTurn;

        public Game(Handler p1, Handler p2)
        {
            player1 = new Player(Server.FIRST_PLAYER, p1);
            player2 = new Player(Server.SECOND_PLAYER, p2);

            //setting the player for the client handlers
            p1.gamePlayer = player1;
            p2.gamePlayer = player2;

            board = new char[12, 12];
        }

        public void SendChat(Handler from, string message)
        {
            Player to = from == player1.handler ? player2 : player1;
            to.SendMessage(Messages.CHAT, message);
        }

        //Go Move
        public void MakeMove(Handler handler, int r1, int c1, params String[] move)
        {
            //if it's not this player's move
            if (handler == player1.handler && currentTurn == 1 || handler == player2.handler && currentTurn == 0)
            {
                AsyncSocketListener.Read(handler);
                return;
            }

            char piece = handler.gamePlayer.amBlack ? 'S' : 's';
            board[r1, c1] = piece;

            Server.consoleWrite("Placed " + piece + " at " + r1 + ", " + c1);

            if (handler == player1.handler)
            {   //now player 2 goes
                player2.SendMessage(Messages.MOVE, move);
                currentTurn = 1;
                AsyncSocketListener.Read(player2.handler);
            }
            else
            { //now player 1 goes
                player1.SendMessage(Messages.MOVE, move);
                currentTurn = 0; 
                AsyncSocketListener.Read(player1.handler);
            }
        }

        //Chess Move From: (r1,c1) To: (r2,c2)
        public void MakeMove(Handler handler, int r1, int c1, int r2, int c2, params String[] move)
        {
            //if it's not this player's move
            if(handler == player1.handler && currentTurn == 1 || handler == player2.handler && currentTurn == 0)
            {
                AsyncSocketListener.Read(handler);
                return;
            }

            char piece = board[r1,c1];
            board [r2, c2] = piece;
            board [r1, c1] = '\0';

            Server.consoleWrite("Moved " + r1 + ", " + c1 + " to " + r2 + ", " + c2);

            if (handler == player1.handler)
            {
                player2.SendMessage(Messages.MOVE, move);
                currentTurn = 1;
                AsyncSocketListener.Read(player2.handler);
            }
            else
            { 
                player1.SendMessage(Messages.MOVE, move);
                currentTurn = 0;
                AsyncSocketListener.Read(player1.handler);
            }
        }

        public Player GetPlayerByNum(int num)
        {
            if (player1.handler.clientNumber == num)
                return player1;
            else
                return player2;
        }
        public bool BothConnected()
        {//FIXME  change to && in production.
            return player1.IsConnected() && player2.IsConnected();
        }
        public void StartGame()
        {
            Random playerTurn = new Random();
            this.currentTurn = playerTurn.Next(Server.FIRST_PLAYER, Server.SECOND_PLAYER + 1);

            gameState = Server.GAME_STATE_START;

            player1.amBlack = currentTurn == 0;
            player2.amBlack = currentTurn == 1;

            player1.SendMessage(Messages.STARTGAME, currentTurn == 0 ? 1 : 0);
            player2.SendMessage(Messages.STARTGAME, currentTurn == 1 ? 1 : 0);

            AsyncSocketListener.Read(currentTurn == 1 ? player2.handler : player1.handler);
        }
        public void EndGame(Player loser)
        {


        }
    }
}
