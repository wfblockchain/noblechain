package app_test

import (
	"os"
	"testing"

	"cosmossdk.io/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/stretchr/testify/require"
	"github.com/wfblockchain/noblechain/v5/app"
	"github.com/wfblockchain/noblechain/v5/cmd"
)

func init() {
	simcli.GetSimulatorFlags()
}

const AppChainID = "app"

type SimApp interface {
	cmd.App
	GetBaseApp() *baseapp.BaseApp
	AppCodec() codec.Codec
	SimulationManager() *module.SimulationManager
	ModuleAccountAddrs() map[string]bool
	Name() string
	LegacyAmino() *codec.LegacyAmino
	BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock
	EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock
	InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain
}

// BenchmarkSimulation run the chain simulation
// Running using starport command:
// `starport chain simulate -v --numBlocks 200 --blockSize 50`
// Running as go benchmark test:
// `go test -benchmem -run=^$ -bench ^BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func BenchmarkSimulation(b *testing.B) {
	simcli.FlagNumBlocksValue = 250
	simcli.FlagBlockSizeValue = 50
	simcli.FlagCommitValue = true
	simcli.FlagVerboseValue = true
	simcli.FlagEnabledValue = true
	config := simcli.NewConfigFromFlags()
	config.ChainID = AppChainID

	db, dir, logger, _, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
	require.NoError(b, err, "simulation setup failed")

	b.Cleanup(func() {
		db.Close()
		err = os.RemoveAll(dir)
		require.NoError(b, err)
	})

	encoding := cmd.MakeEncodingConfig(app.ModuleBasics)

	app := app.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encoding,
		simtestutil.EmptyAppOptions{},
	)

	// simApp, ok := app
	// require.True(b, ok, "can't use simapp")

	// Run randomized simulations
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		app.GetBaseApp(),
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), simapp.SimApp{}.DefaultGenesis()),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(&simapp.SimApp{}, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(&simapp.SimApp{}, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}
