FROM yaccio/orchid

RUN mkdir /download
WORKDIR /download

RUN wget https://github.com/joewalnes/websocketd/releases/download/v0.2.10/websocketd-0.2.10-linux_amd64.zip
RUN unzip websocketd-0.2.10-linux_amd64.zip
RUN chmod u+x websocketd
RUN cp websocketd /bin/websocketd

ADD orchid-server /bin/orchid-server

WORKDIR /project

EXPOSE 3000

ENTRYPOINT ["websocketd", "--port=3000", "--devconsole", "orchid-server"]
