# CEP Weather Service

Este projeto consiste em dois serviços, A e B, para validação de CEP e obtenção de informações meteorológicas com base na localização do CEP. 

## Estrutura do Projeto

- **Serviço A**: Responsável por receber e validar o CEP.
- **Serviço B**: Responsável pela orquestração, validando o CEP, obtendo a localização e retornando informações meteorológicas formatadas.
- Foi utilizado a lib [Zipkin](https://zipkin.io/) para rastreamento de requisições. É uma ferramenta de tracing distribuido que ajuda a coletar, visualizar e analisar dados de rastreamento de solicitações em microsserviços.

## ⚙️ Configuração

Você precisará das seguintes tecnologias abaixo:

- [Docker](https://docs.docker.com/get-docker/) 🐳
- [Docker Compose](https://docs.docker.com/compose/install/) 🐳
- [GIT](https://git-scm.com/downloads)
- [Postman ☄️](https://www.postman.com/downloads/) ou [VS Code](https://code.visualstudio.com/download) com a
  extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) instalada.

## 🛠️ Requisitos

### Serviço A (Responsável pelo Input)
- Receber um input de 8 dígitos via POST no formato JSON: `{ "cep": "29902555" }`
- Validar se o input é válido (contém 8 dígitos e é uma STRING).
- Caso seja válido, encaminhar para o Serviço B via HTTP.
- Caso seja inválido, retornar:
    - Código HTTP: 422
    - Mensagem: `invalid zipcode`

### Serviço B (Responsável pela Orquestração)
- Receber um CEP válido de 8 dígitos.
- Realizar a pesquisa do CEP e encontrar o nome da localização.
- Retornar as temperaturas formatadas em Celsius, Fahrenheit e Kelvin juntamente com o nome da localização.
- Responder adequadamente nos seguintes cenários:
    - Em caso de sucesso:
        - Código HTTP: 200
        - Response Body: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }`
    - Em caso de falha, caso o CEP seja inválido (com formato correto):
        - Código HTTP: 422
        - Mensagem: `invalid zipcode`
    - Em caso de falha, caso o CEP não seja encontrado:
        - Código HTTP: 404
        - Mensagem: `can not find zipcode`

## 🚀 Iniciando

1. Clone o repositório:
    ```sh
    git clone https://github.com/brunoliveiradev/observability-otel-challenge-GoExpertPostGrad.git
    cd observability-otel-challenge-GoExpertPostGrad
    ```

2. Execute o comando abaixo na pasta raiz do projeto para iniciar o ambiente de desenvolvimento:
    ```sh
    docker-compose up -d
    ```
   
    Para parar os serviços:
    ```sh
    docker-compose down
    ```


## 🛠️ Endpoints

Veja abaixo os comportamentos de cada serviço.

### Serviço A

Você pode acessar o serviço A em `http://localhost:8080/` e enviar um cep valido no formato JSON. O arquivo `api/get_temperatures.http` contém exemplos de uso.

Comportamento:
- **POST** `/`
    - Request Body:
      ```json
      {
        "cep": "29902555"
      }
      ```
    - Responses:
        - 200: Encaminha para o Serviço B.
        - 422: `invalid zipcode` caso seja inválido.

### Serviço B

Você pode acessar o serviço B em `http://localhost:8081/{cep}`. O arquivo `api/get_temperatures.http` contém exemplos de uso.
- **GET** `/{cep}`
    - Responses:
        - 200: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }`
        - 404: `can not find zipcode` caso não encontre o CEP.
        - 422: `invalid zipcode` caso o CEP seja inválido.


### Zipkin
Para acessar a telemetria use o seguinte endereço do `zipkin` e após realizar uma requisição clique no botão "`RUN QUERY`":
- `http://localhost:9411/zipkin`


## 🧪 Testes

Após iniciar o ambiente de desenvolvimento, você pode testar com o cURL de exemplo abaixo ou com o arquivo `api/get_temperatures.http`:

```sh
curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -d '{"cep": "01001000"}'
```

```sh
curl http://localhost:8081/29902555
```