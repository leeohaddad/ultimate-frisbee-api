FROM golang:alpine as build

WORKDIR /build

RUN apk add --update make

COPY . .

RUN make setup build

FROM alpine
WORKDIR /disc

# RUN adduser -D ultimate_frisbee_manager && chown -R ultimate_frisbee_manager /disc
# USER ultimate_frisbee_manager

COPY --from=build /build/config/ config/
COPY --from=build /build/infra/database/migrations/ infra/database/migrations/
COPY --from=build /build/bin/ultimate_frisbee_manager .

EXPOSE 80

ENTRYPOINT [ "/disc/ultimate_frisbee_manager" ]
