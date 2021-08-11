build:
	go build -o server.exe main.go

run: build
	./server

watch: build
# 	ulimit -n 1000 #increase the file watch limit, might required on MacOS
# 	reflex -s -r 'server.go' make run
# 	watcher -cmd="make run" -startcmd="./server.exe" main.go -list -keepalive
# 	watcher -cmd="make run" -startcmd="make run" -list -keepalive
# 	watcher -cmd="make run" -startcmd #-list -keepalive
	watcher -cmd="pwsh -command \"make run\"" -startcmd
# 	watcher -dotfiles=false -recursive=false -cmd="./server.exe" main.go ../

