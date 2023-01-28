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

type ClientesFiltro struct {
	Cnpj_cpf string `json:"cnpj_cpf"`
}

type OmieRequestListClient struct {
	Pagina                 int            `json:"pagina"`
	Registros_por_pagina   int            `json:"registros_por_pagina"`
	Apenas_importado_api   string         `json:"apenas_importado_api"`
	Exibir_caracteristicas string         `json:"exibir_caracteristicas"`
	ClientesFiltro         ClientesFiltro `json:"clientesFiltro"`
}

const (
	Yes = "S"
	No  = "N"
)

type OmieRequests[T interface{}] struct {
	Call       string `json:"call"`
	App_key    string `json:"app_key"`
	App_secret string `json:"app_secret"`
	Param      []T    `json:"param"`
}

type RequestsBodyData struct {
	Call       string `json:"call"`
	App_key    string `json:"app_key"`
	App_secret string `json:"app_secret"`
	Cnpj_cpf   string `json:"cnpj_cpf"`
}

type ListClientsResponse struct {
	ClientesCadastro []struct {
		Bairro              string `json:"bairro"`
		BloquearFaturamento string `json:"bloquear_faturamento"`
		Caracteristicas     []struct {
			Campo    string `json:"campo"`
			Conteudo string `json:"conteudo"`
		} `json:"caracteristicas"`
		Cep                     string `json:"cep"`
		Cidade                  string `json:"cidade"`
		CidadeIbge              string `json:"cidade_ibge"`
		CnpjCpf                 string `json:"cnpj_cpf"`
		CodigoClienteIntegracao string `json:"codigo_cliente_integracao"`
		CodigoClienteOmie       int64  `json:"codigo_cliente_omie"`
		CodigoPais              string `json:"codigo_pais"`
		Complemento             string `json:"complemento"`
		DadosBancarios          struct {
			Agencia       string `json:"agencia"`
			CodigoBanco   string `json:"codigo_banco"`
			ContaCorrente string `json:"conta_corrente"`
			DocTitular    string `json:"doc_titular"`
			NomeTitular   string `json:"nome_titular"`
			TransfPadrao  string `json:"transf_padrao"`
		} `json:"dadosBancarios"`
		Email           string `json:"email"`
		Endereco        string `json:"endereco"`
		EnderecoEntrega struct {
		} `json:"enderecoEntrega"`
		EnderecoNumero string `json:"endereco_numero"`
		Estado         string `json:"estado"`
		Exterior       string `json:"exterior"`
		Inativo        string `json:"inativo"`
		Info           struct {
			CImpAPI string `json:"cImpAPI"`
			DAlt    string `json:"dAlt"`
			DInc    string `json:"dInc"`
			HAlt    string `json:"hAlt"`
			HInc    string `json:"hInc"`
			UAlt    string `json:"uAlt"`
			UInc    string `json:"uInc"`
		} `json:"info"`
		InscricaoEstadual  string `json:"inscricao_estadual"`
		InscricaoMunicipal string `json:"inscricao_municipal"`
		NomeFantasia       string `json:"nome_fantasia"`
		PessoaFisica       string `json:"pessoa_fisica"`
		RazaoSocial        string `json:"razao_social"`
		Recomendacoes      struct {
			GerarBoletos string `json:"gerar_boletos"`
		} `json:"recomendacoes"`
		Tags []struct {
			Tag string `json:"tag"`
		} `json:"tags"`
		Telefone1Ddd    string `json:"telefone1_ddd"`
		Telefone1Numero string `json:"telefone1_numero"`
	} `json:"clientes_cadastro"`
	Pagina           int `json:"pagina"`
	Registros        int `json:"registros"`
	TotalDePaginas   int `json:"total_de_paginas"`
	TotalDeRegistros int `json:"total_de_registros"`
}

func FindUserByCpf(requestBodyMessage RequestsBodyData) (int, error, *http.Response) {
	omieRequestUrl := "https://app.omie.com.br/api/v1/geral/clientes/"
	json_request_list_client := &OmieRequests[OmieRequestListClient]{
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
		omieRequestUrl,
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if requestError != nil {
		return 404, requestError, nil
	}
	defer r.Body.Close()
	return 200, nil, r
}

func ParseListClientResponse(r *http.Response) []byte {
	readedBody, _ := io.ReadAll(r.Body)
	var listClientResponse map[string]ListClientsResponse
	json.Unmarshal(readedBody, &listClientResponse)
	body, _ := json.Marshal(&listClientResponse)
	return body
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (LambdaResponse, error) {
	paramsBody := request.Body
	var requestBodyMessage RequestsBodyData
	json.Unmarshal([]byte(paramsBody), &requestBodyMessage)
	var buf bytes.Buffer

	statusCode, errorResp, r := FindUserByCpf(requestBodyMessage)
	if statusCode == 404 {
		return LambdaResponse{StatusCode: 404}, errorResp
	}

	body := ParseListClientResponse(r)
	json.HTMLEscape(&buf, body)
	resp := LambdaResponse{StatusCode: 200, Body: buf.String()}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
