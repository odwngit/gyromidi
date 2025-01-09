openssl req -new -subj "/CN=gyromidi" -newkey rsa:2048 -nodes -keyout gyromidi.key -out gyromidi.csr
openssl x509 -req -days 365 -in gyromidi.csr -signkey gyromidi.key -out gyromidi.crt
