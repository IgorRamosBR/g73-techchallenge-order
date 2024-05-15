# g73-techchallenge-order

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

Este é um microsserviço responsável por criar produtos e gerenciar pedidos online para uma lanchonete. Ele oferece endpoints para criar, atualizar, buscar e listar pedidos, além de interagir com produtos e processar pagamentos.


## Tecnologias Utilizadas

- Linguagem de Programação: Go
- Banco de Dados: PostgreSQL
- Docker



## Requisitos

- go 1.20
- docker
- kubernetes cluster (Docker desktop)
- kubectl



## Funcionalidades
- **Criar Pedido:** Permite que os clientes criem novos pedidos adicionando produtos ao carrinho.
- **Atualizar Pedido:** Permite atualizar o status dos pedidos.
- **Listar Pedidos:** Fornece uma lista de todos os pedidos, filtrada por diferentes critérios.
- **Buscar Pedido por ID:** Permite buscar um pedido específico pelo seu ID.
- **Gerar QR Code de Pagamento:** Gera um QR code para processar pagamentos dos pedidos.



## Como Executar
Para executar este microsserviço, siga estas etapas:

**1.** Certifique-se de ter o Docker instalado em sua máquina.

**2.** Clone este repositório para o seu ambiente local.

**3.** Navegue até o diretório do projeto.

**4.** Execute o seguinte comando para iniciar o contêiner Docker do PostgreSQL:

```bash
docker-compose up -d
```

**5.** Depois que o PostgreSQL estiver em execução, construa a imagem Docker do microsserviço:

```bash
docker build -t order-service .
```

**6.** Por fim, execute o contêiner do microsserviço:

```bash
docker run -p 8080:8080 order-service
```

Agora que o microsserviço está em execução, você pode acessar os endpoints conforme documentado abaixo.



## Endpoints

### Criar pedido

```bash
POST /orders
```

Cria um novo pedido com base nos dados fornecidos.

**Parâmetros**
- **customerCpf:** CPF do cliente.
- **items:** Lista de itens do pedido.

### Atualizar Status do Pedido

```bash
PUT /orders/{orderId}
```

Atualiza o status de um pedido existente.

**Parâmetros**
- **orderId:** ID do pedido.
- **status:** Novo status do pedido.

### Listar pedidos

```bash
GET /orders
```

Retorna uma lista de todos os pedidos.

**Parâmetros**
- **limit:** (Opcional) Número máximo de resultados a serem retornados.
- **offset:** (Opcional) Deslocamento para paginar os resultados.

### Buscar Pedido por ID

```bash
GET /orders/{orderId}
```

Retorna detalhes de um pedido específico com base no ID fornecido.

**Parâmetros**
**orderId:** ID do pedido.

### Gerar QR Code de Pagamento

```bash
POST /payments/qrcode
```

Gera um QR code para pagamento do pedido especificado.

**Parâmetros**
**orderId:** ID do pedido.


## Execução com Kubernetes

Entrar na pasta do Kubernetes
```bash
  cd k8s
```

Criar Persistent Volume
```bash
  kubectl apply -f pv-volume.yaml
```

Criar Persistent Volume Claim
```bash
  kubectl apply -f pv-claim.yaml
```

Criar Postgres Config Map
```bash
  kubectl apply -f postgres-config-map.yaml
```

Criar Postgres Service
```bash
  kubectl apply -f postgres-service.yaml
```

Criar Postgres Deployment
```bash
  kubectl apply -f postgres-deployment.yaml
```

Criar API Service
```bash
  kubectl apply -f api-service.yaml
```

Criar API Deployment
```bash
  kubectl apply -f api-deployment.yaml
```

## Documentação
[Documentation](https://github.com/IgorRamosBR/g73-techchallenge-order/tree/master/docs)


## Arquitetura
Clean Architecture com a estrutura de pastas baseada no [Standard Go Project Layout](https://github.com/golang-standards/project-layout#go-directories) 

```bash
├── cmd
├── configs
├── docs
├── internal
|   |── api
|   |── controllers
|   ├── core
|   │   ├── entities
|   │   ├── usecases
|   ├── infra
|   │   ├── drivers
|   │   ├── gateways
├── k8s
├── migrations
```

