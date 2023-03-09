#!/bin/sh

rm *.pem

openssl req -x509 -newkey rsa:4096 \
  -days 3650 \
  -nodes -keyout ca-key.pem \
  -out ca-cert.pem \
  -subj "/C=IT/ST=Padova/L=Italy/O=fuu/OU=billgatesluamaro/CN=*.nelretrobottega.it"

openssl x509 -in ca-cert.pem -noout -text

openssl req -newkey rsa:4096 \
  -nodes -keyout server-key.pem \
  -out server-req.pem \
  -subj "/C=IT/ST=Padova/L=Italy/O=fuu/OU=billgatesluamaro/CN=*.nelretrobottega.it"

openssl x509 -req \
  -in server-req.pem \
  -days 3650 \
  -CA ca-cert.pem \
  -CAkey ca-key.pem \
  -CAcreateserial \
  -out server-cert.pem \

openssl x509 -in server-cert.pem -noout -text