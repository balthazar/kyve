package types

import (
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx sdk.Context, plan types.Plan) error
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	GetPoolWithError(ctx sdk.Context, poolId uint64) (pooltypes.Pool, error)
	GetPool(ctx sdk.Context, id uint64) (val pooltypes.Pool, found bool)

	IncrementBundleInformation(ctx sdk.Context, poolId uint64, currentHeight uint64, currentKey string, currentValue string)

	GetAllPools(ctx sdk.Context) (list []pooltypes.Pool)
	ChargeFundersOfPool(ctx sdk.Context, poolId uint64, amount uint64) error
}

type StakerKeeper interface {
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
	GetCommission(ctx sdk.Context, stakerAddress string) sdk.Dec
	AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error

	DoesStakerExist(ctx sdk.Context, staker string) bool
	DoesValaccountExist(ctx sdk.Context, poolId uint64, stakerAddress string) bool

	LeavePool(ctx sdk.Context, staker string, poolId uint64)

	IncrementPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (newPoints uint64)
	ResetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (previousPoints uint64)
}

type DelegationKeeper interface {
	GetDelegationAmount(ctx sdk.Context, staker string) uint64
	GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64
	PayoutRewards(ctx sdk.Context, staker string, amount uint64, payerModuleName string) (success bool)
	SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType delegationTypes.SlashType)
}
