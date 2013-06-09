#!/bin/sh

dropdb yobs && createdb yobs && psql yobs < schema.sql
