package enricher

import (
	"github.com/mrizkifadil26/medix/enricher/local"
	"github.com/mrizkifadil26/medix/utils/datawrapper"
)

func Enrich(
	data any,
	config *Config,
) (any, error) {
	// switch config.Provider {
	// case "tmdb":
	// 	return tmdb.Enrich(data, config)
	// default:
	// 	return nil, nil
	// }

	wrappedData := datawrapper.WrapData(data)
	// params := map[string]string{
	// 	"concurrency": "4",
	// }

	// enricher := &tmdb.TMDbEnricher{}
	// data, err := enricher.Enrich(wrappedData, params)
	// if err != nil {
	// 	return nil, err
	// }

	// mediaSourceEnricher := &local.MediaSourceEnricher{}
	// data, err := mediaSourceEnricher.Enrich(wrappedData, nil)
	// if err != nil {
	// 	return nil, err
	// }

	subtitleEnricher := &local.SubtitleEnricher{}
	data, err := subtitleEnricher.Enrich(wrappedData, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
