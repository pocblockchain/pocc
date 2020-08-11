package lcd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"

	"github.com/pocblockchain/pocc/client/context"
	"github.com/pocblockchain/pocc/client/flags"
	"github.com/pocblockchain/pocc/codec"
	keybase "github.com/pocblockchain/pocc/crypto/keys"
	"github.com/pocblockchain/pocc/server"

	// unnamed import of statik for swagger UI support
	_ "github.com/pocblockchain/pocc/client/lcd/statik"
)

//go:generate statik -src=./swagger-ui -c "Run go generate root.go with client/lcd as working dir to generate this file"
// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux     *mux.Router
	CliCtx  context.CLIContext
	KeyBase keybase.Keybase

	log      log.Logger
	listener net.Listener
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *codec.Codec) *RestServer {
	r := mux.NewRouter()
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	return &RestServer{
		Mux:    r,
		CliCtx: cliCtx,
		log:    logger,
	}
}

// Start starts the rest server
func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint) (err error) {
	server.TrapSignal(func() {
		err := rs.listener.Close()
		rs.log.Error("error closing listener", "err", err)
	})

	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	rs.log.Info(
		fmt.Sprintf(
			"Starting application REST service (chain-id: %q)...",
			viper.GetString(flags.FlagChainID),
		),
	)

	return rpcserver.StartHTTPServer(rs.listener, rs.Mux, rs.log, cfg)
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec, registerRoutesFn func(*RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := NewRestServer(cdc)

			registerRoutesFn(rs)
			rs.registerSwaggerUI()

			//	fmt.Printf("ServeCommand, cdc %v, rs:%v", cdc, rs)
			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
			)

			return err
		},
	}

	return flags.RegisterRestServerFlags(cmd)
}

func (rs *RestServer) registerSwaggerUI() {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}

func (rs *RestServer) RegisterSwaggerUI() {
	rs.registerSwaggerUI()
}
