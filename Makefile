test:
	fg go run cmd/slotframe-manager/main.go
	echo "Starting root node"
	fg go run cmd/testclient/main.go run --parentId 0 --id 1 --etx 1
	echo "Starting child node of root named 2"
	fg go run cmd/testclient/main.go run --parentId 1 --id 2 --etx 1
	echo "Starting child node of root named 3"
	fg go run cmd/testclient/main.go run --parentId 1 --id 3 --etx 1
	echo "Starting child node of 2 named 4"
	fg go run cmd/testclient/main.go run --parentId 2 --id 4 --etx 1
