FROM postgres:17.8@sha256:69dddb030ab69d669d8d7c6abf67aeb448178e5270d5f123a21f4f7ac8b46a24

RUN cd /var/lib/postgresql/ && \
    openssl req -new -text -passout pass:abcd -subj /CN=localhost -out server.req -keyout privkey.pem && \
    openssl rsa -in privkey.pem -passin pass:abcd -out server.key && \
    openssl req -x509 -in server.req -text -key server.key -out server.crt && \
    chmod 600 server.key && \
    chown postgres:postgres server.key
CMD ["postgres", "-c", "ssl=on", "-c", "ssl_cert_file=/var/lib/postgresql/server.crt", "-c", "ssl_key_file=/var/lib/postgresql/server.key" ]
