package tsdb

type TsdbClient interface {
	Query(start, end, rawTags, aggregator, metric string, is_relative bool) ([]Result, error)
}
