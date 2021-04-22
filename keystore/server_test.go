package keystore

import (
	"context"
	"crypto/ed25519"
	"log"
	"net"
	"testing"

	pb "github.com/psychobummer/pbrelay/rpc/keystore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewID(t *testing.T) {
	assert.NotEqual(t, newID(), newID())
}

func TestCreateKey(t *testing.T) {
	pubKey, _, _ := ed25519.GenerateKey(nil)

	tests := []struct {
		key     []byte
		errCode codes.Code
	}{
		{
			key:     []byte("thisistooshorttobeakey"),
			errCode: codes.FailedPrecondition,
		},
		{
			key:     []byte("idosaythisisfartoolongtobeaed25519key"),
			errCode: codes.FailedPrecondition,
		},
		{
			key:     pubKey,
			errCode: codes.OK,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewKeystoreServiceClient(conn)
	for _, tt := range tests {
		_, err := client.CreateKey(ctx, &pb.CreateKeyRequest{PublicKey: tt.key})
		e, _ := status.FromError(err)
		assert.Equal(t, tt.errCode, e.Code())
	}
}

// create a key and then verify we can fetch it
func TestGetKey(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	pubKey, _, _ := ed25519.GenerateKey(nil)

	client := pb.NewKeystoreServiceClient(conn)
	createResp, err := client.CreateKey(ctx, &pb.CreateKeyRequest{PublicKey: pubKey})
	assert.Nil(t, err)

	getResp, err := client.GetKey(ctx, &pb.GetKeyRequest{Id: createResp.GetId()})
	assert.Nil(t, err)
	assert.Equal(t, pubKey, ed25519.PublicKey(getResp.GetPublicKey()))
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	keystoreServer := NewServer()

	pb.RegisterKeystoreServiceServer(grpcServer, keystoreServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
