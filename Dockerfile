FROM alpine

COPY jira-exporter /usr/bin/
ENTRYPOINT ["/usr/bin/jira-exporter"]