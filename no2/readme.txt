Step:
-Install docker
-docker ps
-go mod tidy
-docker run -d -p 6379:6379 redis
-go run login.go