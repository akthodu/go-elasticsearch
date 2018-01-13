package elasticsearch

import "testing"

var aggregateClient *Client

func init() {
	var err error
	aggregateClient, err = Open("http://localhost:9200")
	if err != nil {
		panic(err)
	}
	if err := aggregateClient.Ping(); err != nil {
		panic(err)
	}
	aggregateClient.DeleteIndex("testclient_termaggregate")
}

func TestClient_TermAggregate(t *testing.T) {
	aggregateClient.InsertDocument("testclient_updatedocument", "doc", "1", map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}, true)
	aggregateClient.InsertDocument("testclient_updatedocument", "doc", "2", map[string]interface{}{
		"field1": "value1",
		"field2": "value3",
	}, true)
	aggregateClient.InsertDocument("testclient_updatedocument", "doc", "3", map[string]interface{}{
		"field1": "value1",
		"field2": "value4",
	}, true)
	result, err := aggregateClient.TermAggregate("testclient_updatedocument", "doc", nil, NewTermAggregations([]*TermAggregation{
		{Field: "field1", Size: 10},
		{Field: "field2", Size: 10},
	}))
	if err != nil {
		t.Fatalf("could not get aggregations: %s", err)
	}
	field1 := result["field1"].Buckets
	if len(field1) != 1 || field1[0].Key.(string) != "value1" || field1[0].Count != 3 {
		t.Fatalf("wrong field1 aggs: %#v", field1)
	}
	field2 := result["field2"].Buckets
	if len(field2) != 3 || field2[0].Count != 1 || field2[1].Count != 1 || field2[2].Count != 1 {
		t.Fatalf("wrong field2 aggs: %#v", field2)
	}
}
