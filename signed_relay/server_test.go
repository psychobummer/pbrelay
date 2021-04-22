package signedrelay

import (
	"context"
	"crypto/ed25519"
	"io"
	"net"
	"testing"
	"time"

	"github.com/psychobummer/pbrelay/keystore"
	pbk "github.com/psychobummer/pbrelay/rpc/keystore"
	pbr "github.com/psychobummer/pbrelay/rpc/relay"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestGetStream(t *testing.T) {
	dialer := testDialer(t)
	producerClient := makeRelayClient(t, dialer)
	subscriberClient := makeRelayClient(t, dialer)
	keystoreClient := makeKeystoreClient(t, dialer)

	// Store a public key for this producer
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	resp, err := keystoreClient.CreateKey(context.Background(), &pbk.CreateKeyRequest{PublicKey: pubKey})
	assert.Nil(t, err)

	// Before any subscribers can read from a stream, something needs to have created a stream.
	producerStream, err := producerClient.CreateStream(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}

	msg := pbr.StreamMessage{
		Id:        resp.GetId(),
		Data:      []byte("sign me please"),
		Signature: ed25519.Sign(privKey, []byte("sign me please")),
	}

	if err := producerStream.Send(&msg); err != nil {
		assert.Nil(t, err)
		t.Fatal(err.Error())
	}

	done := make(chan bool)
	// Connect a subscriber but don't read until we've written a message to catch
	go func() {
		<-time.After(50 * time.Millisecond) // give the initial message a moment to get processed
		subscriberStream, err := subscriberClient.GetStream(context.Background(), &pbr.GetStreamRequest{
			Id: msg.Id,
		})
		if err != nil {
			log.Fatal().Err(err)
		}

		for {
			recv, err := subscriberStream.Recv()
			if err != nil {
				t.Logf("TestStream: should not have encountered a grpc error: %v", err)
				t.Fail()
				break
			}
			assert.Equal(t, msg.Data, recv.Data)
			break
		}
		done <- true
	}()

	// write a message for the subscriber to consume
waitLoop:
	for {
		select {
		case <-done:
			break waitLoop
		default:
			if err := producerStream.Send(&msg); err != nil {
				log.Fatal().Err(err)
			}
		}
	}
}

func TestGetStreamNoRegisteredProducer(t *testing.T) {
	dialer := testDialer(t)
	subscriberClient := makeRelayClient(t, dialer)

	subscriberStream, err := subscriberClient.GetStream(context.Background(), &pbr.GetStreamRequest{
		Id: "idontexist",
	})
	if err != nil {
		log.Fatal().Err(err)
	}

	for {
		_, err = subscriberStream.Recv()
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				t.Log("Could not parse grpc status code from error?")
				t.Fail()
				break
			}
			assert.Equal(t, codes.NotFound, e.Code())
			break
		}
	}
}

func TestCreateStreamInvalidSignature(t *testing.T) {
	dialer := testDialer(t)

	// Store a public key for this producer
	pubKey, _, _ := ed25519.GenerateKey(nil)
	keystoreClient := makeKeystoreClient(t, dialer)
	keyResp, err := keystoreClient.CreateKey(context.Background(), &pbk.CreateKeyRequest{PublicKey: pubKey})
	assert.Nil(t, err)
	tests := []struct {
		testName string
		msg      *pbr.StreamMessage
		errCode  codes.Code
	}{
		{
			testName: "No public key",
			msg: &pbr.StreamMessage{
				Id:        "idonotexisdfsdfst",
				Data:      []byte(""),
				Signature: []byte(""),
			},
			errCode: codes.NotFound,
		},
		{
			testName: "Invalid signature",
			msg: &pbr.StreamMessage{
				Id:        keyResp.GetId(),
				Data:      []byte("thisisbaddata"),
				Signature: []byte("thisisaninvalidsig"),
			},
			errCode: codes.FailedPrecondition,
		},
	}

	producerClient := makeRelayClient(t, dialer)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			producerStream, err := producerClient.CreateStream(context.Background())
			if err != nil {
				t.Fatal(err.Error())
			}
			for {
				if err := producerStream.Send(tt.msg); err != nil {
					// We only care about errors that we didn't generate by closing the connection
					if err == io.EOF {
						break
					}
					_, actualErr := producerStream.CloseAndRecv()
					e, _ := status.FromError(actualErr)
					assert.Equal(t, tt.errCode, e.Code())
					break
				}

			}

		})
	}
}

type dialFunc func(context.Context, string) (net.Conn, error)

func makeRelayClient(t *testing.T, dialer dialFunc) pbr.RelayServiceClient {
	t.Helper()
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		log.Fatal().Err(err)
	}
	return pbr.NewRelayServiceClient(conn)
}

func makeKeystoreClient(t *testing.T, dialer dialFunc) pbk.KeystoreServiceClient {
	t.Helper()
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		log.Fatal().Err(err)
	}
	return pbk.NewKeystoreServiceClient(conn)
}

func testDialer(t *testing.T) dialFunc {
	t.Helper()

	listener := bufconn.Listen(1024 * 1024)
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	keystoreServer := keystore.NewServer()
	relayServer := NewServer(Config{
		KeystoreClient:            pbk.NewKeystoreServiceClient(conn),
		ValidateProducerSignature: true,
	})

	grpcServer := grpc.NewServer()
	pbr.RegisterRelayServiceServer(grpcServer, relayServer)
	pbk.RegisterKeystoreServiceServer(grpcServer, keystoreServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
