package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"strings"
	"virel-blockchain/blockchain"
	"virel-blockchain/config"
	"virel-blockchain/logger"
)

var Log = logger.New()

func init() {
	blockchain.Log = Log
}

var cpu_profile = flag.String("cpu-profile", "", "write cpu profile to the provided file")

func main() {
	p2p_bind_port := flag.Uint("p2p-bind-port", config.P2P_BIND_PORT, "starts P2P server on this port")
	public_rpc := flag.Bool("public-rpc", false, "required for public RPC nodes: blocks private RPC calls and binds on 0.0.0.0")
	rpc_bind_port := flag.Uint("rpc-bind-port", config.RPC_BIND_PORT, "starts RPC server on this port")
	stratum_bind_ip := flag.String("stratum-bind-ip", "127.0.0.1", "use 0.0.0.0 to expose Stratum server")
	stratum_bind_port := flag.Uint("stratum-bind-port", config.STRATUM_BIND_PORT, "")
	log_level := flag.Uint("log-level", 1, "sets the log level")

	var slavechains_stratums *string
	var stratum_wallet *string
	if config.IS_MASTERCHAIN {
		slavechains_stratums = flag.String("slavechains-stratums", "", "comma-separated list of slave chain stratum IP:PORT")
		stratum_wallet = flag.String("mining-wallet", "", "merge mining stratum address")
	}

	flag.Parse()

	if *cpu_profile != "" {
		f, err := os.Create(*cpu_profile)
		if err != nil {
			Log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			Log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	Log.SetLogLevel(uint8(*log_level))

	Log.Info("Starting", config.NETWORK_NAME, "node")
	if config.NETWORK_NAME != "mainnet" {
		Log.Warn("This is a", strings.ToUpper(config.NETWORK_NAME), "node, only for testing the blockchain.")
		Log.Warn("Be aware that any amount transacted in", config.NETWORK_NAME, "is worthless.")
	}

	bc := blockchain.New()

	if config.IS_MASTERCHAIN {
		if len(*slavechains_stratums) > 0 {
			stratums := strings.Split(*slavechains_stratums, ",")
			if stratum_wallet == nil || len(*stratum_wallet) == 0 {
				Log.Fatal("you must specify your wallet address when merge mining using --mining-wallet")
			}
			for _, v := range stratums {
				go bc.AddStratum(v, *stratum_wallet, true)
			}
		}
	}

	bind_ip := "127.0.0.1"
	if *public_rpc {
		bind_ip = "0.0.0.0"
	}

	go startRpc(bc, bind_ip, uint16(*rpc_bind_port), *public_rpc)
	go bc.StartStratum(*stratum_bind_ip, uint16(*stratum_bind_port))
	go bc.StartP2P(config.SEED_NODES, uint16(*p2p_bind_port))
	go bc.NewStratumJob(true)

	prompts(bc)
}
