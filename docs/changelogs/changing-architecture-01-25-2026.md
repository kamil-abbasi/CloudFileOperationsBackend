# Introduction

In this blog post I'm going to cover the following things:

- current architecture

- new architecture

- why I want to change the architecture

- pros and cons of each architecture

# Current Architecture

As for now user files and their metadata are stored in the container filesystem under a specific path. This directory then gets mounted to an external docker volume. Metadata is stored in SQLITE database. Each user gets their own directory in which the data is stored. The goal was to dynamically create volume for the user upon registration, so files are isolated.

Pros:

- file isolation out of the box (separate volumes)

Cons:

- to dynamically create volumes container must have access to docker on the host (not secure, defeats one of the docker purposes which is isolation)

- coupling with docker

- good for small scale, but storing metadata for large amount of files in a sqlite database simply won't work

- complexity (accessing docker from the container, managing volumes)

- hard to set qoutas (volumes don't natively support quotas)

- almost impossible multi-server deployment (because of volumes)

- migration to cloud is very hard (requires complete architecture change)

# New Architecture

In the new architecture files will be stored inside an object storage (for example Amazon S3 or Minio) and their metadata in SQL database like Postgres. All user files will be stored together and isolation will be achieved using software.

Pros:

- horizontal scaling out of the box (with Amazon S3, with minio requires a little configuration but still easy)

- secure

- storing large amount of file metadata is easy thanks to server with SQL database

- cloud ready (Minio is compatible with Amazon S3)

Cons:

- isolation is not provided out of the box, but still easily achievable using software

# Why

I want to change the architecture because new one offers less complexity, better scaling and is more secure. It is better than current architecture in almost every way.