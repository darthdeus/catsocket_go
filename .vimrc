nmap ,r :call VimuxRunCommand("go run main.go")<cr>
nmap ,t :call VimuxRunCommand("go test")<cr>
nmap ,k :!killall a.out; killall go<cr><cr>

