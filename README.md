<h1>Open telemetry challenge</h1>

<h3>Intruções de execução</h3>

  - Clonar repositório
  - Configurar token para weather api criando ./getCepApp/.env com base no getCepApp/.env.example
  - Executar `docker compose up`

<h3>Intruções de uso</h3>
  
  - Realizar requisições para localhost"8081 
    `curl -X POST http://localhost:8081 -d '{"cep": "96065710" }'`
  - É possível visualizar o trace distribuido entre os serviços utilizando
    - Zipkin: <a>http://localhost:9411 </a>
    - Jaeger: <a>http://localhost:16686 </a>

