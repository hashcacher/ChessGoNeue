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
        public static bool amIWhite;
        public static string secret;
        public static int gameID;
        public static bool hotseat = true;
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
