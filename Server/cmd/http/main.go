package main

import (
	"fmt"
	"log"
	"net/http"

	inmemory "github.com/hashcacher/ChessGoNeue/Server/v2/inmemory"
)

func main() {

	s := inmemory.NewWebService()

	// TODO add http.servermux with metrics/logging middleware
	http.HandleFunc("/", webHandler)
	http.HandleFunc("/ding", dingHandler)
	http.HandleFunc("/v1/getUser", s.GetUser)
	http.HandleFunc("/v1/matchMe", s.MatchMe)
	// http.HandleFunc("/v1/move", s.moveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(poorMansHTML()))
	return
}

func poorMansHTML() string {
	baseURL := "https://storage.cloud.google.com/chessgo/"
	linuxURL := "master/linux_chessgo_master"
	osxURL := "master/osx_chessgo_master.tar.gz"
	webglURL := "master/webgl"

	return fmt.Sprintf(
		`<html lang="en-us">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>Unity WebGL Player | ChessGoNeue</title>
    <link rel="shortcut icon" href="%s/TemplateData/favicon.ico">
    <link rel="stylesheet" href="%s/TemplateData/style.css">
    <script src="%s/TemplateData/UnityProgress.js"></script>
    <script src="%s/Build/UnityLoader.js"></script>
    <script>
      var gameInstance = UnityLoader.instantiate("gameContainer", "%s/Build/webgl-build.json", {onProgress: UnityProgress});
    </script>
  </head>
  <body>
	Welcome to ChessGo. Here are our beta clients to try:<br>
	<a href="%s">ChessGo for Mac OSX</a><br>
	<a href="%s">ChessGo for Linux</a><br><br>

    <div class="webgl-content">
      Play right in the browser!
      <div id="gameContainer" style="width: 960px; height: 600px"></div>
      <div class="footer">
        <div class="fullscreen" onclick="gameInstance.SetFullscreen(1)"></div>
      </div>
    </div>

	Email chessgoinfo@gmail.com for more info.
  </body>
</html>
`,
		baseURL+webglURL, baseURL+webglURL, baseURL+webglURL, baseURL+webglURL, baseURL+webglURL,
		baseURL+osxURL, baseURL+linuxURL)
}

func dingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dong"))
}
