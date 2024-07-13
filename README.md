# labs-auction

## Inicie o mongodb e a aplicacao
```bash
docker compose up -d
```

no .env temos a variavel de ambiente `AUCTION_EXPIRED` que será responsável por definir o tempo de expiração de um leilão. O default é de 20s
também temos `FETCH_EXPIRED_INTERVAL` que é responsável pelo intervalo de tempo que a aplicação irá verificar se algum leilão expirou. O default é de 10s

## Execute os testes
```bash
docker-compose exec app go test ./internal/infra/database/auction -v
```