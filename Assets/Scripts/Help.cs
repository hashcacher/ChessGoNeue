using UnityEngine;
using UnityEngine.Events;
using System.Collections;

namespace ChessGo {
    class Help : MonoBehaviour {
        private Transform canvas;
        public RectTransform helpPanel;
        public UnityEvent helpOpened;
        public UnityEvent helpClosed;

        void Awake() {
            canvas = Camera.main.transform.Find("Canvas").transform;
            helpPanel = canvas.Find("Help Panel").GetComponent<RectTransform>();
            helpPanel.gameObject.SetActive(false);
        }

        public void ShowHelp()
        {
            helpPanel.gameObject.SetActive(true);
            Vector3 startPos = new Vector3(0, 500, -20);
            Vector3 endPos = new Vector3(0, -230, -20);
            helpOpened.Invoke();
            StartCoroutine(Util.SmoothMoveUI(helpPanel, startPos, endPos, .5f));
        }

        public void CloseHelp()
        {
            StartCoroutine(CloseHelpHelper());
        }

        private IEnumerator CloseHelpHelper()
        {
            Vector3 startPos = new Vector3(0, -230, -20);
            Vector3 endPos = new Vector3(0, -1230, -20);
            StartCoroutine(Util.SmoothMoveUI(helpPanel, startPos, endPos, .5f));
            yield return new WaitForSeconds(.5f);
            helpPanel.gameObject.SetActive(false);
            helpClosed.Invoke();
        }
    }
}
