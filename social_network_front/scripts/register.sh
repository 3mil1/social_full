#!/bin/bash

curl -X POST -H "Content-Type: application/json" -d @silver.json http://localhost:8080/user/signup
curl -X POST -H "Content-Type: application/json" -d @niki.json http://localhost:8080/user/signup
curl -X POST -H "Content-Type: application/json" -d @vici.json http://localhost:8080/user/signup
curl -X POST -H "Content-Type: application/json" -d @maik.json http://localhost:8080/user/signup
