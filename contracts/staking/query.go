package staking

import (
	"errors"
	"math/big"

	"github.com/0xPolygon/polygon-sdk/contracts/abis"
	"github.com/0xPolygon/polygon-sdk/state"
	"github.com/0xPolygon/polygon-sdk/types"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
)

var (
	// staking contract address
	AddrStakingContract = types.StringToAddress("1001")
)

func DecodeValidators(method *abi.Method, returnValue []byte) ([]types.Address, error) {
	decodedResults, err := method.Outputs.Decode(returnValue)
	if err != nil {
		return nil, err
	}
	results, ok := decodedResults.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed type assertion from decodedResults to map")
	}
	web3Addresses, ok := results["0"].([]web3.Address)
	if !ok {
		return nil, errors.New("failed type assertion from results[0] to []web3.Address")
	}

	addresses := make([]types.Address, len(web3Addresses))
	for idx, waddr := range web3Addresses {
		addresses[idx] = types.Address(waddr)
	}

	return addresses, nil
}

func QueryValidators(t *state.Transition, from types.Address) ([]types.Address, error) {
	method, ok := abis.StakingABI.Methods["validators"]
	if !ok {
		return nil, errors.New("validators method doesn't exist in Staking contract ABI")
	}

	selector := method.ID()
	res, err := t.Apply(&types.Transaction{
		From:     from,
		To:       &AddrStakingContract,
		Value:    big.NewInt(0),
		Input:    selector[:],
		GasPrice: big.NewInt(0),
		Gas:      100000000,
		Nonce:    t.GetNonce(from),
	})
	if err != nil {
		return nil, err
	}
	if res.Failed() {
		return nil, res.Err
	}

	return DecodeValidators(method, res.ReturnValue)
}