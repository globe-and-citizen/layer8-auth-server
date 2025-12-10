FROM ubuntu:latest
LABEL authors="thaodtp"

ENTRYPOINT ["top", "-b"]