#bin/bash

cat test/deneme.orhun | go run ./... | diff -q test/deneme.orhun.test -
cat test/deneme2.orhun | go run ./... | diff -q test/deneme2.orhun.test -
cat test/deneme3.orhun | go run ./... | diff -q test/deneme3.orhun.test -