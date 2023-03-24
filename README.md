# Synergy Protocol

## Introduction

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

> Media = {content-type,content-channel,references,[{author,capacity}],hash-of-previous version}

If a media object refers to a previous version, the authors of the previous version
must approve the link. The authors of the new version can be appended on the capacity of
co-authors, collaborators or contributors. 

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

