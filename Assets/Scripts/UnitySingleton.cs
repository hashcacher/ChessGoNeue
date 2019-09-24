using UnityEngine;
using System.Collections;

namespace ChessGo
{
    // Stays alive between scenes
    public class UnitySingleton : MonoBehaviour
    {
        private static UnitySingleton instance = null;
        public static UnitySingleton Instance
        {
            get { return instance; }
        }
        public bool amIWhite;
        public string secret;
        public string gameID;
        void Awake()
        {
            if (instance != null && instance != this)
            {
                Destroy(this.gameObject);
                return;
            }
            else
            {
                instance = this;
            }
            DontDestroyOnLoad(this.gameObject);
        }
    }
}
