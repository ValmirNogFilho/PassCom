Protocolos e conceitos estudados:

- Gossip Protocol
- Reliable Broadcast (Total Order)
- Eventual Consistency
- Vector Clocks
- Transações distribuídas
- Raft, Paxos
- Linearização
- Last Writer Wins (alcancável com relógios vetoriais)

# Resumo

# Introdução

No último século, o advento do comércio digital (também chamado de "E-commerce") se propagou de maneira exponencial. Isso ocorreu principalmente nos últimos anos, devido a variáveis como o período pandêmico, que forçou a digitalização do comércio. As plataformas digitais fornecem qualidade de serviço, atendimento prático, rápido e automatizado, proporcionando melhor experiência aos clientes.

O presente relatório visa descrever o desenvolvimento do sistema "Passcom", um sistema de venda de passagens distribuído para três companhias aéreas (Rumos, Giro e Boreal). O software foi desenvolvido utilizando a arquitetura REST, com servidores HTTP individualizados codificados em Go, com acesso a banco de dados SQLite via uso da biblioteca GORM para mapeamento relacional de objetos. O front-end foi desenvolvido com o framework React. Também foi adicionada uma interface CLI acessível via TCP para monitoramento e configuração dos servidores. O sistema é conteinerizado com Docker, permitindo consistência no funcionamento.

O software assegura a sincronização distribuída a partir do uso de diversos algoritmos de consenso e protocolos a serem explicados no relatório. O resultado atendeu os requisitos, pela confecção de um projeto bem arquiteturado, robusto que visa a praticidade do uso para os usuários.

### Equipe:

[![Pierre Machado](https://github.com/pierremachado.png?size=20)](https://github.com/pierremachado) [Pierre Machado](https://github.com/pierremachado)

[![Valmir Nogueira](https://github.com/valmirnogfilho.png?size=20)](https://github.com/valmirnogfilho) [Valmir Nogueira](https://github.com/valmirnogfilho)

# Arquitetura da solução

Stateless

# Protocolo de comunicação

Especifique as APIs de comunicação implementadas entre os componentes desenvolvidos, descrevendo os métodos remotos, parametros e retornos empregados para permitir a compra de passagens entre clientes e servidores, e entre servidores.

HTTP

# Roteamento

Qual o método usado para o cálculo distribuído das rotas entre origem e destino da passagem, e se o sistema consegue mostrar para os usuários todas as rotas possíveis considerando os trechos disponíveis nos servidores de todas as companhias.

Réplicas dos bancos de dados

# Concorrência Distribuída

Especifique conteitualmente a solução empregada para evitar mais vendas de passagens do que a quantidade existente e ou a venda da mesma passagem para clientes distintos.

Broadcast do voo para as réplicas

# Confiabilidade da solução

Desconectando e conectando os servidores das companhias, o sistema continua garantindo a concorrência distribuída e a finalização da compra anteriormente iniciada por um cliente.

Relógios vetoriais

# Avaliação da Solução

Se foi desenvolvido e mantido no Github o código para testar a consistência da solução sob condições críticas e ou avaliar o desempenho do sistema.

# Documentação do código

Se o código do projeto possui comentários explicando as principais classes, e se as funções têm descrições sobre seu propósito, parâmetros, e o retorno esperado.

# Emprego do Docker

O sistema completo foi conteinerizado via uso do Docker. Cada um dos servidores separados possui um Dockerfile, com instruções de diretórios a serem copiados, portas a serem expostas e volumes de persistência de dados a serem considerados (arquivos JSON e os arquivos de database SQLite). Também foram criados contêineres para execução das interfaces React, e a comunicação entre front-end e back-end pelas APIs foram asseguradas pelas networks criadas. Os Dockerfiles dos servidores também expõem as portas para acesso ao server CLI de monitoramento dos servidores REST.

Assim, o arquivo `docker-compose.yaml` une a execução dos contêineres, permitindo o build e execução dos componentes de cada uma das companhias aéreas a partir do comando:

```
docker compose up --build
```

# Conclusão
