FROM golang:latest
 
RUN mkdir /ascii-web-docker
 
# Copy all files from the current directory to the app directory
COPY . /ascii-web-docker
 
# Set working directory
WORKDIR /ascii-web-docker
 

RUN go build -o server . 
 
# Run the server executable
CMD [ "/ascii-web-docker/server" ]

