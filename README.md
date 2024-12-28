Hello, This is go implementation of async server using event loop

Since the code uses linux kernel functions use Dockerfile to build and run else just the bash `runserver.sh`

To run code using Dockerfile use following steps
1. `docker build -t go-async-server .`
2. `docker run --rm -it -p 8080:8080 go-async-server`

Since this is a TCP server you can use tool `netcat` as follows
`echo -e "PING" | nc 127.0.0.1 8080`