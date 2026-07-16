# Arquitetura atual

Exemplo genérico:

```text
Aplicação legada -> Banco relacional -> Serviço externo
```

Características comuns:

- forte acoplamento;
- chamadas síncronas;
- regras distribuídas entre aplicação e banco;
- baixa observabilidade;
- evolução com alto risco.
