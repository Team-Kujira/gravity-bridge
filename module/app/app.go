package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"cosmossdk.io/errors"
	// "cosmossdk.io/simapp"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	// distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/gorilla/mux"
	gravityparams "github.com/peggyjv/gravity-bridge/module/v4/app/params"
	v2 "github.com/peggyjv/gravity-bridge/module/v4/app/upgrades/v2"
	v3 "github.com/peggyjv/gravity-bridge/module/v4/app/upgrades/v3"
	v4 "github.com/peggyjv/gravity-bridge/module/v4/app/upgrades/v4"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity"
	gravityclient "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/client"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/keeper"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmos "github.com/tendermint/tendermint/libs/os"
	// unnamed import of statik for swagger UI support
	// _ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

const (
	appName = "app"

	// MaxAddrLen is the maximum allowed length (in bytes) for an address.
	//
	// NOTE: In the SDK, the default value is 255.
	MaxAddrLen = 20
)

var (
	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome string

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				// distrclient.ProposalHandler,
				// upgradeclient.LegacyProposalHandler,
				// upgradeclient.LegacyCancelProposalHandler,
				gravityclient.ProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		gravity.AppModuleBasic{},
	)

	// module account permissions
	// NOTE: We believe that this is giving various modules access to functions of the supply module? We will probably need to use this.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}

	// verify app interface at compile time
	// _ simapp.App              = (*Gravity)(nil)
	_ servertypes.Application = (*Gravity)(nil)
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
	vesting.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

// Gravity extended ABCI application
type Gravity struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tKeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	accountKeeper         authkeeper.AccountKeeper
	bankKeeper            bankkeeper.Keeper
	capabilityKeeper      *capabilitykeeper.Keeper
	stakingKeeper         *stakingkeeper.Keeper
	slashingKeeper        slashingkeeper.Keeper
	mintKeeper            mintkeeper.Keeper
	distrKeeper           distrkeeper.Keeper
	govKeeper             govkeeper.Keeper
	crisisKeeper          *crisiskeeper.Keeper
	upgradeKeeper         *upgradekeeper.Keeper
	consensusParamsKeeper consensusparamkeeper.Keeper
	paramsKeeper          paramskeeper.Keeper
	ibcKeeper             *ibckeeper.Keeper
	evidenceKeeper        evidencekeeper.Keeper
	transferKeeper        ibctransferkeeper.Keeper
	gravityKeeper         keeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// Module Manager
	ModuleManager      *module.Manager
	BasicModuleManager module.BasicManager

	// configurator
	configurator module.Configurator

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".gravity")
}

func NewGravityApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig gravityparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Gravity {

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		gravitytypes.StoreKey,
	)
	tKeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	var app = &Gravity{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
		memKeys:           memKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.consensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authority,
		runtime.EventService{},
	)
	bApp.SetParamStore(&app.consensusParamsKeeper.ParamsStore)

	app.capabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should
	// call `Seal` after creating their scoped modules in the app via
	// `ScopeToModule`.
	app.capabilityKeeper.Seal()

	app.accountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(gravityparams.Bech32PrefixAccAddr),
		gravityparams.Bech32PrefixAccAddr,
		authority,
	)

	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.accountKeeper,
		app.BlockedAddrs(),
		authority,
		logger,
	)

	app.stakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.accountKeeper,
		app.bankKeeper,
		authority,
		authcodec.NewBech32Codec(gravityparams.Bech32PrefixValAddr),
		authcodec.NewBech32Codec(gravityparams.Bech32PrefixConsAddr),
	)

	app.mintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		app.stakingKeeper,
		app.accountKeeper,
		app.bankKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.stakingKeeper,
		authority,
	)

	app.crisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.bankKeeper,
		authtypes.FeeCollectorName,
		authority,
		app.accountKeeper.AddressCodec(),
	)

	app.upgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authority,
	)

	app.gravityKeeper = keeper.NewKeeper(
		appCodec,
		keys[gravitytypes.StoreKey],
		app.GetSubspace(gravitytypes.ModuleName),
		app.accountKeeper,
		app.stakingKeeper,
		app.bankKeeper,
		app.slashingKeeper,
		app.distrKeeper,
		sdk.DefaultPowerReduction,
		app.ModuleAccountAddressesToNames([]string{}),
		app.ModuleAccountAddressesToNames([]string{distrtypes.ModuleName}),
	)

	app.stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.gravityKeeper.Hooks(),
		),
	)

	app.ibcKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.stakingKeeper,
		app.upgradeKeeper,
		scopedIBCKeeper,
		authority,
	)

	app.transferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper,
		app.ibcKeeper.ChannelKeeper,
		app.ibcKeeper.PortKeeper,
		app.accountKeeper,
		app.bankKeeper,
		scopedTransferKeeper,
		authority,
	)

	transferModule := ibctransfer.NewAppModule(app.transferKeeper)
	transferIBCModule := ibctransfer.NewIBCModule(app.transferKeeper)

	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.stakingKeeper,
		app.slashingKeeper,
		app.accountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	app.evidenceKeeper = *evidenceKeeper

	govRouter := govtypesv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.ibcKeeper.ClientKeeper))

	// Are defaults ok here?
	govConfig := govtypes.DefaultConfig()
	app.govKeeper = *govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		app.distrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authority,
	)

	app.setupUpgradeStoreLoaders()

	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.ModuleManager = module.NewManager(
		genutil.NewAppModule(
			app.accountKeeper,
			app.stakingKeeper,
			app,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(
			appCodec,
			app.accountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		vesting.NewAppModule(
			app.accountKeeper,
			app.bankKeeper,
		),
		bank.NewAppModule(
			appCodec,
			app.bankKeeper,
			app.accountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(
			appCodec,
			*app.capabilityKeeper,
			false,
		),
		crisis.NewAppModule(
			app.crisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
		gov.NewAppModule(
			appCodec,
			&app.govKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(
			appCodec,
			app.mintKeeper,
			app.accountKeeper,
			nil,
			app.GetSubspace(minttypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.slashingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
			app.interfaceRegistry,
		),
		distr.NewAppModule(
			appCodec,
			app.distrKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(
			appCodec,
			app.stakingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(app.upgradeKeeper, app.accountKeeper.AddressCodec()),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
		gravity.NewAppModule(
			app.gravityKeeper,
			app.bankKeeper,
		),
	)

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.ModuleManager,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(app.LegacyAmino())
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.ModuleManager.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	app.ModuleManager.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		gravitytypes.ModuleName,
	)
	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		gravitytypes.ModuleName,
	)
	app.ModuleManager.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		gravitytypes.ModuleName,
	)

	app.ModuleManager.RegisterInvariants(app.crisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.ModuleManager.RegisterServices(app.configurator)

	app.setupUpgradeHandlers()

	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(
			app.appCodec,
			app.accountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.accountKeeper,
			BankKeeper:      app.bankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  nil,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create ante handler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App
func (app *Gravity) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *Gravity) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.ModuleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *Gravity) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.ModuleManager.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *Gravity) InitChainer(
	ctx sdk.Context,
	req *abci.RequestInitChain,
) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	if err := app.upgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
		panic(err)
	}

	return app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *Gravity) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Gravity) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ModuleAccountNames returns a map of module account address to module name
func (app *Gravity) ModuleAccountAddressesToNames(moduleAccounts []string) map[string]string {
	modAccNames := make(map[string]string)
	for _, acc := range moduleAccounts {
		modAccNames[authtypes.NewModuleAddress(acc).String()] = acc
	}

	return modAccNames
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *Gravity) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	// allow the following addresses to receive funds
	delete(blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	delete(blockedAddrs, authtypes.NewModuleAddress(authtypes.FeeCollectorName).String())

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *Gravity) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tKeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Gravity) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
func (app *Gravity) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *Gravity) SimulationManager() *module.SimulationManager {
	return app.sm
}

// API server.
func (app *Gravity) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// TODO: build the custom gravity swagger files and add here?
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterSwaggerAPI registers swagger route with API Server
// TODO: build the custom gravity swagger files and add here?
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *Gravity) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Gravity) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

func (app *Gravity) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}

	return modAccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)

	return paramsKeeper
}

func VerifyAddressFormat(bz []byte) error {
	if len(bz) == 0 {
		return errors.Wrap(sdkerrors.ErrUnknownAddress, "invalid address; cannot be empty")
	}
	if len(bz) != MaxAddrLen {
		return errors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"invalid address length; got: %d, max: %d", len(bz), MaxAddrLen,
		)
	}

	return nil
}

func (app *Gravity) setupUpgradeStoreLoaders() {
	_, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	// if upgradeInfo.Name matches a plan name with a module being added, renamed, or deleted,
	// create a storetypes.StoreUpgrades struct and
	// app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	// see also:
	// https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/upgrade.md#add-storeupgrades-for-new-modules
}

func (app *Gravity) setupUpgradeHandlers() {
	app.upgradeKeeper.SetUpgradeHandler(
		v2.UpgradeName,
		v2.CreateUpgradeHandler(
			app.ModuleManager,
			app.configurator,
			app.bankKeeper,
		),
	)

	app.upgradeKeeper.SetUpgradeHandler(
		v3.UpgradeName,
		v3.CreateUpgradeHandler(
			app.ModuleManager,
			app.configurator,
		),
	)

	app.upgradeKeeper.SetUpgradeHandler(
		v4.UpgradeName,
		v4.CreateUpgradeHandler(
			app.ModuleManager,
			app.configurator,
		),
	)
}
