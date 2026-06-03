# Figurinhas  Ponderada 3

API REST desenvolvida em Go para gerenciar figurinhas da Copa do Mundo, construída com os princípios de Clean Code e arquitetura em camadas.

---

## O que foi proposto

Desenvolver uma API para cadastro e gerenciamento de figurinhas, aplicando os conceitos de Clean Code vistos em aula. A API deveria separar responsabilidades em camadas distintas (Domain, Repository, Service, Handler), usar interfaces para desacoplar essas camadas, aplicar injeção de dependência via construtores, nomear erros de domínio e mapeá-los para status HTTP, e conectar a aplicação a um banco de dados local SQLite.

---

## O que foi implementado

### Estrutura em camadas

O projeto foi dividido em quatro camadas:

- Domain: define a entidade `Figurinha`, os tipos enumerados (`FigurinhaType`, `FigurinhaPosition`), os mapas de validação e os DTOs de request
- Repository: interface `FigureRepository` e sua implementação concreta com GORM, responsável exclusivamente pelo acesso ao banco
- Service: interface `FigureService` e sua implementação, onde residem todas as regras de negócio e validações
- Handler: camada HTTP com Gin, responsável apenas por receber requests, delegar ao service e mapear erros para status HTTP

### Injeção de dependência

Cada camada recebe suas dependências via construtor e expõe apenas a interface, nunca o tipo concreto:

```go
func NewFigureRepository(db *gorm.DB) FigureRepository { ... }
func NewFigureService(repo repository.FigureRepository) FigureService { ... }
func NewFigureHandler(svc service.FigureService) *FigureHandler { ... }
```

Isso permite que cada camada seja testada de forma isolada, substituindo a dependência real por um mock sem alterar nada além do ponto de construção.

### Enums e mapas de validação no domain

Uma decisão que tomamos foi tipar os campos `Tipo` e `Posicao` como tipos próprios (`FigurinhaType` e `FigurinhaPosition`) em vez de usar `string` diretamente, e criar mapas de validação para cada um:

```go
var ValidType = map[FigurinhaType]bool{
    TypeCommun:       true,
    TypeShine:        true,
    TypeLegendGold:   true,
    TypeLegendCopper: true,
}

var ValidPosition = map[FigurinhaPosition]bool{
    PositionGoalkeeper: true,
    PositionDefender:   true,
    PositionMidfielder: true,
    PositionForward:    true,
}
```

A dificuldade aqui foi entender onde esses mapas deveriam viver. A ideia inicial era colocá-los no service, junto com as validações, mas percebemos que os valores válidos de um tipo são uma propriedade do próprio domínio e não da lógica de negócio. Mover os mapas para o domain deixou o service mais limpo e garantiu que qualquer camada que precise validar um tipo possa fazer isso sem depender do service.

### Validação do número via regex no service

Optamos por remover a tag `min=6` do campo `Numero` nos DTOs e centralizar toda a validação de formato no service, por meio de um regex compilado uma única vez no topo do código:

```go
var numeroRegex = regexp.MustCompile(`^[A-Za-z]{3} \d{2}$`)
```

A razão dessa decisão é que o `min=6` é uma regra de negócio disfarçada de validação de binding. O binding deve apenas garantir que o JSON está bem formado e que os campos obrigatórios vieram preenchidos, ele não deve se preocupar se o conteúdo faz sentido para o domínio. Além disso, o regex já garante o comprimento implicitamente, um padrão `XXX 00` tem exatamente 6 caracteres por definição, tornando o `min=6` redundante. Concentrar a lógica no service significa que ela pode ser testada diretamente, sem precisar montar um contexto HTTP.

### Extração do validateNumero

A validação do número aparecia de forma idêntica em `Create` e `Update`, o que vai contra o princípio de não repetição do clean code. Extraímos para uma função privada:

```go
func validateNumero(numero string) error {
    if !numeroRegex.MatchString(numero) {
        return ErrInvalidNumber
    }
    return nil
}
```

Além de eliminar a duplicação, essa extração deixa explícito que existe uma regra com nome próprio, `validateNumero` comunica ao leitor melhor do que um bloco `if` inline.

### created_at definido no service

O campo `CreatedAt` é preenchido com `time.Now()` dentro do service, não pelo banco e não pelo cliente:

```go
figurinha := &domain.Figurinha{
    Numero:    req.Numero,
    Tipo:      req.Tipo,
    Posicao:   req.Posicao,
    UpdateAt:  time.Now(),
    CreatedAt: time.Now(),
}
```

Isso é importante porque o service é a camada que conhece as regras de negócio e para nós a data de criação é o momento em que o registro foi aceito pelo sistema então, consequentemente, é uma regra de negócio. Deixar essa responsabilidade cair para o banco (via `DEFAULT CURRENT_TIMESTAMP`) acoplaria uma regra de domínio à infraestrutura. Deixar o cliente definir seria um problema de segurança, então prefirimos definir no service para manter a regra em um lugar ideal e auditável.

### Mapeamento de erros para HTTP

Os erros de domínio são nomeados no service e mapeados para status HTTP exclusivamente no handler, através de uma função `mapErrorToStatus`:

```go
func mapErrorToStatus(err error) int {
    switch {
    case errors.Is(err, service.ErrFigureNotFound):
        return http.StatusNotFound
    case errors.Is(err, service.ErrInvalidNumber),
        errors.Is(err, service.ErrInvalidType),
        errors.Is(err, service.ErrInvalidPosition):
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

O service não conhece HTTP, ele apenas diz o que deu errado em termos de negócio. Quem traduz isso para códigos de status é o handler.

---

## Como rodar localmente

```bash
git clone <url-do-repositorio>
cd ponderada-3
go mod init ponderada-3
go mod tidy
go run main.go
```

O servidor sobe em `http://localhost:8080`.

---

## Contrato da API

| Método | Rota | Corpo | Respostas |
|---|---|---|---|
| POST | `/figurinha` | `CreateFigureRequest` (JSON) | 201, 400 |
| GET | `/figurinha` | `?tipo=` e/ou `?posicao=` (opcional) | 200, 400 |
| GET | `/figurinha/:id` | — | 200, 404 |
| PUT | `/figurinha/:id` | `UpdateFigureRequest` (JSON) | 200, 400, 404 |
| DELETE | `/figurinha/:id` | — | 204, 404 |

Valores válidos para `tipo`: `comum`, `brilhante`, `legends_ouro`, `legends_bronze`

Valores válidos para `posicao`: `Goleiro`, `Zagueiro`, `Meio-campista`, `Atacante`

Formato de `numero`: padrão `XXX 00`, três letras, espaço, dois dígitos (ex: `BRA 15`, `FWC 02`)

---

## Critérios atendidos

| Critério | Status |
|---|---|
| Camadas separadas com responsabilidades distintas | ✅ |
| Interfaces desacoplando camadas | ✅ |
| Injeção de dependência via construtores | ✅ |
| Erros de domínio mapeados para status HTTP no handler | ✅ |
| Validações de negócio na camada service | ✅ |
| Banco de dados SQLite conectado | ✅ |
| CRUD completo (POST, GET, GET por ID, PUT, DELETE) | ✅ |