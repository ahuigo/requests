# Generating CA

    openssl genrsa -out rootCA.key 4096
    openssl req -x509 -new -key rootCA.key -days 3650 -out rootCA.crt
# Generate certificate for local.com signed with created CA

    openssl genrsa -out local.com.key 2048
    openssl req -new -key local.com.key -out local.com.csr
    #In answer to question `Common Name (e.g. server FQDN or YOUR name) []:` you should set `local.com` (your real domain name)
    openssl x509 -req -in local.com.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -days 365 -out local.com.crt
