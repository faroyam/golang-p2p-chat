mkdir certs 2>/dev/null
rm certs/* 2>/dev/null
echo "creating server cert"
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 3650 -subj "/C=RU/ST=MSK/L=MSK/O=Chat Company/OU=IT/CN=www.chat.com/emailAddress=example@mail.com"
echo "creating client cert"
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 3650 -subj "/C=RU/ST=MSK/L=MSK/O=Chat Company/OU=IT/CN=www.chat.com/emailAddress=example@mail.com"
echo "done"