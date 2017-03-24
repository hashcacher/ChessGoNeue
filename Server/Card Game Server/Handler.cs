using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net.Sockets;
using System.Threading;
using System.IO;
using System.Diagnostics;
using MySql.Data.MySqlClient;

namespace CardGameServer
{
    public class Handler
    {
        const int CLIENT_STATE_LOGIN        = 1;
        const int CLIENT_STATE_LOADING      = 2;
        const int CLIENT_STATE_LANDING      = 3;
        const int CLIENT_STATE_DECKBUILDER  = 4;
        const int CLIENT_STATE_SHOP         = 5;
        const int CLIENT_STATE_LOOKING      = 6;
    
        public Socket clientSocket;
        public int clientNumber;
        string clientIP;

        public int ID;
        public string nickname;
        public int clientState;

        public Game currentGame;

        public Player gamePlayer;

        MySqlConnection mConnection;
        MySqlDataReader mData;

        public Handler(MySqlConnection myConnection, Socket mySock, int num)
        {
            mConnection = myConnection;

            clientState = Server.CLIENT_STATE_LOGIN;
            clientState = CLIENT_STATE_LOGIN;

            clientSocket = mySock;

            clientNumber = num;
        }

        public void OnReceiveClientMessage(int message, params string[] parameters)
        {
            Console.WriteLine("Handing message " + message + " from client in state " + clientState);

            foreach (string str in parameters)
                Console.Write(str + ", ");
            Console.WriteLine();

            switch (message)
            {
                case Messages.DISCONNECT:
                    {
                        Disconnect();
                        break;
                    }
                case Messages.FINDMATCH: //params (difficulty, ID, nickname)
                    {
                        clientState = Server.CLIENT_STATE_LOOKING;

                        int difficulty = int.Parse(parameters[0]);
                        int ID = int.Parse(parameters[1]);//not using right now
                         
                        for (var listIndex = 0; listIndex < Server.queueList.Count; listIndex++)
                        {
                            if (Server.queueList[listIndex][0] == this)
                            {
                                // leave queue
                                Console.WriteLine("removed client # " + clientNumber + " from queue");
                                Server.queueList.RemoveAt(listIndex);
                                SendMessage(Messages.FINDMATCH);
                                AsyncSocketListener.Read(clientSocket, this);
                                return;
                            }
                        }

                        // add to queue
                        try
                        {
                            Server.addToQueue(new object[] { this, ID, });
                            AsyncSocketListener.Read(clientSocket, this);
                            break;
                        }
                        catch (Exception ex)
                        {
                            Server.serverLog("error_log.txt", ex.ToString());
                        }

                        break;
                    }
                case Messages.STARTGAME: //the client has loaded the game
                    {//FIXME: the client currently doesn't receive this message.
                        Server.consoleWrite("Client " + clientNumber + " has LOADGAME");
                        currentGame.GetPlayerByNum(clientNumber).HasConnected();
                        if (currentGame.BothConnected())
                        {
                            currentGame.StartGame();
                            Console.WriteLine("Both connected, sending STARTGAME");
                        }
                        break;
                    }
                case Messages.MOVE:
                    {
                        String[] fromMove = Server.escape(parameters[0].ToString()).Split(',');
                        
                        //check if client state is in play mode?
                        
                        //if player placed a Go stone
                        if (parameters.Length == 1)
                        {
                            currentGame.MakeMove(this, Int32.Parse(fromMove[0]), Int32.Parse(fromMove[1]), (string[]) parameters);
                        }
                        else
                        {
                            String[] toMove = Server.escape(parameters[1].ToString()).Split(',');

                            //updates the game state and tells the other player to make the next move. then listens to that player
                            currentGame.MakeMove(this, Int32.Parse(fromMove[0]), Int32.Parse(fromMove[1]), 
                                Int32.Parse(toMove[0]), Int32.Parse(toMove[1]), (string[])parameters);
                        }
                        break;
                    }
                case Messages.CHAT:
                    {
                        String msg = Server.escape(parameters[0].ToString());
                        currentGame.SendChat(this, msg);
                        break;
                    }
                case Messages.WIN:
                    {
                        break;
                    }
                case Messages.LOSE:
                    {
                        currentGame.EndGame(this.gamePlayer);
                        break;
                    }
            }
        }

        public void FoundMatchLoadGame()
        {
            SendMessage(Messages.LOADGAME);
            AsyncSocketListener.Read(clientSocket, this);
        }

        public void SendMessage(int message, params object[] parameters)
        {
            AsyncSocketListener.Send(clientSocket, message, parameters);
        }

        public bool IsConnected(Socket clientSocket)
        {
            try
            {
                return !(clientSocket.Poll(1000000, SelectMode.SelectRead) && clientSocket.Available == 0);
            }
            catch (SocketException) { return false; }
        }

        private void Disconnect()
        {
            clientSocket.Shutdown(SocketShutdown.Both);
            clientSocket.Close();
            Program.DisconnectClient(this);
        }
    }
}