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

type LambdaResponse events.APIGatewayProxyResponse

type Caracteristicas struct {
	Campo    string `json:"campo"`
	Conteudo string `json:"conteudo"`
}
type OmieRequestUpSertClient struct {
	CodigoClienteIntegracao string            `json:"codigo_cliente_integracao"`
	Email                   string            `json:"email"`
	Telefone1Ddd            string            `json:"telefone1_ddd"`
	Telefone1Numero         string            `json:"telefone1_numero"`
	RazaoSocial             string            `json:"razao_social"`
	NomeFantasia            string            `json:"nome_fantasia"`
	Caracteristicas         []Caracteristicas `json:"caracteristicas"`
	Tags                    []string          `json:"tags"`
}

type OmieRequests[T interface{}] struct {
	Call       string `json:"call"`
	App_key    string `json:"app_key"`
	App_secret string `json:"app_secret"`
	Param      []T    `json:"param"`
}

type RequestsBodyData struct {
	Call                    string            `json:"call"`
	App_key                 string            `json:"app_key"`
	App_secret              string            `json:"app_secret"`
	CodigoClienteIntegracao string            `json:"codigo_cliente_integracao"`
	Email                   string            `json:"email"`
	Telefone1Ddd            string            `json:"telefone1_ddd"`
	Telefone1Numero         string            `json:"telefone1_numero"`
	RazaoSocial             string            `json:"razao_social"`
	NomeFantasia            string            `json:"nome_fantasia"`
	Caracteristicas         []Caracteristicas `json:"caracteristicas"`
}

type UpsertClientResponse struct {
	CodigoClienteOmie       int64  `json:"codigo_cliente_omie"`
	CodigoClienteIntegracao string `json:"codigo_cliente_integracao"`
	CodigoStatus            string `json:"codigo_status"`
	DescricaoStatus         string `json:"descricao_status"`
}

func UpSertClient(requestBodyMessage RequestsBodyData) *http.Response {
	omieRequestUrl := "https://app.omie.com.br/api/v1/geral/clientes/"
	json_request_create_client := &OmieRequests[OmieRequestUpSertClient]{
		Call:       requestBodyMessage.Call,
		App_key:    requestBodyMessage.App_key,
		App_secret: requestBodyMessage.App_secret,
		Param: []OmieRequestUpSertClient{
			{
				CodigoClienteIntegracao: requestBodyMessage.CodigoClienteIntegracao,
				Email:                   requestBodyMessage.Email,
				Telefone1Ddd:            requestBodyMessage.Telefone1Ddd,
				Telefone1Numero:         requestBodyMessage.Telefone1Numero,
				RazaoSocial:             requestBodyMessage.RazaoSocial,
				NomeFantasia:            requestBodyMessage.NomeFantasia,
				Caracteristicas:         requestBodyMessage.Caracteristicas,
			},
		},
	}
	json_data, _ := json.Marshal(json_request_create_client)
	r, _ := http.Post(
		omieRequestUrl,
		"application/json",
		bytes.NewBuffer(json_data),
	)
	return r
}

func ParseUpSertResponse(r *http.Response) ([]byte, UpsertClientResponse) {
	readedBody, _ := io.ReadAll(r.Body)
	var res UpsertClientResponse
	json.Unmarshal(readedBody, &res)
	body, _ := json.Marshal(&res)
	return body, res
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (LambdaResponse, error) {
	paramsBody := request.Body
	var requestBodyMessage RequestsBodyData

	err := json.Unmarshal([]byte(paramsBody), &requestBodyMessage)
	if err != nil {
		panic("Campos inválidos")
	}
	var buf bytes.Buffer
	upSertClientResponse := UpSertClient(requestBodyMessage)
	defer upSertClientResponse.Body.Close()
	body, _ := ParseUpSertResponse(upSertClientResponse)
	json.HTMLEscape(&buf, body)
	resp := LambdaResponse{
		StatusCode: 200,
		Body:       buf.String(),
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
