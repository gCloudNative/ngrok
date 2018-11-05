cd ngrok
export GOBIN=""
rm -rf server* base* bin pkg
rm -f assets/client/tls/ngrokroot.crt  assets/server/tls/*.crt assets/server/tls/*.key

export NGROK_DOMAIN="ngrok.XXX.cn";
openssl genrsa -out base.key 2048;
openssl req -new -x509 -nodes -key base.key -days 100000 -subj "/CN=$NGROK_DOMAIN" -out base.pem;
openssl genrsa -out server.key 2048;
openssl req -new -key server.key -subj "/CN=$NGROK_DOMAIN" -out server.csr;
openssl x509 -req -in server.csr -CA base.pem -CAkey base.key -CAcreateserial -days 100000 -out server.crt;
cp base.pem assets/client/tls/ngrokroot.crt;
#cp server.crt assets/server/tls/snakeoil.crt;
#cp server.key assets/server/tls/snakeoil.key;

#make release-server release-client
GOOS=linux GOARCH=amd64 make release-server
GOOS=darwin GOARCH=amd64 make release-client
GOOS=windows GOARCH=amd64 make release-client


[SERVER]
./bin/ngrokd -tlsKey=server.key -tlsCrt=server.crt -domain="ngrok.XXX.cn" -tunnelAddr=":4443" -httpAddr=":8081" -httpsAddr=":8082" -apiAddr=":8083"


[CLIENT]
echo "xx.xx.xx.xx ngrok.XXX.cn" >>/etc/hosts

cat <<EOF > ngrok.cfg
server_addr: ngrok.XXX.cn:4443
trust_host_root_certs: true
EOF


./ngrok -proto=tcp -config=ngrok.cfg -log=stdout     unix:///var/run/docker.sock
./ngrok -proto=tcp -config=ngrok.cfg -subdomain=vnc  5900

