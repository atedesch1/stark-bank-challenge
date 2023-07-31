# Stark Bank Backend Trial
This app sends 8 to 12 invoices to random people and listens for an invoice webhook.
The webhook sends a transfer equal to the invoice amount to Stark Bank.

## Deploy
This app was deployed to GCP and is listening for the invoice webhook on [this url](https://stark-bank-test-3fzbrnzjyq-rj.a.run.app/invoicehook).

## Continuous deployment
On push to master this repository triggers a build for the app's docker image which is pushed to GCP and served immediately.
