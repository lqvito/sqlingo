<img src="https://raw.githubusercontent.com/lqs/sqlingo/master/logo.png" width="236" height="106">

[![Travis CI](https://travis-ci.org/lqs/sqlingo.svg?branch=master)](https://travis-ci.org/lqs/sqlingo)
[![Go Report Card](https://goreportcard.com/badge/github.com/lqs/sqlingo)](https://goreportcard.com/report/github.com/lqs/sqlingo)
[![codecov](https://codecov.io/gh/lqs/sqlingo/branch/master/graph/badge.svg)](https://codecov.io/gh/lqs/sqlingo)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

**sqlingo** is a SQL DSL library in Go. It generates code from the database and lets you write SQL easily and correctly.

## Tutorial

### Install and use sqlingo code generator
In order to generate code, sqlingo requires your tables are already created in the database.

```
$ go get -u github.com/lqs/sqlingo/sqlingo-gen
$ mkdir -p generated/sqlingo
$ sqlingo-gen root:123456@/database_name >generated/sqlingo/database_name.dsl.go
```

### Start using sqlingo
Create `main.go` to use the generated code
```go
package main

import (
    "github.com/lqs/sqlingo"
    . "./generated/sqlingo"
)

func main() {
    db, err := sqlingo.Open("mysql", "root:123456@/database_name")
    if err != nil {
        panic(err)
    }
    
    // insert some rows
    customer1 := &CustomerModel{name: "Customer One"}
    customer2 := &CustomerModel{name: "Customer Two"}
    db.InsertInto(Customer).
        Values(customer1, customer2).
        Execute()

    // do some queries
    var customers []*CustomerModel
    db.SelectFrom(Customer).
        Where(Customer.Id.In(1, 2)).
        FetchAll(&customers)

    // more examples
    var customerId int64
    var orderId int64
    db.Select(Customer.Id, Order.Id).
        From(Customer, Order).
        Where(Customer.Id.Equals(Order.CustomerId), Order.Id.Equals(1)).
        FetchFirst(&customerId, &orderId)
}
```
