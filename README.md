# Go Legacy Integration Platform

Projeto de referência para modernização gradual de sistemas legados por meio de APIs REST, processamento assíncrono, auditoria e integração com armazenamento compatível com S3.

> O projeto utiliza nomes, dados e regras genéricas. Não contém código, credenciais ou informações internas de organizações reais.

## Objetivos

- Demonstrar desenvolvimento backend com Go.
- Expor uma API REST para documentos.
- Representar integração entre aplicações legadas e serviços modernos.
- Aplicar idempotência, auditoria e rastreabilidade.
- Preparar evolução para RabbitMQ, MinIO/S3, Redis e Oracle.
- Documentar decisões arquiteturais e estratégia de migração gradual.

## Arquitetura do MVP

```text
Aplicação legada
      |
      v
API REST em Go
      |
      +--> Serviço de documentos
      +--> Repositório em memória
      +--> Auditoria estruturada
      |
      v
Worker de processamento
```

O MVP executável usa repositório em memória para facilitar testes locais. As portas de infraestrutura permitem substituir essa implementação por Oracle, MinIO/S3 e mensageria sem alterar o domínio.

## Executar localmente

Pré-requisitos:

- Go 1.23 ou superior

```bash
go test ./...
go run ./cmd/api
```

A API estará disponível em:

```text
http://localhost:8080
```

### Endpoints

```text
GET    /health
POST   /api/v1/documents
GET    /api/v1/documents/{id}
DELETE /api/v1/documents/{id}
GET    /api/v1/documents/{id}/history
```

### Criar documento

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: carga-001" \
  -d '{"name":"documento.pdf","content_type":"application/pdf"}'
```

## Estrutura

```text
cmd/
  api/
  worker/
internal/
  application/
  config/
  domain/
  infrastructure/
  interfaces/
docs/
  wiki/
.github/
  workflows/
```

## Próximas evoluções

1. Persistência Oracle.
2. Armazenamento MinIO/S3.
3. RabbitMQ e worker assíncrono.
4. Redis para cache e idempotência distribuída.
5. JWT/OAuth.
6. OpenAPI/Swagger.
7. Métricas, tracing e logs estruturados.
8. Testes de integração com containers.
