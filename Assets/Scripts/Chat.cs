using UnityEngine;

namespace ChessGo {
    class Chat : MonoBehaviour {
        int failedChatConnections = 0;
        public Text chatBox;

        // Prefabs
        public SpeechBubble chatBubble;

        void Awake() {
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
                var host = GetServerHost(); // Post to our api
                using (UnityWebRequest www = GoodPost(host + "/v1/getChat", msg))
                {
                    yield return www.SendWebRequest();

                    if (www.isNetworkError) {
                        // Exponential backoff
                        Debug.LogError("ReceiveMoves network error: " + www.error);
                        this.failedChatConnections++;
                        yield return new WaitForSeconds(Mathf.Pow(2f, this.failedChatConnections) / 10f * Random.Range(.5f, 1.0f));

                        if (this.failedChatConnections >= 10) {
                            errorMessage.text = "Server is experiencing technical difficulties. Please try again later.";
                            OnBackPress();
                        } else {
                            // TODO give up on the server
                            // Go back to main menu
                            StartCoroutine(MatchMe());
                        }
                    } else if (www.isHttpError) {
                        var response = JsonUtility.FromJson<ChatResponse>(www.downloadHandler.text);
                        if (response != null) {
                            Debug.Log("Chat error: " + response.err);
                        }
                    } else {
                        this.failedChatConnections = 0;

                        var response = JsonUtility.FromJson<ChatResponse>(www.downloadHandler.text);
                        if (response != null) {
                            if (response.chat) {
                                // Move received!
                                AppendNewChat(response.chat);
                            }
                        }
                    }
                }
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
            GameObject myKing = IAmBlack ? GameObject.Find("BlackKing(Clone)") : GameObject.Find("WhiteKing(Clone)");
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
            if (usingServer)
            {
                AsyncServerConnection.Send(Messages.CHAT, inputChat.text);
                AsyncServerConnection.Receive();
            }

            //reset the input box
            inputChat.text = "";
            inputChat.Select();

        }

    }
}
