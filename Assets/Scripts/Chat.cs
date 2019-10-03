using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;
using System.Collections;

namespace ChessGo {
    class Chat : MonoBehaviour {
        int failedChatConnections = 0;
        public Text chatBox;
        public InputField inputChat;

        // Prefabs
        public SpeechBubble chatBubble;

        void Awake() {
            // TODO chatBox = 
            // TODO inputChat = 
            StartCoroutine("ReceiveChat"); // TODO cleanup
        }

        void FixedUpdate() {
            if (Input.GetKeyDown("return")) { SubmitChat(); }
        }

        private void AppendNewChat(string msg)
        {
            chatBox.text += "\n" + msg;
        }


        IEnumerator ReceiveChat() {
            while (true) {
                var request = new MoveRequest(); // Same as the moveRequest, but different endpoint
                request.secret = UnitySingleton.secret;
                request.gameID = UnitySingleton.gameID;
                var msg = JsonUtility.ToJson(request);
                var host = Utilities.GetServerHost(); // Post to our api
                using (UnityWebRequest www = Utilities.GoodPost(host + "/v1/getChat", msg))
                {
                    yield return www.SendWebRequest();

                    if (www.isNetworkError) {
                        // Exponential backoff
                        Debug.LogError("ReceiveMoves network error: " + www.error);
                        this.failedChatConnections++;
                        yield return new WaitForSeconds(Mathf.Pow(2f, this.failedChatConnections) / 10f * Random.Range(.5f, 1.0f));

                        if (this.failedChatConnections >= 10) {
                            // TODO OnBackPress();
                        } else {
                        }
                    } else if (www.isHttpError) {
                        Debug.Log("Chat error: " + www.downloadHandler.text);
                    } else {
                        this.failedChatConnections = 0;

                        /* TODO
                        var response = JsonUtility.FromJson<ChatResponse>(www.downloadHandler.text);
                        if (response != null) {
                            if (response.chat) {
                                // Move received!
                                AppendNewChat(response.chat);
                            }
                        }
                        */
                    }
                }

                yield return new WaitForSeconds(1f);
            }
        }

        public void SubmitChat()
        {
            if (inputChat.text.Trim() == "")
            {
                Debug.Log("empty string");
                return;
            }

            chatBox.text += "\n" + inputChat.text;

            //speech bubbles
            GameObject myKing = null; // TODO: IAmBlack ? GameObject.Find("BlackKing(Clone)") : GameObject.Find("WhiteKing(Clone)");
            SpeechBubble b;
            if (myKing.transform.childCount == 0)
            {
                b = Instantiate(chatBubble) as SpeechBubble;
                b.transform.parent = myKing.transform;
                b.transform.position = myKing.transform.position;
                //b.transform.position = new Vector3()
            }
            else
            {
                b = myKing.transform.GetComponentInChildren<SpeechBubble>() as SpeechBubble;

            }

            if (b == null)
                Debug.LogError("Failed to get chat bubble");
            else
            {
                b.gameObject.SetActive(true);
                b.SetText(inputChat.text);
            }

            //send the chat message
            /* TODO
            if (usingServer)
            {
                AsyncServerConnection.Send(Messages.CHAT, inputChat.text);
                AsyncServerConnection.Receive();
            }
            */

            //reset the input box
            inputChat.text = "";
            inputChat.Select();

        }

    }
}
