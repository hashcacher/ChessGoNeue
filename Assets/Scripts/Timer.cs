using UnityEngine;
using UnityEngine.Events;
using UnityEngine.UI;
using System;
using System.Collections;

namespace ChessGo {
    class Timer : MonoBehaviour {
        public Text text;
        public UnityEvent timeout;
        public bool amWhite;
        private bool ourTurn;
        private TimeSpan timeLeft;
        private TimeSpan timeLeftAtStartOfTurn;
        private DateTime turnStarted;

        void Awake() {
            text = GetComponent<Text>();
        }

        void Start() {
            if (false && UnitySingleton.hotseat) {
                this.gameObject.SetActive(false);
            } else {
                // Set initial clocks with a fake move response
                MoveResponse init = new MoveResponse();
                init.blackLeft = UnitySingleton.match.duration;
                init.whiteLeft = UnitySingleton.match.duration;
                //init.blackLeft = 5l * 1000 * 1000 * 1000 * 60;
                //init.whiteLeft = 5l * 1000 * 1000 * 1000 * 60;
                DateTime now = DateTime.Now;
                init.blackTurnStarted = now.ToString();
                init.whiteTurnStarted = now.ToString();
                UnitySingleton.lastMove = init;
                UpdateClocks();

                if (amWhite) {
                    OnTurnStart();
                }
            }
        }

        public void OnTurnStart() {
            ourTurn = true;
            UpdateClocks();
            StartCoroutine("tick");
        }

        public void OnTurnEnd() {
            ourTurn = false;
            StopCoroutine("tick");
        }

        IEnumerator tick() {
            while (ourTurn) {
                timeLeft = timeLeftAtStartOfTurn - (DateTime.Now - turnStarted);
                if (timeLeft.TotalMilliseconds <= 0) {
                    Debug.Log("nope " + timeLeft.TotalMilliseconds);
                    break;
                } else {
                    Debug.Log("tick");
                    text.text = timeLeft.Minutes + ":" + timeLeft.Seconds;
                }
                yield return new WaitForSeconds(1f);
            }
        }

        // Go server gives us time in ns(?)
        void UpdateClocks() {
            long myTimeLeft = amWhite ? UnitySingleton.lastMove.whiteLeft : UnitySingleton.lastMove.blackLeft;
            long mLeft =  myTimeLeft / 1000 / 1000 / 1000 / 60;
            float sLeft = (myTimeLeft - (mLeft * 1000 * 1000 * 1000 * 60)) / 1000f;
            TimeSpan.TryParse("0:"+ mLeft + ":" + sLeft, out timeLeftAtStartOfTurn);

            Debug.Log(myTimeLeft);
            Debug.Log(mLeft);
            Debug.Log(sLeft);
            Debug.Log(timeLeftAtStartOfTurn);

            string myStarted = amWhite ? UnitySingleton.lastMove.whiteTurnStarted : UnitySingleton.lastMove.blackTurnStarted;
            bool success = DateTime.TryParse(myStarted, out turnStarted);
        }

    }
}
