FROM alpine
RUN apk add --no-cache libc6-compat
ADD api-gateway-service /api-gateway-service
ENTRYPOINT [ "/api-gateway-service" ]
