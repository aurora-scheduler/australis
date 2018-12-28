FROM ubuntu:16.04

RUN apt-get update -y && \
 apt-get install -y build-essential devscripts dh-exec dh-make git lintian wget && \
 wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz -O /tmp/go.tar.gz

RUN tar -C /usr/local -xzf /tmp/go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

COPY ./entrypoint.sh /

ENTRYPOINT [ "/entrypoint.sh" ]
