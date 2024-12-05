# ANA Pay importer for Pocketsmith

A simple CLI to import ANA Pay transactions into Pocketsmith.

## Usage


go run main.go -username=YOUR_ANAPAY_USERNAME -password=YOUR_ANAPAY_PASSWORD -token=YOUR_POCKETSMITH_TOKEN


Or set environment variables:


export ANAPAY_USERNAME=YOUR_ANAPAY_USERNAME
export ANAPAY_PASSWORD=YOUR_ANAPAY_PASSWORD
export POCKETSMITH_TOKEN=YOUR_POCKETSMITH_TOKEN

go run main.go


### Run with docker (recommended)

`docker run -e ANAPAY_USERNAME=xxx -e ANAPAY_PASSWORD=xxx -e POCKETSMITH_TOKEN=xxx dvcrn/pocketsmith-anapay`


### Optional flags

- `-num-transactions`: Number of transactions to fetch (default: 100)

## Features

- Automatically creates ANA Pay institution and account in Pocketsmith if they don't exist
- Updates account balance
- Imports transactions with proper categorization for:
  - Charges
  - Cashback
  - Credit card payments
  - Apple Pay
  - Auto-charges
  - Virtual prepaid card
  - VISA touch payments
  - iD touch payments
- Prevents duplicate transactions
