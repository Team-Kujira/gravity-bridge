package types

import (
	context "context"
	"time"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
	GetLastValidatorPower(ctx context.Context, operator sdk.ValAddress) (int64, error)
	GetLastTotalPower(ctx context.Context) (sdkmath.Int, error)
	IterateValidators(context.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error
	ValidatorQueueIterator(ctx context.Context, endTime time.Time, endHeight int64) (corestore.Iterator, error)
	GetParams(ctx context.Context) (stakingtypes.Params, error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	IterateBondedValidatorsByPower(context.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error
	IterateLastValidators(context.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error
	Validator(context.Context, sdk.ValAddress) (stakingtypes.ValidatorI, error)
	ValidatorByConsAddr(context.Context, sdk.ConsAddress) (stakingtypes.ValidatorI, error)
	Slash(context.Context, sdk.ConsAddress, int64, int64, math.LegacyDec) (sdkmath.Int, error)
	Jail(context.Context, sdk.ConsAddress) error
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx context.Context, denom string) (bank.Metadata, bool)
}

type SlashingKeeper interface {
	GetValidatorSigningInfo(ctx context.Context, address sdk.ConsAddress) (info slashingtypes.ValidatorSigningInfo, err error)
}

// AccountKeeper defines the interface contract required for account
// functionality.
type AccountKeeper interface {
	GetSequence(ctx context.Context, addr sdk.AccAddress) (uint64, error)
}
