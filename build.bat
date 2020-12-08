SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386
go build -ldflags "-s -w" -buildmode=c-shared -o OneBot-YaYa.XQ.dll
pause