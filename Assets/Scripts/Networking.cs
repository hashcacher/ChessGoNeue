using System.Text;
using UnityEngine;
using UnityEngine.Networking;


namespace ChessGo
{
    //Static methods for computing Go surrounds, valid moves, etc.
    public static class Net
    {
        public static UnityWebRequest GoodPost(string url, string body) {
           byte[] bytes = Encoding.ASCII.GetBytes(body);
           var request             = new UnityWebRequest(url);
           request.uploadHandler   = new UploadHandlerRaw(bytes);
           request.downloadHandler = new DownloadHandlerBuffer();
           request.method          = UnityWebRequest.kHttpVerbPOST;
           request.timeout = 120;
           return request;
        }


        public static string GetServerHost() {
            if (Application.isEditor) {
                return "localhost:8080";
            }

            return "https://chessgo.xyz";
        }
    }
}
