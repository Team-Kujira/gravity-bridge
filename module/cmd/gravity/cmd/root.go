package cmd

import (
	"context"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	// snapshotstypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"cosmossdk.io/log"
	tmcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/peggyjv/gravity-bridge/module/v4/app"
	"github.com/peggyjv/gravity-bridge/module/v4/app/params"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	app.SetAddressConfig()

	encodingConfig := app.MakeEncodingConfig()
	tempApp := app.NewGravityApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		nil,
		"",
		0,
		encodingConfig,
		simtestutil.NewAppOptionsWithFlagHome(tempDir()),
	)

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome)

	rootCmd := &cobra.Command{
		Use:   "gravity",
		Short: "Stargate Gravity App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig, tempApp.BasicModuleManager, tempApp.ModuleManager)

	return rootCmd, encodingConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	cfg.P2P.Seeds = ""

	// these values put a higher strain on node memory
	cfg.P2P.MaxNumInboundPeers = 320
	cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	// WASMConfig defines configuration for the wasm module.
	type WASMConfig struct {
		// This is the maximum sdk gas (wasm and storage) that we allow for any x/wasm "smart" queries
		QueryGasLimit uint64 `mapstructure:"query_gas_limit"`

		// Address defines the gRPC-web server to listen on
		LruSize uint64 `mapstructure:"lru_size"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	customAppTemplate := serverconfig.DefaultConfigTemplate
	return customAppTemplate, srvCfg
}

// Execute executes the root command.
func Execute(rootCmd *cobra.Command) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	srvCtx := server.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, srvCtx)

	rootCmd.PersistentFlags().String("log_level", "info", "The logging level in the format of <module>:<level>,...")

	executor := tmcli.PrepareBaseCmd(rootCmd, "", app.DefaultNodeHome)
	return executor.ExecuteContext(ctx)
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig, basicModuleManager module.BasicManager, mm *module.Manager) {
	gentxModule := mm.Modules[genutiltypes.ModuleName].(genutil.AppModule)
	txConfig := encodingConfig.TxConfig

	rootCmd.AddCommand(
		genutilcli.InitCmd(basicModuleManager, app.DefaultNodeHome),
		CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, gentxModule.GenTxValidator, txConfig.SigningContext().ValidatorAddressCodec()),
		GenTxCmd(mm, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, txConfig.SigningContext().ValidatorAddressCodec(), app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
	)

	a := appCreator{encodingConfig}
	server.AddCommands(rootCmd, app.DefaultNodeHome, a.newApp, a.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(),
		Commands(app.DefaultNodeHome),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.QueryTxCmd(),
		authcmd.QueryTxsByEventsCmd(),
		rpc.QueryEventForTxCmd(),
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		server.QueryBlockResultsCmd(),
		server.QueryBlocksCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	// var cache sdk.MultiStorePersistentCache

	// if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
	// 	cache = store.NewCommitKVStoreCacheManager()
	// }

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	// snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	// snapshotDB, err := tmdb.NewDB("metadata", server.GetAppDBBackend(appOpts), snapshotDir)
	// if err != nil {
	// 	panic(err)
	// }
	// snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	// if err != nil {
	// 	panic(err)
	// }

	// snapshotOpts := snapshotstypes.SnapshotOptions{
	// 	Interval:   cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
	// 	KeepRecent: cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
	// }
	return app.NewGravityApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		// baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		// baseapp.SetSnapshot(snapshotStore, snapshotOpts),
	)
}

func (a appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {

	var gravity *app.Gravity
	if height != -1 {
		gravity = app.NewGravityApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1), a.encCfg, appOpts)

		if err := gravity.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		gravity = app.NewGravityApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1), a.encCfg, appOpts)
	}

	return gravity.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)

}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "gravityd")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	return dir
}
