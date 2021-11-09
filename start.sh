#!/bin/bash

# Preliminary set up
mkdir -p certs static templates
touch talks.db

docker-compose -f docker-compose.dev.yml up --build
