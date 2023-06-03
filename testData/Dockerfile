FROM debian:latest

RUN apt-get update -y
RUN apt-get install -y python3
CMD python3 -m http.server 80 --bind 0.0.0.0
EXPOSE 8888
