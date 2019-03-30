FROM golang:alpine AS build
ADD . /go/src/github.com/JosephSalisbury/twitter-cleanup
RUN cd /go/src/github.com/JosephSalisbury/twitter-cleanup && go build -o twitter-cleanup

FROM alpine
COPY --from=build /go/src/github.com/JosephSalisbury/twitter-cleanup/twitter-cleanup /twitter-cleanup
ENTRYPOINT ["/twitter-cleanup"]
