
# Backend | Desafio Técnico - Resolução

# Introdução
Este projeto consiste em uma API criada para a resolução de um case para o processo seletivo da Ideal CVTM

# Bibliotecas do Go
- Foi utilizada a biblioteca Fiber para a criação dessa API
    - https://github.com/gofiber/fiber
    - https://docs.gofiber.io/
- Para buscar as informações das ações e stocks foi utlilizada a biblioteca **finance-go**. Esta biblioteca utiliza as informaçõs do Yahoo finance.
    - https://pkg.go.dev/github.com/piquette/finance-go
    - https://piquette.io/projects/finance-go/
- Para a camada de persistencia foi utilizada a biblioteca Mongo Drive
    - https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
    - https://github.com/mongodb/mongo-go-driver

# Banco de dados
Neste projeto foi utilizado o banco de dados MongoDB. A motivação para utilizar este banco foi devido à sua facilidade de utilização, livre de custos para a sua utilização e formato das informações a serem persistidas.
    - Usuario: usuario
	- Senha: senha123
    - conexão: MONGOURI=mongodb+srv://usuario:senha123@cluster0.po0gmgr.mongodb.net/?retryWrites=true&w=majority
 
# Funcionamento
Esta API tem a função de cadastrar valores de ativos associados à um CPF e buscar o valor atualizado destes ativos na Yahoo Finance.
Os valores apresentados dos ativos estão em dolar.

    1. Primeiramente deve-se cadastrar um usuario e suas informações básicas de CPF, Nome e E-mail;
    2. As informações do usuario de Nome e E-mail podem ser editadas pós cadastro inicial;
    3. Depois deve-se atribuir nomes de ativos para serem associados à esse CPF;
    4. Os ativos podem ser adicionados ou deletados individualmente ou em blocos
    5. A ordem dos ativos também pode ser alterada por tipo e ordem
        - Tipos: Nome, Preço e Lista
        - Ordem Crescente = 1 e Decrescente = 0

# Rotas
**Invetidor:**

	- Criar perfil de Investidor:
        - POST: http://localhost:6000/investor
	- Traz a informação de um Investidor por CPF:
        - GET: http://localhost:6000/investor/:cpf
	-  Traz as informações de todos os Investidores:
        - GET: http://localhost:6000/investors
	- Edita as informações de um Investidor por CPF:
        - PUT: http://localhost:6000/investor/:cpf
	- Deleta as informações de um investidor por CPF:
        - Delete: http://localhost:6000/investor/:cpf

**Assets:**

    - Traz o preço de um ou mais ativos com o valor em Dolar:
        - POST: http://localhost:6000/
	- Adiciona um Asset à lista de Assets de um Investidor:
	    - POST: http://localhost:6000/asset/asset/:cpf
	- Remove um Asset à lista de Assets de um Investidor:
	    - POST: http://localhost:6000/asset/:cpf/remove
	- Ordena os Assets da lista de um Investidor por criteriod e Nome, Preço ou Lista 
	    - POST: http://localhost:6000/asset/order/:type/:asc/:cpf
            - type = "name" ou "type" ou "list"
            - asc = 0 (Crescente)
            - asc = 1 (Decrescente)

