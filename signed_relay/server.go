package signedrelay

import (
	"context"
	"crypto/ed25519"
	"io"

	"github.com/google/uuid"
	pbk "github.com/psychobummer/pbrelay/rpc/keystore"
	pbr "github.com/psychobummer/pbrelay/rpc/relay"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Config configures a signed realy server
type Config struct {
	KeystoreClient            pbk.KeystoreServiceClient // a configured keyseystore grpc client
	ValidateProducerSignature bool                      // should we validate the signature of producer messages.
}

// The Server contains the GRPC service implementations and a broker
// for broacasting messages from a producer to all subscribed clients.
type Server struct {
	broker *Broker
	config Config
}

// NewServer returns a new initailized *Server
func NewServer(config Config) *Server {
	return &Server{
		broker: NewBroker(),
		config: config,
	}
}

// CreateStream will, on the first received message, register this stream with the broker.
// From there on, we'll simply broadcast any message we receive, and the broker
// will determine to whom they go.
// If we encounter any errors in the producer grpc stream, or it disconnects,
// we deregister from the broker, which will disconnect all subscribers to that stream.
func (s *Server) CreateStream(stream pbr.RelayService_CreateStreamServer) error {
	validated := false
	registered := false

	for {
		message, err := stream.Recv()
		log.Trace().Msgf("server received message from %v", message.GetId())
		if err == io.EOF {
			return status.Errorf(codes.OK, "client disconnected")
		}
		if err != nil {
			log.Error().Msgf("error reading from producer: %s", err.Error())
			return status.Errorf(codes.OK, "%v", err)
		}

		// verify the signature of the first message in the stream. clients will verify the
		// integrity of every message they receive -- the server only want to make sure that people
		// can't accidentally start streaming data that will never verify.
		producerID := message.GetId()
		if s.config.ValidateProducerSignature && !validated {
			keyResp, err := s.config.KeystoreClient.GetKey(context.Background(), &pbk.GetKeyRequest{Id: producerID})
			if err != nil {
				log.Error().Msg(err.Error())
				return status.Errorf(codes.NotFound, "could not find producer public key: %v", err)
			}
			if !ed25519.Verify(keyResp.GetPublicKey(), message.GetData(), message.GetSignature()) {
				log.Error().Msgf("signature mismatch for %s", producerID)
				return status.Error(codes.FailedPrecondition, "a public key was found, but it does appear to have signed this data")
			}
			validated = true
		}

		// register this producer if they aren't already registered
		if !registered {
			if err := s.broker.RegisterProducer(producerID); err != nil {
				log.Error().Msgf("error registering producer: %s", err.Error())
				return status.Errorf(codes.Internal, "%v", err)
			}
			registered = true
			defer func() {
				if err := s.broker.DeregisterProducer(producerID); err != nil {
					log.Error().Msg(err.Error())
				}
			}()
		}

		// all good
		s.broker.Broadcast(message)
	}
}

// GetStream registers this client as a subscriber to the requested stream.
// On error or disconnect, we dregister the subscription.
func (s *Server) GetStream(request *pbr.GetStreamRequest, stream pbr.RelayService_GetStreamServer) error {
	clientID := uuid.New().String()
	log.Trace().Msgf("subscriber connection from client %s", clientID)
	queue, err := s.broker.RegisterSubscriber(request.GetId(), clientID)
	if err != nil {
		log.Error().Msgf("could not subscribe %s to requested stream %s: %s", clientID, request.GetId(), err.Error())
		return status.Errorf(codes.NotFound, "%v", err)
	}
	defer s.broker.DeregisterSubscriber(request.GetId(), clientID)
	for message := range queue {
		if err := stream.Send(message); err != nil {
			log.Error().Msg(err.Error())
			return err
		}
	}
	return nil
}
