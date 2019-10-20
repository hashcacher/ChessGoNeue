using UnityEngine;
using UnityEngine.UI;


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
            title.text = "Game " + data.id;
            opponent.text = data.opponent;

            timeLeft.text = data.timeLeft;
            myTurn.text = data.myTurn ? "My turn" : opponent + "'s turn";
        }
    }
}
