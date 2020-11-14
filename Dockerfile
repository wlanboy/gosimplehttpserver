FROM busybox:1.31

RUN mkdir /home
COPY ./gosimplehttpserver /home/
COPY ./.env /home/
EXPOSE 7000

CMD ["/home/gosimplehttpserver"]
