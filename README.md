# Rinha de Backend 2024 Q1 · [![ci](https://github.com/vitorsalgado/rinha-2024-q1-custom-lb/actions/workflows/ci.yml/badge.svg)](https://github.com/vitorsalgado/rinha-2024-q1-custom-lb/actions/workflows/ci.yml) · ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vitorsalgado/rinha-2024-q1-custom-lb) · ![GitHub License](https://img.shields.io/github/license/vitorsalgado/rinha-2024-q1-custom-lb)

Proposta de implementação da **[Rinha de Backend 2024 Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)** com um **Load Balancer** implementatdo em Go.  
Código do Load Balancer [aqui](./cmd/load_balancer).  
Os resultados dos testes são publicados automaticamente neste **[site](https://vitorsalgado.github.io/rinha-2024-q1-custom-lb/)**.  
Submissão: [aqui](https://github.com/zanfranceschi/rinha-de-backend-2024-q1/tree/main/participantes/vitorsalgado-custom-lb).

## Tech

- Go
- Postgres
- Load Balancer Próprio (Go)
- PgBouncer

## Sobre

A idéia era criar um projeto bem simples, com o mínimo possível de libs e frameworks e que também fosse fácil de replicar em outras linguagens.  
Em relação a **performance**, aqui algumas idéias que guiaram o projeto, mais ou menos em uma ordem de prioridade:  

- menos **round trips** possíveis ao banco de dados. para isso, usei uma **function** no Postgres para as transações e uma query única para o obter o extrato bancário.

- gestão eficiente de conexões com o banco. esse ponto é um complemente do anterior, conexões com o banco de dados são "caras" e aqui demorei para achar o setup ideal. desde o início a solução contava com um **pool** de conexões e no começo esse pool girou em torno de ~100 - ~300 de máx. conexões. depois de vários experimentos, encontrei uma ferramente interessante para o pool de conexões, **PgBouncer**. com o PgBouncer integrado, o setup ideal acabou sendo: __pool=5__ nas apis e __pool=20__ no PgBouncer, um número muito menor do que os experimentos inicias sem esse componente.

- essa submissão em específico implementa um load balancer próprio, implementado em Go também. o load balancer é bastante simpler, basicamente um proxy TCP com um round robin simples que distribui o tráfico entre as apis igualmente independente do estado delas. [código fonte aqui](cmd/load_balancer). 

- **threads**: dadas as limitações do ambiente em relação a CPU e memória, experimentei diferentes setups de threads para as aplicações. no final, os testes se sairam melhor configurando tudo com o mínimo de processos possível. as apis e o load balancer definem **GOMAXPROCS=1**. no final, usei a lib _automaxprocs_ para a definição adequada e automática da variável GOMAXPROCS para evitar surpresas no ambiente de testes da rinha.

- com a solução final pronta, utilizei um recurso chamado **PGO (Performance Guided Optimization)** para gerar binários mais eficientes das aplicações em Go. basicamente rodei um profiling das apis e do load balancer também, gerando no final um arquivo **default.pgo**. esse arquivo é então submetido ao **build** posteriormente.

## Executando

Para executar o projeto completo em um **docker compose** local, execute no seu terminal:
```
make up
```

## Testes de Carga

Para executar os testes de carga contidos no repositório original da rinha, 
primeiro execute o comando de preparação:
```
make prepare
```

O comando `make prepare` clona o repositório da rinha e instala a ferramente Gatling.  
**Ele deve ser executado apenas uma vez.**  
Para rodar os testes, execute o comando:
```
make test
```
