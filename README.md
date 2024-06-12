# Rate Limiter

## Executando a aplicação

Primeiramente deve criiar um arquivo .env tendo como exemplo o arquivo .env.example dentro da pasta cmd;
Depois pode executar o comando para subir as aplicações utilizando docker:
docker compose up --build

## Testando a aplicação

Para testar a aplicação basta que faça as requisições para a rota determinada
[http://localhost:8080]

## Testes automatizados

Para testar a aplicação primeiramente deve ingresar na pasta cmd;, depois executar o comando:
go test ./...
