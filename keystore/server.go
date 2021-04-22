package keystore

import (
	"context"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	pb "github.com/psychobummer/pbrelay/rpc/keystore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KeymasterServer contains the implementation for the keymaster RPCs
type Server struct {
	store Backend
}

// NewServer returns a new server with memory-backed storage
func NewServer() *Server {
	return &Server{
		store: NewMemoryStorage(),
	}
}

// CreateKey accepts a key, stores it at a randomly-generated id, and returns the id to the caller
func (s *Server) CreateKey(ctx context.Context, request *pb.CreateKeyRequest) (*pb.CreateKeyResponse, error) {
	id := newID()
	for {
		if s.store.HasKey(id) {
			// QUESTION: exponential backoff in case our entropy is somehow.. not very.. good?
			id = newID()
			continue
		}
		break
	}

	artifact, err := NewKeyArtifact(request.GetPublicKey())
	if err != nil {
		return &pb.CreateKeyResponse{}, status.Errorf(codes.FailedPrecondition, "refused key: %v", err)
	}

	if s.store.Set(id, artifact) != nil {
		return &pb.CreateKeyResponse{}, status.Error(codes.Internal, "could not store key")
	}
	response := &pb.CreateKeyResponse{
		Id: id,
	}
	return response, nil
}

// GetKey returns the public key stored at `GetKeyRequest.Id`
func (s *Server) GetKey(ctx context.Context, request *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	artifact, err := s.store.Get(request.GetId())
	if err != nil {
		return &pb.GetKeyResponse{}, status.Errorf(codes.NotFound, "could not find key for requested id: %s", request.GetId())
	}
	response := &pb.GetKeyResponse{
		PublicKey: artifact.Data(),
	}
	return response, nil
}

// QUESTION: do we want to use UUIDs or something cutesy?
func newID() string {
	fake := gofakeit.NewCrypto()
	components := []string{
		fake.Adjective(),
		fake.Word(),
		fake.Word(),
		fake.Animal(),
	}
	id := strings.Join(components, "-")
	return id
}
