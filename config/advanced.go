package config

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"time"
)

// This file holds advanced config options. You shouldn't edit these options unless you really know what you
// are doing.

const MAX_MERGE_MINED_CHAINS = 16
const DEFAULT_CHECKPOINT_INTERVAL = 32
const SEEDHASH_DURATION = 4 * (60 * 60 * 24) // seed hash changes once every 4 days

const BLOCKS_PER_DAY = 60 * 60 * 24 / TARGET_BLOCK_TIME

const STRATUM_READ_TIMEOUT = 90 * time.Second
const STRATUM_JOBS_HISTORY = MINIDAG_ANCESTORS

var ATOMIC = math.Round(math.Log10(COIN))

const WALLET_PREFIX = "s" // Wallet prefix should be the same for all merge mined chains

// True if the current chain is a Master Chain; if false, the node will try connecting to the masterchain
// node to send Merge Mining jobs
const IS_MASTERCHAIN = NETWORK_ID == 0x4af15cf1542ba49a // do not change this

const PARALLEL_BLOCKS_DOWNLOAD = 50

var BinaryNetworkID = make([]byte, 8)

func init() {
	binary.LittleEndian.PutUint64(BinaryNetworkID, NETWORK_ID)
}
func AssertHexDec(s string) []byte {
	dat, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return dat
}
