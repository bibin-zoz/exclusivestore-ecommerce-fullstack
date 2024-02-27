# Build Stage
FROM golang:1.21.1-alpine3.18 AS build-stage

WORKDIR /app
COPY ./ /app

RUN mkdir -p /app/build
RUN go mod download
RUN go build -v -o /app/build/api .

# Final Stage
FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /app/build/api /api
COPY --from=build-stage /app/templates /templates
COPY --from=build-stage /app/static /static
COPY --from=build-stage /app/.env /

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["/api"]
