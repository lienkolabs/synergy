# Synergy Protocol

## Chamado

De tempos em tempos indivíduos inquietos, inconformados decidem agrupar-se e 
estabelecerem entre si novos meios de coordenação para ação coletiva, numa
iniciativa de afronta ao status quo. 

Foi assim que, por exemplo, a chamada Reública das Letras pariu, através do 
estabelecimento de sociedades e periódicos, a revolução científica. Foi assim 
que a ARPANET, através da consolidação do mecanismo do RFC, pariu a internet
em padrões abertos. Foi assim que riot grrrrl pariu, através dos zines, a 
terceira onda do feminismo. Os exemplos abundam. Todos contemplam um núcleo
comum: experimentalismo ao invés de dogma, reputação ao invés de autoridade, 
autonomia ao invés de controle. 

Acreditamos que é chegado o momento em que os inquietos, os inconformados, hão
de agrupar-se e estabelecerem meios de coordenação para ação coletiva, numa 
iniciativa de afronta à mídia social e as plataformas de tecnologia. 

Acreditamos que estes novos meios tem que, obrigatoriamente, ser construídos 
fora das plataformas vigentes. Acreditamos que a tarefa de reinventar a 
internet, desta vez como internet verdadeiramente pessoal, é uma tarefa que 
requer múltiplos talentos e diversidade. Não é mais admissível o reinado dos 
interesses do capital, ou o viés dos tecnólogos. A construção da nova internet
é uma tarefa para todos nós.

## Visão



## Introduction

Synergy é um novo protocolo social concebido para coordenar a ação coletiva
através do estabelecimento de padrões de cooperação. 

Ele é uma expansão das funcionalidades básicas fornecidas pelo protocolo social
axé. 

Coletivo é uma entidade nomeada que consiste em cada momento num conjunto de 
indivíduos em acordo com uma política de ação conjunta em vista de um objetivo
comum explicitado. Ela fornece ao mesmo tempo uma possibilidade de reputação
coletiva, uma vez que suas ações promovem ou depreciam o nome da colatividade.

Uma assembléia consiste num grupo de indivíduos que tem acertado entre si, de 
comum e livre acordo, regras de consenso para a ação coletiva. Noutras palavras,
estabelecem requesitos mínimos para que todos estejam de acordo com que ações 

Um coletivo é um assembléia associada a um nome:

Coletivo = {Nome,Assembléia}

Havendo nome, as ações autorizadas pela assembléia do coletivo podem ao longo do 
tempo promoder ou depreciar a reputação do mesmo. Noutras palavras, espera-se 
que ao contrário da assembléia, em que a responsabilidade das ações coletivas 
recaem essencialmente sobre os indivíduos, num coletivo espera-se que a 
responsabilidade 



Coletivo = {{Executiva,Maioria},{}}

Coletivo = {Membros,Maioria,Super-maioria}

Coletivo Nomeado = {Nome, Coletivo}

Coletivos podem ser formados em propósito genérico, quando na lógica do 
protocolo, formada as condições de consenso eles podem realizar quaisquer ações
protocolares possíveis aos indivíduos, excessão feita à votação (não está 
prevista Coletivos de Coletivos, por exemplo).

Coletivos podem ser formados para a exclusiva 


Its is built on top of authoriship funcionalities provided by the axé base 
social protocol. Besides individual


Synergy is a social protocol to organize collaborative patterns to advance the
boundaries of collective knowledge of different subjects. 

It is constructed on top of a key basic objects which are built on top of 
primitive functionalitiers of axé social protocol (id and stages). Any object can
have an unique handle and a smal textual description. 

> Collective = {{Owners,Majority},{Moderators, Majority},Members}

A collective is any collection of individuals, some of which are designated owners
with the exclusive capacity to appoint and remove moderators. Moderators are 
grante the right to approve or remove members and perform all the other actions
prescribe to Colletive objects in the protocol. 

Majority is a rule governing the capacities of a subset of the designated authorities

Collectives can manage subjects 

> Subject = {collective,[parent subjects]}

If parent subjects are not *nil*, then their approval are required on the creation
of a subject and they reserve the right to reject the relationship at any time.

Information on synergy is abstracted on a media object 

> Media = {content-nature,content-type,content-channel,references,[{author,capacity}],hash-of-previous version}

If a media object refers to a previous version, the authors of the previous version must approve the link. The authors of the new version can be appended on the capacity of co-authors, collaborators or contributors. Content-nature can be: question, denounce, recomend, propose, release, pin. 

A hash-linked sequence of media objects is a draft. A draft can be submitted as 
pre-print on a number of sujects. Moderators mujst approve the draft, but this 
should be understood not as an approval on the content of the draft, but rather on
the its form. At a very broad analysis it seems to have the required style and
relevance to the subject matter. 

In order to get approval by peer-review a pre-print must be sent to a Journal

> Journal = {collective,[subjects],review-policy}

On the collective, moderators could be interpreted as editors, and members as 
peer-reviewers. If a journal accept a submitted pre-print it is telling the 
community that it has looked carefully on the content of the media and in good
faith attests that it is correct as by the standards of the subject. 

Besides media, that is the basic assyncronous communication primitive, there is a synchronous primitive in synergy

> Event = {venue,start, duration, scope}

A meeting is an event associated to a collective

> Meeting = {collective,event}

A festival is a list of events associated to a subject 

> Festival = {subject, [events]}

The system runs on reputation systems. Reputation system ranks different authors
based on their titles (belonging to collectives), their publications and their
participation on events. 

