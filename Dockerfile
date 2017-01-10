FROM centurylink/ca-certs

COPY bin/statspout /

WORKDIR /

EXPOSE 8080

ENTRYPOINT ["/statspout"]

