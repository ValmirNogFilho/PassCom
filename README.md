# PassCom - Sistema distribuído de venda de passagens aéreas

## Introdução

No último século, o advento do comércio digital (também chamado de "E-commerce") se propagou de maneira exponencial. Isso ocorreu principalmente nos últimos anos, devido a variáveis como o período pandêmico, que forçou a digitalização do comércio. As plataformas digitais fornecem qualidade de serviço, atendimento prático, rápido e automatizado, proporcionando melhor experiência aos clientes.

O presente relatório visa descrever o desenvolvimento do sistema "Passcom", um sistema de venda de passagens distribuído para três companhias aéreas (Rumos, Giro e Boreal). O software foi desenvolvido utilizando a arquitetura REST, com servidores HTTP individualizados codificados em Go, com acesso a banco de dados SQLite via uso da biblioteca GORM para mapeamento relacional de objetos. O front-end foi desenvolvido com o framework React. Também foi adicionada uma interface CLI acessível via TCP para monitoramento e configuração dos servidores. O sistema é conteinerizado com Docker, permitindo consistência no funcionamento.

O software assegura a sincronização distribuída a partir do uso de diversos algoritmos de consenso e protocolos a serem explicados no relatório. O resultado atendeu os requisitos, pela confecção de um projeto bem arquiteturado, robusto que visa a praticidade do uso para os usuários.

### Equipe:

[![Pierre Machado](https://github.com/pierremachado.png?size=20)](https://github.com/pierremachado) [Pierre Machado](https://github.com/pierremachado)

[![Valmir Nogueira](https://github.com/valmirnogfilho.png?size=20)](https://github.com/valmirnogfilho) [Valmir Nogueira](https://github.com/valmirnogfilho)

## Arquitetura da solução

Stateless. Os servidores foram projetados para operar de forma independente e escalável. A conexão entre os servidores é feita utilizando um cliente TCP, e, assim que os servidores se conectam, o sistema verifica periodicamente o status da conexão através de heartbeats.

De forma independente, o diagrama de sequência do sistema mostra a comunicação entre o cliente e o servidor. Nessa perspectiva, é possível visualizar o sistema distribuído como se fosse apenas um servidor centralizado.

Quando o servidor de duas companhias estão conectados, eles trocam o banco de dados local entre si, tornando-se réplicas de si mesmo. Essa conexão é exemplificada pelo diagrama de sequência a seguir.

A sincronização de dados é feita através de broadcast do voo quando é editado em seu servidor de origem.

Para implementar o sistema PassCom, o projeto foi feito em Go e a interface gráfica em React.

O sistema distribuído foi projetado tendo em mente as seguintes assunções:

- Assume-se que qualquer um dos nós pode falhar a qualquer momento de forma permanente (por falha de hardware, desligamento), ou que a latência entre os servidores impediria a comunicação devida. O sistema se reorganizaria para continuar funcionando sem o nó ausente.
- Tempo parcialmente-síncrono: o sistema assume que a comunicação entre os servidores é rápida e não há latência significativa. Entretanto, pode ocorrer que as mensagens cheguem de forma atrasada.
- Teorema CAP (CAP Theorem): teorema fundamental de sistemas distribuídos que dita que estes não podem garantir a consistência, disponibilidade e a partição de rede do sistema simultaneamente. Logo, a disponibilidade foi sacrificada - se um dos nós falharem, mesmo que os outros servidores continuem funcionando normalmente e tenha os dados deste servidor, a venda não será efetuada por nenhum nó.

## Protocolo de comunicação

HTTP. Os servidores das companhias possuem endpoints que utilizam para se conectar e conectar com os seus clientes. São eles:

(tabela com endpoints)

Através de solicitações get, post, put e delete, são capazes de organizar a compra de passagens entre clientes e servidores.

O diagrama de sequência a seguir mostra o fluxo de comunicação entre os servidores quando o cliente deseja comprar uma passagem.

## Roteamento

Como cada servidor mantém uma réplica do banco de dados do outro quando eles se conectam, é possível determinar rotas entre os trechos dos servidores das companhias que estão conectados.

O algoritmo realiza uma BFS entre rotas e retorna uma rota contendo as rotas possíveis entre os trechos dos servidores das companhias que estão conectados.

## Concorrência Distribuída

Quando um voo é editado, todos os servidores das companhias que estão conectados com os outros servidores do sistema notificam aos nós conectados essa alteração através de um broadcast; essa operação é idempotente, pois o servidor não pede para decrementar em um a quantidade de assentos, e sim envia o estado atual do voo e pede para os nós conectados substituirem as informações atuais. Portanto, se a mensagem chegar múltiplas vezes, será alterado apenas uma vez.

Além disso, o servidor não permite a venda da passagem de outro servidor que esteja offline, pois parte do pressuposto que não é possível determinar se o problema está localizado na rede ou se o servidor caiu.

Algoritmos de consenso que possuem como alicerce a eleição de nós lideres, como Paxos e Raft, foram cogitados para o projeto. Todavia, a implementação destes foi descartada. Se deve ao fato de que algoritmos de consenso dessa forma impediria que os servidores pudessem operar de forma independente assim que não fosse possível se conectar a um quorum de servidores operando e recebendo mensagens.

Portanto, o sistema permite que as informações sobre os vôos de outras companhias, como a quantidade de assentos disponíveis, estejam desatualizadas eventualmente caso a mensagem de broadcast não chegue ao servidor destinatário. Entretanto, nenhum servidor pode vender a passagem de um servidor que esteja offline.

Para tratar a eventual concorrência de dois clientes tentando comprar o mesmo assento, o sistema implementa locks otimistas. O lock acontece apenas no momento da transação ou no envio de uma mensagem, e, caso resulte em erro, a transação é cancelada e o cliente é notificado.

## Confiabilidade da solução

No momento atual, o sistema PassCom possui algumas vulnerabilidades. Atualmente, não há um algoritmo de consenso confiável implementado para o sistema. Isso faz com que, caso as informações cheguem de forma inconsistente, os dados dos outros servidores podem aparecer desatualizados para o cliente: um assento de outro servidor pode estar marcado como disponível para um cliente local, mas os assentos do outro servidor podem estar marcados como indisponíveis para o cliente do servidor em questão. Em ambos os casos, a transação resultará em um erro. 

Além disso, se um dos servidores falhe durante uma transação, é possível que a quantidade de assentos disponíveis seja decrementada, mas o servidor remetente não receba a mensagem de confirmação da transação. Para resolver esse problema, um algoritmo de consenso confiável ou de transações distribuídas, como Two-Phase Commit (2PC), poderia ser implementado para garantir a confiabilidade da solução. Por fim, cada servidor dispõe de relógios vetoriais para determinar a causalidade e ordem dos eventos, mas não os utilizam.

## Avaliação da Solução

O código foi mantido no GitHub para o seu desenvolvimento e testado através da interface gráfica do projeto.

## Documentação do código

O código possui a documentação dos endpoints HTTP que ditam a comunicação do servidor com o cliente.

## Emprego do Docker

O sistema completo foi conteinerizado via uso do Docker. Cada um dos servidores separados possui um Dockerfile, com instruções de diretórios a serem copiados, portas a serem expostas e volumes de persistência de dados a serem considerados (arquivos JSON e os arquivos de database SQLite). Também foram criados contêineres para execução das interfaces React, e a comunicação entre front-end e back-end pelas APIs foram asseguradas pelas networks criadas. Os Dockerfiles dos servidores também expõem as portas para acesso ao server CLI de monitoramento dos servidores REST.

Assim, o arquivo `docker-compose.yaml` une a execução dos contêineres, permitindo o build e execução dos componentes de cada uma das companhias aéreas a partir do comando:

```
docker compose up --build
```

## Conclusão

Conclui-se que o sistema PassCom, apesar de possuir vulnerabilidades, é funcional e requeriu o estudo de uma série de técnicas de arquitetura, protocolo de comunicação, roteamento, concorrência distribuída e confiabilidade. Algumas melhorias podem ser feitas, como a implementação de algoritmos de consenso confiáveis e o reforço da robustez do sistema em caso de falhas, como, por exemplo, armazenar os logs das transações pendentes em disco.

Ademais, alguns tópicos podem ser revisitados futuramente, como:


 Espera-se que a implementação do sistema PassCom contribua positivamente para a compreensão e aperfeiçoamento dos sistemas distribuídos no futuro.