package fixture

type Query struct {
	query      string
	parameters []interface{}
}

func (query Query) GetQuery() string {
	return query.query
}

func (query Query) GetParameters() []interface{} {
	return query.parameters
}

func GenerateCustomQuery(query string, parameters ...interface{}) Query {
	return Query{query: query, parameters: parameters}
}

func MergeQueries(queries ...[]Query) []Query {
	var allQueries []Query
	for _, queryArray := range queries {
		allQueries = append(allQueries, queryArray...)
	}

	return allQueries
}
