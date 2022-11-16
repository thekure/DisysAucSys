package main

import (
	auction "github.com/thekure/DisysAucSys/grpc"
	"golang.org/x/net/context"
)

type Server struct {
	auction.UnimplementedAuctionServer
	currentHighestBid int32
}

func (s *Server) Bid(
	ctx context.Context,
	RequestBid *auction.RequestBid,
) (*auction.Ack, error) {

	if s.CheckIfBidIsHighest(RequestBid.Amount) {
		return &auction.Ack{
			Message: "Success",
		}, nil
	} else {
		return &auction.Ack{
			Message: "Failed",
		}, nil
	}
}

func (s *Server) CheckIfBidIsHighest(bid int32) bool {
	if bid > s.currentHighestBid {
		return true
	} else {
		return false
	}
}
