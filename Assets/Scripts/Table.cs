using UnityEngine;
using System.Collections;

namespace ChessGo
{
    public class Table : MonoBehaviour
    {
        public Board gameboard;

        // Use this for initialization
        void Start()
        {

        }

        // Update is called once per frame
        void Update()
        {

        }

        //remove surround when the go stone falls to the ground.
        void OnCollisionEnter(Collision collision)
        {
            Debug.Log("collision table!");
            foreach (ContactPoint contact in collision.contacts)
            {
                if (contact.otherCollider.tag.Equals("Stone"))
                    gameboard.CheckSurrounded(gameboard.ClosestPoint(contact.point));
            }
        }
    }
}