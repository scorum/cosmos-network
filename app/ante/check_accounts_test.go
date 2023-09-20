package ante

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	aviatrixtypes "github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func Test_extractAddresses(t *testing.T) {
	addr1, addr2, addr3 := sample.AccAddress(), sample.AccAddress(), sample.AccAddress()

	tt := []struct {
		Msg  sdk.Msg
		Addr []sdk.AccAddress
	}{
		{
			Msg:  &scorumtypes.MsgBurn{Supervisor: addr1.String()},
			Addr: []sdk.AccAddress{addr1},
		},
		{
			Msg: &aviatrixtypes.MsgCreatePlane{
				Supervisor: addr1.String(),
				Owner:      addr2.String(),
			},
			Addr: []sdk.AccAddress{addr1, addr2},
		},
		{
			Msg: &banktypes.MsgMultiSend{
				Inputs: []banktypes.Input{
					{Address: addr1.String()},
					{Address: addr2.String()},
				},
				Outputs: []banktypes.Output{
					{Address: addr3.String()},
				},
			},
			Addr: []sdk.AccAddress{addr1, addr2, addr3},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(reflect.Indirect(reflect.ValueOf(tc.Msg)).Type().Name(), func(t *testing.T) {
			t.Parallel()

			require.ElementsMatch(t, tc.Addr, extractAddresses(tc.Msg))
		})
	}
}
