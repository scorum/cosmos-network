package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func TestMsgServer_UpdatePlaneName(t *testing.T) {
	supervisorAddr := sample.AccAddress()
	addr := sample.AccAddress()
	id := uuid.New().String()

	testCases := []struct {
		name string

		owner string
		msg   *types.MsgUpdatePlaneName

		isError bool
	}{
		{
			name:  "not allowed",
			owner: sample.AccAddress().String(),
			msg: &types.MsgUpdatePlaneName{
				Supervisor: sample.AccAddress().String(),
				Id:         id,
				Name:       "new name",
			},
			isError: true,
		},
		{
			name:  "name is busy",
			owner: sample.AccAddress().String(),
			msg: &types.MsgUpdatePlaneName{
				Supervisor: supervisorAddr.String(),
				Id:         id,
				Name:       "busy",
			},
			isError: true,
		},
		{
			name:  "by owner",
			owner: addr.String(),
			msg: &types.MsgUpdatePlaneName{
				Supervisor: addr.String(),
				Id:         id,
				Name:       "new name",
			},
			isError: true,
		},
		{
			name:  "by supervisor",
			owner: supervisorAddr.String(),
			msg: &types.MsgUpdatePlaneName{
				Supervisor: supervisorAddr.String(),
				Id:         id,
				Name:       "new name",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			set, ctx := setupKeeper(t)
			set.scorumKeeper.SetParams(ctx.Context, scorumtypes.Params{
				Supervisors:                       []string{supervisorAddr.String()},
				GasLimit:                          sdk.IntProto{Int: sdk.NewInt(1000)},
				GasUnconditionedAmount:            sdk.IntProto{Int: sdk.NewInt(500)},
				GasAdjustCoefficient:              sdk.DecProto{Dec: sdk.NewDec(1)},
				RegistrationSPDelegationAmount:    sdk.IntProto{Int: sdk.NewInt(5)},
				SpWithdrawalTotalPeriods:          1,
				SpWithdrawalPeriodDurationSeconds: 1,
			})

			s := keeper.NewMsgServerImpl(set.keeper)
			_, err := s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
				Id:         id,
				Supervisor: supervisorAddr.String(),
				Owner:      addr.String(),
				Meta: &types.PlaneMeta{
					Name:       "name",
					Color:      "white",
					Experience: 10,
				},
			})
			require.NoError(t, err)
			_, err = s.CreatePlane(ctx.Context, &types.MsgCreatePlane{
				Id:         uuid.New().String(),
				Supervisor: supervisorAddr.String(),
				Owner:      sample.AccAddress().String(),
				Meta: &types.PlaneMeta{
					Name:       "busy",
					Color:      "white",
					Experience: 10,
				},
			})
			require.NoError(t, err)

			_, err = s.UpdatePlaneName(ctx.Context, tc.msg)

			if tc.isError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			act, err := set.keeper.GetPlane(ctx.Context, id)
			require.NoError(t, err)
			require.Equal(t, &types.Plane{
				Id:    id,
				Owner: addr.String(),
				Meta: &types.PlaneMeta{
					Name:       tc.msg.Name,
					Color:      "white",
					Experience: 10,
				},
			}, act)
		})
	}
}
