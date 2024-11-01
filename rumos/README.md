<h1 align="center">VendePass</h1>
<h2 align="center">Sistema de Vendas de Passagens Aéreas Baseado na Arquitetura TCP Cliente-Servidor</h2>

![presentation](docs/presentation.png)

### Equipe

[![Pierre Machado](https://github.com/pierremachado.png?size=20)](https://github.com/pierremachado) [Pierre Machado](https://github.com/pierremachado)

[![Valmir Nogueira](https://github.com/valmirnogfilho.png?size=20)](https://github.com/valmirnogfilho) [Valmir Nogueira](https://github.com/valmirnogfilho)

### Introdução

É evidente que houve uma crescente demanda para que as empresas se adaptem ao rápido crescimento do comércio virtual nos últimos anos, já que os sistemas de compras online oferecem conforto aos clientes que buscam agilidade e preços acessíveis. Um dos setores que mais se destaca nesse cenário é o das companhias aéreas, que, através de plataformas automatizadas, proporciona interações rápidas e práticas para a compra de passagens. Em vista disso, uma empresa aérea identificou a necessidade de implementar um sistema mais eficiente para pesquisar rotas e vender passagens, de modo a otimizar a experiência do cliente e aumentar a competitividade no mercado. 

Neste contexto, este relatório visa detalhar o desenvolvimento de um sistema de compra e venda de passagens aéreas pela internet, utilizando o protocolo TCP/IP e a arquitetura cliente-servidor. A solução proposta baseia-se no uso dos pacotes da linguagem de programação Go para a implementação de um servidor TCP e a comunicação com o cliente, a criação de uma API baseada no protocolo de solicitações HTTP e a framework React do JavaScript para o front-end. O sistema foi dockerizado e testado extensivamente, tratando requisições de clientes com paralelismo. A solução atendeu aos requisitos propostos e resultou em um projeto funcional e robusto. 

### Arquitetura da solução

A solução adotada utilizou a arquitetura de comunicação cliente-servidor. No modelo cliente-servidor, o cliente é o programa que faz solicitações, enquanto o servidor é o responsável por responder a essas solicitações. A comunicação entre ambos ocorre por meio de protocolos, como o Transmission Control Protocol (TCP), que garante a entrega confiável de dados na internet, mesmo em redes instáveis. O uso do protocolo TCP garante a integridade e a ordem dos pacotes, algo crucial em um sistema de compra de passagens, onde a precisão dos dados trocados é essencial para evitar inconsistências, como reservas duplicadas ou falhas de processamento. 

O servidor do projeto foi desenvolvido utilizando as bibliotecas da linguagem de programação Go, que oferece ferramentas nativas para a criação de servidores eficientes. A linguagem de programação Go é particularmente vantajosa para lidar com múltiplas requisições, graças às goroutines e aos channels. As goroutines são threads leves que permitem executar múltiplas tarefas simultaneamente, otimizando o desempenho do servidor através do paralelismo. Já os channels atuam como canais de comunicação entre goroutines, garantindo a sincronização e evitando condições de corrida. Isso é fundamental em situações em que vários clientes tentam reservar os mesmos assentos, onde transações atômicas são necessárias para assegurar que cada reserva seja processada de forma única e consistente. 

Além disso, foi implementado um banco de dados utilizando o padrão Data Access Object (DAO). Embora a linguagem Go não seja orientada a objetos por natureza, é possível simular a orientação a objetos por meio de estruturas (structs). Com base em um diagrama de classes, foram criadas estruturas para representar entidades reais, como cliente, aeroporto, voo e sessão, e seus respectivos relacionamentos. Por exemplo, a estrutura de voo armazena dados sobre o aeroporto de origem, aeroporto de destino e a quantidade de assentos disponíveis, permitindo o gerenciamento eficaz de reservas e rotas. 

![fig1](docs/0.jpg)

Figura 1. Diagrama de classes inicial do projeto

O cliente foi desenvolvido utilizando o framework React no JavaScript, que é amplamente utilizado para a construção de interfaces de usuário interativas e dinâmicas. O React facilita a criação de componentes reutilizáveis e a atualização eficiente da interface com base nas interações do usuário, permitindo o desenvolvimento ágil do front-end do projeto. Através de uma interface gráfica projetada no Figma e desenvolvida no framework React junto ao CSS, o cliente é capaz de realizar login, escolher rotas, adicionar as rotas ao carrinho e confirmar a compra. 

![fig2](docs/1.jpg)

Figura 2. Tela de login 

![fig3](docs/2.jpg)

Figura 3. Tela inicial após login 

A comunicação entre o front-end e o back-end ocorre por meio de uma Interface de Programação de Aplicativo (Application Programming Interface, API) desenvolvida em Go, utilizando a biblioteca net. O cliente envia requisições HTTP para a API, que por sua vez converte essas requisições para um formato que o servidor TCP possa processar e devolver uma resposta, que é recebida pela API HTTP e retornada para o "front-end". Isso permite que operações como a pesquisa de rotas e a compra de passagens sejam realizadas de forma eficiente e segura. 

![fig4](docs/3.png)

Figura 4. Fluxograma de vendas do serviço 

### Paradigma de comunicação

Existem dois tipos de paradigma de comunicação: stateful e stateless. No paradigma de comunicação stateful, o sistema mantém o registro das interações do usuário ao longo de suas sessões, permitindo que cada nova solicitação seja tratada com base no histórico de solicitações anteriores. Isso possibilita uma experiência mais interativa e personalizada, como quando o sistema “lembra” de itens colocados no carrinho ou o estado de uma transação em andamento. Já no paradigma stateless, cada solicitação é tratada de forma isolada, sem considerar o histórico de interações do cliente. O servidor não mantém informações sobre o estado da sessão entre uma requisição e outra, o que pode simplificar a escalabilidade do sistema, mas torna necessário que todas as informações relevantes sejam transmitidas a cada nova solicitação. 

Além disso, uma diferença importante entre esses paradigmas está na maneira como as threads são gerenciadas. Em sistemas stateful, há uma maior complexidade no gerenciamento de threads, pois elas precisam lidar com a continuidade do estado da sessão, garantindo que as interações anteriores sejam levadas em conta. No entanto, em sistemas stateless, cada thread trata uma requisição individualmente, sem necessidade de armazenar ou consultar o estado anterior, o que pode facilitar a execução paralela de múltiplas solicitações. 

O paradigma de comunicação adotado para este projeto foi o stateless, pois o cliente precisa enviar um token de autenticação associado à sua sessão a cada requisição. Entretanto, o servidor registra as sessões ativas no DAO (Data Access Object) assim que o cliente faz login. Esse mecanismo é útil, por exemplo, para adicionar uma rota ao carrinho de compras. Associado a sessão do cliente, o cliente tem 30 minutos para confirmar a compra, e, nesse período, a quantidade de assentos disponíveis para o voo é temporariamente reduzida, garantindo a reserva até que a transação seja concluída, o tempo expire ou o cliente se desconecte.

No servidor, cada requisição é tratada de forma independente por meio de goroutines. Quando uma solicitação chega, uma nova goroutine é criada para processá-la, garantindo eficiência e paralelismo. A goroutine verifica se o token enviado pelo cliente é válido, o que permite ao servidor identificar e autorizar a operação, sem a necessidade de manter uma sessão ativa entre as requisições. Após o processamento da requisição e a geração da resposta, a goroutine é finalizada.

### Protocolo de comunicação

O software faz uso de dois protocolos de comunicação. O primeiro foi desenvolvido pela equipe, de modo a organizar e assegurar os dados na comunicação através da API Socket Básica, que não dispõe um protocolo de troca de mensagens pronto.

O protocolo desenvolvido é modelado em cima de "requests" e "responses", inspirado no padrão HTTP, no entanto, com uma proposta mais simples, voltada à escala do projeto. A comunicação servidor-cliente se dá pelo cliente criando uma conexão TCP com o servidor e enviando uma "request". O servidor cria uma "goroutine" para processar a solicitação da "request", gerando ao final uma "response" com os resultados e fechando a "goroutine".

!["Requests" e "Responses"](https://github.com/ValmirNogFilho/VendePass/blob/master/docs/umlpbleredes.jpg)

Figura 5. Estruturas de dados das "Requests" e "Responses".

A "request" é composta pelos campos:

- "Action", indicando a ação solicitada do cliente para o servidor;
- "Auth", enviando um token de sessão de um usuário específico como string, para permitir que os dados solicitados serão processados para o usuário correto, de maneira autenticada;
- "Data", enviando qualquer tipo de dados que seja incluído em um objeto "JSON", em pares chave-valor com chave de tipo string e valor de tipo genérico.

A requisição de login não necessita do envio da autenticação, visto que é a solicitação que pede um token de autenticação do servidor, para permitir que o usuário faça outras "requests" posteriormente.

As "responses" retornam respostas nos campos:

- "Error", indicando possíveis erros que impossibilitaram o feitio da "request" (credenciais inválidas de usuário, erros internos, regras de negócio);
- "Data", enviando qualquer tipo de resposta em "JSON", como já explicado anteriormente.

Para a integração da aplicação React com o servidor TCP, foi necessária a adição do protocolo de comunicação Hypertext Transfer Protocol (HTTP). O servidor permanece em um loop infinito, utilizando um listener que "ouve" as requisições dos clientes e as processa conforme chegam. Quando o cliente faz uma solicitação ao servidor, ela é recebida como uma requisição HTTP pela API, que atua como um middleware, intermediando a comunicação entre o cliente e o servidor.  

Os métodos HTTP mais comuns incluem: 

- GET: utilizado para solicitar dados do servidor. Por exemplo, ao pesquisar voos disponíveis, o cliente envia uma requisição GET para receber informações sobre rotas e horários. 

- POST: utilizado para enviar dados ao servidor. Por exemplo, quando um cliente realiza uma compra, as informações da compra são enviadas através de uma requisição POST. 

- PUT: usado para atualizar informações já existentes. 

- DELETE: usado para excluir recursos no servidor, como o cancelamento de uma reserva. 

Após a API receber a solicitação HTTP, ela a traduz para um formato que o servidor TCP possa entender, ou seja, transforma os dados recebidos em uma "request" no padrão definido pela equipe. A API HTTP atua como o cliente na conexão TCP. Ao receber uma requisição HTTP da aplicação REACT, a API HTTP valida os dados e garante que a requisição é bem formada. Em seguida, a API encaminha a solicitação para o servidor, que processa a lógica necessária, e lhe retorna a "response". Após o recebimento, a API transmite a "response" ao front-end com os dados necessários.

![Protocolos de comunicação](https://github.com/ValmirNogFilho/VendePass/blob/master/docs/protocol.jpg)

Figura 6. Diagrama de comunicação entre front-end e back-end.

| Método               | Descrição                                                                                           |
|----------------------|-----------------------------------------------------------------------------------------------------|
| `login`              | Processa o login do usuário, verificando as credenciais e criando uma nova sessão se bem-sucedido.    |
| `getUserBySessionToken` | Recupera o usuário associado a um token de sessão válido.                                           |
| `logout`             | Realiza o logout do usuário, encerrando a sessão e liberando reservas associadas.                    |
| `AllRoutes`          | Retorna uma lista de todas as rotas disponíveis após validar o token de autenticação.                |
| `Route`              | Recupera uma rota entre duas cidades, validando o token de autenticação e verificando se a rota existe. |
| `Flights`            | Retorna detalhes de voos com base nos IDs fornecidos, após validar o token de autenticação.          |
| `Reservation`        | Cria uma reserva para os voos solicitados, verificando a disponibilidade e adicionando à sessão do cliente. |
| `CancelReservation`   | Cancela uma reserva de voo específica, liberando o assento e removendo a reserva da sessão do cliente. |
| `GetCart`            | Retorna o carrinho do cliente, com as reservas atuais, baseado no token de autenticação.             |
| `BuyTicket`          | Finaliza a compra de um bilhete, validando a reserva e atualizando os dados do voo e do cliente.     |
| `CancelBuy`          | Cancela a compra de um bilhete, atualizando os dados do cliente e do voo, se autorizado.            |
| `GetTickets`         | Retorna todos os bilhetes comprados pelo cliente autenticado.                                        |

Tabela 1. Métodos responsáveis por lidar com diferentes requisições de um cliente.

Esse ciclo garante que a comunicação entre o cliente e o servidor ocorra de forma eficiente e organizada, com a API atuando como uma ponte que gerencia o fluxo de requisições e respostas. 

### Formatação e tratamento de dados

Os dados trocados entre o cliente e a API, e entre a API e o servidor, foram formatados em JavaScript Object Notation (JSON). O formato JSON é amplamente utilizado por ser leve, de fácil leitura e compatível com a maioria das linguagens de programação, incluindo Go e JavaScript. Ele permite a transmissão de estruturas de dados complexas, como objetos e arrays, de forma simples e eficiente. 

No contexto da linguagem de programação Go, o JSON pode ser manipulado através das funções Marshal() e Unmarshal(). A função Marshal() é responsável por converter estruturas de dados, como structs, em um formato JSON. Isso é útil quando o servidor precisa enviar uma resposta para o cliente ou API, pois permite transformar os dados da aplicação em um formato legível para a comunicação HTTP. Já a função Unmarshal() converte dados no formato JSON em estruturas de dados Go. Isso é essencial quando o servidor recebe uma solicitação JSON e precisa interpretá-la, convertendo-a em variáveis manipuláveis pelo código. 

No contexto da implementação do protocolo HTTP, essas funções são extremamente úteis, pois facilitam a conversão de dados entre os formatos utilizados pelo servidor e pela API. Sempre que o cliente envia uma solicitação HTTP com um corpo em JSON, o servidor usa Unmarshal() para decodificar essa solicitação e tratá-la de forma adequada. Da mesma forma, ao retornar uma resposta ao cliente, o servidor utiliza Marshal() para converter os dados da resposta em JSON, que é então enviado via HTTP de forma estruturada e legível. 

Esse processo simplifica a comunicação entre as diferentes camadas do sistema, garantindo que os dados sejam sempre transmitidos de maneira consistente e eficiente entre o cliente, a API e o servidor.

### Tratamento de conexões simultâneas

Com o uso de goroutines, o servidor é capaz de lidar com requisições simultâneas dos clientes de forma eficiente. Quando um cliente realiza login, uma sessão única é criada e associada ao ID desse cliente no DAO de sessões. Isso permite que, em toda solicitação subsequente, uma nova goroutine seja iniciada para processar a requisição de forma paralela. Assim que a operação no servidor é concluída, a goroutine encerra seu processo. 

Para validar as requisições do cliente, é gerada uma string de identificação única que serve como token de autenticação. Esse token é utilizado pelo servidor para garantir que cada solicitação seja proveniente de uma sessão válida e autenticada. Quando o usuário se desconecta, a sessão correspondente ao seu ID é removida do DAO, o que também cancela todas as reservas parciais feitas durante aquela sessão. 

Além disso, o sistema não permite que um mesmo usuário se conecte simultaneamente em múltiplos dispositivos, garantindo que apenas uma sessão ativa esteja vinculada ao ID do cliente a qualquer momento.

### Tratamento de concorrência

Na linguagem de programação Go, as goroutines são sincronizadas por meio de channels. Um channel é um mecanismo de comunicação que permite que goroutines troquem dados entre si de maneira segura e sincronizada, evitando "condições de corrida" (situações prejudiciais onde diferentes threads consumam os mesmos recursos ao mesmo tempo). Eles são utilizados para passar informações de uma goroutine para outra, garantindo que o acesso a recursos compartilhados seja controlado e ordenado. Quando combinados com filas e mutexes, os channels permitem gerenciar a concorrência e evitar problemas como condições de corrida, onde múltiplas goroutines tentam acessar o mesmo recurso simultaneamente.

Sem a sincronização adequada entre as goroutines, podem ocorrer erros em cenários de concorrência. Por exemplo, se dois clientes tentarem comprar o último assento disponível em um voo, sem o controle adequado, ambos poderiam acabar conseguindo reservar a mesma passagem, gerando uma inconsistência no sistema. 

Para evitar isso, foram utilizados channels para implementar filas de requisições e mutexes para garantir a exclusão mútua ao acessar recursos compartilhados, como a quantidade de assentos disponíveis. Quando um cliente faz uma solicitação para reservar um assento, o channel controla a ordem das requisições e o mutex garante que apenas uma goroutine por vez possa acessar e modificar os dados do voo. Assim, a concorrência é tratada de maneira segura, prevenindo que dois clientes reservem o mesmo assento ao mesmo tempo. 

Além disso, foi implementada uma etapa extra para aumentar a segurança da reserva, utilizando o conceito de reserva parcial. Quando o cliente adiciona uma passagem ao carrinho, o sistema temporariamente reserva essa passagem para ele até a sua sessão associada expirar ou por até 30 minutos. Durante esse período, a quantidade de assentos disponíveis é subtraída, garantindo que o assento não seja oferecido a outro cliente. Se o cliente desistir da compra ou se desconectar, o assento é liberado, e a quantidade de assentos disponíveis é incrementada novamente, permitindo que outros clientes façam a reserva. Se a compra for efetivada dentro do prazo, a reserva se torna permanente e a quantidade de assentos permanece subtraída. 

![fig5](docs/4.png)

Figura 7. Carrinho de um dos clientes.

![fig6](docs/5.jpg)

Figura 8. Indicação de que não é possível reservar um voo já reservado por outro cliente.

### Desempenho e avaliação

Para aumentar o desempenho do sistema, múltiplas goroutines são criadas e executadas em paralelo. Essas goroutines são responsáveis por várias tarefas, como o gerenciamento de reservas, controle de sessões, além de "ouvir" e responder às solicitações dos clientes com consultas aos bancos de dados. O uso de goroutines permite que o servidor lide com muitas requisições simultaneamente, aproveitando o paralelismo e distribuindo as tarefas de maneira eficiente.

Para coordenar essas goroutines e garantir a comunicação segura entre elas, são utilizados channels. Os channels funcionam como um meio de comunicação entre as goroutines, permitindo que dados sejam enviados de uma para outra de forma sincronizada. No contexto do sistema, channels são usados para passar informações, como requisições de reserva ou atualizações de sessão. Dessa forma, o uso de channels e mutexes garante que a execução paralela das goroutines maximize o desempenho e a capacidade de resposta do servidor, tornando-o mais eficiente e escalável. 

### Confiabilidade da solução

O servidor implementa vários mecanismos à prova de falhas para garantir a robustez do sistema. Por exemplo, ele associa um período de inatividade a todas as sessões e registra o último tempo ativo de cada uma. Se o tempo atual ultrapassar o último registro de atividade mais o tempo máximo de inatividade permitido, o cliente é automaticamente desconectado, garantindo que recursos não sejam mantidos ocupados por sessões inativas.

No entanto, o servidor não foi projetado para lidar com encerramentos abruptos. Durante a execução, ele utiliza o DAO (Data Access Object) sem persistência de memória, o que significa que todos os dados em memória são voláteis. Em caso de uma interrupção inesperada, como uma falha de energia ou um erro no sistema, todos os dados em uso são perdidos. Para mitigar essa limitação, o servidor é configurado para recuperar dados mockados do diretório de stubs após ser reiniciado. Essa abordagem permite que o sistema retorne a um estado funcional, embora com dados temporários, o que ajuda a minimizar o impacto da perda de informações durante um encerramento inesperado. 

Embora essa estratégia ajude na recuperação do sistema, é recomendável que, em versões futuras, sejam implementadas soluções de persistência de dados mais robustas, como bancos de dados que garantam a integridade e a persistência dos dados mesmo em casos de falhas no servidor.

### Documentação do código

O código do projeto foi devidamente organizado em diretórios distintos, seguindo o padrão Model-View-Controller (MVC), e documentado em inglês. Essa separação por camadas permite que a lógica de negócio (Model), a interface do usuário (View) e o controle de fluxo de dados (Controller) fiquem isolados, facilitando a manutenção e a evolução do sistema. A documentação das funções foi feita de maneira detalhada, e pode ser consultada com o comando go doc nomeDoPacote para verificar as funcionalidades de cada pacote, ou com o comando go doc nomeDoPacote.nomeDaFuncao para acessar a descrição individual de cada função.

Essa abordagem melhora significativamente a modularização do código. Ao separar as responsabilidades em diferentes pacotes, o código se torna mais organizado e compreensível, facilitando a reutilização de componentes e a adição de novas funcionalidades sem comprometer outras partes do sistema.

### Emprego do Docker

Para garantir portabilidade dos serviços do sistema em diferentes sistemas operacionais e computadores, o software conta com a conteinerização dos seus componentes. Foi utilizado o Docker como ferramenta para criar as imagens do sistema, com suas dependências necessárias (compiladores, bibliotecas e demais recursos). O software está dividido em três imagens, e duas delas estão definidas no Dockerfile da pasta raiz do projeto, contendo:
- A aplicação do servidor com a API Socket Básica para comunicação TCP;
- A API HTTP utilizada como comunicador intermediário entre a aplicação React e o servidor TCP.

Ambas são feitas com base na imagem do compilador "Go", versão 22.

A terceira imagem Docker está descrita pelo Dockerfile na pasta `ui`, com base na imagem do Node 18, para tornar possível a execução do servidor Vite para a aplicação React.

Para tornar mais prática a conteinerização das imagens, o arquivo `docker-compose.yaml` executa o build dos arquivos Dockerfile e cria um "network" para tornar a comunicação entre os três conteineres possível, através de uma rede virtual. A "network" segue o padrão de driver "bridge".

Para executar a conteinerização e execução do sistema, é necessário ter a ferramenta Docker instalada no computador e digitar no terminal, na pasta raiz do projeto, o comando `docker compose up --build`. Para encerrar a execução, deve ser executado `docker compose down`.

Para a execução para testes, sem uso do docker, o arquivo shell `start.sh`, na pasta raiz do projeto, contém as instruções de execução sem conteinerização do projeto. Para usá-las, deve ser executado os comandos `chmod +x start.sh`, para dar permissão de execução ao arquivo shell e `./start.sh`, na pasta raiz do projeto.

### Conclusão

O software desenvolvido atende aos requisitos definidos pelo Aprendizado Baseado em Problemas. O sistema permite a execução de um servidor capaz de receber diversas requisições de diferentes clientes em simultâneo. As estruturas de channels e mutexes em go assegura que as estruturas de dados sejam acessadas por apenas um cliente de cada vez, evitando "condições de corrida" por recursos de memória. A solução é eficiente por permitir a concorrência de goroutines, habilitando pseudo-paralelismo das tarefas do servidor, acelerando o "throughput". O software também conta com uma interface gráfica intuitiva, como solicitado no problema, além de possuir design moderno, seguindo princípios de UX Design.

Conforme a plausível evolução do projeto, possíveis melhorias incluiriam:
- A estruturação de todo servidor exclusivamente no protocolo HTTP (o que não foi feito, devido as exigências do projeto em relação ao uso do TCP);
- Melhorias na estilização da interface gráfica, voltadas a design responsivo para diferentes telas;
- Maiores baterias de testes, visando analisar casos de erros mais aprofundados ou em maior densidade de usuários simultâneos;

A experiência da produção do sistema foi proveitosa, por trazer a elucidação de como ferramentas modernas de desenvolvimento web back-end funcionam "por baixo dos panos". O aprendizado prático sobre algoritmos de sincronização e testagem para lidar com sistemas distribuídos foi determinante. Além disso, outro detalhe importante para a evolução profissional/acadêmica dos desenvolvedores foi a oportunidade do uso de ferramentas modernas de "front-end", no caso do React, e "DevOps" para conteinerização, com Docker, altamente utilizados na atualidade.

### Referências

Holz, Leonhard. Well-structured Logic: A Golang OOP Tutorial. Toptal. Disponível em: https://www.toptal.com/golang/golang-oop-tutorial. Acesso em: 22 set. 2024.

MDN Web Docs. HTTP: Overview. Disponível em: https://developer.mozilla.org/pt-BR/docs/Web/HTTP/Overview. Acesso em: 22 set. 2024.

MDN Web Docs. Client-Server Overview. Disponível em: https://developer.mozilla.org/pt-BR/docs/Learn/Server-side/First_steps/Client-Server_overview. Acesso em: 22 set. 2024.

Red Hat. Stateful vs Stateless. Disponível em: https://www.redhat.com/pt-br/topics/cloud-native-apps/stateful-vs-stateless. Acesso em: 22 set. 2024.
