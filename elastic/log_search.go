package elastic

import (
	"context"

	"bitbucket.org/firstrow/logvoyage/models"
	"bitbucket.org/firstrow/logvoyage/shared/config"
	"gopkg.in/olivere/elastic.v5"
)

// How many log records display per-page.
const pageSize = 10

// LogRecord fetched from storage
type LogRecord struct {
	Source   string
	Datetime int64
}

// SearchLogsResult contains logs and total number of records in storage.
type SearchLogsResult struct {
	Logs  []string `json:"logs"`
	Total int64    `json:"total"`
}

// SearchLogs sends query to elastic search index
// types - ealstic types to search on.
// queryString - user provided data.
func SearchLogs(user *models.User, project *models.Project, types []string, queryString string, page int) (SearchLogsResult, error) {
	ctx := context.Background()
	es, _ := elastic.NewClient(elastic.SetURL(config.Get("elastic.url")))

	q := buildQuery(queryString)
	s := es.Search().
		Index(project.IndexName()).
		Type(types...).
		From(page*pageSize).Size(pageSize).
		Sort("_datetime", false).
		Query(q)

	searchResult, err := s.Do(ctx)

	if err != nil {
		// Index not found. That's ok, user didn't sent any data for now.
		// Otherwise error should be handled.
		if elastic.IsNotFound(err) {
			return SearchLogsResult{}, nil
		}
		return SearchLogsResult{}, err
	}

	if searchResult.Hits.TotalHits > 0 {
		var result = make([]string, len(searchResult.Hits.Hits))
		for i, hit := range searchResult.Hits.Hits {
			result[i] = string(*hit.Source)
		}

		return SearchLogsResult{
			Logs:  result,
			Total: searchResult.Hits.TotalHits,
		}, nil
	}

	return SearchLogsResult{}, nil
}

// If queryString is empty - return all records.
// Else, use query string dsl.
func buildQuery(queryString string) elastic.Query {
	if len(queryString) == 0 {
		return elastic.NewMatchAllQuery()
	}
	return elastic.NewQueryStringQuery(queryString).DefaultField("msg")
}
