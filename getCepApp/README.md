## Instruções de execução

### Dependências
- Docker

### Variáveis de ambiente
- Crie um arquivo .env no diretório root
- Siga o mesmo formato do .env.example
- Substitua seu token da Weather API 

### Comando para iniciar a aplicação
- `$ docker compose up`

### Requisição
- `$ curl -X POST localhost:8080 -d '{"zipcode":"01153000"}'`


### Comando para executar os testes
- `$ go test ./...`

### Requisição para o cloud run
- `$ curl -X POST https://fc-cloud-run-4wzrxu7gwq-uc.a.run.app -d '{"zipcode":"96030610"}'`