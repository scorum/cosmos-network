package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/gogo/protobuf/proto"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePlane{}, "aviatrix/CreatePlane", nil)
	cdc.RegisterConcrete(&MsgUpdatePlaneExperience{}, "aviatrix/UpdatePlaneExperience", nil)
	cdc.RegisterConcrete(&MsgAdjustPlaneExperience{}, "aviatrix/AdjustPlaneExperience", nil)
	cdc.RegisterConcrete(&PlaneMeta{}, "aviatrix/PlaneMeta", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePlane{},
		&MsgUpdatePlaneExperience{},
		&MsgAdjustPlaneExperience{},
	)

	registry.RegisterImplementations((*proto.Message)(nil), &PlaneMeta{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
