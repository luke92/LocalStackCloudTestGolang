# LocalStackCloudTestGolang
Tests with LocalStackCloud and Golang with Amazon S3

# LocalStack Cloud
- https://localstack.cloud/ 
- A fully functional local cloud stack
- Develop and test your cloud and serverless apps offline!

## Documentation
- https://docs.localstack.cloud/get-started/

## Requisites
- python (Python 3.6 up to 3.10 is supported)
- pip (Python package manager)
- Docker

## Install LocalStack
- `python3 -m pip install localstack`

## Install GUI - LocalStack Cockpit (Beta)
- You need request the link for download in https://localstack.cloud/products/cockpit/
- You will get a link similar to this https://api.localstack.cloud/download/cockpit/v0.1.1/LocalStack%20Cockpit%200.1.1.exe?token=9e3ac160dbc35bf57dfd3bd979458820

## Other GUI
- https://app.localstack.cloud/ (You need create an account)

## Starting LocalStack with the LocalStack CLI
- `localstack start`
- If not work localstack (maybe by PATH) run `python3 -m localstack.cli.main start`

## Starting LocalStack with Docker
- `docker run --rm -it -p 4566:4566 -p 4510-4559:4510-4559 localstack/localstack`

## Starting LocalStack with Docker-Compose
- `docker-compose up` (you need a `docker-compose.yml` file)

## Check if running
- You will see in the command line that LocalStack works in PORT `4566`
- Enter in http://localhost:4566/ and you will see `{"status": "running"}`


