# Synergy Protocol

## Overview
Synergy protocol was designed as a digital framework for collaboration and
collective construction. Protocol functionalities were inspired by the - arguably
- most successful form of human cooperation: scientific publishing.
Synergy was built on top of axé social protocol, which provides primitives
for identity and stage management.

## General structure
The most fundamental functionalities for collaborative project development 
derive from two basic constructions: groups of individuals with a common goal,
which are called COLLECTIVE by this protocol, and these groups’ ability to
perform actions as a single entity. For COLLECTIVE actions, the protocol
provides a voting mechanism called POOL.

### COLLECTIVE
COLLECTIVE are an association of members that have a common goal. Each
COLLECTIVE is created with a name that is unique within the network and
must not be a handle (cannot be an @), and with a description of its goal.

A COLLECTIVE goal might have a deadline, but that is not mandatory.
All members that wish to take part in the goal can apply to become part
of the COLLECTIVE. The COLLECTIVE has a list of members to which new
member’s signatures are appended once they are accepted.

Any of the actions prescribed by the Synergy protocol can be submitted as
being on behalf of a COLLECTIVE. That means COLLECTIVE can act as a
unit in the network. They may create, react, index and review content, submit
content to peer review, create events and pools.

Upon being created, a COLLECTIVE must provide its choice mechanism
for pooling, which will dictate the voting mechanism for the actions that will
be performed in its name.
Whenever an instruction is posted on behalf of a COLLECTIVE the pooling
mechanism is automatically activated and its result is accepted as the decision
of the COLLECTIVE.

If an instruction is posted on behalf of a COLLECTIVE and the instruction’s
pool results in non-acceptance of the instruction, the instruction is discarded.
The pooling mechanism of a COLLECTIVE may updated by the result of a
POOL with this purpose.

###  POOL

POOLs are voting mechanisms for any group of members to collectively decide
or express their opinion.

Each COLLECTIVE created must point to a POOL mechanism previously
created.

POOL instructions prescribe a list of participating members, a list of options
to be chosen from, a number of points to be distributed amongst the list of
options, a maximum number of points that each option might get, a counting
method and a winner criteria.

POOL instructions also may, or may not, have a deadline. For pools with a
deadline, POOL result is achieved by the end of the POOL period. For pools
with no deadline, POOL result is achieved once every member of the prescribed
list has voted.

## Information exchange instructions

Given the basic functionalities of COLLECTIVE and POOL, the protocol also
provides a list of functionalities for information exchange.

Each instruction can be signed by only a single network member and may, or
may not, be on behalf of a COLLECTIVE. If the COLLECTIVE field is filled
with the COLLECTIVE name, a POOL is automatically activated.

If pool results in acceptance of the instruction, instruction in considered valid
and on behalf of the pointed COLLECTIVE. If POOL results in non-acceptance
of the instruction, instruction is discarded.

To account for basic elements for information exchange, the protocol provides
the following instructions:

### DRAFT

It is the basic element used for publishing ideas and contributions that are in
progress. DRAFTs are to be used for public idea elaboration and collaboration,
amongst network members.

By creating a DRAFT the member is sharing a starting proposition with
the community, and anyone who wishes to contribute to the forming of it can
apply for an new version of the DRAFT by adding or removing info from the
original DRAFT publishing and referencing to the old version as deprecated.
Any member can propose a new DRAFT as an individual contribution or
on behalf of a COLLECTIVE.

Besides the actual DRAFT content, DRAFT instructions must include a
title for easy identification (it does not need to be unique), a brief description
of the content, and a list of keywords the content is related to. It may, or
may not include a list of internal references used (content previously posted on
the network), a list of external references and, if the DRAFT posted is a new
version of a previously posted DRAFT, it must contain the hash of the previous
DRAFT as its predecessor.

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