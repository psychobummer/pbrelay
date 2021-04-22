package cmd

import (
	"net"

	"github.com/psychobummer/pbrelay/keystore"
	pbk "github.com/psychobummer/pbrelay/rpc/keystore"
	pbr "github.com/psychobummer/pbrelay/rpc/relay"
	signedrelay "github.com/psychobummer/pbrelay/signed_relay"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Start the server but a longer description",
	Run:   doServer,
}

func init() {
	serverCmd.Flags().StringP("listen", "l", "", "ip:port to server should bind to (e.g: 0.0.0.0:9999) (required)")
	serverCmd.Flags().BoolP("validate", "s", true, "validate signature of producer messages (default: true)")
	serverCmd.MarkFlagRequired("listen")
	rootCmd.AddCommand(serverCmd)
}

func doServer(cmd *cobra.Command, args []string) {
	listenOn, _ := cmd.Flags().GetString("listen")
	validateSig, _ := cmd.Flags().GetBool("validate")

	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		log.Fatal().Msgf("couldn't create tcp listener: %v", err)
	}
	conn, err := grpc.Dial(listenOn, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("couldnt dial server: %v", err)
	}
	defer conn.Close()

	keystoreServer := keystore.NewServer()
	relayServer := signedrelay.NewServer(signedrelay.Config{
		KeystoreClient:            pbk.NewKeystoreServiceClient(conn),
		ValidateProducerSignature: validateSig,
	})
	grpcServer := grpc.NewServer()
	pbr.RegisterRelayServiceServer(grpcServer, relayServer)
	pbk.RegisterKeystoreServiceServer(grpcServer, keystoreServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Error().Msg(err.Error())
	}
}
