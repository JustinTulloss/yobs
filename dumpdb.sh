#!/bin/sh

pg_dump -s yobs > schema.sql
pg_dump --column-inserts --data-only --table=users --table=transactions yobs > seed.sql
