using UnityEngine;
using UnityEngine.UI;
using UnityEngine.SceneManagement;
using System.Collections;

namespace ChessGo
{
    public class LandingScreen : MonoBehaviour
    {
        public InputField welcomeLabel;
        public InputField queueTimeLabel;
        public Text findMatchText;

        float queueTime;
        public bool inQueue;

        void Start()
        {
            inQueue = false;
            AsyncServerConnection.Send(Messages.LOADING);
            AsyncServerConnection.Receive();
        }

        void Update()
        {
            if (AsyncServerConnection.messageQueue.Count > 0)
                OnReceiveServerMessage(AsyncServerConnection.messageQueue.Dequeue());

            if (inQueue)
            {
                queueTime += Time.deltaTime;
                queueTimeLabel.text = "Searching....\n" + queueTime.ToString("F1") + " seconds in queue.";
            }

        }

        public void startQueueTimer()
        {
            if (inQueue)
            {
                inQueue = !inQueue;
                stopQueueTimer();

            }
            else
            {
                queueTimeLabel.enabled = true;
                queueTime = 0;
                findMatchText.text = "Cancel search";
                inQueue = true;
            }
        }

        public void stopQueueTimer()
        {
            findMatchText.text = "Find Match";
            queueTimeLabel.enabled = false;
        }

        void OnReceiveServerMessage(Message msg) //changed from object[] params
        {
            int message = msg.message;
            string[] parameters = msg.parameters;

            Debug.Log("Got a message!");
            switch (message)
            {
                case Messages.LOADING:
                    {
                        //this is loading too fast - find a way to fix this.
                        SetUsernameText(AsyncServerConnection.FixParam(parameters[0]));
                        break;
                    }
                case Messages.LOADGAME:
                    {
                        SceneManager.LoadScene("GameBoard");
                        break;
                    }
                case Messages.FINDMATCH: //not using this anymore
                {
                    AsyncServerConnection.Receive();
                    //wait for the LOADGAME
                    break;
                }
            }
        }


        public void onPressFindMatch()
        {
            //GUI.enabled = false;

            AsyncServerConnection.Send(Messages.FINDMATCH);
            AsyncServerConnection.Receive();
            startQueueTimer();
            Debug.Log("Pressed find match");
        }

        public void onPressPracticeBots()
        {
            //				GUI.enabled = false;

            //ServerConnection.sendMessage(Messages.FINDMATCH);
            Debug.Log("Pressed practice bots");
        }

        /*public void onPressDeckBuilder()
        {   
            //				GUI.enabled = false;

            Debug.Log("Pressed DeckBuilder");
            AsyncServerConnection.Send(Messages.DECKBUILDER);
            AsyncServerConnection.Receive();
        }*/

        public void SetUsernameText(string username)
        {
            GameObject.Find("Welcome Label").GetComponent<InputField>().text = "Welcome " + username;
            Client.username = username;
        }
    }
}