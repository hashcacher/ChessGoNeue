using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using UnityEngine.UI;
using UnityEngine.SceneManagement;

namespace ChessGo
{
    public class Board : MonoBehaviour
    {
        public Camera camera;

        //GameObjects placed to mark the max/min X/Z coordinated of the grid.
        public GameObject TopLeft;
        public GameObject BottomRight;

        public GameObject table;
        public GameObject cameraLight;
        public GameObject blackStone;
        public GameObject whiteStone;
        public GameObject blackPawn;
        public GameObject whitePawn;
        public GameObject blackRook;
        public GameObject whiteRook;
        public GameObject blackKnight;
        public GameObject whiteKnight;
        public GameObject blackBishop;
        public GameObject whiteBishop;
        public GameObject blackQueen;
        public GameObject whiteQueen;
        public GameObject blackKing;
        public GameObject whiteKing;

        //speech bubbles
        public SpeechBubble chatBubble;

        //2D chess pieces
        public GameObject chess2D;

        //Highlights
        public GameObject validMove;
        public GameObject redSquareHighlight;
        public GameObject greenSquareHighlight;

        //Camera perspectives
        public GameObject cameraHotspots;
        public GameObject topDownSpot;

        //GUI
        public Button TopDownToggleButton;
        public Text myTurnText;
        public InputField inputChat;
        public Text chatBox;
        public Text playerNamesText;
        public GameObject winScreen;

        GameObject[,] pieces;
        GameObject[,] pieces2D;
        char[,] board;
        GameObject[,] highlights;

        GameObject mouseHighlight; //keeps track of which point or piece is being highlighted
        Point mouseHighlightPoint;
        Color mouseHighlightColor; //the color of the piece before it got highlighted

        MeshRenderer pieceHighlight;

        const int nRows = 8; //

        float maxX = 45.381f;
        float maxZ = 45.394f;
        float minX = -45.381f;
        float minZ = -45.394f;

        float w, h; //width and height of the board
        float smallW, smallH; // width and height of a square on the board

        bool myTurn = true;
        bool IAmBlack = true;

        //Dragging:
        public const int boardLayerBitmask = (1 << 8); // layer 8

        Transform grabbed;
        Vector3 grabbedPrevPos; //for putting it back if invalid drag
        Point grabbedInitialPoint;

        int grabLayerMask;
        Vector3 grabOffset; //delta between transform transform position and hit point

        //Camera rotation
        int curHotspot = 6;
        bool topdown = false;

        //if this is false, we're just debugging locally
        bool usingServer = false;
	bool gameOver = false;

        public Sprite whiteTurnImage, blackTurnImage;

        private Transform canvas;
        private RectTransform helpPanel, winPanel;
        private bool preventMoves;

        // Use this for initialization
        void Start()
        {
            usingServer = PlayerPrefs.GetInt("Hotseat") != 1 && 
                            AsyncServerConnection.clientSocket != null &&
                            AsyncServerConnection.clientSocket.Connected;

            maxX = TopLeft.transform.position.x;
            maxZ = TopLeft.transform.position.z;
            minX = BottomRight.transform.position.x;
            minZ = BottomRight.transform.position.z;

            //Find objects
            canvas = Camera.main.transform.Find("Canvas").transform;
            helpPanel = canvas.Find("Help Panel").GetComponent<RectTransform>();
            winPanel = canvas.Find("Win Panel").GetComponent<RectTransform>();
            helpPanel.gameObject.SetActive(false);
            winPanel.gameObject.SetActive(false);


            // Send a start game message to the server
            if (usingServer)
            {
                AsyncServerConnection.Send(Messages.STARTGAME);
                AsyncServerConnection.Receive();

                playerNamesText.text = Client.username + " vs ";
            }
            else
            {
                if (!IAmBlack)
                    StartCoroutine(Rotate180(camera.transform));
            }

            //set the width and height of the board
            w = getW();
            h = getH();
            smallW = w / (nRows - 1);
            smallH = h / (nRows - 1);

            //moar gravity
            Physics.gravity = new Vector3(0, -100, 0);

            //initialize stuff
            board = new char[nRows, nRows];
            pieces = new GameObject[nRows, nRows];
            pieces2D = new GameObject[nRows, nRows];
            highlights = new GameObject[nRows, nRows];

            Setup2D(); //shoots rays at the 2D pieces to get references.
            SetupPieces(pieces); //places the initial chess pieces

            StartTurn();
        }

	IEnumerator VictoryAnimation()
	{
	    gameOver = true;
	    // Spam Stones
	    for (int n = 1; n <= 4; n++) {
		for (int x = 0; x < 12; x++) {
		    for (int y = 0; y < 12; y++) {
			if(board[x,y] == '\0')	PlaceStone (new Point (x, y),IAmBlack ? 'S' : 's');
		    }
		    yield return null;
		}

	    }

	}

        // Spins an object (the camera) to the other side of the board
        IEnumerator<GameObject> Rotate180(Transform o)
        {
            float t = 0.0f;
            float speed = 125;
            while (t + speed * Time.deltaTime <= 180)
            {
                t += speed * Time.deltaTime;
                o.transform.RotateAround(table.transform.position, Vector3.up, speed * Time.deltaTime);

                //rotate the 2d pieces
                foreach (Transform piece in chess2D.transform)
                    piece.Rotate(new Vector3(0, 0, 1), speed * Time.deltaTime);

                yield return null;
            }
            float timeLeft = 180 - t;
            o.transform.RotateAround(table.transform.position, Vector3.up, timeLeft);

            //rotate the 2d pieces
            foreach (Transform piece in chess2D.transform)
                piece.Rotate(new Vector3(0, 0, 1), timeLeft);
        }

        public void QuitGame()
        {
            SceneManager.LoadScene("MainMenu");
        }

        public void SubmitChat()
        {
            if (inputChat.text.Trim() == "")
            {
                Debug.Log("empty string");
                return;
            }

            chatBox.text += "\n" + inputChat.text;

            //speech bubbles
            GameObject myKing = IAmBlack ? GameObject.Find("BlackKing(Clone)") : GameObject.Find("WhiteKing(Clone)");
            SpeechBubble b;
            if (myKing.transform.childCount == 0)
            {
                b = Instantiate(chatBubble) as SpeechBubble;
                b.transform.parent = myKing.transform;
                b.transform.position = myKing.transform.position;
                //b.transform.position = new Vector3()
            }
            else
            {
                b = myKing.transform.GetComponentInChildren<SpeechBubble>() as SpeechBubble;

            }

            if (b == null)
                Debug.LogError("Failed to get chat bubble");
            else
            {
                b.gameObject.SetActive(true);
                b.SetText(inputChat.text);
            }

            //send the chat message
            if (usingServer)
            {
                AsyncServerConnection.Send(Messages.CHAT, inputChat.text);
                AsyncServerConnection.Receive();
            }

            //reset the input box
            inputChat.text = "";
            inputChat.Select();

        }

        //shoots rays all the 2D chess models and 
        void Setup2D()
        {
            for(int x = 0; x < 8; x++)
            {
                for(int y = 0; y < 2; y++)
                {
                    //Black piece
                    Point p = new Point(x, y);
                    Vector3 target = GetWorldAtPoint(p, table);

                    Vector3 direction = target - camera.transform.position;
                    RaycastHit hit;
                    if (Physics.Raycast(camera.transform.position, direction, out hit, Mathf.Infinity, (1 << LayerMask.NameToLayer("Board"))))
                    {
                        pieces2D[x, y] = hit.transform.gameObject;
                        hit.transform.gameObject.SetActive(false);
                        if(hit.transform.name == "Table") {
                            Debug.Log(x + " " + y);
                        }
                    }
                    else
                        Debug.LogError("Missed the 2D piece: " + p);

                    //White piece
                    p = new Point(x, 7-y);
                    target = GetWorldAtPoint(p, table);

                    direction = target - camera.transform.position;
                    Physics.Raycast(camera.transform.position, direction, out hit, Mathf.Infinity, (1 << LayerMask.NameToLayer("Board")));

                    pieces2D[x, 7-y] = hit.transform.gameObject;
                    hit.transform.gameObject.SetActive(false);
                        if(hit.transform.name == "Table") {
                            Debug.Log(x + " " + (7-y));
                        }
                }
            }
        }

        public void TogglePerspective(Button b)
        {
            //if currently in topdown
            if (topdown)
            {
                Show3DPieces();
                Hide2DPieces();

                RotateCamera(GetCameraHotspot(IAmBlack ? "DefaultBlack" : "DefaultWhite"));
                curHotspot = IAmBlack ? 0 : 6; //update this if we change DefaultWhite's array location
                TopDownToggleButton.image.color = Color.white;
            }
            else
            {
                Show2DPieces();
                Hide3DPieces();

                RotateCamera(GetCameraHotspot(IAmBlack ? "Topdown" : "TopdownWhite"));
                curHotspot = IAmBlack ? 1 : 5;
                TopDownToggleButton.image.color = new Color(1, .67f,.67f);
            }
            topdown = !topdown;
        }

        void Hide2DPieces()
        {
            foreach (Transform o in chess2D.transform)
            {
                StartCoroutine(Utilities.FadeOut(o.gameObject, 1f)); //fade it out over 1 second
            }
        }

        void Show3DPieces()
        {
            foreach (GameObject o in pieces)
            {
                if(o != null)
                    StartCoroutine(Utilities.FadeIn(o, 1f)); //fade it out over 1 second
            }
        }

        void Show2DPieces()
        {
            chess2D.SetActive(true);
            foreach (Transform o in chess2D.transform)
            {
                StartCoroutine(Utilities.FadeIn(o.gameObject, 1f)); //fade it out over 1 second
            }
        }

        void Hide3DPieces()
        {
            foreach (GameObject o in pieces)
            {
                if (o != null && o.tag != "Stone")
                    StartCoroutine(Utilities.FadeOut(o, 1f)); //fade it out over 1 second
            }
        }

        //"Default", "Topdown", etc...
        Transform GetCameraHotspot(string s)
        {
            Transform t = cameraHotspots.transform.Find(s);
            if (t == null)
                Debug.LogError("No camera hotspot named " + s);
            return t;
        }

        void MouseOver()
        {
            if (Input.GetAxis("Mouse X") != 0 || Input.GetAxis("Mouse Y") != 0)
            {
                Point p = CurrentlyDraggingChessPiece() ? ClosestPoint(grabbed.transform.position) : GetSquareUnderMouse();
                //if(Input.GetKey(KeyCode.D))
                //    Debug.Log("Mouseing over square: " + p);

                if (p.col != -1)
                {
                    //if we're still mousing over the same square
                    if (p.Equals(mouseHighlightPoint))
                        return;
                    else
                    {
                        //reset the old highlight
                        ResetMouseOver();

                        mouseHighlightPoint = p;
                    }


                    if (CurrentlyDraggingChessPiece())
                    {
                        if (Utilities.IsValidMove(grabbedInitialPoint, p, board) && !Utilities.IsEmptyAt(p, board))
                        {
                            mouseHighlight = GetPieceAtPoint(p).gameObject;
                            Renderer rend = mouseHighlight.GetComponent<Renderer>();
                            mouseHighlightColor = rend.material.GetColor("_EmissionColor");
                            rend.material.SetColor("_EmissionColor", Color.red);
                        }
                    }
                    else
                    {
                        if (Utilities.IsEmptyAt(p, board))
                        {
                            mouseHighlight = CreateHighlight(validMove, p.row, p.col);
                        }
                        else
                        {
                            GameObject go = GetPieceAtPoint(p).gameObject;
                            if (IsMyPiece(go))
                            {
                                pieceHighlight = go.GetComponent<MeshRenderer>();
                                mouseHighlightColor = Color.black;
                                pieceHighlight.material.EnableKeyword("_EMISSION");
                                pieceHighlight.material.SetColor("_EmissionColor", new Color(.4f, 1f, .2f, 1f) * .5f);
                            }
                        }
                    }
                }
                else
                    ResetMouseOver();
            }
        }

        void ResetMouseOver()
        {
            if (mouseHighlight != null)
            {
                if (mouseHighlight.name == "ValidMove(Clone)") //the little red circle
                {
                    mouseHighlight.SetActive(false);
                }
                else {
                    mouseHighlight.GetComponent<Renderer>().material.SetColor("_EmissionColor", mouseHighlightColor);
                }
                mouseHighlight = null;
            }
            if (pieceHighlight != null)
            {
                //Debug.Log("emission was " + pieceHighlight.material.GetColor("_EmissionColor"));
                pieceHighlight.material.SetColor("_EmissionColor", mouseHighlightColor);

            }
        }

        bool IsChessPiece(Point p)
        {
            if (board[p.row, p.col] != '\0' && board[p.row, p.col] != 's' && board[p.row, p.col] != 'S')
                return true;
            else
                return false;
        }
        
        bool CurrentlyDraggingChessPiece()
        {
            return grabbed != null;
        }

        Point GetSquareUnderMouse()
        {
            Ray ray = Camera.main.ScreenPointToRay(Input.mousePosition);
            RaycastHit hit;
            if (Physics.Raycast(ray, out hit, Mathf.Infinity, boardLayerBitmask)) //ray only hits the board layer
            { //cast a ray and save hit to hit var
                if (ClickedWithinGrid(hit.point)) //inside the grid
                {
                    Point hitsquare = ClosestPoint(hit.point);
                    return hitsquare;
                }
                else return new Point(-1,-1);
            }
            else
                return new Point(-1,-1);
        }

        GameObject CreateHighlight(GameObject o, int row, int col)
        {
            if (highlights[row, col] == null || highlights[row,col].tag != o.tag)
                return highlights[row, col]
                 = Instantiate(o, GetWorldAtPoint(new Point(row, col), o), Quaternion.identity)
                   as GameObject;
            else
            {
                highlights[row, col].SetActive(true);
                return highlights[row, col];
            }
        }

        GameObject CreatePiece(GameObject o, int row, int col)
        {
            board[row, col] = Utilities.GetCharForPiece(o);
            pieces[row, col] = Instantiate(o, GetWorldAtPoint(new Point(row, col), o), Quaternion.identity) as GameObject;
            return pieces[row, col];
        }

        // Update is called once per frame
        void Update()
        {
            if (usingServer) {
		if (AsyncServerConnection.messageQueue.Count > 0)
		    OnReceiveServerMessage (AsyncServerConnection.messageQueue.Dequeue ());
	    } else {

	    }
            if (myTurn && !preventMoves)
            {
                if(!UpdateToggleDrag()) //if we didn't pick something up or drop something.
                    if(!grabbed && Input.GetMouseButtonDown(0))
                        StartCoroutine(PlaceGoStone());

                MouseOver(); //maybe call less frequently
            }
            if (Input.GetKeyDown("down")) { RotateCamera(); }
            if (Input.GetKeyDown("return")) { SubmitChat(); }
        }

        public void RotateCamera()
        {
            int nextHotspot = (IAmBlack ? (++curHotspot) : (--curHotspot)) % (cameraHotspots.transform.childCount - 1);
            Transform t = cameraHotspots.transform.GetChild(Mathf.Abs(nextHotspot));
            StartCoroutine(Utilities.SmoothMove(camera.transform, t, 1.0f));
        }

        void RotateCamera(Transform t)
        {
            StartCoroutine(Utilities.SmoothMove(camera.transform, t, 1.0f));
        }

        void OnReceiveServerMessage(Message msg) //changed from object[] params
        {
            int message = msg.message;
            string[] parameters = msg.parameters;

            Debug.Log("Message received in Board");

            switch (message)
            {
                case Messages.STARTGAME:
                    {
                        IAmBlack = AsyncServerConnection.FixParam(parameters[0]).Equals("1");
                        if (!IAmBlack)
                        {
                            curHotspot = 6;
                            myTurn = false;
                            myTurnText.enabled = false;
                            Debug.Log("I am White!");
                            StartCoroutine(Rotate180(camera.transform));
                            AsyncServerConnection.Receive();
                        }
                        else
                        {
                            curHotspot = 0;
                            StartTurn();

                        }
                        break;
                    }
                //Execute's opponent's chess move.
                case Messages.MOVE:
                    {
                        //escape the parameters
                        string[] fromMove = parameters[0].ToString().Split(',');

                        //get board coords
                        Point p1 = new Point(int.Parse(fromMove[0]), int.Parse(fromMove[1]));
                        //the move was a Go placement
                        if (parameters.Length == 1)
                        {
                            PlaceGoStone(p1);
                        }
                        //the move was a chess move
                        else
                        {
                            string[] toMove = parameters[1].ToString().Split(',');
                            Point p2 = new Point(int.Parse(toMove[0]), int.Parse(toMove[1]));

                            //checks if move is valid, if so it updates board[][]
                            if (MovePieceBoard(p1, p2))
                            {
                                Debug.Log("moving the actual table piece");
                                //move gameobject to p2

                                Transform piece = pieces[p1.row, p1.col].transform;
                                MovePieceTable(piece, p1, p2);

                                piece = pieces2D[p1.row, p1.col].transform;
                                MovePieceTable2D(piece, p1, p2);

                                CheckSurrounded(p2);
                            }
                            else
                            {
                                Debug.LogError("Server sent us an invalid move");
                            }
                        }
                        StartTurn();
                        break;
                    }
                case Messages.CHAT:
                    {
                        //escape the parameters
                        string chat = parameters[0].ToString();
                        ReceiveChat(chat);
                        break;
                    }
            }
        }

        private void StartTurn()
        {
            myTurn = true;
            myTurnText.enabled = true;
            if (!usingServer)
                myTurnText.text = IAmBlack ? "Black's Turn" : "White's Turn";
            StartCoroutine(SmoothTextOpen(myTurnText, 1f));
            FlipTurnButtonColor();
        }

        IEnumerator SmoothTextOpen(Text text, float seconds)
        {
            float t = 0.0f;
            float origWidth = text.rectTransform.sizeDelta.x;
            Vector2 textSize = new Vector2(0,text.rectTransform.sizeDelta.y);
            Vector2 futureTextSize = new Vector2(origWidth, text.rectTransform.sizeDelta.y);

            Color a = new Color(0, 1f, 0, 0);
            Color b = new Color(0, 1f, 0, 1f);
            while (t <= seconds)
            {
                t += Time.deltaTime / seconds;
                text.rectTransform.SetSizeWithCurrentAnchors(RectTransform.Axis.Horizontal,Mathf.Lerp(textSize.x, futureTextSize.x, Mathf.SmoothStep(0.0f, 1.0f, t)));
                text.color = Color.Lerp(a, b, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }
            t = 0f;
            while (t <= seconds)
            {
                t += Time.deltaTime / seconds/2;
                text.color = Color.Lerp(b, a, Mathf.SmoothStep(0.0f, 1.0f, t));
                yield return null;
            }


        }

        private void ReceiveChat(string msg)
        {
            chatBox.text += "\n" + msg;
        }


        // Toggles drag with mouse click
        //returns true if we picked something up or put something down.
        bool UpdateToggleDrag()
        {
            if (Input.GetMouseButtonDown(0))
            {
                return Grab();
            }
            else            {
                if (grabbed)
                {
                    Drag();
                }

                return false;
            }
        }

        void HideHighlights()
        {
            foreach (GameObject o in highlights)
                if (o != null)
                    o.SetActive(false);
        }
        void ResetGrab()
        {
            if (grabbed)
            {
                //restore the original layermask
                //grabbed.gameObject.layer = grabLayerMask;
                HideHighlights();
            }
            grabbed = null;
        }

        void ReleaseSuccessfulGrab()
        {
            if (grabbed)
            {
                //restore the original layermask
                grabbed.gameObject.layer = grabLayerMask;
                HideHighlights();

                Point prevPoint = ClosestPoint(grabbedPrevPos);//change this back in Grab()

                Point p = ClosestPoint(grabbed.transform.position);

                
                if (MovePieceBoard(prevPoint, p))
                {
                    MovePieceTable2D(pieces2D[grabbedInitialPoint.row,grabbedInitialPoint.col].transform, prevPoint, p);
                    MovePieceTable(grabbed, prevPoint, p);
                    CheckSurrounded(p);

                    // TODO SendMoveToServer(prevPoint, p);

                    if (usingServer)
                        EndTurn();
                    else
                        EndTurnHotseat();
                }
                else //put it back
                    StartCoroutine(Utilities.SmoothMove(grabbed, grabbedPrevPos, .2f));
            }
            grabbed = null;
        }

        void SendGoMoveToServer(Point p1)
        {
            if (!usingServer)
                return;

            AsyncServerConnection.Send(Messages.MOVE, p1.ToString());
            AsyncServerConnection.Receive();
        }
        void SendMoveToServer(Point p1, Point p2)
        {
            if (!usingServer)
                return;

            Debug.Log("Sending move " + p1 + " to " + p2);
            AsyncServerConnection.Send(Messages.MOVE, p1.ToString(), p2.ToString());
            AsyncServerConnection.Receive();
        }

        bool Grab()
        {
            if (grabbed)
            {
                //check if the release was on the board
                RaycastHit hit;
                Ray ray = Camera.main.ScreenPointToRay(Input.mousePosition);
                if (Physics.Raycast(ray, out hit))
                {
                    if (ClickedWithinGrid(hit.point))
                    {
                        ReleaseSuccessfulGrab();
                        return true;
                    }
                    else
                        return false;
                }
                else
                    return false;
            }
            else
            {
                Point p = GetSquareUnderMouse();
                if (p.col == -1 || Utilities.IsEmptyAt(p, board))
                {
                    return false;
                }
                else
                {
                    grabbed = GetPieceAtPoint(p);
                    grabbedPrevPos = CloneVector3(grabbed.position);

                    if (grabbed.parent)
                    {
                        //grabbed = grabbed.parent.transform;
                    }

                    //can't grab the board or wrong color
                    if (!IsMyPiece(grabbed.gameObject) ||
                       grabbed.gameObject.tag.Equals("Stone"))
                    {
                        Debug.Log("Reset uh oh");

                        ResetGrab();
                        return false;
                    }
                    //set the object to ignore raycasts
                    grabLayerMask = grabbed.gameObject.layer;
                    grabbed.gameObject.layer = 2;
                    //now immediately do another raycast to calculate the offset
                    RaycastHit hit;
                    Ray ray = Camera.main.ScreenPointToRay(Input.mousePosition);
                    if (Physics.Raycast(ray, out hit, 1 << LayerMask.NameToLayer("Board")))
                    {
                        grabOffset = grabbed.position - hit.point;
                        grabOffset.y += hit.point.y;

                        grabbedInitialPoint = ClosestPoint(grabbed.transform.position);
                        CreateHighlight(greenSquareHighlight, grabbedInitialPoint.row, grabbedInitialPoint.col);
                        List<Point> pointsToHighlight = Utilities.GetValidDestinations(grabbedInitialPoint, board);
                        HighlightPoints(pointsToHighlight);
                        return true;
                    }
                    else
                    {
                        Debug.Log("Uh oh2");
                        return false;
                    }
                }
            }
        }

        bool IsMyPiece(GameObject go)
        {
            return go.tag.Equals(IAmBlack ? "Black" : "White");
        }

        void HighlightPoints(List<Point> points)
        {
            foreach (Point p in points)
            {
                CreateHighlight(validMove, p.row, p.col);
            }
        }


        void Drag()
        {
            RaycastHit hit;
            Ray ray = Camera.main.ScreenPointToRay(Input.mousePosition);
            if (Physics.Raycast(ray, out hit, Mathf.Infinity, boardLayerBitmask))
            {
                grabbed.position = hit.point + grabOffset;
            }
        }

        //returns [row, column] of the point closest to the one clicked
        public Point ClosestPoint(Vector3 p)
        {
            if (!ClickedWithinGrid(p))
                return new Point(-1, -1);

            int x = 0;
            int z = 0;
            for (float i = minX; i <= maxX; i += smallW) //can do this is constant time with mod?
            {
                if (Mathf.Abs(p.x - i) < smallW / 2)
                    break;
                x++;
            }
            for (float i = minZ; i <= maxZ; i += smallH)
            {
                if (Mathf.Abs(p.z - i) < smallH / 2)
                    break;
                z++;
            }
            return new Point(x, z);
        }

        float getW()
        {
            return maxX - minX;
        }
        float getH()
        {
            return maxZ - minZ;
        }

        bool ClickedWithinGrid(Vector3 p)
        {
            if (p.x > minX - smallW / 2 && p.x < maxX + smallW / 2
                && p.z > minX - smallH / 2 && p.z < maxZ + smallH / 2)
                return true;
            return false;
        }

        Transform GetPieceAtPoint(Point p)
        {
            return pieces[p.row, p.col].transform;
        }

        //A version that keep's o's y position.
        Vector3 GetWorldAtPoint(Point p, GameObject o)
        {
            Vector3 v = new Vector3(minX + smallW * (p.row), o.transform.position.y, minZ + smallH * (p.col));
            return v;
        }


        Vector3 GetWorldAtPoint(Point p)
        {
            return new Vector3(minX + smallW * (p.row), 35f, minZ + smallH * (p.col));
        }

        //when the server sends us opponent's go move
        void PlaceGoStone(Point p)
        {
            //if they clicked in an empty spot
            if (GetBoardPiece(p) == '\0')
            {
                //place either white or black Stone.
                PlaceStone(p, IAmBlack ? 's' : 'S');
            }
            else
                Debug.LogError("Server told me to place a Go stone where the was already a piece");
        }

        //when a player clicks an empty place
        IEnumerator PlaceGoStone()
        {
            Ray ray = Camera.main.ScreenPointToRay(Input.mousePosition);
            RaycastHit hit;
            if (Physics.Raycast(ray, out hit))
            { //cast a ray and save hit to hit var
                if (hit.transform.position.y == 0) //if we hit the board
                {
                    if (ClickedWithinGrid(hit.point)) //inside the grid
                    {
                        //get the row and column they clicked
                        Point p = ClosestPoint(hit.point);

                        //if they clicked in an empty spot
                        if (GetBoardPiece(p) == '\0')
                        {
                            //place either white or black Stone.
                            PlaceStone(p, IAmBlack ? 'S' : 's');
                            // TODO SendGoMoveToServer(p);

                            if (usingServer)
                                EndTurn();
                            else
                            {
                                preventMoves = true;
                                yield return new WaitForSeconds(.8f); //wait for the stone to drop
                                EndTurnHotseat();
                                preventMoves = false;
                            }
                        }
                    }
                }
            }
        }



        GameObject GetTablePiece(Point p)
        {
            return pieces[p.row, p.col];
        }

        char GetBoardPiece(Point p)
        {
            return board[p.row, p.col];
        }

        void PlaceStone(Point p, char piece)
        {
            GameObject obj = Instantiate(piece == 'S' ? blackStone : whiteStone,
                              GetWorldAtPoint(p), Quaternion.identity)
                              as GameObject;

            board[p.row, p.col] = piece;
            if (pieces[p.row, p.col] != null)
            {
                Debug.LogError("trying to put a Go stone in non-empty square");
            }
            pieces[p.row, p.col] = obj;
        }

        void EndTurn()
        {
            myTurnText.enabled = false;
            myTurn = false;
            FlipTurnButtonColor();
        }

        private void FlipTurnButtonColor()
        {
            TopDownToggleButton.image.sprite = (myTurn && IAmBlack || !myTurn && !IAmBlack) ? blackTurnImage : whiteTurnImage;
        }

        void EndTurnHotseat()
        {
			if (!gameOver) {
				IAmBlack = !IAmBlack;
				StartTurn ();
				StartCoroutine (Rotate180 (camera.transform));
            }
        }

        bool MovePieceTable(Transform o, Point p1, Point p2)
        {
            Transform chessPiece = GetPieceAtPoint(p1);
            Vector3 targetLocation = GetWorldAtPoint(p2);
            targetLocation.y = chessPiece.GetComponent<Collider>().bounds.extents.y;

            pieces[p2.row, p2.col] = pieces[p1.row, p1.col];
            pieces[p1.row, p1.col] = null;

            StartCoroutine(Utilities.SmoothMove(chessPiece, targetLocation, .2f));

            return true;
        }

        bool MovePieceTable2D(Transform o, Point p1, Point p2)
        {
            Transform chessPiece = pieces2D[p1.row, p1.col].transform;
            Vector3 targetLocation = GetWorldAtPoint(p2);
            targetLocation.y = o.transform.position.y;

            pieces2D[p2.row, p2.col] = pieces2D[p1.row, p1.col];
            pieces2D[p1.row, p1.col] = null;

            StartCoroutine(Utilities.SmoothMove(chessPiece, targetLocation, .2f));

            return true;
        }


        bool MovePieceBoard(Point p1, Point p2)
        {
            //TODO
            //Initial placement of pieces
            if (board[p1.row, p1.col] == '\0')
            {
                Debug.Log("no piece here");
                return false;
            }
            else if (Utilities.IsValidMove(p1, p2, board))
            {
                Destroy(pieces[p2.row, p2.col]);
                Destroy(pieces2D[p2.row, p2.col]);
               	
				// Was a king destroyed?
				if(board[p2.row, p2.col] == 'k' || board[p2.row, p2.col] == 'K') {
					Debug.Log ("The Old King is dead, long live the King!");
					//Remove from board and in-game piece as well.
					Destroy(pieces[p2.row, p2.col]);
					Destroy(pieces2D[p2.row, p2.col]);
					board[p2.row, p2.col] = '\0';
					StartCoroutine(VictoryAnimation());

                    winPanel.gameObject.SetActive(true);
                    if (board[p2.row, p2.col] == 'k' && IAmBlack || board[p2.row, p2.col] == 'K' && !IAmBlack)
                    {
                        AsyncServerConnection.Send(Messages.WIN);
                    }
                    else
                    {
                        AsyncServerConnection.Send(Messages.LOSE);
                    }
				}
                board[p2.row, p2.col] = board[p1.row, p1.col];
                board[p1.row, p1.col] = '\0';

                return true;
            }
            else
            {
                return false;
            }
        }


        void KillAtPoint(Point p)
		{
			//Kills piece at point p
			
			if(board[p.row, p.col] == 'k' || board[p.row, p.col] == 'K') {
				Debug.Log ("The Old King is dead, long live the King!");
				//Remove from board and in-game piece as well.
				Destroy(pieces[p.row, p.col]);
				Destroy(pieces2D[p.row, p.col]);
				board[p.row, p.col] = '\0';
				StartCoroutine(VictoryAnimation());
			}
			else if (board [p.row, p.col] == '\0') {
				Debug.Log ("Attempted to Kill Empty Position");
			} else {
				//Remove from board and in-game piece as well.
                Destroy(pieces[p.row, p.col]);
                Destroy(pieces2D[p.row, p.col]);
				board[p.row, p.col] = '\0';
			}
		}

		//Checks the adjacent points, removing their group from game if the group is surrounded
		public void CheckSurrounded(Point p)
		{
			if (!gameOver) {
				HashSet<Point> affected = Utilities.GetAdjacentPoints (p);
				foreach (Point a in affected) {
					//See if the group has now been surrounded
					if (Utilities.IsGroupDead (a, board)) {
						//Kill the group
						HashSet<Point> targets = Utilities.GetGroup (a, board);
						foreach (Point t in targets) {
							KillAtPoint (t);
						}
					}
				}
			}
		}

        Vector3 CloneVector3(Vector3 v)
        {
            return new Vector3(v.x, v.y, v.z);
        }

        //pieces[x,y] 0 <= x <= 12 going from left to right
        //            0 <= y <= 12 going from top down
        void SetupPieces(GameObject[,] pieces)
        {
            //pawns
            CreatePiece(whitePawn, 0, 6);
            CreatePiece(whitePawn, 1, 6);
            CreatePiece(whitePawn, 2, 6);
            CreatePiece(whitePawn, 3, 6);
            CreatePiece(whitePawn, 4, 6);
            CreatePiece(whitePawn, 5, 6);
            CreatePiece(whitePawn, 6, 6);
            CreatePiece(whitePawn, 7, 6);

            CreatePiece(blackPawn, 0, 1);
            CreatePiece(blackPawn, 1, 1);
            CreatePiece(blackPawn, 2, 1);
            CreatePiece(blackPawn, 3, 1);
            CreatePiece(blackPawn, 4, 1);
            CreatePiece(blackPawn, 5, 1);
            CreatePiece(blackPawn, 6, 1);
            CreatePiece(blackPawn, 7, 1);

            //rooks
            CreatePiece(whiteRook, 0, 7);
            CreatePiece(whiteRook, 7, 7);

            CreatePiece(blackRook, 0, 0);
            CreatePiece(blackRook, 7, 0);

            //knights
            CreatePiece(whiteKnight, 1, 7).transform.Rotate(0, 180, 0);
            CreatePiece(whiteKnight, 6, 7).transform.Rotate(0, 180, 0);

            CreatePiece(blackKnight, 1, 0);
            CreatePiece(blackKnight, 6, 0);

            //bishops
            CreatePiece(whiteBishop, 2, 7);
            CreatePiece(whiteBishop, 5, 7);

            CreatePiece(blackBishop, 2, 0);
            CreatePiece(blackBishop, 5, 0);

            //Queen and King
            CreatePiece(whiteQueen, 3, 7);
            CreatePiece(whiteKing, 4, 7);

            CreatePiece(blackQueen, 3, 0);
            CreatePiece(blackKing, 4, 0);
        }

        public void ShowHelp()
        {
            helpPanel.gameObject.SetActive(true);
            //Vector2 startPos = new Vector2(0, -300f);
            //Vector2 endPos = new Vector2(0, 632f);
            Vector3 startPos = new Vector3(0, 500, -20);
            Vector3 endPos = new Vector3(0, -230, -20);
            preventMoves = true;
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

    }
}
