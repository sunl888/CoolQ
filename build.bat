SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386
go build -ldflags "-s -w" -buildmode=c-shared -o app.dll

cqcfg.exe ./

SET DIR=%cd%
MOVE %DIR%\app.dll D:\cq\dev\com.ypdan.ypdan
MOVE %DIR%\app.json D:\cq\dev\com.ypdan.ypdan
COPY %DIR%\config\*.json D:\cq\data\app\com.ypdan.ypdan