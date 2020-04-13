FROM busybox:1.31

COPY ./ip /home/
EXPOSE 8000

CMD ["/home/ip"]
