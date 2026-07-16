# API de documentos

Responsabilidades:

- registrar metadados;
- garantir idempotência;
- consultar documento;
- excluir logicamente;
- manter histórico;
- futuramente integrar upload e download com MinIO/S3.

## Contratos iniciais

```text
POST   /api/v1/documents
GET    /api/v1/documents/{id}
DELETE /api/v1/documents/{id}
GET    /api/v1/documents/{id}/history
```
