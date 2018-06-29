package main

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CloudWatch Response Object
type CloudWatchResponse struct {
	TimeStamp   string  `json:"timestamp"`
	MetricValue float64 `json:"metric_value"`
	Unit        string  `json:"unit"`
}

// Get Required metrics based on metric input and instance id
func getMetrics(metricName string, instanceID string, namspace string, unit string) []CloudWatchResponse {

	// Load session from shared config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)

	now := time.Now()
	then := now.Add(-24 * time.Hour) //duration
	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(now),          // Required
		MetricName: aws.String(metricName), // Required
		Namespace:  aws.String(namspace),   // Required
		Period:     aws.Int64(3600),        // Required //interval
		StartTime:  aws.Time(then),         // Required
		Statistics: []*string{ // Required
			aws.String("Maximum"), // Required
			// More values...
		},
		Dimensions: []*cloudwatch.Dimension{
			{ // Required
				Name:  aws.String("InstanceId"), // Required
				Value: aws.String(instanceID),   // Required
			},
			// More values...
		},
		Unit: aws.String(unit),
	}
	resp, err := svc.GetMetricStatistics(params)
	if err != nil {
		panic(err)
	}
	response := []CloudWatchResponse{}

	// Retrive the data points in the response
	for _, key := range resp.Datapoints {
		temp := key.Timestamp.UnixNano() / int64(time.Millisecond)
		// Append a Datapoint to the JSONArray
		slot := CloudWatchResponse{strconv.FormatInt(temp, 10), *key.Maximum, *key.Unit}
		response = append(response, slot)
	}
	return response
}
