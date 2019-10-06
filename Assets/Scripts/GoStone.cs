using System.Collections;
using UnityEngine;

public class GoStone : MonoBehaviour {
    public float secondsTillOff = .25f;

    void Awake() {
    }

    void OnCollisionEnter(Collision _) {
        StartCoroutine(TurnOffRB());
    }

    IEnumerator TurnOffRB() {
        yield return new WaitForSeconds(secondsTillOff);
       Destroy(GetComponent<Rigidbody>()); 
    }
}
