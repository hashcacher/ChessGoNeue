using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;
using UnityEngine.EventSystems;
using System.Collections;
using System.Text;
using System.Linq;
using UnityEngine.SceneManagement;

namespace ChessGo
{
    public class MainMenu : MonoBehaviour
    {
        public Canvas canvas;

        // First panel
        public InputField nickname;
        public Button playOnline;
        public Button playHotseat;
        public Text errorMessage;

        public Text countdown;
        public Text searching;
        public Button backButton;

        private float queueTime;
        private bool inQueue;

        public ToggleGroup toggleGroup;

        public RectTransform [] menuPanels;

        private MarkovNameGenerator generator;

        public delegate void ConnectDelegate(bool success);
        private int failedConnections = 0;
        private string playerID = "null";
        private bool matching = false;

        void Awake() {
            // Get / generate player ID
            if (PlayerPrefs.HasKey("ID")) {
                this.playerID = PlayerPrefs.GetString("ID");
            } else {
                this.playerID = RandomString(20);
                PlayerPrefs.SetString("ID", this.playerID);
            }

            canvas = GameObject.Find("Canvas").GetComponent<Canvas>();
            menuPanels = new RectTransform[3];
            int count = 0;
            for (int i = 0; i < canvas.transform.childCount; i++) {
                var child = canvas.transform.GetChild(i).GetComponent<RectTransform>();
                if (child) {
                    menuPanels[count++] = child; 
                }
            }

            toggleGroup = GameObject.FindObjectOfType<ToggleGroup>();
            nickname = GameObject.Find("Nickname Input").GetComponent<InputField>();
            playOnline = GameObject.Find("Online").GetComponent<Button>();
            errorMessage = GameObject.Find("Error Text").GetComponent<Text>();

            countdown = GameObject.Find("Countdown").GetComponent<Text>();
            searching = GameObject.Find("Searching Header").GetComponent<Text>();
        }

        void Start() {
        }

        public static string RandomString(int length)
        {
            const string chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
            return new string(Enumerable.Repeat(chars, length)
              .Select(s => s[Random.Range(0, s.Length)]).ToArray());
        }

        IEnumerator MatchMe() {
            var request = new MatchRequest();
            request.secret = playerID;
            var msg = JsonUtility.ToJson(request);
            var host = "https://chessgo.xyz";
            if (Application.isEditor) {
                host = "localhost:8080";
            }

            // Post to our api
            using (UnityWebRequest www = GoodPost(host + "/v1/matchMe", msg))
            {
                matching = true;
                yield return www.SendWebRequest();

                if (www.isNetworkError) {
                    // Exponential backoff
                    Debug.LogError("MatchMe network error: " + www.error);
                    this.failedConnections++;
                    yield return new WaitForSeconds(Mathf.Pow(2f, this.failedConnections) / 10f * Random.Range(.5f, 1.0f));

                    if (this.failedConnections >= 5) {
                        errorMessage.text = "Server is experiencing technical difficulties. Please try again later.";
                        OnBackPress();
                    } else {
                        StartCoroutine(MatchMe());
                    }
                } else if (www.isHttpError) {
                    OnBackPress();
                    var response = JsonUtility.FromJson<MatchResponse>(www.downloadHandler.text);
                    if (response != null) {
                        Debug.Log("MatchMe error: " + response.err);
                    }
                } else {
                    // Found! Start game
                    matching = false;
                    var response = JsonUtility.FromJson<MatchResponse>(www.downloadHandler.text);
                    if (response != null) {
                        if (response.haveMatch) {
                            Debug.Log("Found match!");
                            StartCoroutine(Countdown());
                        } else {
                            Debug.Log("MatchMe 200 doesnt have match wtf error: " + response.err);
                        }
                    } else {
                        Debug.Log("Couldnt parse server response: " + www.downloadHandler.text);
                    }
                }
            }
        }

        IEnumerator Countdown() {
            StopCoroutine("Dots");

            searching.transform.Translate(-10, 0, 0);
            searching.text = "Match Found";
            for (int i = 3; i >= 0; i--) {
                countdown.text = i.ToString();
                yield return new WaitForSecondsRealtime(1);
            }

            PlayerPrefs.SetInt("Hotseat", 0);
            SceneManager.LoadScene("GameBoard");
        }

        IEnumerator Dots() {
            var dotCount = 3;
            while (true) {
                searching.text = "Searching" + string.Concat(Enumerable.Repeat(".", dotCount).ToArray());
                yield return new WaitForSeconds(0.3f);
                dotCount = dotCount % 5 + 1;
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
            this.failedConnections = 0;
            StartCoroutine(MatchMe());
            StartCoroutine("Dots");
        }

        public void OnRulesPress()
        {
            SlidePanels(2);
        }

        void Unmatchme() {

        } 

        public void OnBackPress() {
            if (matching) {
                Unmatchme();
                matching = false;
            }
            int delta = (int)menuPanels[0].localPosition.x / 800;
            SlidePanels(delta);
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
