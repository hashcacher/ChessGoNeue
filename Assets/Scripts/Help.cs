using UnityEngine;

namespace ChessGo {
    class Help : MonoBehaviour {
        public RectTransform helpPanel;

        void Awake() {
            helpPanel = canvas.Find("Help Panel").GetComponent<RectTransform>();
            helpPanel.gameObject.SetActive(false);
        }

        public void ShowHelp()
        {
            helpPanel.gameObject.SetActive(true);
            Vector3 startPos = new Vector3(0, 500, -20);
            Vector3 endPos = new Vector3(0, -230, -20);
            preventMoves = true; // Eventsystem
            StartCoroutine(Utilities.SmoothMoveUI(helpPanel, startPos, endPos, .5f));
        }

        public void CloseHelp()
        {
            StartCoroutine(CloseHelpHelper());
        }

        private IEnumerator CloseHelpHelper()
        {
            Vector3 startPos = new Vector3(0, -230, -20);
            Vector3 endPos = new Vector3(0, -1230, -20);
            StartCoroutine(Utilities.SmoothMoveUI(helpPanel, startPos, endPos, .5f));
            yield return new WaitForSeconds(.5f);
            helpPanel.gameObject.SetActive(false);
            preventMoves = false;
        }

        public interface IChatEvent : IEventSystemHandler
        {
            // functions that can be called via the messaging system
            void ChatOpened();
            void ChatClosed();
        }
    }
}
