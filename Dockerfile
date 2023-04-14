FROM golang:1.20.3-alpine

# Create an move to the working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum form3/go.mod ./

# Download all the dependencies
RUN go mod download

# Copy the source from the current directory to the working directory
COPY . .

CMD ["go", "test", "-v", "./..."]