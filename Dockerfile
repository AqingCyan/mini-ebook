FROM ubuntu:20.04
COPY mini-book /app/mini-book
WORKDIR /app
CMD ["/app/mini-book"]