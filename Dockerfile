FROM scratch
COPY tcalc /usr/bin/tcalc
ENV HOME=/home/user
ENTRYPOINT ["/usr/bin/tcalc"]
