package sources

import (
	"testing"

	"github.com/NibiruChain/nibiru/x/common/set"
	"github.com/NibiruChain/pricefeeder/types"
	"github.com/stretchr/testify/require"
)

func TestBitfinexSource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rawPrices, err := BitfinexPriceUpdate(set.New[types.Symbol]("tBTCUSD", "tETHUSD"))
		require.NoError(t, err)
		require.Equal(t, 2, len(rawPrices))
		require.NotZero(t, rawPrices["tBTCUSD"])
		require.NotZero(t, rawPrices["tETHUSD"])
	})

}
