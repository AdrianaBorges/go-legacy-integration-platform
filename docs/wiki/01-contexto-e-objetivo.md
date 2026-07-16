# Contexto e objetivo

Sistemas corporativos legados frequentemente concentram regras de negócio críticas e não podem ser reescritos de uma só vez. Este projeto propõe uma modernização incremental, adicionando uma camada de serviços em Go entre as aplicações existentes, bancos relacionais e serviços externos.

## Objetivos

- Reduzir acoplamento.
- Criar contratos REST claros.
- Introduzir idempotência e auditoria.
- Permitir processamento assíncrono.
- Preparar integração com Oracle, MinIO/S3, RabbitMQ e Redis.
