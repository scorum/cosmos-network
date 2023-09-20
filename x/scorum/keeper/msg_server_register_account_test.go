package keeper_test

import (
	"fmt"
	"testing"

	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_RegisterAccount(t *testing.T) {
	set, ctx := setupKeeper(t)

	s := keeper.NewMsgServer(set.keeper)

	addr := sample.AccAddress()

	_, err := s.RegisterAccount(ctx, types.NewMsgRegisterAccount(set.supervisor.String(), addr.String()))
	require.NoError(t, err)

	require.Equal(
		t,
		fmt.Sprintf("%d%s", set.keeper.GetParams(ctx.Context).RegistrationSPDelegationAmount.Int.Int64(), types.SPDenom),
		set.bankKeeper.GetBalance(ctx.Context, addr, types.SPDenom).String(),
	)
}

func TestMsgServer_RegisterAccount_DoubleRegistration(t *testing.T) {
	set, ctx := setupKeeper(t)

	s := keeper.NewMsgServer(set.keeper)

	addr := sample.AccAddress()

	_, err := s.RegisterAccount(ctx, types.NewMsgRegisterAccount(set.supervisor.String(), addr.String()))
	require.NoError(t, err)

	_, err = s.RegisterAccount(ctx, types.NewMsgRegisterAccount(set.supervisor.String(), addr.String()))
	require.Error(t, err, "conflict")
}

func TestMsgServer_RegisterAccount_NotSupervisor(t *testing.T) {
	set, ctx := setupKeeper(t)

	s := keeper.NewMsgServer(set.keeper)

	addr := sample.AccAddress()

	_, err := s.RegisterAccount(ctx, types.NewMsgRegisterAccount(addr.String(), addr.String()))
	require.Error(t, err)
}
