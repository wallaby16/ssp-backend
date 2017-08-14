FROM golang:1.8

WORKDIR /usr/ssp/

# Download the sources and UI from github
ADD https://github.com/oscp/cloud-selfservice-portal/releases/download/v1.0.9/self-service-portal.tar.gz self-service-portal.tar.gz

# Extract the content
RUN tar xfvz self-service-portal.tar.gz

EXPOSE 8080

CMD ["/usr/ssp/dist/server"]