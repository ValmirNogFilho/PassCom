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

O servidor desenvolvido utilizou uma abordagem de arquitetura "stateless". Tanto as comunicações entre clientes das companhias e o servidor, quanto a comunicação entre servidores, consistem em um sistema de requisições e respostas.

A arquitetura consiste na transação de mensagens curtas, com solicitações pequenas para operações nos servidores. A escolha adotada traz mais escalabilidade, além de portabilidade para mais dispositivos acessarem os servidores em simultâneo.

Entretanto, embora o paradigma seja stateless, que prevê a ausência de estados, são armazenadas informações temporárias de sessão dos usuários. A abordagem transmite um ID de "token" da sessão, para garantir segurança em acessos e permissões privadas para os usuários do sistema, além de criar funcionalidades personalizadas, como carrinho de compras temporário para os usuários.

# Protocolo de comunicação

**Especifique as APIs de comunicação implementadas entre os componentes desenvolvidos, descrevendo os métodos remotos, parametros e retornos empregados para permitir a compra de passagens entre clientes e servidores, e entre servidores.**

A equipe foi orientada a implementar o protocolo de comunicação HTTP, que acarretou na decisão do paradigma stateless. O protocolo HTTP é feito em cima do paradigma stateless, além de ser construido em cima do protocolo TCP, que assegura a entrega de informações na conexão.

# Roteamento

A decisão adotada para o projeto faz com que cada servidor da PassCom possua, além de seus dados, uma réplica do banco de dados dos outros servidores. Os algoritmos de consenso e roteamento permitem a sincronização das compras de forma segura. Quaisquer operações sobre as passagens dos vôos fazem com que todas as réplicas sejam alteradas, independente das operações serem locais ou referentes aos outros servidores.

Tendo em vista que todos esses dados são compartilhados, se torna mais prático fazer algoritmos de grafos sobre os dados para formular as rotas, visto que a tabela de vôos forma um supergrafo das três companhias. A decisão de adicionar uma permanência de dados a partir de um banco de dados relacional auxiliou na associatividade dos vôos, a partir da relação entre aeroportos de origem e destino. 

Como trata-se de um protótipo, foi utilizado um banco de dados SQLite, que é mais simples e possui os mesmos princípios SQL de bancos mais complexos. O acesso aos dados a partir do padrão Data Access Object (DAO) de forma centralizada permite a mudança para um banco de dados mais escalável e seguro com poucas mudanças nas configurações de drivers. Além permitir acesso aos dados do banco pelos models do projeto, a biblioteca GORM abstrai o acesso a banco de dados relacionais, tornando essa adaptação ainda mais simples.

Um algoritmo de Bread-First-Search forma o caminho mais curto a partir das rotas distribuídas. Essas informações são expostas na interface gráfica a partir das passagens individualmente compráveis e do mapa, que ilustra o caminho das rotas, com as cores das rotas simbolizando as cores temáticas das três companhias. As passagens também expoem as logomarcas de suas respectivas companhias.

# Concorrência Distribuída

Para assegurar a consistência dos dados distribuídos de maneira descentralizada, foi utilizado o "gossip protocol". Trata-se de um algoritmo de consenso peer-to-peer para sistemas distribuídos focado em manter o estado do seus nós (no caso, os servidores).

A lógica do protocolo consiste no envio periódico em \"broadcast" dos dados de compra e cancelamento de vôos para as três companhias, independentemente da operação ser relacionada ao servidor local ou aos servidores remotos. Isso garante a estabilidade e consistência de todas réplicas do banco de dados, de maneira descentralizada, evitando que a existência de um servidor central que possa se desconectar atrapalhe o sistema como um todo.

# Confiabilidade da solução

A solução da equipe aplicou o conceito de "heartbeat" e relógios vetoriais na comunicação entre os servidores, para asssegurar a confiabilidade dos dados após a possível desconexão de um dos servidores. 

O "heartbeat" trata-se de um algoritmo que envia mensagens periódicas para os servidores, a fim de apenas checar se estão ativos. Caso contrário, o servidor desconectado é desconsiderado para operações de consultas, até que possa talvez se reconectar novamente. Para isso, o heartbeat persiste lhe mandando sinais, a espera de um possível retorno. A proposta de algoritmo não causa grande peso nos servidores, por mandar mensagens leves e em um período de tempo razoável.

 Os relógios vetoriais armazenam três contadores de processos relativos aos  respectivos três servidores. Após um processo de um servidor, seu contador é incrementado em cada cópia do relógio de cada servidor. O incremento dos contadores após cada processo assegura que o sistema saiba a ordem causal dos eventos, a partir da visualização das cópias e a ordem que seus contadores são incrementados.

Assim, se um servidor se desconecta por um período e se reconecta posteriormente, pode recuperar os dados perdidos após descobrir que seus contadores estão reduzidos em relação aos demais relógios. Após a desconexão de qualquer um dos servidores, seu relógio vetorial é serializado e armazenado no seu arquivo `systemvars.json`, na sua pasta root, além de outros dados importantes para a sincronização, como registros de conexões, seus horários, endereços de server, logs e informações de identificação do próprio server.

# Avaliação da Solução

Cada um dos servidores possui uma pasta `test`, com testes de sincronização entre servidores, a partir da consulta dos relógios vetoriais. Os testes funcionam plenamente, demonstrando a confiabilidade das abordagens adotadas, inclusive em possíveis desconexões de partes do sistema.

# Documentação do código

As funções e métodos do projeto relativas a lógica de negócios, endpoints da API e componentes da lógica interna de comunicação distribuída estão documentadas, permitindo melhor visualização dos parâmetros a serem passados e o retorno das operações.

# Emprego do Docker

O sistema completo foi conteinerizado via uso do Docker. Cada um dos servidores separados possui um Dockerfile, com instruções de diretórios a serem copiados, portas a serem expostas e volumes de persistência de dados a serem considerados (arquivos JSON e os arquivos de database SQLite). Também foram criados contêineres para execução das interfaces React, e a comunicação entre front-end e back-end pelas APIs foram asseguradas pelas networks criadas. Os Dockerfiles dos servidores também expõem as portas para acesso ao server CLI de monitoramento dos servidores REST.

Assim, o arquivo `docker-compose.yaml` une a execução dos contêineres, permitindo o build e execução dos componentes de cada uma das companhias aéreas a partir do comando:

```
docker compose up --build
```

# Conclusão

A equipe produziu uma solução confiável, robusta, e de apresentação moderna e intuitiva para usuários leigos com interesse na compra de passagens aéreas. A arquitetura organizacional do projeto foi produzida visando separar bem funções em arquivos e diretórios com propósitos em comum, de acordo com o Single Responsability Principle. Além disso, atribuiu os padrões MVC para dividir diretórios relativos a dados, lógica de negócios e exibição, e DAO para centralizar localmente o acesso ao banco de dados, tornando o consumo de dados fácilmente reconfigurável e evitando repetição inconsistente de locais.

Além disso, a escolha mais robusta dos algoritmos de consenso priorizam a consistência das informações, acima da decisão de economia de espaço de armazenamento, visto que cada servidor possui réplicas dos dados dos outros.
