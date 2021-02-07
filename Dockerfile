FROM golang:alpine as build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch as run
WORKDIR /app
COPY --from=build /build/main .
COPY ./www ./www
EXPOSE 9844
ENV PORT 9844
ENV PROJECTOR 192.168.12.241:4998
ENV PROJECTOR_IR_PORT 2
ENV RECEIVER 192.168.12.40:23
ENV ROKU 192.168.12.238:8060
ENV LIGHTS 192.168.12.3:8899
CMD ["./main"]