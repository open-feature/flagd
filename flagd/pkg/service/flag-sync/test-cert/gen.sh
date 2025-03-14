:'
This script can be used to recreate the SSL certificates that are used in the sync service test

Warning: there might be issues running the script on Windows with the -subj argument
 -> workaround: run the commands manually without the -subj argument and provide info when asked by the console output
'
rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 9999 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/CN=flagD test certificate"

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/CN=flagD test server PR and CSR"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 9999 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

