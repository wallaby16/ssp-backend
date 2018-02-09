FROM golang:1.8

WORKDIR /usr/ssp/

# Download the sources and UI from github
ADD https://github.com/oscp/cloud-selfservice-portal-backend/releases/download/v1.2.3/self-service-portal-backend.tar.gz self-service-portal-backend.tar.gz

# Extract the content
RUN tar xfvz self-service-portal-backend.tar.gz &&mv dist/* .

EXPOSE 8080

CMD ["/usr/ssp/server"]