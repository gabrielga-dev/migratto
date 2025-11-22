# 🧭 Summary/ Sumário

- [🇺🇸 English](#en)
- [🇧🇷 Português](#ptbr)

# 🇺🇸 En

## 🧩 Migratto — SQL Migration Runner for Go Projects

**Migratto** is a simple and straightforward tool for applying SQL migrations to **PostgreSQL** databases within **Go** projects.  
It automates the sequential execution of `.sql` scripts from a configured directory and keeps an internal record of applied migrations.

---

### 🚀 Features

- 📂 Automatically reads `.sql` files from the configured directory.  
- 🧠 Automatically creates the `migratto_migration_history` table (if it doesn’t exist).  
- ⚙️ Executes only the new migrations that haven’t been applied yet.  
- 🪶 Simple and direct configuration through a Go struct.  
- 🔍 Optional logging of executed operations.  

> ⚠️ **Important:** Migratto currently tracks **only the number of applied migrations**, and **not the file content**, unlike tools such as Flyway.  
> Changes to previously executed scripts are not detected — this functionality will be implemented in future versions.

---

### 🧱 Configuration Structure

Configuration is defined through the `ConfigDTO` struct:

```go
type ConfigDTO struct {
	DatabaseDriver   string
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	Sslmode          string
	MigrationsDir    string
	Log              bool
}
```

### 📁 Migrations Directory Structure

The folder defined in MigrationsDir must contain only .sql files following the naming pattern:

`V<migration_number>_<description>.sql`

Example:

```
migrations/
├── V01_create_users_table.sql
├── V02_insert_initial_data.sql
├── V03_create_orders_table.sql
```

---

### ⚙️ How It Works

1) The script checks whether the table migratto_migration_history exists.
    - If it doesn’t, it will be created automatically.

2) Then, it reads all .sql files in the directory specified by MigrationsDir.

3) It compares the number of migrations already applied with the number of files found.

4) If new migrations exist, only those that have not yet been applied are executed in sequence.

---

### 🧩 Practical Example

Suppose you initially have:

```
V01_script.sql
V02_create_table.sql
```

After running Migratto, both scripts will be executed and recorded in the migratto_migration_history table.

Later, if you change the file names and contents to:

```
V01_script_database.sql
V02_create_tables.sql
V03_add_columns.sql
```

Migratto will execute only V03_add_columns.sql, ignoring changes in the previous files, since it only tracks the migration count and not the content.

---

###🗃️ Internal Control

Migration history is managed through the migratto_migration_history table, which is automatically created if it does not exist.
It stores the number of the last executed migration, preventing duplicates and maintaining consistent progress across runs.

---

### 🧠 Execution Example

Import Migratto running:

> go get https://github.com/gabrielga-dev/migratto

Suppose your main.go file looks like this:

```
package main

import (
	"fmt"

	dto "github.com/gabrielga-dev/migratto/dto"
	migration_service "github.com/gabrielga-dev/migratto/service/migration"
)

func main() {
		config := dto.ConfigDTO{
		DatabaseDriver:   "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "migratto",
		DatabaseUsername: "user",
		DatabasePassword: "pass",
		Sslmode:          "disable",
		MigrationsDir:    "./migrations",
		Log:              true,
	}
	migration_service.Migrate(config)
}
```
Then, simply run in your terminal:

> go run main.go

---

### 🧩 Next Steps

 - [x] Validate migration content (not just the count).

 - [x] Add support for other database types.


# 🇧🇷 Ptbr

## 🧩 Migratto — Aplicador de migrações SQL para projetos Go

**Migratto** é uma ferramenta simples e direta para aplicar migrações SQL em bancos de dados **PostgreSQL** dentro de projetos **Go**.  
Ele automatiza a execução sequencial de scripts `.sql` a partir de um diretório configurado e mantém um histórico interno das migrações já aplicadas.

---

### 🚀 Funcionalidades

- 📂 Leitura automática de arquivos `.sql` no diretório configurado.  
- 🧠 Criação automática da tabela `migratto_migration_history` (caso não exista).  
- ⚙️ Execução apenas das novas migrações ainda não aplicadas.  
- 🪶 Configuração simples e direta via struct Go.  
- 🔍 Log opcional das operações executadas.  

> ⚠️ **Importante:** atualmente o Migratto controla **apenas a quantidade de migrações aplicadas**, e **não o conteúdo dos arquivos**, como o Flyway faz.  
> Alterações em scripts antigos não são detectadas — essa funcionalidade será implementada futuramente.

---

### 🧱 Estrutura de Configuração

A configuração é feita por meio da struct `ConfigDTO`:

```go
type ConfigDTO struct {
	DatabaseDriver   string
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	Sslmode          string
	MigrationsDir    string
	Log              bool
}
```

---

### 📁 Estrutura de diretório de migrações

A pasta configurada em MigrationsDir deve conter apenas arquivos .sql que seguem o padrão:

```V<numero_da_migracao>_<descricao>.sql```

Exemplo:

```
    migrations/
    ├── V01_cria_tabela_usuarios.sql
    ├── V02_insere_dados_iniciais.sql
    ├── V03_cria_tabela_pedidos.sql
```

---

### ⚙️ Como funciona

1) O script verifica se existe a tabela migratto_migration_history.
    -   Caso não exista, ela é criada automaticamente.

2) Em seguida, ele lê todos os arquivos .sql no diretório definido em MigrationsDir.

3) Ele compara a quantidade de migrações já aplicadas com a quantidade de arquivos encontrados.
4) Caso existam novas migrações, apenas as que ainda não foram aplicadas são executadas em sequência.

---

### 🧩 Exemplo prático

Suponha que você tenha inicialmente:

```
V01_script.sql
V02_cria_tabela.sql
```


Após rodar o Migratto, ambas serão aplicadas e registradas na tabela migratto_migration_history.

Mais tarde, se você alterar o conteúdo e nomes para:

```
V01_script_banco.sql
V02_cria_tabelas.sql
V03_cria_colunas.sql
```


O Migratto executará somente V03_cria_colunas.sql, ignorando as mudanças nos arquivos anteriores, pois ele observa apenas a contagem, e não o conteúdo das migrações.

---

### 🗃️ Controle interno

O controle de histórico é feito por meio da tabela migratto_migration_history, criada automaticamente no banco de dados caso não exista.
Ela armazena o número da última migração executada, evitando repetições e mantendo o progresso do projeto.

---

### 🧠 Exemplo de execução

Import Migratto running:

> go get https://github.com/gabrielga-dev/migratto

Suponha que você tenha o arquivo main.go assim:

```go
package main

import (
	"fmt"

	dto "github.com/gabrielga-dev/migratto/dto"
	migration_service "github.com/gabrielga-dev/migratto/service/migration"
)

func main() {
		config := dto.ConfigDTO{
		DatabaseDriver:   "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "migratto",
		DatabaseUsername: "user",
		DatabasePassword: "pass",
		Sslmode:          "disable",
		MigrationsDir:    "./migrations",
		Log:              true,
	}
	migration_service.Migrate(config)
}

```

E no terminal, basta executar:

> go run main.go

---

### 🧩 Próximos passos

- [x] Validar conteúdo das migrações (não apenas a contagem);
- [x] Suporte a outros tipos de bancos de dados.