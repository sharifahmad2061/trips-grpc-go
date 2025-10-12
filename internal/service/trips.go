package service

import (
	"context"
	"time"

	apiv1 "github.com/sharifahmad2061/trip-grpc-go/api/gen/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TripsServiceImpl struct {
	apiv1.UnimplementedTripsServer
}

func (s *TripsServiceImpl) CreateTrip(
	ctx context.Context,
	req *apiv1.CreateTripRequest,
) (*apiv1.CreateTripResponse, error) {
	return &apiv1.CreateTripResponse{
		Id:        1,
		Name:      req.GetName(),
		MemberId:  req.GetMemberId(),
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (s *TripsServiceImpl) GetTripById(
	ctx context.Context,
	req *apiv1.GetTripByIdRequest,
) (*apiv1.Trip, error) {
	return &apiv1.Trip{
		Id:       req.GetId(),
		Name:     "Sample Trip",
		MemberId: 12345,
		StartDate: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   0,
		},
		EndDate: &timestamppb.Timestamp{
			Seconds: time.Now().Add(24 * time.Hour).Unix(), // Example timestamp
			Nanos:   0,
		},
	}, nil
}
