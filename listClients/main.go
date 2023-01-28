package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type ClientesFiltro struct {
	Cnpj_cpf string `json:"cnpj_cpf"`
}

const (
	Yes = "S"
	No  = "N"
)

type OmieRequestListClient struct {
	Pagina                 int            `json:"pagina"`
	Registros_por_pagina   int            `json:"registros_por_pagina"`
	Apenas_importado_api   string         `json:"apenas_importado_api"`
	Exibir_caracteristicas string         `json:"exibir_caracteristicas"`
	ClientesFiltro         ClientesFiltro `json:"clientesFiltro"`
}

type OmieRequests struct {
	Call       string                  `json:"call"`
	App_key    string                  `json:"app_key"`
	App_secret string                  `json:"app_secret"`
	Param      []OmieRequestListClient `json:"param"`
}

type RequestsBodyData struct {
	Call       string `json:"call"`
	App_key    string `json:"app_key"`
	App_secret string `json:"app_secret"`
	Cnpj_cpf   string `json:"cnpj_cpf"`
}

// {
// 	"call": "ListarClientes",
// 	"app_key": "2131745201586",
// 	"app_secret": "541ca5d21847e02e7142fc35db84f0d1",
// 	"param": [
// 		{
// 			"pagina": 1,
// 			"registros_por_pagina": 2,
// 			"apenas_importado_api": "N",
// 			"exibir_caracteristicas": "S",
// 			"clientesFiltro": {
// 				"cnpj_cpf": "023.082.445-50"
// 			}
// 		}
// 	]
// }

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	paramsBody := request.Body
	var requestBodyMessage RequestsBodyData

	err := json.Unmarshal([]byte(paramsBody), &requestBodyMessage)
	if err != nil {
		panic("Campos inválidos")
	}
	fmt.Println("requestBodyMessage", requestBodyMessage)
	var buf bytes.Buffer
	// request.Body
	json_request_list_client := &OmieRequests{
		Call:       requestBodyMessage.Call,
		App_key:    requestBodyMessage.App_key,
		App_secret: requestBodyMessage.App_secret,
		Param: []OmieRequestListClient{
			{
				Pagina:                 1,
				Registros_por_pagina:   2,
				Apenas_importado_api:   No,
				Exibir_caracteristicas: Yes,
				ClientesFiltro: ClientesFiltro{
					Cnpj_cpf: requestBodyMessage.Cnpj_cpf,
				},
			},
		},
	}
	json_data, _ := json.Marshal(json_request_list_client)
	r, requestError := http.Post(
		"https://app.omie.com.br/api/v1/geral/clientes/",
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if requestError != nil {
		return Response{StatusCode: 404}, requestError
	}
	defer r.Body.Close()

	readedBody, _ := io.ReadAll(r.Body)
	var res map[string]interface{}
	json.Unmarshal(readedBody, &res)

	body, err := json.Marshal(&res)
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
