FROM golang:1.14

WORKDIR /workspace
VOLUME ["/workspace"]
CMD ["go", "build", "-o", "bin/trial"]
