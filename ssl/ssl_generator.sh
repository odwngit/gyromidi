openssl req -new -subj "/CN=starwheel" -newkey rsa:2048 -nodes -keyout starwheel.key -out starwheel.csr
openssl x509 -req -days 365 -in starwheel.csr -signkey starwheel.key -out starwheel.crt
