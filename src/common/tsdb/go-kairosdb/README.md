#go-kairosdb
A go client for [KairosDB](http://kairosdb.github.io/).

## Introduction
*go-kairosdb* is a package that aims to provide wrapper APIs around the REST APIs exposed by KairosDB.
These APIs can be used to store the metrics collected and query them back in different ways.

The code structure/APIs are based on the original KairosDB client that has been implemented in Java.
More info about the original client can be found [here](https://github.com/kairosdb/kairosdb-client).

Like the Java based client, *go-kairosdb* also employs a builder pattern to construct the JSON objects
that are sent to KairosDB using an http client.

NOTE: THIS IS STILL WORK IN PROGRESS.

## Getting go-kairosdb
*go-kairosdb* is written in golang and hence it's assumed that the user has setup the golang programming environment
correctly. In order to get access to the package the following commands need to be executed:

```
go get github.com/ajityagaty/go-kairosdb
go get github.com/stretchr/testify/assert
```

[assert](https://github.com/stretchr/testify) is a package that's available under testify project. It
provides a very friendly API to perform validations. This package has been employed in the unit tests.

## Usage
The following sections describe the way the APIs can be used.

### Sending Metrics
The MetricBuilder interface is primarily used to stitch the metrics together and send them to KairosDB.
One can add metrics, their associated tags and the data points using the builder.

```
// Instantiate the MetricBuilder
mb := NewMetricBuilder()

// Add a metric along with tags and datapoints.
mb.AddMetric("m1").
	AddDataPoint(1234, int64(1000)).
	AddDataPoint(1235, int64(304)).
	AddTag("t1", "v1").
	AddTag("t2", "v2")

// Add another metric with its data points.
mb.AddMetric("m2").
	AddDataPoint(1236, int64(320)).
	AddDataPoint(1237, 201.3).
	AddTag("t3", "v3")

// Get an instance of the http client
cli := client.NewHttpClient("http://localhost:1234")
pushResp, _ := cli.PushMetrics(mb)
```

### Querying Metrics
The QueryBuilder is used to build the query. Every query requires a date range wherein the start date
is mandatory while the end date defaults to NOW. A specific metric can be queried for by specifying the
metric's name and tags can be added to narrow down the search.

```
// Instantiate a QueryBuilder
qb := builder.NewQueryBuilder()

// Set a relative start time of 4 years and specify the metric name.
// Set Limit can be used to specify how many data points need to be returned.
qb.SetRelativeStart(2, utils.HOURS).
	AddMetric("m1").
	SetLimit(100)

// Add another metric with Absolute start time set to 3 days ago.
qb.SetAbsoluteStart(time.Now().AddDate(0, 0, -3)).
	AddMetric("m2").
    SetLimit(100)

// Get an instance of the http client
cli := client.NewHttpClient("http://localhost:1234")
queryResp, _ := cli.Query(qb)
```

### Query Metric Names
One can get a list of all the metric names stored in KairosDB.

```
// Get an instance of the client.
cli := client.NewHttpClient("http://localhost:1234")

// Get all the metric names.
resp, _ := cli.GetMetricNames()

// Print all the metrics.
for _, metric := range resp.Results {
	fmt.Println(metric)
}
```

### Query Tag Names and Values
Similarly one can get a list of all tag names and values stored in KairosDB.

```
// Get an instance of the client.
cli := client.NewHttpClient("http://localhost:1234")

// Get all the tag names.
tagNamesResp, _ := cli.GetTagNames()

// Get all the tag values.
tagValuesResp, _ := cli.GetTagValues()

// Print all the tag names.
for _, tagName := range resp.Results {
	fmt.Println(tagName)
}

// Print all the tag values.
for _, tagVal := range resp.Results {
	fmt.Println(tagVal)
}
```

### Delete Metric
One can delete a metric and all its associated data points from KairosDB.
On success - *StatusNoContent* is returned.
On failure - *StatusBadRequest* or *StatusInternalServerError* is returned.

```
// Get an instance of the client.
cli := client.NewHttpClient("http://localhost:1234")

// Delete a metric.
delResp, _ := cli.DeleteMetric("m1")

if delResp.GetStatusCode() == http.StatusNoContent {
	fmt.Println("Delete Metric succeeded")
} else {
	fmt.Println("Delete Metric failed")
}
```

### Health Check
One can also query the health of the KairosDB server. An HTTP *StatusNoContent*
return code indicates that all is well. If there is any problem then an HTTP
*StatusInternalServerError* will be returned.

```
// Get an instance of the client.
cli := client.NewHttpClient("http://localhost:1234")

healthResp, _ := cli.HealthCheck()

if healthResp.GetStatusCode() == http.StatusNoContent {
	fmt.Println("All is well")
} else {
	fmt.Println("Internal error")
}

```