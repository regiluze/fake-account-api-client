FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Ruben Eguiluz <regiluze@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN make deps
