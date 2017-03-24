using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using MySql.Data.MySqlClient;

namespace CardGameServer
{
    class Program
    {
        public const int SERVER_PORT = 8888;
        public const int MAX_USERS = 100;
        public static int usersConnected = 0;
        public static List<Handler> clientList = new List<Handler>();
        public static MySqlConnection mConnection;

        // Thread signal.
        public static int Main(String[] args)
        {
            try
            {
                mConnection = Server.MySQL_Connect();
                mConnection.Open();
            }
            catch (Exception ex)
            {
                Server.consoleWrite("Failed to connect to MySQL (" + ex.ToString() + ")");
                Console.ReadKey();
                return 1;
            }
            AsyncSocketListener.StartListening();
            return 0;
        }

        public static bool DisconnectClient(Handler handler)
        {
            return clientList.Remove(handler);
        }
    }
}
