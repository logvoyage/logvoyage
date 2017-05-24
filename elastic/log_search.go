package elastic

import (
	"context"
	"reflect"

	"github.com/logvoyage/logvoyage/models"
	"github.com/logvoyage/logvoyage/shared/config"
	"gopkg.in/olivere/elastic.v5"
)

// How many log records display per-page.
const pageSize = 10

// LogRecord loaded from elastic
type LogRecord struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

// SearchLogsResult contains logs and total number of records in storage.
type SearchLogsResult struct {
	Logs  []LogRecord `json:"logs"`
	Total int64       `json:"total"`
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
		var result = make([]LogRecord, len(searchResult.Hits.Hits))
		for i, hit := range searchResult.Hits.Hits {
			r := LogRecord{
				Source: string(*hit.Source),
				Type:   string(hit.Type),
			}
			result[i] = r
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

func GetIndexTypes(p *models.Project) []string {
	ctx := context.Background()
	es, _ := elastic.NewClient(elastic.SetURL(config.Get("elastic.url")))
	mapping, _ := es.GetMapping().Index(p.IndexName()).Do(ctx)

	keys := []string{}
	if indexMapping, ok := mapping[p.IndexName()]; ok {
		m := indexMapping.(map[string]interface{})
		for _, v := range reflect.ValueOf(m["mappings"]).MapKeys() {
			keys = append(keys, v.String())
		}
	}

	return keys
}
