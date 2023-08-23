package clob_test

import (
	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPlaceOrder_EquityTierLimit(t *testing.T) {
	tests := map[string]struct {
		firstOrder                   clobtypes.Order
		secondOrder                  clobtypes.Order
		equityTierLimitConfiguration clobtypes.EquityTierLimitConfiguration
		cancellation                 *clobtypes.MsgCancelOrder
		advanceBlock                 bool
		expectError                  bool
	}{
		"Short-term order would exceed max open short-term orders in same block": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			expectError: true,
		},
		"Long-term order would exceed max open stateful orders in same block": {
			firstOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			expectError: true,
		},
		"Conditional order would exceed max open stateful orders in same block": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			expectError: true,
		},
		"Short-term order would exceed max open short-term orders across blocks": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Long-term order would exceed max open stateful orders across blocks": {
			firstOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Conditional order would exceed max open stateful orders across blocks": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Order cancellation prevents exceeding max open short-term orders for short-term order in same block": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderShortTerm(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.GetGoodTilBlock(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation prevents exceeding max open stateful orders for long-term order in same block": {
			firstOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation prevents exceeding max open stateful orders for conditional order in same block": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation prevents exceeding max open short-term orders for short-term order across blocks": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderShortTerm(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.GetGoodTilBlock(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order cancellation prevents exceeding max open stateful orders for long-term order across blocks": {
			firstOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order cancellation prevents exceeding max open stateful orders for conditional order across blocks": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() types.GenesisDoc {
				genesis := testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
					state.Subaccounts = []satypes.Subaccount{
						constants.Alice_Num0_10_000USD,
					}
				})
				testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
					state.EquityTierLimitConfig = tc.equityTierLimitConfiguration
				})
				return genesis
			}).Build()

			ctx := tApp.InitChain()

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.firstOrder)) {
				require.True(t, tApp.CheckTx(tx).IsOK())
			}

			if tc.advanceBlock {
				ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			}

			if tc.cancellation != nil {
				for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *tc.cancellation) {
					require.True(t, tApp.CheckTx(tx).IsOK())
				}
			}

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.secondOrder)) {
				response := tApp.CheckTx(tx)
				if tc.expectError {
					require.False(t, response.IsOK())
					require.Contains(t, response.Log, "Opening order would exceed equity tier limit of 1.")
				} else {
					require.True(t, response.IsOK())
				}
			}

			// Ensure that any succesful transactions can be delivered.
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
		})
	}
}

func TestPlaceOrder_EquityTierLimit_OrderExpiry(t *testing.T) {
	tests := map[string]struct {
		firstOrder                   clobtypes.Order
		secondOrder                  clobtypes.Order
		equityTierLimitConfiguration clobtypes.EquityTierLimitConfiguration
		advanceToBlockAndTime        uint32
		expectError                  bool
	}{
		"Short-term order has not expired": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 14,
			expectError:           true,
		},
		"Short-term order has expired": {
			firstOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 15,
		},
		"Stateful order has not expired": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 4,
			expectError:           true,
		},
		"Stateful order has expired": {
			firstOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
				testapp.DefaultGenesis(),
			),
			secondOrder: MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() types.GenesisDoc {
				genesis := testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
					state.Subaccounts = []satypes.Subaccount{
						constants.Alice_Num0_10_000USD,
					}
				})
				testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
					state.EquityTierLimitConfig = tc.equityTierLimitConfiguration
				})
				return genesis
			}).Build()

			ctx := tApp.InitChain()

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.firstOrder)) {
				require.True(t, tApp.CheckTx(tx).IsOK())
			}

			ctx = tApp.AdvanceToBlock(
				tc.advanceToBlockAndTime,
				testapp.AdvanceToBlockOptions{BlockTime: time.Unix(int64(tc.advanceToBlockAndTime), 0)},
			)

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.secondOrder)) {
				response := tApp.CheckTx(tx)
				if tc.expectError {
					require.False(t, response.IsOK())
					require.Contains(t, response.Log, "Opening order would exceed equity tier limit of 1.")
				} else {
					require.True(t, response.IsOK())
				}
			}

			// Ensure that any succesful transactions can be delivered.
			tApp.AdvanceToBlock(lib.MustConvertIntegerToUint32(tApp.GetBlockHeight()+1), testapp.AdvanceToBlockOptions{})
		})
	}
}