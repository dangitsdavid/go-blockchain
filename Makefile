# default block data
BLOCK_DATA="hello world!"

add-block:
	go run main.go add -block "$(BLOCK_DATA)"

print:
	go run main.go print
