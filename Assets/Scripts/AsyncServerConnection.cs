using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net;
using System.Net.Sockets;
using System.Threading;

//using DeckBuilder = Card_Game.DeckBuilder;
using UnityEngine;

namespace ChessGo
{

    // State object for receiving data from remote device.
    public class StateObject
    {
        // Client socket.
        public Socket workSocket = null;
        // Size of receive buffer.
        public const int BufferSize = 1024;
        // Receive buffer.
        public byte[] buffer = new byte[BufferSize];
        // Received data string.
        public StringBuilder sb = new StringBuilder();
    }

    class AsyncServerConnection
    {
        // The port number for the remote device.
        private const int port = Client.SERVER_PORT;

        // ManualResetEvent instances signal completion.
        private static ManualResetEvent connectDone =
            new ManualResetEvent(false);
        private static ManualResetEvent sendDone =
            new ManualResetEvent(false);
        private static ManualResetEvent receiveDone =
            new ManualResetEvent(false);

        public static Queue<Message> messageQueue;

        //this might be problematic
        public static Socket clientSocket;
        // The response from the remote device.
        private static String response = String.Empty;

        public static void StartClient(MainMenu.ConnectDelegate callback)
        {
            // Connect to a remote device.
            try
            {
                Debug.Log("ServerConnection, Starting client");
                // Establish the remote endpoint for the socket.
                // The name of the 
                // remote device is "host.contoso.com".
                IPHostEntry ipHostInfo = Dns.GetHostEntry(Client.SERVER_IP);
                Debug.Log("ServerConnection, got host entry");
                IPAddress ipAddress = ipHostInfo.AddressList[0];
                IPEndPoint remoteEP = new IPEndPoint(ipAddress, port);
                //IPEndPoint[] remoteEPs = Enumerable.Select(ipHostInfo.AddressList, (x) => new IPEndPoint(x, port)).ToArray();
                Debug.Log("ServerConnection, made IPEndpoint" + ipAddress.ToString() + port);

                // Create a TCP/IP socket.
                clientSocket = new Socket(ipAddress.AddressFamily, SocketType.Stream, ProtocolType.Tcp);

                Debug.Log("ServerConnection, made new socket");

                // Connect to the remote endpoint.
                clientSocket.BeginConnect(remoteEP,
                    new AsyncCallback(ConnectCallback), clientSocket);
                //Debug.Log("ServerConnection, connected to socket");

                int numTries = 0;
                while (!clientSocket.Connected && numTries < 3)
                {
                    connectDone.WaitOne(1000);
                    numTries++;
                    Debug.Log("ServerConnection, waited one, socket connected? : " + clientSocket.Connected);
                }

                if (clientSocket.Connected)
                {
                    callback(true);
                }
                else
                {
                    callback(false);
                    clientSocket.Close();
                    Application.Quit();
                }
                
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        private static void ConnectCallback(IAsyncResult ar)
        {
            try
            {
                // Retrieve the socket from the state object.
                Socket client = (Socket)ar.AsyncState;

                // Complete the connection.
                client.EndConnect(ar);

                Debug.Log("Socket connected to " +
                    client.RemoteEndPoint.ToString());

                messageQueue = new Queue<Message>();

                // Signal that the connection has been made.
                connectDone.Set();
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        public static void Receive(Socket client)
        {
            try
            {
                // Create the state object.
                StateObject state = new StateObject();
                state.workSocket = client;

                // Begin receiving the data from the remote device.
                client.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                    new AsyncCallback(ReceiveCallback), state);
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        public static void Receive()
        {
            try
            {
                // Create the state object.
                StateObject state = new StateObject();
                state.workSocket = clientSocket;

                // Begin receiving the data from the remote device.
                clientSocket.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                    new AsyncCallback(ReceiveCallback), state);
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        private static void ReceiveCallback(IAsyncResult ar)
        {
            try
            {
                // Retrieve the state object and the client socket 
                // from the asynchronous state object.
                StateObject state = (StateObject)ar.AsyncState;
                Socket client = state.workSocket;

                // Read data from the remote device.
                int bytesRead = client.EndReceive(ar);
                state.sb.Append(Encoding.ASCII.GetString(state.buffer, 0, bytesRead));

                Debug.Log("bytes read: " + bytesRead);
                if (bytesRead <= 0)
                {
                    Debug.Log("reading more " + state.sb.ToString());
                    // Get the rest of the data.
                    client.BeginReceive(state.buffer, 0, StateObject.BufferSize, 0,
                        new AsyncCallback(ReceiveCallback), state);
                }
                else
                {
                    // All the data has arrived; put it in response.
                    if (state.sb.Length > 1)
                    {
                        Debug.Log("all arrived " + response);
                        List<string> receivedMessage = System.Text.Encoding.ASCII.GetString(state.buffer).Split(';').ToList();
                        int messageID = Convert.ToInt32(receivedMessage[0]);

                        receivedMessage.RemoveAt(0);

                        messageQueue.Enqueue(new Message(messageID, receivedMessage.ToArray()));
                    }
                    // Signal that all bytes have been received.
                    receiveDone.Set();
                }
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        public static void Send(int message, params String[] parameters)
        {
            if (clientSocket != null && clientSocket.Connected)
            {
                string sendingMessage = message + ";" + String.Join(";", parameters);

                byte[] byteData = System.Text.Encoding.ASCII.GetBytes(sendingMessage + "$");

                // Begin sending the data to the remote device.
                Debug.Log("Sending, client socket is: " + clientSocket.ToString());
                clientSocket.BeginSend(byteData, 0, byteData.Length, 0,
                    new AsyncCallback(SendCallback), clientSocket);

                sendDone.WaitOne(); //new
            }
        }

        private static void Send(Socket client, int message, params String[] parameters)
        {
            string sendingMessage = message + ";" + String.Join(";", parameters);

            byte[] byteData = System.Text.Encoding.ASCII.GetBytes(sendingMessage + "$");

            // Begin sending the data to the remote device.
            client.BeginSend(byteData, 0, byteData.Length, 0,
                new AsyncCallback(SendCallback), client);
        }

        private static void SendCallback(IAsyncResult ar)
        {
            try
            {
                // Retrieve the socket from the state object.
                Socket client = (Socket)ar.AsyncState;

                // Complete sending the data to the remote device.
                int bytesSent = client.EndSend(ar);
                Debug.Log("Sent " + bytesSent + " bytes to server.");

                // Signal that all bytes have been sent.
                sendDone.Set();
            }
            catch (Exception e)
            {
                Debug.Log(e.ToString());
            }
        }

        public static void Shutdown()
        {
            Send(Messages.DISCONNECT);

            // Release the socket.
            clientSocket.Shutdown(SocketShutdown.Both);
            clientSocket.Close();
        }

        public static string FixParam(string p)
        {
            return new string(p.Where(c => char.IsLetterOrDigit(c)).ToArray());
        }
    }
}
