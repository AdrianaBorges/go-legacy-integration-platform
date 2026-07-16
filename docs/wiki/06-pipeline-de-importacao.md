# Pipeline de importação

Fluxo proposto:

```text
Recebimento
  -> Validação
  -> Persistência de estágio
  -> Publicação de evento
  -> Processamento pelo worker
  -> Armazenamento
  -> Atualização de status
  -> Auditoria
```

Princípios:

- idempotência;
- retentativas controladas;
- dead-letter queue;
- correlation ID;
- rastreabilidade ponta a ponta.
