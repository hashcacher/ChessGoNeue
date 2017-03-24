using UnityEngine;
using System.Collections;

namespace ChessGo
{
    public class UnitySingleton : MonoBehaviour
    {
        private static UnitySingleton instance = null;
        public static UnitySingleton Instance
        {
            get { return instance; }
        }
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

        public void OnApplicationQuit()
        {
            Debug.Log("Closing Socket");
            // Release the socket.
            AsyncServerConnection.Shutdown();
        }

        // any other methods you need
        //http://answers.unity3d.com/questions/11314/audio-or-music-to-continue-playing-between-scene-c.html
    }
}