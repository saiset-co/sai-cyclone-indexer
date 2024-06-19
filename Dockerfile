ARG SERVICE="sai-cyclone-indexer"

FROM golang as BUILD

ARG SERVICE

WORKDIR /src/

COPY ./ /src/

RUN go build -o sai-cyclone-indexer -buildvcs=false

FROM ubuntu

ARG SERVICE

WORKDIR /srv

COPY --from=BUILD /src/sai-cyclone-indexer /srv/sai-cyclone-indexer
COPY ./config.yml /srv/config.yml
COPY ./addresses.json /srv/addresses.json

RUN chmod +x /srv/sai-cyclone-indexer

CMD /srv/sai-cyclone-indexer start
