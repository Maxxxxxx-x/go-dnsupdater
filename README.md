# go-dynamicdns

A code that aims to automatically update DNS A records for those using dynamic IPs, reducing the amount of headaches that are caused by IP rotation.
Currently only works with Cloudflare


# Why?
I want to write one

# Prerequisite
- go v1.23.1
- cloudflare account
- cloudflare token with DNS read and edit access

# config
check .env.example

# build
make build/prod

# TODO
- daemon? (or it can be a cronjob)
