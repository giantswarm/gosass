FROM jbergknoff/sass

RUN cp /usr/bin/sass /usr/bin/sassc
COPY gosass /usr/bin/

ENTRYPOINT ["gosass"]

# sudo docker build -t local/gosass .

# cd $GITHUB/giantswarm/giantswarmio-webapp
# sudo docker run -ti --rm -v $PWD:$PWD -w $PWD local/gosass -input webapp/scss/ -output webapp/static/css/ -style compressed
