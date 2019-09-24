from  gableroux/unity3d:2018.4.9f1-mac

RUN apt update && apt install makeself

COPY build/license.ulf /root/.local/share/unity3d/Unity/Unity_lic.ulf

WORKDIR /app
COPY . .

# Build linux binary
RUN mkdir -p /app/build/artifacts && xvfb-run --auto-servernum --server-args='-screen 0 640x480x24' \
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

# folder to zip up
RUN makeself /app/build/artifacts/linux \
  # name to give the self-extracting archive
  linux_chessgo_master \
  ChessGo \
  # Executable to launch
  ./linux-chessgo-master \
    # Replace the separate stuff with the combined archive
    && rm -rf /app/build/artifacts/linux/* \
    && mv linux_chessgo_master /app/build/artifacts/linux


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
  -logFile /dev/stdout \

  && tar cvzf /app/build/artifacts/osx/osx_chessgo_master.tar.gz \
              /app/build/artifacts/osx/osx-chessgo-master.app/ --remove-files 
