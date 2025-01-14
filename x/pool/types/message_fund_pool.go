package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFundPool{}

func NewMsgFundPool(creator string, id uint64, amount uint64) *MsgFundPool {
	return &MsgFundPool{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgFundPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgFundPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if util.ValidateNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	return nil
}
