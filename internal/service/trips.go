package service

import (
	"context"
	"database/sql"

	apiv1 "github.com/sharifahmad2061/trip-grpc-go/api/gen/go"
	queries "github.com/sharifahmad2061/trip-grpc-go/internal/db/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TripsServiceImpl struct {
	apiv1.UnimplementedTripsServer
	Db *sql.DB
}

func (s *TripsServiceImpl) CreateTrip(
	ctx context.Context,
	req *apiv1.CreateTripRequest,
) (*apiv1.CreateTripResponse, error) {
	dbHandle := queries.New(s.Db)

	trip, err := dbHandle.CreateTrip(ctx, queries.CreateTripParams{
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
	return &tripResponse, nil
}

func (s *TripsServiceImpl) GetTripById(
	ctx context.Context,
	req *apiv1.GetTripByIdRequest,
) (*apiv1.Trip, error) {
	dbHandle := queries.New(s.Db)

	trip, err := dbHandle.GetTripByID(ctx, int64(req.GetId()))
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
