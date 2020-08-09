FROM alpine
RUN apk add --no-cache libc6-compat
ADD stock-service /stock-service
ENTRYPOINT [ "/stock-service" ]
