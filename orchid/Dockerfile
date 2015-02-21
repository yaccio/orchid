FROM debian

RUN apt-get update
RUN apt-get install -y wget curl unzip

ADD orchid /bin/orchid

RUN mkdir /project
WORKDIR /project

ENTRYPOINT ["orchid"]
