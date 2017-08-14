FROM golang:1.8

WORKDIR /go/src/github.com/oscp/cloud-selfservice-portal/server

COPY server/ /go/src/github.com/oscp/cloud-selfservice-portal

# Build the backend: Golang
RUN go get gopkg.in/gin-gonic/gin.v1 \
    && go get gopkg.in/appleboy/gin-jwt.v2 \
    && go get gopkg.in/dgrijalva/jwt-go.v3 \
    && go get github.com/jtblin/go-ldap-client \
    && go get github.com/Jeffail/gabs

RUN go install -v

# Build the frontend: npm
RUN curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash - \
    && sudo apt-get install -y nodejs \
    && cd ui && npm install && npm run build \
    && mv dist ../ && ls

EXPOSE 8080

CMD ["server"]