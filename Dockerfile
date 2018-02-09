FROM golang:1.8

MAINTAINER Reto Lehmann <reto.lehmann@sbb.ch>

WORKDIR /usr/ssp/

# Download the sources and UI from github
RUN apt-get update && apt-get install -y wget curl \
  && curl -s https://api.github.com/repos/oscp/cloud-selfservice-portal-backend/releases/latest -k \
     | grep "browser_download_url" | cut -d : -f 2,3 | tr -d \" | wget -qi - \
  && tar xfvz self-service-portal-backend.tar.gz \
  && mv dist/* . \
  && apt-get purge -y --auto-remove wget curl

EXPOSE 8080

CMD ["/usr/ssp/server"]