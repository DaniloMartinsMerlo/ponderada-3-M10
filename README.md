# Figurinhas - Ponderada 3

API REST desenvolvida em Go para gerenciar figurinhas da Copa do Mundo, construída com os princípios de Clean Code e arquitetura em camadas.

## O que foi proposto

Desenvolver uma API para cadastro e gerenciamento de figurinhas, aplicando os conceitos de Clean Code vistos em aula. A API deveria separar responsabilidades em camadas distintas (Domain, Repository, Service, Handler), usar interfaces para desacoplar essas camadas, aplicar injeção de dependência via construtores, nomear erros de domínio e mapeá-los para status HTTP, e conectar a aplicação a um banco de dados local SQLite.

## O que foi implementado

### Tecnologia

Decidimos trabalhar em go, e a nossa escolha se baseia em afinidade e também porque já haviamos utilizado esse linguagem recentemente, assim, estando com ela facil na memória. Juntamente, temos a questão de ser uma linguagem que podemos utilizar de modo runtime enquanto desenvolvemos, e também podemos compila-lá para melhor otimização de performance.

### Definição de variáveis e funções

Durante o desenvolvimento tivemos dúvidas sobre usar nomes em português ou inglês. Como os endpoints e a interface externa são em português e os DTOs internos, informados pela ponderada, estão em inglês, optamos por um meio-termo: implementações e estruturas internas em inglês, e a interface exposta ao cliente (endpoints e JSON) em português. Mantivemos a entidade em português para que o JSON retornado ao usuário esteja em português. Em alguns casos usamos termos em português por restrições da linguagem, por exemplo, `type`.

### Estrutura em camadas

O projeto foi dividido em quatro camadas:

- Domain: define a entidade `Figurinha`, os tipos enumerados (`FigureType`, `FigurePosition`), os mapas de validação e os DTOs de request
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

Uma decisão que tomamos foi tipar os campos `Tipo` e `Posicao` como tipos próprios (`FigureType` e `FigurePosition`) em vez de usar `string` diretamente, e criar mapas de validação para cada um:

```go
var ValidType = map[FigureType]bool{
    TypeCommun:       true,
    TypeShine:        true,
    TypeLegendGold:   true,
    TypeLegendCopper: true,
}

var ValidPosition = map[FigurePosition]bool{
    PositionGoalkeeper: true,
    PositionDefender:   true,
    PositionMidfielder: true,
    PositionForward:    true,
}
```

A escolha do mapa foi feita, pois, como vemos que são coisas que não vão mudar, não acreditamos que fosse necessário um `for`, que aumentaria a complexidade para algo que não precisa de flexibilidade.
A dificuldade aqui foi entender onde esses mapas deveriam viver. A ideia inicial era colocá-los no service, junto com as validações, mas percebemos que os valores válidos de um tipo são uma propriedade do próprio domínio e não da lógica de negócio.

### Validação do número via regex no service

Inicialmente, começamos trabalhando com a tag `min=6` no campo `Numero` nos DTOs, mas com o desenvolvimento da aplicação optamos por removê-la e centralizar toda a validação de formato no service, por meio de um regex compilado uma única vez no topo do código:

```go
var numeroRegex = regexp.MustCompile(`^[A-Za-z]{3} \d{2}$`)
```

A razão dessa decisão é que a tag `min=6` é uma regra de negócio disfarçada de validação de binding. O binding deve apenas garantir que o JSON está bem formado e que os campos obrigatórios vieram preenchidos; ele não deve se preocupar se o conteúdo faz sentido para o domínio. Além disso, o regex já garante o comprimento implicitamente: um padrão `XXX 00` (considere X = letra e 0 = número) tem exatamente 6 caracteres por definição, tornando a tag `min=6` redundante. Concentrar a lógica no service significa que ela pode ser testada diretamente, sem precisar montar um contexto HTTP.

### Extração do validateNumero

Ao terminar de desenvolver percebemos que a validação do número aparecia de forma idêntica em `Create` e `Update`, o que vai contra o princípio de não repetição do clean code. Devido a isso decidimos extraí-la para uma função privada:

```go
func validateNumero(numero string) error {
    if !numeroRegex.MatchString(numero) {
        return ErrInvalidNumber
    }
    return nil
}
```

Além de eliminar a duplicação, essa extração deixa explícito que existe uma regra com nome próprio, `validateNumero`, o que comunica ao leitor melhor do que um bloco `if` inline.

### campos definidos no service

Enquanto trabalhávamos no projeto, tivemos muita dúvida de onde preencheríamos as informações dos campos `ID`, `UpdatedAt` e `CreatedAt`, mas decidimos concentrá-las dentro do service, não pelo domain, repository ou handler, pois acreditamos que essas informações são regras de negócio. Porém, ao pesquisarmos um pouco mais, descobrimos que não era necessário fazer isso para o ID, pois o tipo `uint` já se auto incrementa. Dessa forma, ficando a lógica para o `UpdatedAt` e `CreatedAt` assim:

```go
figurinha := &domain.Figurinha{
    Numero:    req.Numero,
    Tipo:      req.Tipo,
    Posicao:   req.Posicao,
    UpdateAt:  time.Now(),
    CreatedAt: time.Now(),
}
```

Isso é importante porque o service é a camada que conhece as regras de negócio e para nós a data de criação e edição dizem respeito o momento em que o registro foi aceito pelo sistema então, consequentemente, é uma regra de negócio. Deixar essa responsabilidade cair para o banco (via `DEFAULT CURRENT_TIMESTAMP`) acoplaria uma regra de domínio à infraestrutura. Deixar o cliente definir seria um problema de segurança, então prefirimos definir no service para manter a regra em um lugar ideal e auditável.

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

## Como rodar localmente

```bash
git clone <url-do-repositorio>
cd ponderada-3
go mod init ponderada-3
go mod tidy
go run main.go
```

O servidor sobe em `http://localhost:8080`.

## Contrato da API

| Método | Rota | Corpo | Respostas |
|---|---|---|---|
| POST | `/figurinha` | `CreateFigureRequest` (JSON) | 201, 400 |
| GET | `/figurinha` | `?tipo=` e/ou `?posicao=` (opcional) | 200, 400 |
| GET | `/figurinha/:id` | — | 200, 404 |
| PUT | `/figurinha/:id` | `UpdateFigureRequest` (JSON) | 200, 400, 404 |
| DELETE | `/figurinha/:id` | — | 204, 404 |

Valores válidos para `tipo`: `comum`, `brilhante`, `legends_ouro`, `legends_bronze`

Valores válidos para `posicao`: `goleiro`, `zagueiro`, `meio-campista`, `atacante`

Formato de `numero`: padrão `XXX 00`, três letras, espaço, dois dígitos (ex: `BRA 15`, `FWC 02`)

Retorno esperado: 
    - POST(/figurinha), GET(/figurinha/:id) e PUT(/figurinha/:id) retornam uma figurinha, completa ou erro
    - GET(/figurinha) retorna uma lista de figurinhas, completa ou erro
    - DELETE(/figurinha) retorna somente o erro -> entendemos que essa rota não faça sentido retornar a figurinha deletada

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