package priceprovider

import (
	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/price-feeder/feeder/types"
	"github.com/rs/zerolog"
)

func NewAggregatePriceProvider(exchangesToPairToSymbolMap map[string]map[common.AssetPair]string, log zerolog.Logger) types.PriceProvider {
	log = log.With().Str("component", "aggregate-price-provider").Logger()
	providers := make([]types.PriceProvider, 0, len(exchangesToPairToSymbolMap))
	for exchangeName, pairToSymbolMap := range exchangesToPairToSymbolMap {
		providers = append(providers, NewPriceProvider(exchangeName, pairToSymbolMap, log.With().Str("exchange", exchangeName).Logger()))
	}
	return newAggregatePriceProvider(providers, log)
}

// NewAggregatePriceProvider instantiates a new AggregatePriceProvider instance
// given multiple PriceProvider.
func newAggregatePriceProvider(providers []types.PriceProvider, log zerolog.Logger) types.PriceProvider {
	a := AggregatePriceProvider{log, make(map[int]types.PriceProvider, len(providers))}
	for i, p := range providers {
		a.providers[i] = p
	}
	return a
}

// AggregatePriceProvider aggregates multiple price providers
// and queries them for prices.
type AggregatePriceProvider struct {
	log       zerolog.Logger
	providers map[int]types.PriceProvider // we use a map here to provide random ranging (since golang's map range is unordered)
}

// GetPrice fetches the first available and correct price from the wrapped PriceProviders.
// Iteration is exhaustive and random.
// If no correct PriceResponse is found, then an invalid PriceResponse is returned.
func (a AggregatePriceProvider) GetPrice(pair common.AssetPair) types.Price {
	// iterate randomly, if we find a valid price, we return it
	// otherwise we go onto the next PriceProvider to ask for prices.
	for _, p := range a.providers {
		price := p.GetPrice(pair)
		if price.Valid {
			return price
		}
	}

	// if we reach here no valid symbols were found
	a.log.Warn().Str("pair", pair.String()).Msg("no valid prices")
	return types.Price{
		Pair:   pair,
		Price:  0,
		Source: "aggregate",
		Valid:  false,
	}
}

func (a AggregatePriceProvider) Close() {
	for _, pp := range a.providers {
		pp.Close()
	}
}