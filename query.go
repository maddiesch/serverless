package serverless

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// QueryPaginatedBuilder is the function called by QueryPaginated to generate a query
type QueryPaginatedBuilder = func() (*dynamodb.QueryInput, error)

// QueryPaginatedOutput The results of the paginated query
type QueryPaginatedOutput struct {
	// Records contains the slice of results
	Records []map[string]*dynamodb.AttributeValue

	// RequestCount is the total number of requests made to DynamoDB
	RequestCount uint64
}

// QueryPaginated perform a paginated dynamodb query
//
// This will repeatedly call the handler function until there are no more results from DynamoDB. If there are multiple
// results they will all be collected into a single slice and returned.
//
// If the handler returns an error, no results will be returned and the error will be passed through.
//
// If the handler returns nil for the query and no error, the loop will be aborted the currently collected results will
// be returned.
func QueryPaginated(db *DB, handler QueryPaginatedBuilder) (*QueryPaginatedOutput, error) {
	results := []map[string]*dynamodb.AttributeValue{}
	exclusiveStartKey := map[string]*dynamodb.AttributeValue{}
	var count uint64

	for {
		query, err := handler()
		if err != nil {
			return nil, err
		}
		if query == nil {
			break
		}

		if query.TableName == nil {
			query.SetTableName(db.TableName)
		}

		if len(exclusiveStartKey) > 0 {
			query.SetExclusiveStartKey(exclusiveStartKey)
		}

		count++

		result, err := db.Client.Query(query)
		if err != nil {
			return nil, err
		}

		for _, res := range result.Items {
			results = append(results, res)
		}

		if len(result.LastEvaluatedKey) == 0 {
			break
		}

		exclusiveStartKey = result.LastEvaluatedKey
	}

	return &QueryPaginatedOutput{Records: results, RequestCount: count}, nil
}
