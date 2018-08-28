FROM scratch

EXPOSE 8080
ADD rabbit_consumer_groups_linux /srv/app

CMD ["/srv/app"]