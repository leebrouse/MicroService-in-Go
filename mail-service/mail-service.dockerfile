# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
COPY templates /templates

CMD [ "/app/mailApp" ]
