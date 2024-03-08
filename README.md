# Rinha de Backend 2024 Q1 · [![ci](https://github.com/vitorsalgado/rinha-2024-q1-custom-lb/actions/workflows/ci.yml/badge.svg)](https://github.com/vitorsalgado/rinha-2024-q1-custom-lb/actions/workflows/ci.yml) · ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vitorsalgado/rinha-2024-q1-custom-lb) · ![GitHub License](https://img.shields.io/github/license/vitorsalgado/rinha-2024-q1-custom-lb)

Proposta de implementação da **[Rinha de Backend 2024 Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)** com um **Load Balancer** implementatdo em Go.  
Código do Load Balancer [aqui](./cmd/load_balancer).  
Os resultados dos testes são publicados automaticamente neste **[site](https://vitorsalgado.github.io/rinha-2024-q1-custom-lb/)**.

## Tech

- Go
- Postgres
- Load Balancer Próprio (Go)
- PgBouncer

## Sobre

Alguns pontos sobre o projeto:  

- api com std lib;
- **load balancer** implementado em Go. a implementação é bem simples: basicamente um proxy TCP com um round-robin simples. [código fonte aqui](cmd/load_balancer). 
- as operações de débito, crédito e extrato são feitas com apenas uma chamada ao banco, reduzindo o número de idas e vindas ao mesmo. No caso das operações de débito e crédito, foi utilizada uma function no Postgres que concentra a regra de negócio.
- PGO (Performance Guided Optimization) para gerar um binário mais eficiente.
- uso do componente PgBouncer para uma gestão mais eficiente de conexões com o banco.
- _Envoy_ como load balancer.
- GOMAXPROCS=1 (definido pela lib _automaxprocs_).

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
