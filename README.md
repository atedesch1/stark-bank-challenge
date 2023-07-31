# Stark Bank Backend Trial
This app sends 8 to 12 invoices to random people and listens for an invoice webhook on a separate go routine.
The webhook sends a transfer equal to the invoice amount to Stark Bank.

## Deploy
This app was deployed to GCP and is listening for the invoice webhook on [this url](https://stark-bank-test-3fzbrnzjyq-rj.a.run.app/invoicehook).

## Continuous deployment
On push to master this repository triggers a build for the app's docker image which is pushed to GCP and served immediately.

## Running locally
Add a .env file in the root directory with the following environment variables:
- PROJECT_ID (the stark bank project id)
- PRIVATE_KEY (the stark bank private key)
- ENV=sandbox

After setting the environment variables build and run the application with make:
```
make run
```
The invoices will be sent correctly however stark bank's webhook won't be able to reach your machine as you don't have any exposed ports.
The application being served on GCP will handle the webhooks.