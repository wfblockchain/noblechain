package network

import (
	"fmt"
	"testing"
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	genutil "github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	cctptypes "github.com/wfblockchain/noble-cctp/x/cctp/types"
	paramauthoritytypes "github.com/wfblockchain/noble-paramauthority/x/params/types/proposal"
	paramauthorityupgradetypes "github.com/wfblockchain/noble-paramauthority/x/upgrade/types"
	"github.com/wfblockchain/noblechain/v5/app"
	"github.com/wfblockchain/noblechain/v5/cmd"
	"github.com/wfblockchain/noblechain/v5/testutil/sample"
)

type (
	Network = network.Network
	Config  = network.Config
)

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...network.Config) *network.Network {
	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = configs[0]
	}
	//
	net, err := network.New(t, "", cfg)
	if err != nil {
		// handle the error
		fmt.Println("Error:", err)
		// return or exit the function depending on your use case
	}

	t.Cleanup(net.Cleanup)
	return net
}

func DefaultConfig() network.Config {
	// app doesn't have this modules anymore, but we need them for test setup, which uses gentx and MsgCreateValidator
	app.ModuleBasics[genutiltypes.ModuleName] = genutil.AppModuleBasic{}
	app.ModuleBasics[stakingtypes.ModuleName] = staking.AppModuleBasic{}

	encoding := cmd.MakeEncodingConfig(app.ModuleBasics)
	cfg := network.Config{
		Codec:             encoding.Marshaler,
		TxConfig:          encoding.TxConfig,
		LegacyAmino:       encoding.Amino,
		InterfaceRegistry: encoding.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			return app.New(
				val.GetCtx().Logger, tmdb.NewMemDB(), nil, true, map[int64]bool{}, val.GetCtx().Config.RootDir, 0,
				encoding,
				simtestutil.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			)
		},
		GenesisState:  app.ModuleBasics.DefaultGenesis(encoding.Marshaler),
		TimeoutCommit: 2 * time.Second,
		ChainID:       "chain-" + tmrand.NewRand().Str(6),
		// Some changes are introduced to make the tests run as if Noble is a standalone chain.
		// This will only work if NumValidators is set to 1.
		NumValidators:   1,
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}

	// Authority needs to be present to pass genesis validation
	params := paramauthoritytypes.DefaultGenesis()
	params.Params.Authority = sample.AccAddress()
	cfg.GenesisState[paramstypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(params)

	// Authority needs to be present to pass genesis validation
	upgrade := paramauthorityupgradetypes.DefaultGenesis()
	upgrade.Params.Authority = sample.AccAddress()
	cfg.GenesisState[upgradetypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(upgrade)

	cctp := cctptypes.DefaultGenesis()
	cctp.Owner = sample.AccAddress()
	cctp.AttesterManager = sample.AccAddress()
	cctp.Pauser = sample.AccAddress()
	cctp.TokenController = sample.AccAddress()
	cfg.GenesisState[cctptypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(cctp)

	return cfg
}
