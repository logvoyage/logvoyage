package elastic

import (
	"context"

	"bitbucket.org/firstrow/logvoyage/models"
	"bitbucket.org/firstrow/logvoyage/shared/config"
	"gopkg.in/olivere/elastic.v5"
)

// LogRecord fetched from storage
type LogRecord struct {
	Source   string
	Datetime int64
}

// SearchLogs sends query to elastic search index
func SearchLogs(user *models.User, project *models.Project, types []string, queryString string) ([]string, error) {
	ctx := context.Background()
	es, _ := elastic.NewClient(elastic.SetURL(config.Get("elastic.url")))

	q := buildQuery(queryString)
	s := es.Search().
		Index(project.IndexName()).
		Type(types...).
		Sort("_datetime", false).
		Query(q)

	searchResult, err := s.Do(ctx)

	if err != nil {
		// Index not found. That's ok, user didn't sent any data for now.
		if elastic.IsNotFound(err) {
			return []string{}, nil
		}
		return []string{}, err
	}

	if searchResult.Hits.TotalHits > 0 {
		var result = make([]string, searchResult.Hits.TotalHits)
		for i, hit := range searchResult.Hits.Hits {
			result[i] = string(*hit.Source)
		}
		return result, nil
	}

	return []string{}, nil
}

// If queryString is empty - return all records.
// Else, use query string.
func buildQuery(queryString string) elastic.Query {
	if len(queryString) == 0 {
		return elastic.NewMatchAllQuery()
	}
	return elastic.NewQueryStringQuery(queryString).DefaultField("msg")
}
