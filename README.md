Redis-db is an implementation of redis server. It includes support for the basic commands along with support for persistance.

## Commands Supported
- PING
- SET
- GET
- HSET
- HGET
- HGETALL

## Running the Server
- Try installing redis on your machine if not already. We'll need the redis-cli to send requests to our server.
- Then when `redis` is installed, try opening up the redis-cli using `redis-cli` command. Do a `PING` to see if the server responds with `PONG`.
- Do the following steps to run the server:
    - Before running the server, try `go mod tidy` to tidy up the imports.
    - Then run the server using the `go run .` command.
    - If the port is already being used, run `sudo systemctl stop redis`. This will stop the redis server which started when you installed redis.

- If the connection in the `redis-cli` drops, try running the cli again. 
- After connection, you'd see a log in the server terminal indicating that the server got a connection.
- Now, try to run the commands as you'd run for the native Redis server.
