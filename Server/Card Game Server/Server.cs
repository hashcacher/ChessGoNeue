using System;
using System.IO;
using System.Text;
using System.Collections.Generic;
using System.Security.Cryptography;
using System.Diagnostics;
using MySql.Data.MySqlClient;

namespace CardGameServer
{
    public class Server
    {
        public const string MYSQL_HOST = "localhost";
        public const string MYSQL_USERNAME = "root";
        public const string MYSQL_PASSWORD = "CHANGE";
        public const string MYSQL_DATABASE = "chessgo";

        public const int CLIENT_STATE_LOGIN = 1;
        public const int CLIENT_STATE_LOADING = 2;
        public const int CLIENT_STATE_LANDING = 3;
        public const int CLIENT_STATE_DECKBUILDER = 4;
        public const int CLIENT_STATE_SHOP = 5;
        public const int CLIENT_STATE_LOOKING = 6;
        public const int CLIENT_STATE_INGAME = 7;

        public const int GAME_STATE_START = 1;

        public const int MAX_CARDS = 40;

        public static List<object[]> queueList = new List<object[]>();

        public static List<Game> liveGames = new List<Game>();

        public const int FIRST_PLAYER = 0;
        public const int SECOND_PLAYER = 1;

        
        //player[] is [handler, int accountID, int elo]
        public static void addToQueue(object[] player)
        {
            Handler player1 = (Handler)player[0];
            
            if (queueList.Count > 0)
            {
                player1.FoundMatchLoadGame();
                
                Handler player2 = (Handler)queueList[0][0];
                player2.FoundMatchLoadGame();

                queueList.RemoveAt(0);

                Server.StartMatch(player1, player2);
            }
            else
            {
                queueList.Add(player);
                consoleWrite("Client #" + player1.clientNumber + " added to queue.");
                //AsynchronousSocketListener.Send(player1.clientSocket, Messages.FINDMATCH);
            }
        }

        public static void StartMatch(Handler player1, Handler player2)
        {
            consoleWrite("Clients " + player1.clientNumber + " and " + player2.clientNumber + " loading games.");
            
            Game createdGame = new Game(player1, player2);

            player1.clientState = CLIENT_STATE_INGAME;
            player2.clientState = CLIENT_STATE_INGAME;

            player1.currentGame = createdGame;
            player2.currentGame = createdGame;

            liveGames.Add(createdGame);
        }

        public static void file_put_contents(string path, string contents)
        {
            File.AppendAllLines(path, new[] { contents });
        }

        public static void serverLog(string file, string log)
        {
            consoleWrite(log);
        }

        public static void consoleWrite(string message)
        {
            Console.WriteLine(" >> " + message);
        }

        public static string getRunningPath()
        {
            return AppDomain.CurrentDomain.BaseDirectory;
        }

        public static string escape(string text)
        {
            return MySql.Data.MySqlClient.MySqlHelper.EscapeString(text.ToString()).Replace("_", "\\_").Replace("%", "\\%");
        }

        public static MySqlConnection MySQL_Connect()
        {
            return new MySqlConnection("Server=" + MYSQL_HOST + ";" + "Database=" + MYSQL_DATABASE + ";" + "Uid=" + MYSQL_USERNAME + ";" + "Pwd=" + MYSQL_PASSWORD + ";");
        }

        public static void UpdateElo()
        {

        }

        public static string password(string pass)
        {
            SHA512 alg = SHA512.Create();

            byte[] result = alg.ComputeHash(Encoding.UTF8.GetBytes(pass));
 
            return BitConverter.ToString(result).Replace("-", "");
        }

        public static int GetUnixTimestamp()
        {
            return (Int32)(DateTime.UtcNow.Subtract(new DateTime(1970, 1, 1))).TotalSeconds;
        }
    }
}
