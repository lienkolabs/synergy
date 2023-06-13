From time to time groups of worried and bothered individuals take action for 
themselves to establish new media for the coordination of action and sharing ideas on a 
commitment to face the status quo.

Examples abound, of greater or lesser consequence. For example, the so called 
Republic of Letters gave birth, through the culmination of scientific societies
and their periodicals, to the scientific revolution. ARPANET gave birth, through 
the internet task force and the RFC publication, to the open standards internet.
More recently, riot grrrl gave birth, through their zine, to the third wave
of feminism. All of these share a common theme: experimentalism rather than dogma,
reputation rather than authority, autonomy rather than control.  

We believe that time has come for the worried and bothered of our era to stand 
against technology plataforms. And we are proposing a new media, the Synergy 
Social Protocol, in order to facilitate coordination of action and sharing of 
ideas.

We strongly believe that any social media worthy of its name must be 
invented in itself, and must be invented socially. We cannot continue to accept
that the terms of our digital experince are to be dictate by a tiny group of 
individuals with their peculiar cultural biases, worldview and motivations. 
We all must build a new, personal, internet!

# Synergy Protocol

## Overview

Synergy protocol was designed as a digital framework for collaboration and
collective construction. Protocol functionalities were inspired by the - arguably - most successful form of human cooperation: scientific publishing.

Synergy was built on top of axé social protocol, which provides primitives
for identity and stage management.

## General structure

The most fundamental functionalities for collaborative project development 
derive from two basic constructions: groups of individuals with a common goal,
which are called COLLECTIVE, and content creation dynamics, that this protocol
enables by the evolution of DRAFTs. 

A COLLECTIVE's ability to perform actions as a single entity makes it easy for
groups of people to act as a unity. For every action taken on behalf of a 
COLLECTIVE, the protocol automatically triggers a voting mechanism.

A DRAFT's evolution is acconted for by its EDIT, RELEASE and STAMP, this last action
is performed in the sense of a peer review, and can be taken on behalf of a COLLECTIVE.

### Collective action

COLLECTIVE are an association of members that have a common goal. Each
COLLECTIVE is created with a name that is unique within the network and
must not be a handle (cannot be an @), and with a description of its goal. 

Besides the name and description, upon being created a COLLECTIVE must provide 
its choice policy for pooling, which will dictate the voting mechanism 
for the actions that will be performed in its name. It must also provide a "super" 
policty for policy update, meaning the policy for changing the COLLECTIVE voting policy.

Every action taken on behalf of the collective triggers a voting mechanism
according to the COLLECTIVE policy. Whenever an instruction is posted on behalf 
of a COLLECTIVE the pooling mechanism is automatically triggered and its result 
is accepted as the decision of the COLLECTIVE.

If pool results in acceptance of the instruction, instruction in considered valid
and on behalf of the pointed COLLECTIVE. If an instruction is posted on behalf of 
a COLLECTIVE and the instruction’s pool results in non-acceptance of the instruction,
the instruction is discarded. 

The pooling mechanism of a COLLECTIVE may be updated.

Any of the actions prescribed by the protocol can be submitted as
being on behalf of a COLLECTIVE. That means COLLECTIVE can act as a
unit in the network. 

Everyone can apply to join a COLLECTIVE. Upon applying, a voting amongst COLLECTIVE 
members is taken to either accept or deny the request.


## Information dynamics

To account for basic elements for information exchange, the protocol provides
the following instructions:

### DRAFT

It is the basic element used for publishing ideas and contributions that are in
progress. DRAFTs are to be used for public idea elaboration and collaboration,
amongst network members.

By creating a DRAFT the member, or group of members, is sharing a starting proposition with
the community, and anyone who wishes to contribute to the forming of it can
apply for an new version of the DRAFT by adding or removing info from the
original DRAFT publishing and referencing to the old version as deprecated.
Any member can propose a new DRAFT as an individual contribution or
on behalf of a COLLECTIVE.

Besides the actual DRAFT content, DRAFT instructions must include a
title for easy identification (it does not need to be unique), a brief description
of the content, and a list of keywords the content is related to. It may, or
may not include a list of internal references used (content previously posted on
the network), if the DRAFT posted is a new version of a previously posted DRAFT, 
it must contain the hash of the previous DRAFT as its predecessor.

DRAFT instructions are necessarily public. All information published as a
DRAFT can be viewed and revised by the whole community.

### JOURNAL

JOURNALs are created for members to endorse content they have reviewed and
wish to promote. Any single member or COLLECTIVE can create a JOURNAL,
therefore being responsible for managing the content that is accepted by it.

JOURNALs are a sort of stamp to a content, and that makes it easier for
community members to know which content has been reviewed and is being
vouched by peers.

If a single member created the JOURNAL, there is no need to point to a
POOL mechanism for content acceptance, since the member will directly decide.
For JOURNALs created on behalf of a COLLECTIVE, content acceptance will
follow the COLLECTIVE POOL mechanism.

JOURNALs must have a unique name/title (it must not be an @) and a
brief description of its purpose, so it is easy for members to decide upon which
journal to follow.

###  POST

POSTs are DRAFTs that have reached a final form and are, therefore, applying
for the endorsement of a JOURNAL.

They differ from DRAFTs only on two points: POSTs must include a list of
JOURNALs they are applying to (the list must contain at least one JOURNAL)
and POSTs have no history track, since they are a final version and not a new
version of a DRAFT.

If accepted by at least one of the JOURNALs it applied to, the POST is
considered valid. Otherwise, it is considered invalid.

### REACTION

To be used as both a ranking tool for DRAFT or POST contents and a means
for the community to express its interest.

Members can react either positively or negatively to both DRAFTs and
POSTs instructions.

REACTIONs can be either signed by a single member, as an individual, or
on behalf of a COLLECTIVE. If signed on behalf of a COLLECTIVE, a POOL
according to the COLLECTIVE’s pooling mechanism is automatically generated. 
The REACTION is only considered valid if the pool results in acceptance
of the instruction.

### BOARD

BOARDs are instructions that provide a keyword, or unique group of keywords,
to index DRAFTs that were posted referencing the BOARD’s dedicated keyword
or group of keywords.

They are both a means for the community to easily find DRAFTs by their
keywords, and a place for content that has not yet been reviewed by peers to
get visibility.

Content indexed by a BOARD has not necessarily been reviewed by board’s
creator(s).

BOARDs can be created by either a single member or on behalf of a COLLECTIVE.
Each keyword or set of keywords can be associated to a single BOARD, so
the indexing is kept organized.

### EVENT

Synchronous exchange of information is sometimes needed for collective creation.
EVENT is a digital space for live discussion on any given topic where there
may, or may not, be an audience that does not interact with discussion.
Given the synchronous aspect of the EVENTs, they must contain a starting
and finishing time, besides a name for easy identification, a brief description
of their goal. For private EVENTs, the instruction must contain the list of
members with interaction privileges. If the EVENT allows for an audience to
watch, it may contain either a list of non-interacting allowed members or a
public key if it is open for public audience.

EVENTs can only be created by a COLLECTIVE.