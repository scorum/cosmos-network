package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MintProposal{}, "scorum/MintProposal", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "scorum/Burn", nil)
	cdc.RegisterConcrete(&MsgRegisterAccount{}, "scorum/RegisterAccount", nil)
	cdc.RegisterConcrete(&MsgConvertSCR2SP{}, "scorum/ConvertSCR2SP", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&MintProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgBurn{},
		&MsgRegisterAccount{},
		&MsgConvertSCR2SP{},
	)
	registry.RegisterInterface(
		"cosmos.gov.v1beta1.Content",
		(*govtypes.Content)(nil),
		&MintProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
