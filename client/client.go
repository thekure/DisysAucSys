package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	auction "github.com/thekure/DisysAucSys/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name              string
	currentHighestBid int32
	servers           map[int32]auction.AuctionClient
}

func main() {

	fmt.Println("Please enter a name")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)

	// Set client name:
	tempName := ""
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			break
		} else {
			tempName += scanner.Text()
		}
	}

	// flag.Parse()
	setLog()

	client := &Client{
		name:              tempName,
		currentHighestBid: 0,
		servers:           make(map[int32]auction.AuctionClient),
	}

	go handleClient(client)

	for {

	}

}

func (client *Client) makeBid(amount int32) {
	bid := &auction.RequestBid{
		Name:    client.name,
		Message: client.name + " has made the following bid: " + strconv.Itoa(int(amount)),
		Amount:  amount,
	}

	resultList := make([]string, 0, 3)

	for port, server := range client.servers {
		ack, err := server.Bid(context.Background(), bid)
		if err != nil {
			delete(client.servers, port)
			log.Printf(client.name + "lost connection to a server, operating number of servers are now " + strconv.Itoa(int(len(resultList))))
			// fmt.Printf("something went wrong in bid method: %v", err)
		} else {
			resultList = append(resultList, ack.GetMessage())
		}

	}
	fmt.Printf(resultList[0] + "\n")
}

func (client *Client) requestHighestBid() {
	reqBid := &auction.HighestBidRequest{
		Message: client.name + " requested the current bid value.",
	}
	// outcome, err := client.connection.Result(context.Background(), reqBid)
	// if err != nil {
	// 	log.Printf(err.Error())
	// }

	resultList := make([]string, 0, 3)

	for port, server := range client.servers {
		outcome, err := server.Result(context.Background(), reqBid)
		if err != nil {

			delete(client.servers, port)
			log.Printf(client.name + "lost connection to a server, operating number of servers are now " + strconv.Itoa(int(len(resultList))))
			// fmt.Printf("something went wrong in requestHB method: %v", err)
		} else {
			resultList = append(resultList, outcome.GetStatus())
		}

	}
	fmt.Printf(resultList[0] + "\n")
}

// handles clientinput
func (client *Client) sendMessage() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		//Checks if the client input is of type string og integer (Tries to convert/parse, if an error occurs = string is not parseable)
		amount, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			//Tell client that the current bid is at ... and to make a bid, type an integer
			client.requestHighestBid()
		} else {
			client.makeBid(int32(amount))
		}
	}
}

func handleClient(client *Client) {
	client.getServerConnection()

	go client.sendMessage()

	for {

	}

}

// hardcoded method that connects to three different servers to ensure active replications
func (client *Client) getServerConnection() {

	for i := 0; i < 3; i++ {

		port := int32(5001) + int32(i)
		var conn *grpc.ClientConn

		fmt.Printf("Trying to dial: %v\n", port)
		insecure := insecure.NewCredentials()
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		fmt.Printf("--- "+client.name+" succesfully dialed to %v\n", port)
		log.Printf("--- "+client.name+" succesfully dialed to %v\n", port)

		// defer conn.Close()
		c := auction.NewAuctionClient(conn)
		client.servers[port] = c
	}

}

// Sets log output to file in project dir
func setLog() *os.File {
	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
