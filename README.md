# fr-challenge-v2
This project is a freight quote simulator that uses an API from Frete Rápido to perform quote calculations and return the available carriers for delivery along with their pricing.

The information related to the quotes is saved in a database and used to calculate metrics (such as the highest price), which are then displayed to the user.
## Tree
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
docker-compose up --build -d
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
        {
            "name": "BOX DELIVERY",
            "service": "Rodoviário",
            "deadline": 0,
            "price": 0
        },
        {
            "name": "BOX DELIVERY",
            "service": "Rodoviário",
            "deadline": 0,
            "price": 0
        },
        {
            "name": "AZUL CARGO",
            "service": "Aéreo",
            "deadline": 2,
            "price": 41.82
        },
        {
            "name": "AZUL CARGO",
            "service": "Aéreo",
            "deadline": 0,
            "price": 41.82
        },
        {
            "name": "PRESSA FR (TESTE)",
            "service": "Rodoviário",
            "deadline": 0,
            "price": 58.95
        },
        {
            "name": "FR EXPRESS (TESTE)",
            "service": "Rodoviário",
            "deadline": 3,
            "price": 74.95
        },
        {
            "name": "BTU BRASPRESS",
            "service": "Rodoviário",
            "deadline": 5,
            "price": 93.35
        },
        {
            "name": "CORREIOS",
            "service": "Rodoviário",
            "deadline": 5,
            "price": 103.71
        },
        {
            "name": "CORREIOS - SEDEX",
            "service": "Rodoviário",
            "deadline": 6,
            "price": 121.03
        },
        {
            "name": "CORREIOS",
            "service": "Rodoviário",
            "deadline": 6,
            "price": 121.03
        },
        {
            "name": "BRASPRESS",
            "service": "Rodoviário",
            "deadline": 4,
            "price": 133.58
        },
        {
            "name": "CORREIOS",
            "service": "Rodoviário",
            "deadline": 1,
            "price": 168.43
        },
        {
            "name": "CORREIOS - SEDEX",
            "service": "Rodoviário",
            "deadline": 2,
            "price": 185.75
        },
        {
            "name": "CORREIOS",
            "service": "Rodoviário",
            "deadline": 2,
            "price": 185.75
        }
    ]
}
```

### [GET] .../metrics?last_quotes={n}

Calculates metrics using information from stored quotes in the database (where `n` specifies the number of quotes in descending order) and then displays the results for the user.

#### Parameters
* `last_quotes` (optional) indicates the amount of quotes used to calculate metrics

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
                "AZUL CARGO": 12,
                "BOX DELIVERY": 12,
                "BRASPRESS": 6,
                "BTU BRASPRESS": 6,
                "CORREIOS": 24,
                "CORREIOS - SEDEX": 12,
                "FR EXPRESS (TESTE)": 6,
                "PRESSA FR (TESTE)": 6
            },
            "total_price_per_carrier": {
                "AZUL CARGO": 501.84,
                "BOX DELIVERY": 0,
                "BRASPRESS": 801.48,
                "BTU BRASPRESS": 560.09,
                "CORREIOS": 3473.5,
                "CORREIOS - SEDEX": 1840.65,
                "FR EXPRESS (TESTE)": 449.7,
                "PRESSA FR (TESTE)": 353.7
            },
            "avg_price_per_carrier": {
                "AZUL CARGO": 41.82,
                "BOX DELIVERY": 0,
                "BRASPRESS": 133.58,
                "BTU BRASPRESS": 93.34,
                "CORREIOS": 144.72,
                "CORREIOS - SEDEX": 153.38,
                "FR EXPRESS (TESTE)": 74.95,
                "PRESSA FR (TESTE)": 58.95
            },
            "cheapest_freight": {
                "AZUL CARGO": 41.82,
                "BOX DELIVERY": 0,
                "BRASPRESS": 133.58,
                "BTU BRASPRESS": 93.35,
                "CORREIOS": 103.71,
                "CORREIOS - SEDEX": 121.03,
                "FR EXPRESS (TESTE)": 74.95,
                "PRESSA FR (TESTE)": 58.95
            },
            "priciest_freight": {
                "AZUL CARGO": 41.82,
                "BOX DELIVERY": 0,
                "BRASPRESS": 133.58,
                "BTU BRASPRESS": 93.35,
                "CORREIOS": 185.75,
                "CORREIOS - SEDEX": 185.75,
                "FR EXPRESS (TESTE)": 74.95,
                "PRESSA FR (TESTE)": 58.95
            }
        }
    ]
}
```

## License
This project is licensed under the MIT License.