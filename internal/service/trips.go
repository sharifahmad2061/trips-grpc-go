package service

import (
	"context"

	apiv1 "github.com/sharifahmad2061/trip-grpc-go/api/gen/go"
	queries "github.com/sharifahmad2061/trip-grpc-go/internal/db/generated"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TripsServiceImpl struct {
	apiv1.UnimplementedTripsServer
	Query *queries.Queries
}

func (s *TripsServiceImpl) CreateTrip(
	ctx context.Context,
	req *apiv1.CreateTripRequest,
) (*apiv1.CreateTripResponse, error) {
	trip, err := s.Query.CreateTrip(ctx, queries.CreateTripParams{
		Name:      req.GetName(),
		MemberID:  int64(req.GetMemberId()),
		StartDate: req.GetStartDate().AsTime(),
		EndDate:   req.GetEndDate().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	tripResponse := apiv1.CreateTripResponse{
		Id:       uint64(trip.ID),
		Name:     trip.Name,
		MemberId: uint64(trip.MemberID),
		StartDate: &timestamppb.Timestamp{
			Seconds: trip.StartDate.Unix(),
			Nanos:   int32(trip.StartDate.Nanosecond()),
		},
		EndDate: &timestamppb.Timestamp{
			Seconds: trip.EndDate.Unix(),
			Nanos:   int32(trip.EndDate.Nanosecond()),
		},
	}
	zap.L().Info("Created trip successfully", zap.Uint64("trip_id", tripResponse.Id))
	return &tripResponse, nil
}

func (s *TripsServiceImpl) GetTripById(
	ctx context.Context,
	req *apiv1.GetTripByIdRequest,
) (*apiv1.Trip, error) {
	trip, err := s.Query.GetTripByID(ctx, int64(req.GetId()))
	if err != nil {
		return nil, err
	}

	tripResponse := &apiv1.Trip{
		Id:       uint64(trip.ID),
		Name:     trip.Name,
		MemberId: uint64(trip.MemberID),
		StartDate: &timestamppb.Timestamp{
			Seconds: trip.StartDate.Unix(),
			Nanos:   int32(trip.StartDate.Nanosecond()),
		},
		EndDate: &timestamppb.Timestamp{
			Seconds: trip.EndDate.Unix(),
			Nanos:   int32(trip.EndDate.Nanosecond()),
		},
	}

	return tripResponse, nil
}
