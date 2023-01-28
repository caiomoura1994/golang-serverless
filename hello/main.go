package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer
	r, requestError := http.Get("https://echoof.me/json")
	if requestError != nil {
		return Response{StatusCode: 404}, requestError
	}
	defer r.Body.Close()
	readedBody, _ := io.ReadAll(r.Body)
	var res map[string]interface{}
	json.Unmarshal(readedBody, &res)
	// print(res["ip"].(string))

	body, err := json.Marshal(map[string]interface{}{
		"message": res["ip"],
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: 200,
		Body:       buf.String(),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
