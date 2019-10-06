using System.Collections.Generic;
using UnityEngine;

namespace ChessGo
{
    public static class Util {
        public static IEnumerator<GameObject> FadeOut(GameObject o, float duration)
        {
            float t = 0f;

            Renderer r = o.GetComponent<Renderer>();
            Color newColor = r.material.color;
            newColor.a = 0;

            while (t < 1)
            {
                r.material.color = Color.Lerp(r.material.color, newColor, t);
                t += Time.deltaTime / duration;
                yield return null;
            }
            o.SetActive(false);
        }

        public static IEnumerator<GameObject> FadeIn(GameObject o, float duration)
        {
            float t = 0f;

            Renderer r = o.GetComponent<Renderer>();
            Color newColor = r.material.color;
            newColor.a = 1.0f;

            o.SetActive(true);
            while (t < 1)
            {
                r.material.color = Color.Lerp(r.material.color, newColor, t);
                t += Time.deltaTime / duration;
                yield return null;
            }

        }

        // Moves a on object o smoothly
        public static IEnumerator<GameObject> SmoothMove(Transform o, Transform end, float seconds)
        {
            float t = 0.0f;
            Vector3 startpos = o.transform.position;
            Quaternion startrot = o.transform.rotation;
            while (t <= 1.0f)
            {
                t += Time.deltaTime / seconds;
                o.position = Vector3.Lerp(startpos, end.position, Mathf.SmoothStep(0.0f, 1.0f, t));
                o.rotation = Quaternion.Lerp(startrot, end.rotation, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }
        }

        // Moves a on object o smoothly
        public static IEnumerator<GameObject> SmoothMove(Transform o, Vector3 endpos, float seconds)
        {
            float t = 0.0f;
            Vector3 startpos = o.transform.position;
            while (t <= 1.0f)
            {
                t += Time.deltaTime / seconds;
                o.position = Vector3.Lerp(startpos, endpos, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }
        }

        internal static IEnumerator<GameObject> SmoothMoveUI(RectTransform rt, Vector3 startPos, Vector3 endPos, float time)
        {
            float elapsed = 0;
            while (elapsed < time)
            {
                elapsed += Time.deltaTime;
                rt.anchoredPosition3D = Vector3.Lerp(startPos, endPos, Mathf.SmoothStep(0, 1, elapsed/time));
                yield return null;
            }
        }
    }
}
