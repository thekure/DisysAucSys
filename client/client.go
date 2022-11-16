package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	auction "github.com/thekure/DisysAucSys/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name              string
	currentHighestBid int32
	connection        *grpc.ClientConn
}

func main() {

	fmt.Println("Please enter a name")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)

	tempName := ""
	for scanner.Scan() {
		tempName += scanner.Text()
	}

	flag.Parse()

	client := &Client{
		name:              tempName,
		currentHighestBid: 0,
	}

	go handleClient(client)

	for {

	}
}

func handleClient(client *Client) {

}

func (client *Client) getServerConnection() (auction.AuctionClient, error) {

	//dial options
	//the server is not using TLS, so we use insecure credentials
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("192.168.0.153:5400", opts...) //smusHjemme

	if err != nil {
		log.Fatalln("Could not dial")
	}
	client.connection = conn
	log.Printf("--- Client is connected to the server ---")

	return auction.NewAuctionClient(conn), err
}

// func (client *Client) Bid(ctx context.Context) (*auction.Ack, error) {

// }
