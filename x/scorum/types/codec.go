package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMint{}, "scorum/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "scorum/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgConvertSCR2SP{}, "scorum/MsgConvertSCR2SP", nil)
	cdc.RegisterConcrete(&MsgWithdrawSP{}, "scorum/MsgWithdrawSP", nil)
	cdc.RegisterConcrete(&MsgStopSPWithdrawal{}, "scorum/MsgStopSPWithdrawal", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgBurn{},
		&MsgConvertSCR2SP{},
		&MsgWithdrawSP{},
		&MsgStopSPWithdrawal{},
		&MsgMint{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
