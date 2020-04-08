FROM scratch
ADD bin/main /
ADD etc/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]
