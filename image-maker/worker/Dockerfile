FROM alpine:latest
# check if need to install curl
RUN apk update && apk add git && apk add podman
WORKDIR /root/
COPY ./job.sh ./job.sh
COPY ./gradio ./gradio
CMD ./job.sh
