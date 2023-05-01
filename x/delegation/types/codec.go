package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegate{}, "kyve/Delegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "kyve/Undelegate", nil)
	cdc.RegisterConcrete(&MsgRedelegate{}, "kyve/Redelegate", nil)
	cdc.RegisterConcrete(&MsgWithdrawRewards{}, "kyve/WithdrawRewards", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "kyve/UpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUndelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRedelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
