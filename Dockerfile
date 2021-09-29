FROM scratch
#ca-certificates.crt get from /etc/ssl/certs
COPY ca-certificates.crt /etc/ssl/certs/
ADD main /
CMD ["/main"]
EXPOSE 8080
