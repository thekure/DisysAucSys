package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	auction "github.com/thekure/DisysAucSys/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	auction.UnimplementedAuctionServer
	currentHighestBid       int32
	port                    int
	currentHighestBidholder string
	remainingTime           int32
	isAuctionRunning        bool
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5001

	setLog()
	flag.Parse()

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := &Server{
		currentHighestBid:       0,
		port:                    int(ownPort),
		currentHighestBidholder: "the initial value",
		isAuctionRunning:        false,
	}

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, server)

	go server.handleTime()

	for {

		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}
}
func (s *Server) handleTime() {

	for {
		s.StartAuctionTimer()
		s.StopAuctionAndAnnounceWinner()
	}
}

func (s *Server) StartAuctionTimer() {

	DurationOfTime := time.Duration(10) * time.Second
	startAuction := func() {
		s.resetAuction()
		s.isAuctionRunning = true
		fmt.Printf("\nThe Auction has now begun and will run for -15- seconds")
		log.Printf("\nThe Auction has now begun and will run for -15- seconds")
	}

	fmt.Println("waiting for auction to begin...")
	Timer1 := time.AfterFunc(DurationOfTime, startAuction)
	time.Sleep(time.Second * 15)
	defer Timer1.Stop()

}

func (s *Server) StopAuctionAndAnnounceWinner() {

	DurationOfTime := time.Duration(10) * time.Second

	stopAuction := func() {
		s.isAuctionRunning = false
		fmt.Printf("\nThe auction is now over with the winner being %v. \n Next Auction will begin in -10- seconds", s.currentHighestBidholder)
		log.Printf("\nThe auction is now over with the winner being %v. \n Next Auction will begin in -10- seconds", s.currentHighestBidholder)
	}

	Timer1 := time.AfterFunc(DurationOfTime, stopAuction)
	time.Sleep(time.Second * 10)
	defer Timer1.Stop()

}

func (s *Server) resetAuction() {
	s.currentHighestBidholder = "the initial value"
	s.currentHighestBid = 0
}

func (s *Server) Bid(ctx context.Context, RequestBid *auction.RequestBid) (*auction.Ack, error) {

	log.Printf(RequestBid.Message)

	incomingBid := strconv.Itoa(int(RequestBid.Amount))
	previousBid := strconv.Itoa(int(s.currentHighestBid))

	if !s.isAuctionRunning {
		return &auction.Ack{
			Message: "The auction is currently over and setting up for the next auction. The last auction was won with by " + s.currentHighestBidholder + "a bid of: " + previousBid,
			Amount:  s.currentHighestBid,
		}, nil
	}

	if s.CheckIfBidIsHighest(RequestBid.Amount) {

		previousHolder := s.currentHighestBidholder

		// sets the auction leader stats:
		s.currentHighestBid = RequestBid.Amount
		s.currentHighestBidholder = RequestBid.Name

		return &auction.Ack{
			Message: RequestBid.Name + " has beaten previous bid " + previousBid + " from " + previousHolder + ".\n Current highest bid is now: " + incomingBid,
			Amount:  s.currentHighestBid,
		}, nil
	} else {
		return &auction.Ack{
			Message: RequestBid.Name + " did not succeed in outbidding on the value of " + s.currentHighestBidholder + " with their bid of " + previousBid,
			Amount:  s.currentHighestBid,
		}, nil
	}
}

func (s *Server) Result(ctx context.Context, RequestBid *auction.HighestBidRequest) (*auction.Outcome, error) {
	log.Printf(RequestBid.Message)
	currentHighestBid := strconv.Itoa(int(s.currentHighestBid))

	if !s.isAuctionRunning {
		return &auction.Outcome{
			Amount: s.currentHighestBid,
			Status: "The auction is over with the winner being: " + s.currentHighestBidholder + " with the highest bid of: " + currentHighestBid,
		}, nil
	}

	return &auction.Outcome{
		Amount: s.currentHighestBid,
		Status: "auction is running where " + s.currentHighestBidholder + " has the highest bid at: " + currentHighestBid,
	}, nil
}

func (s *Server) CheckIfBidIsHighest(bid int32) bool {
	return (bid > s.currentHighestBid)
	// if bid > s.currentHighestBid {
	//  return true
	// } else {
	//  return false
	// }
}

// Sets log output to file in project dir
func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
