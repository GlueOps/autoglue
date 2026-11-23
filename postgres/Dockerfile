FROM postgres:18.7@sha256:f5cbc081c25690710f3fc6f4c61a7565382bae01ca0d0f2028a53d7926b85a20

RUN cd /var/lib/postgresql/ && \
    openssl req -new -text -passout pass:abcd -subj /CN=localhost -out server.req -keyout privkey.pem && \
    openssl rsa -in privkey.pem -passin pass:abcd -out server.key && \
    openssl req -x509 -in server.req -text -key server.key -out server.crt && \
    chmod 600 server.key && \
    chown postgres:postgres server.key
CMD ["postgres", "-c", "ssl=on", "-c", "ssl_cert_file=/var/lib/postgresql/server.crt", "-c", "ssl_key_file=/var/lib/postgresql/server.key" ]