package keeper

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestDetectMaliciousSupply(t *testing.T) {
	input := CreateTestEnv(t)

	// set supply to maximum value
	var testBigInt big.Int
	testBigInt.SetBit(new(big.Int), 256, 1).Sub(&testBigInt, big.NewInt(1))
	bigCoinAmount := math.NewIntFromBigInt(&testBigInt)

	err := input.GravityKeeper.DetectMaliciousSupply(input.Context, "stake", bigCoinAmount)
	require.Error(t, err, "didn't error out on too much added supply")
}
