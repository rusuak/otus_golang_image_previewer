FROM golang:1.18-buster as base

ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./vendor /app/vendor
COPY ./go.mod /app/go.mod
COPY ./.env /app/.env
COPY ./Makefile /app/Makefile

RUN make build

FROM alpine:3.15
WORKDIR /ipreviewer
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /app/.env /ipreviewer/.env
COPY --from=base /app/server /ipreviewer/app

RUN chown -R appuser:appuser /ipreviewer
USER appuser:appuser

CMD ["./app"]
