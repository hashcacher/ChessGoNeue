using UnityEngine;
using UnityEngine.UI;
using System;


namespace ChessGo
{
    public class MyGame : MonoBehaviour {
        public Text title, opponent, timeLeft, myTurn;

        public void Awake() {
            title = this.transform.Find("Title").GetComponent<Text>();
            opponent = this.transform.Find("Opponent").GetComponent<Text>();
            timeLeft = this.transform.Find("TimeLeft").GetComponent<Text>();
            myTurn = this.transform.Find("MyTurn").GetComponent<Text>();
        }

        public void SetData(GamePublic data) {
            title.text = "Active game";
            opponent.text = "vs. " + data.opponent;
            myTurn.text = data.myTurn ? "My turn" : "Opponent's turn";

            long mLeft =  data.timeLeft / (long)60e9;
            float sLeft = (data.timeLeft - mLeft * (long)60e9) / 1e9f;
            TimeSpan timeSpan;
            TimeSpan.TryParse("0:" + mLeft + ":" + sLeft, out timeSpan);
            timeLeft.text = timeSpan.Minutes + ":" + timeSpan.Seconds + " left";
        }
    }
}
