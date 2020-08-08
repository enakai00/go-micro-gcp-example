FROM alpine
RUN apk add --no-cache libc6-compat
ADD purchase-service /purchase-service
ENTRYPOINT [ "/purchase-service" ]
