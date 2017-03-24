using System;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Collections.Generic;
using System.Threading;
using System.Linq;

namespace CardGameServer
{
    // State object for reading client data asynchronously
    public class StateObject
    {
        // Client  socket.
        public Socket workSocket = null;
        // Size of receive buffer.
        public const int BufferSize = 1024;
        // Receive buffer.
        public byte[] buffer = new byte[BufferSize];
        // Received data string.
        public StringBuilder sb = new StringBuilder();

        public Handler messageHandler;
    }

    public class AsyncSocketListener
    {
        public static ManualResetEvent allDone = new ManualResetEvent(false);

        public AsyncSocketListener()
        {
        }

        public static void StartListening()
        {
            // Data buffer for incoming data.
            byte[] bytes = new Byte[1024];

            IPEndPoint localEndPoint = new IPEndPoint((Dns.Resolve(IPAddress.Any.ToString())).AddressList[0], 8383);

            // Create a TCP/IP socket.
            Socket listener = new Socket(AddressFamily.InterNetworkV6, SocketType.Stream, ProtocolType.Tcp);
            listener.SetSocketOption(SocketOptionLevel.IPv6, (SocketOptionName)27, 0);

            // Bind the socket to the local endpoint and listen for incoming connections.
            try
            {
                listener.Bind(new IPEndPoint(IPAddress.IPv6Any, 8383));
                //listener.Bind(localEndPoint);
                listener.Listen(100);

                while (true)
                {
                    // Set the event to nonsignaled state.
                    allDone.Reset();

                    // Start an asynchronous socket to listen for connections.
                    Console.WriteLine("Waiting for a connection on local IP " + listener.LocalEndPoint.ToString());
                    listener.BeginAccept(
                        new AsyncCallback(AcceptCallback),
                        listener);
                    // Wait until a connection is made before continuing.
                    allDone.WaitOne();
                }

            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }

            Console.WriteLine("\nPress ENTER to continue...");
            Console.Read();

        }

        public static void AcceptCallback(IAsyncResult ar)
        {
            Console.WriteLine("Accepting");
            // Signal the main thread to continue.
            allDone.Set();

            // Get the socket that handles the client request.
            Socket listener = (Socket)ar.AsyncState;
            Socket handler = listener.EndAccept(ar);

            // Create the state object.
            StateObject state = new StateObject();
            state.workSocket = handler;
            handler.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                new AsyncCallback(ReadCallback), state);

            Server.consoleWrite("[IP:" + (handler.RemoteEndPoint as IPEndPoint).Address.ToString() + "] has connected\n");

            state.messageHandler = new Handler(Program.mConnection, handler, ++Program.usersConnected);
            state.messageHandler.SendMessage(Messages.LOBBY_PLAYERS, PlayerNicknameList().ToArray());
            Program.clientList.Add(state.messageHandler);
        }

        private static IEnumerable<string> PlayerNicknameList()
        {
            return Program.clientList.Select((x) => x.nickname);
        }
        
        public static void Read(Handler handler)
        {
            Read(handler.clientSocket, handler);
        }
        public static void Read(Socket clientSocket, Handler messageHandler)
        {
            StateObject state = new StateObject();
            state.workSocket = clientSocket;
            state.messageHandler = messageHandler;
            clientSocket.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                                               new AsyncCallback(ReadCallback), state);
            Console.Write("Reading...");
        }

        public static void ReadCallback(IAsyncResult ar)
        {
            String content = String.Empty;

            // Retrieve the state object and the handler socket
            // from the asynchronous state object.
            StateObject state = (StateObject)ar.AsyncState;
            Socket handler = state.workSocket;

            // Read data from the client socket. 
            try
            {
                int bytesRead = handler.EndReceive(ar);

                if (bytesRead > 0)
                {
                    // There might be more data, so store the data received so far.
                    state.sb.Append(Encoding.ASCII.GetString(
                        state.buffer, 0, bytesRead));

                    // Check for end-of-file tag. If it is not there, read 
                    // more data.
                    content = state.sb.ToString();
                    if (content.IndexOf("$") > -1)
                    {
                        content = content.TrimEnd('$');

                        List<string> receivedMessage = new List<string>(content.Split(';'));
                        int messageID = Convert.ToInt32(receivedMessage[0]);

                        receivedMessage.RemoveAt(0);

                        state.messageHandler.OnReceiveClientMessage(messageID, receivedMessage.ToArray());

                        Console.WriteLine("Read {0} bytes from socket client {1}. Message: {2}\n", 
                            content.Length, state.messageHandler.clientNumber, messageID);
                    }
                }
                else
                {
                    // Not all data received. Get more.
                    handler.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                    new AsyncCallback(ReadCallback), state);

                    //this might not work, but we gotta track when the client times out
                    Shutdown(handler);
                }
            }
            catch (Exception e)
            {
                Console.WriteLine(e.Message);
                handler.Close();
            }
        }

        public static void Send(Socket handler, int message, params object[] parameters)
        {
            string sendingMessage = message + ";" + String.Join(";", parameters);

            // Convert the string data to byte data using ASCII encoding.
            byte[] byteData = Encoding.ASCII.GetBytes(sendingMessage);

            // Begin sending the data to the remote device.
            handler.BeginSend(byteData, 0, byteData.Length, 0,
                new AsyncCallback(SendCallback), handler);

            Console.WriteLine("Sending message {0} to client", message);
        }

        private static void SendCallback(IAsyncResult ar)
        {
            try
            {
                // Retrieve the socket from the state object.
                Socket handler = (Socket)ar.AsyncState;

                // Complete sending the data to the remote device.
                int bytesSent = handler.EndSend(ar);
                Console.WriteLine("Sent {0} bytes to client.\n", bytesSent);
            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }
        }

        private static void Shutdown(Socket handler)
        {
            handler.Shutdown(SocketShutdown.Both);
            handler.Close();
        }


    }
}