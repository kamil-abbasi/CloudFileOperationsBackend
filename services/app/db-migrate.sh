#!/bin/bash

migrate -database postgresql://admin:admin@database/app?sslmode=disable -path ./db/migrations up