
# Project Add To Cart

A brief description of what this project does and who it's for


## Feature

 - Get Add List Product
 - Create Product
 - Get User Cart
 - Checkout


## Setup And Installation

Clone the project

```bash
  git clone https://github.com/nuriansyahmalik/atc.git
```

Go to the project directory

```bash
  cd atc
```

Setup Up Database Migration On
```bash
  cd atc/migrations/domain/init.sql
```

Start The Server

```bash
  make dev 
```


```bash
  make run 
```


## Documentation

Swagger 

```bash
  http://localhost:8080/swagger/index.html
```

### Endpoint
```bash
  http://localhost:8080/v1/product/
```
```bash
  http://localhost:8080/v1/product?limit=10&page=1&category=laptop
```
```bash
  http://localhost:8080/v1/cart/add
```
```bash
  http://localhost:8080/v1/cart/checkout
```
```bash
  http://localhost:8080/v1/cart/72af6db9-4cbc-4214-839c-a05a0de951f1
```

