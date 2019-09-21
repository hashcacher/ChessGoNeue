using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;
using UnityEngine.EventSystems;
using System.Collections;
using System.Text;
using UnityEngine.SceneManagement;

namespace ChessGo
{
    public class MainMenu : MonoBehaviour
    {
        EventSystem system;

        //deprecated. we don't login anymore
        public InputField username;
        public InputField password;

        public InputField nickname;

        public Canvas canvas;

        public Text errorMessage;

        public Text queueTimeText;
        public Text playButtonText;

        private float queueTime;
        private bool inQueue;

        public ToggleGroup toggleGroup;

        public RectTransform [] menuPanels;

        private MarkovNameGenerator generator;

        public delegate void ConnectDelegate(bool success);

        void Start()
        {
            system = EventSystem.current;// EventSystemManager.currentSystem;

            try
            {
            }
            catch
            {
                errorMessage.text = "The server is currently offline. Please try again later.";
                return;
            }
            StartCoroutine(MatchMe());


            toggleGroup = GameObject.FindObjectOfType<ToggleGroup>();
        }

        IEnumerator MatchMe()
        {
            var msg = string.Format("{{ \"clientID\": \"{0}\"}}", "bob");
            Debug.Log(msg);
            using (UnityWebRequest www = GoodPost("http://localhost:8080/v1/matchMe", msg))
            {
                yield return www.Send();
                Debug.Log(www.downloadHandler.isDone);
                //Debug.Log(www.downloadHandler.);

                if (www.isNetworkError || www.isHttpError)
                {
                    Debug.LogError(www.error);
                }
                else
                {
                    //Debug.Log("Match Me complete");
                }
            }
        }

        UnityWebRequest GoodPost(string url, string body) {
           byte[] bytes = Encoding.ASCII.GetBytes(body);
           var request             = new UnityWebRequest(url);
           request.uploadHandler   = new UploadHandlerRaw(bytes);
           request.downloadHandler = new DownloadHandlerBuffer();
           request.method          = UnityWebRequest.kHttpVerbPOST;
           request.timeout = 120;
           return request;
        }


        void ConnectCallback(bool success)
        {
            canvas = GameObject.Find("Canvas").GetComponent<Canvas>();
            if (!success)
            {
                Transform mainPanel = canvas.transform.Find("Main Panel").transform;
                mainPanel.Find("Error Text").GetComponent<Text>().text = "Can't connect to server";
                mainPanel.Find("Online").GetComponent<Button>().interactable = false;

            }

        }

        void SetLobbyPlayers(string[] clients)
        {
            Text lobbyPlayers = GameObject.Find("Canvas").GetComponentInChildren<Text>();
            lobbyPlayers.text = string.Join("\n", clients);
            Debug.Log(clients.Length + " players in lobby");
        }

        // Update is called once per frame
        void Update()
        {
            if (inQueue)
            {
                queueTime += Time.deltaTime;
                queueTimeText.text = queueTime.ToString("F1") + " seconds in queue.";
            }
        }


        void OnReceiveServerMessage(Message msg) //changed from object[] params
        {
            int message = msg.message;
            string[] parameters = msg.parameters;

            switch (message)
            {
                case Messages.FINDMATCH:
                    {
                        if(parameters.Length > 0)
                        {
                            PlayerPrefs.SetInt("ID", System.Convert.ToInt32(parameters[0]));
                        }
                        break;
                    }
                case Messages.LOADGAME:
                {
                    PlayerPrefs.SetInt("Hotseat", 0);
                    SceneManager.LoadScene("GameBoard");
                    break;
                }
                case Messages.LOBBY_PLAYERS:
                {
                    SetLobbyPlayers(parameters);
                    break;
                }
            }
        }

        public void OnHotseatPress()
        {
            PlayerPrefs.SetInt("Hotseat", 1);
            SceneManager.LoadScene("GameBoard");
        }

        public void SlidePanels(int delta)
        {
            foreach(RectTransform rt in menuPanels)
            {
                StartCoroutine(SmoothMoveRectTransform(rt, new Vector3(rt.localPosition.x - 800 * delta, 0, 0), .5f));
            }
        }

        public void OnOnlinePress()
        {
            SlidePanels(1);
        }

        public void OnRulesPress()
        {
            SlidePanels(2);
        }

        public void OnBackPress()
        {
            int delta = (int)menuPanels[0].localPosition.x / 800;
            SlidePanels(delta);
        }

        public void OnPlayPress()
        {
            //difficulty
            System.Collections.Generic.IEnumerable<Toggle> toggles = toggleGroup.ActiveToggles();
            int difficulty = 0;
            foreach(Toggle t in toggles)
                if(t.name == "Beginner")
                    difficulty = 1;
                else if(t.name == "Intermediate")
                    difficulty = 2;
                else
                    difficulty = 3;

            //id
            int myID = PlayerPrefs.GetInt("ID");

            AsyncServerConnection.Send(Messages.FINDMATCH, difficulty.ToString(), myID.ToString(), nickname.text);
            AsyncServerConnection.Receive();
            StartQueueTimer();
            Debug.Log("Pressed find match");
        }

        public void StartQueueTimer()
        {
            if (inQueue)
            {
                inQueue = !inQueue;
                StopQueueTimer();
            }
            else
            {
                queueTimeText.enabled = true;
                queueTime = 0;
                playButtonText.text = "Looking...";
                inQueue = true;
            }
        }


        public void StopQueueTimer()
        {
            playButtonText.text = "Play";
            queueTimeText.enabled = false;
        }

        // Moves a on object o smoothly
        public static IEnumerator SmoothMoveRectTransform(RectTransform t, Vector2 endpos, float seconds)
        {
            float time = 0.0f;
            Vector2 startpos = t.localPosition;
            while (time <= seconds)
            {
                time += Time.deltaTime;
                t.localPosition = Vector2.Lerp(startpos, endpos, Mathf.SmoothStep(0.0f, 1.0f, time/seconds));
                yield return null;
            }
        }

        public void OnGeneratePress()
        {
            GenerateName();
        }

        private void GenerateName()
        {
            if (generator == null)
                generator = new MarkovNameGenerator();
            nickname.text = generator.NextName;
        }
    }
}
