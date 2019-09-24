from  gableroux/unity3d:2018.4.9f1-mac

RUN apt update && apt install makeself

COPY build/license.ulf /root/.local/share/unity3d/Unity/Unity_lic.ulf

WORKDIR /app
COPY . .

# Build linux binary
RUN mkdir -p build/artifacts && xvfb-run --auto-servernum --server-args='-screen 0 640x480x24' \
  /opt/Unity/Editor/Unity \
  -batchmode \
  -quit \
  -projectPath . \
  -buildTarget StandaloneLinux64 \
  -customBuildTarget StandaloneLinux64 \
  -customBuildName linux-chessgo-master \
  -customBuildPath build/artifacts/linux/ \
  -executeMethod BuildCommand.PerformBuild \
  -logFile /dev/stdout

RUN makeself build/artifacts/linux \ # folder to zip up
  linux_chessgo_master \ # name to give the self-extracting archive
  ChessGo \ # Another name
  ./linux-chessgo-master # Executable to launch

# Build OSX binary
RUN xvfb-run --auto-servernum --server-args='-screen 0 640x480x24' \
  /opt/Unity/Editor/Unity \
  -batchmode \
  -quit \
  -projectPath . \
  -buildTarget StandaloneOSX \
  -customBuildTarget StandaloneOSX \
  -customBuildName osx-chessgo-master \
  -customBuildPath build/artifacts/osx/ \
  -executeMethod BuildCommand.PerformBuild \
  -logFile /dev/stdout


