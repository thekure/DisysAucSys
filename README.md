# DisysAucSys

Mandatory Hand-In 5
== A Distributed Auction System ==

------- INSTRUCTIONS ------------

To run the program, open 3 separate terminals at the project directory.

If you have MAKE installed (otherwise, see below)

Run the following commands, one in each terminal:

    make server0

    make server1

    make server2

If you don't have MAKE:

Run the following commands, one in each terminal:

    go run server/server.go 0

    go run server/server.go 1

    go run server/server.go 2

It is now possible to connect to the three servers by running the following in a seperate terminal:

    go run client/client.go     
