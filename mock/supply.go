package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type SupplyKeeper interface {
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
}

type supplyKeeper struct {
	sendCoinsFromAccountToModule func(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	sendCoinsFromModuleToAccount func(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
}

func (s *supplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return s.sendCoinsFromAccountToModule(ctx, addr, moduleName, coins)
}

func (s *supplyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
	return s.sendCoinsFromModuleToAccount(ctx, moduleName, addr, coins)
}

type SupplyKeeperMock struct {
	s *supplyKeeper
}

func (s *SupplyKeeperMock) SetSendCoinsFromAccountToModule(f func(sdk.Context, sdk.AccAddress, string, sdk.Coins) error) {
	s.s.sendCoinsFromAccountToModule = f
}

func (s *SupplyKeeperMock) SetSendCoinsFromModuleToAccount(f func(sdk.Context, string, sdk.AccAddress, sdk.Coins) error) {
	s.s.sendCoinsFromModuleToAccount = f
}

func (s *SupplyKeeperMock) Mock() SupplyKeeper {
	return s.s
}

func (s *SupplyKeeperMock) WithDefaultsBalances(balances map[string]sdk.Coins) *SupplyKeeperMock {
	send := func(src, dest sdk.AccAddress, coins sdk.Coins) error {
		if balances[src.String()].IsAllGTE(coins) {
			return sdkerrors.ErrInsufficientFunds
		}
		balances[src.String()].Sub(coins)
		balances[dest.String()].Add(coins...)
		return nil
	}

	// set default
	s.SetSendCoinsFromAccountToModule(func(_ sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return send(addr, authtypes.NewModuleAddress(moduleName), coins)
	})

	s.SetSendCoinsFromModuleToAccount(func(_ sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
		return send(authtypes.NewModuleAddress(moduleName), addr, coins)
	})
	return s
}

func NewSupplyKeeper() *SupplyKeeperMock {
	balances := make(map[string]sdk.Coins)
	mock := &SupplyKeeperMock{s: &supplyKeeper{}}
	// set no-ops
	mock.SetSendCoinsFromAccountToModule(func(_ sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return nil
	})

	mock.SetSendCoinsFromModuleToAccount(func(_ sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
		return nil
	})
	return mock.WithDefaultsBalances(balances)
}
