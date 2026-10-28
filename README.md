# 💰 Finance Control API

API para gerenciamento financeiro pessoal — áreas, categorias, orçamentos e transações. Backend em Go com Echo v5 e PostgreSQL, servido via Docker Compose.

---

## 🛠️ Stack

- **Linguagem:** Go 1.27
- **Framework HTTP:** [Echo v5](https://echo.labstack.com/)
- **Banco de dados:** PostgreSQL 18
- **Migrações:** [Goose](https://github.com/pressly/goose)
- **Autenticação:** JWT (HS256, 24h de validade)
- **Validação:** [go-playground/validator](https://github.com/go-playground/validator)
- **Swagger:** [swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)
- **Testes:** [testcontainers-go](https://github.com/testcontainers/testcontainers-go) (PostgreSQL efêmero)

---

## 📋 Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose
- [Go](https://go.dev/dl/) 1.27+ (apenas para rodar testes localmente)

---

## ⚙️ Configuração

Copie o arquivo de exemplo e ajuste as variáveis:

```bash
cp .env.exemple .env
```

| Variável | Descrição | Valor padrão |
|---|---|---|
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_USER` | Usuário do PostgreSQL | — |
| `DB_PASSWORD` | Senha do PostgreSQL | — |
| `DB_NAME` | Nome do banco de dados | — |
| `JWT_SECRET` | Chave de assinatura JWT | — |
| `ENV` | Ambiente (`development` ou `production`) | `development` |
| `FRONTEND_URL` | URL do frontend (obrigatório em `production`) | — |

---

## 🐳 Rodando com Docker

```bash
docker compose up --build
```

O comando inicia três serviços na ordem:

1. 🐘 **postgres** — banco de dados com healthcheck
2. 🔄 **migrate** — aplica todas as migrações Goose e encerra
3. 🚀 **api** — servidor HTTP na porta `1323`

A API estará disponível em `http://localhost:1323`. O Swagger UI fica acessível em `/swagger/*` quando `ENV=development`.

Para parar:

```bash
docker compose down
```

Para reiniciar limpando volumes:

```bash
docker compose down -v && docker compose up --build
```

---

## 🧪 Testes

```bash
cd api && go test ./tests/ -v
```

Os testes usam **testcontainers** para criar um container PostgreSQL efêmero. Não é necessário ter um banco rodando. O schema é carregado de `postgres/init.sql` e o `JWT_SECRET` é lido do `.env` (ou variável de ambiente).

---

## 📡 Endpoints da API

### 🔐 Autenticação (público)

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/auth/register` | Criar novo usuário |
| `POST` | `/auth/login` | Autenticar e obter token JWT |

### 💓 Health Check (público)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Verifica conectividade com o banco |

### 📁 Áreas (requer JWT)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/areas/list?page=1&limit=50&name=` | Listar áreas (paginado, filtro por nome) |
| `GET` | `/areas/:id` | Obter área por ID |
| `POST` | `/areas/create` | Criar nova área |
| `PUT` | `/areas/:id` | Atualizar área |
| `DELETE` | `/areas/:id` | Excluir área (soft delete) |

### 🏷️ Categorias (requer JWT)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/categories/list?page=1&limit=50&name=` | Listar categorias (paginado, filtro por nome) |
| `GET` | `/categories/:id` | Obter categoria por ID |
| `POST` | `/categories/create` | Criar nova categoria |
| `PUT` | `/categories/:id` | Atualizar categoria |
| `DELETE` | `/categories/:id` | Excluir categoria (soft delete) |

### 📊 Orçamentos (requer JWT)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/budgets/list?page=1&limit=50&year=&month=&area_id=` | Listar orçamentos (paginado, filtros) |
| `GET` | `/budgets/:id` | Obter orçamento por ID |
| `POST` | `/budgets/create` | Criar novo orçamento |
| `PUT` | `/budgets/:id` | Atualizar orçamento |
| `DELETE` | `/budgets/:id` | Excluir orçamento (soft delete) |

### 💸 Transações (requer JWT)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/transactions/list?page=1&limit=50&type=&category_id=&area_id=&from=&to=` | Listar transações (paginado, filtros) |
| `GET` | `/transactions/:id` | Obter transação por ID |
| `POST` | `/transactions/create` | Criar nova transação |
| `PUT` | `/transactions/:id` | Atualizar transação |
| `DELETE` | `/transactions/:id` | Excluir transação (soft delete) |

### 📈 Relatórios (requer JWT)

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/reports/kpis?year=&month=` | Entradas, saídas e saldo do mês (`income - expense`). Se `year/month` omitidos, usa o mês atual. |
| `GET` | `/reports/line-chart?year=&month=&type=` | Um ponto por dia do mês com `value` e `median` (últimos 6 meses, mesmo dia). `type` padrão `expense`. |
| `GET` | `/reports/radar-chart?year=&month=&type=&area_ids=` | Valor e mediana por área ativa. `area_ids` opcional, separado por vírgula. `type` padrão `expense`. |
| `GET` | `/reports/budgets?year=&month=` | Progresso de orçamentos das áreas ativas no mês (só despesas, estouro explícito). |

### 📄 Paginação

Todas as rotas `list` aceitam parâmetros de query:

- `page` — página atual (padrão: `1`)
- `limit` — itens por página (padrão: `50`, máximo: `100`)

Resposta paginada:

```json
{
  "data": [],
  "total": 0,
  "page": 1,
  "limit": 50,
  "total_pages": 0
}
```

### 📝 Exemplos de Request

**🔐 Login:**

```bash
curl -X POST http://localhost:1323/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "meulogin", "password": "minhasenha"}'
```

**📁 Criar área (com token):**

```bash
curl -X POST http://localhost:1323/areas/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Moradia"}'
```

**💸 Listar transações do mês (com filtros):**

```bash
curl "http://localhost:1323/transactions/list?from=2026-09-01&to=2026-09-30&type=expense" \
  -H "Authorization: Bearer <token>"
```

---

## 🗂️ Estrutura do Projeto

```
finance-control-api/
├── api/                        # 🐹 Aplicação Go (módulo: finance-control-api)
│   ├── main.go                 # 🏁 Entrypoint — env check → DB pool → Echo server
│   ├── go.mod / go.sum         # 📦 Dependências
│   ├── Dockerfile              # 🐋 Build multi-stage (golang → distroless)
│   ├── database/               # 🗄️ Camada de repositório (SQL bruto)
│   │   ├── connection.go       # 🔌 Pool global *sql.DB
│   │   ├── errors.go           # ⚠️ Erros sentinela mapeados do PostgreSQL
│   │   ├── pagination.go       # 📑 Helper de paginação genérico
│   │   ├── area_repository.go
│   │   ├── budget_repository.go
│   │   ├── category_repository.go
│   │   ├── transaction_repository.go
│   │   └── user_repository.go
│   ├── models/                 # 📐 Structs de domínio
│   │   ├── area.go
│   │   ├── budget.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   └── user.go
│   ├── routes/                 # 🛣️ Handlers e registro de rotas (Echo v5)
│   │   ├── router.go           # 🧭 BindRoutes — monta grupos e middleware JWT
│   │   ├── validation.go       # ✅ Helpers compartilhados (bind, validação, paginação)
│   │   ├── area_router.go
│   │   ├── auth_router.go
│   │   ├── budget_router.go
│   │   ├── category_router.go
│   │   └── transaction_router.go
│   ├── middlewares/            # 🔧 Middleware globais
│   │   ├── setup.go            # ⚡ BindMiddlewares (logger, recover, ratelimit)
│   │   ├── jwt.go
│   │   ├── auth_context.go
│   │   ├── cors.go
│   │   └── ratelimit.go
│   ├── helpers/                # 🧰 Utilitários
│   │   └── check_envs.go       # 🔍 Validação de variáveis de ambiente (viper)
│   ├── tests/                  # 🧪 Testes de integração (testcontainers)
│   │   ├── setup_test.go       # 🏗️ TestMain — container, schema, usuário de teste
│   │   ├── auth_test.go
│   │   ├── area_test.go
│   │   ├── budget_test.go
│   │   ├── category_test.go
│   │   ├── transaction_test.go
│   │   ├── soft_delete_test.go
│   │   ├── unique_test.go
│   │   └── validation_test.go
│   └── docs/                   # 📚 Swagger gerado (swag init)
├── migrations/                 # 🔄 Migrações Goose (SQL)
│   ├── Dockerfile              # 🐋 Container com goose
│   └── *.sql                   # 📝 Migrações numeradas sequencialmente
├── postgres/
│   └── init.sql                # 🏗️ Schema baseline para testes
├── docker-compose.yml
├── .env.exemple
└── .gitignore
```

---

## 🔄 Migrações

Gerenciadas com [Goose](https://github.com/pressly/goose). Novos arquivos de migração ficam em `migrations/` seguindo a nomenclatura Goose (`YYYYMMDDHHMMSS_descricao.sql`).

Executar manualmente:

```bash
docker compose run migrate
```

Ou deixar o `docker compose up` aplicar automaticamente na inicialização.

---

## 📖 Swagger

A documentação Swagger é gerada a partir das anotações nos arquivos de rota e em `main.go`. Para regenerar:

```bash
cd api && swag init
```

Acesse a UI em `http://localhost:1323/swagger/index.html` quando `ENV=development`.

---

## 📏 Convenções

- 🗑️ **Soft delete:** Todos os registros usam `StatusId = 3` para exclusão lógica. Apenas `StatusId = 1` (ativo) é retornado nas consultas.
- 👤 **Unicidade por usuário:** As constraints de unicidade são parciais (`WHERE "StatusId" = 1`) e incluem `UserId` — dois usuários podem ter registros com o mesmo nome, e nomes podem ser reutilizados após exclusão.
- ⚠️ **Erros sentinela:** `database/errors.go` define erros mapeados (`ErrNotFound`, `ErrDuplicate*`, etc.) que os handlers traduzem para códigos HTTP (404, 409, 500).
- 🟩 **Mensagens em português:** Mensagens de erro e sucesso da API são em português brasileiro.
- 🔗 **Sem injeção de dependência:** Repositórios usam o pool global `database.Pool`.
- ⏱️ **Rate limiting:** 10 requests/s com burst de 30, expira em 1 minuto, por IP.
