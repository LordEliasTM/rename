rm -rf ./build/*
env GOOS=windows GOARCH=amd64 go build -ldflags "-s" -o build/rename-win64.exe
env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o build/rename-linux64
env GOOS=darwin GOARCH=amd64 go build -ldflags "-s" -o build/rename-darwin64
