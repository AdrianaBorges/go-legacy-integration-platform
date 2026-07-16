# Arquitetura proposta

```text
Aplicação legada
      |
      v
API Go
      |
      +--> Domínio
      +--> Serviços de aplicação
      +--> Portas de infraestrutura
              |
              +--> Oracle
              +--> MinIO/S3
              +--> RabbitMQ
              +--> Redis
```

O domínio permanece independente das tecnologias externas. Essa separação facilita testes, manutenção e substituição de componentes.
