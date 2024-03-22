#!/bin/bash
export CONN="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/pggo"
cd ./sql/schema
goose postgres "$CONN" down
