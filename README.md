# fr-challenge-v2

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)

## Description
This project is a freight quote simulator that uses an API from Frete Rápido to perform quote calculations and return the available carriers for delivery along with their charges.

The information related to the quotes is saved in a database and used to calculate metrics (such as the highest price), which are then displayed to the user.

### Tree
```text
.
├── fr
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── handlers_test.go
│   │       ├── helpers.go
│   │       ├── helpers_test.go
│   │       ├── main.go
│   │       ├── models.go
│   │       ├── routes.go
│   │       └── routes_test.go
│   ├── data
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── test-models.go
│   ├── fr.dockerfile
│   ├── go.mod
│   └── go.sum
├── LICENSE
├── project
│   ├── docker-compose.yml
│   └── Makefile
└── README.md
```

## Installation & Usage
### Requisites
* Docker
* Docker Compose
* Make

### Usage
* Clone this repository
```bash
git clone https://github.com/mtrdgs/fr-challenge-v2
```

* Navigate to the project directory
```bash
cd fr-challenge-v2/project
```

* Start containers
```bash
make up_build
```
or
```bash
docker compose up --build -d
```

Obs.: This application runs at `http://localhost:8080/`.

## Endpoints
### [POST] .../quote

Receives data from the user and performs a quote using the Frete Rápido API.

#### Request
```bash
curl --location 'http://localhost:8080/quote' \
--header 'Content-Type: application/json' \
--data '{
    "recipient": {
        "address": {
            "zipcode": "01311000"
        }
    },
    "volumes": [
        {
            "category": 7,
            "amount": 1,
            "unitary_weight": 5,
            "price": 349,
            "sku": "abc-teste-123",
            "height": 0.2,
            "width": 0.2,
            "length": 0.2
        },
        {
            "category": 7,
            "amount": 2,
            "unitary_weight": 4,
            "price": 556,
            "sku": "abc-teste-527",
            "height": 0.4,
            "width": 0.6,
            "length": 0.15
        }
    ]
}'
```

#### Response
```json
{
    "carrier": [
        // ...
        {
            "name": "BTU BRASPRESS",
            "service": "Rodoviário",
            "deadline": 4,
            "price": 78.63
        },
        {
            "name": "CORREIOS",
            "service": "Rodoviário",
            "deadline": 3,
            "price": 332.53
        },

        {
            "name": "CORREIOS - SEDEX",
            "service": "Rodoviário",
            "deadline": 3,
            "price": 332.53
        }
        /// ...
    ]
}
```

### [GET] .../metrics?last_quotes={n}

Calculates metrics using information from stored quotes in the database (where `n` specifies the number of quotes in descending order) and then displays the results for the user.

#### Parameters
* `last_quotes`: Indicates the number of quotes (n) used to calculate metrics

#### Request
```bash
curl --location 'http://localhost:8080/metrics?last_quotes=6'
```

#### Response
```json
{
    "metrics": [
        {
            "results_per_carrier": {
                "AZUL CARGO": 4,
                "BOX DELIVERY": 6,
                "BRASPRESS": 2,
                "BTU BRASPRESS": 2,
                "CORREIOS": 11,
                "CORREIOS - SEDEX": 9,
                "FR EXPRESS (TESTE)": 2,
                "PRESSA FR (TESTE)": 2,
                "RAPIDÃO FR (TESTE)": 2
            },
            "total_price_per_carrier": {
                "AZUL CARGO": 94.64,
                "BOX DELIVERY": 0,
                "BRASPRESS": 120.7,
                "BTU BRASPRESS": 157.26,
                "CORREIOS": 1688.44,
                "CORREIOS - SEDEX": 1553.92,
                "FR EXPRESS (TESTE)": 200.96,
                "PRESSA FR (TESTE)": 117.9,
                "RAPIDÃO FR (TESTE)": 353.16
            },
            "avg_price_per_carrier": {
                "AZUL CARGO": 23.66,
                "BOX DELIVERY": 0,
                "BRASPRESS": 60.35,
                "BTU BRASPRESS": 78.63,
                "CORREIOS": 153.49,
                "CORREIOS - SEDEX": 172.65,
                "FR EXPRESS (TESTE)": 100.48,
                "PRESSA FR (TESTE)": 58.95,
                "RAPIDÃO FR (TESTE)": 176.58
            },
            "cheapest_freight": {
                "AZUL CARGO": 23.66,
                "BOX DELIVERY": 0,
                "BRASPRESS": 60.35,
                "BTU BRASPRESS": 78.63,
                "CORREIOS": 67.27,
                "CORREIOS - SEDEX": 81.05,
                "FR EXPRESS (TESTE)": 100.48,
                "PRESSA FR (TESTE)": 58.95,
                "RAPIDÃO FR (TESTE)": 176.58
            },
            "priciest_freight": {
                "AZUL CARGO": 23.66,
                "BOX DELIVERY": 0,
                "BRASPRESS": 60.35,
                "BTU BRASPRESS": 78.63,
                "CORREIOS": 332.53,
                "CORREIOS - SEDEX": 332.53,
                "FR EXPRESS (TESTE)": 100.48,
                "PRESSA FR (TESTE)": 58.95,
                "RAPIDÃO FR (TESTE)": 176.58
            }
        }
    ]
}
```