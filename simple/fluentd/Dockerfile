FROM fluent/fluentd:latest
MAINTAINER cqshinn92@gmail.com <Cong Quan>
WORKDIR /home/fluent
ENV PATH /home/fluent/.gem/ruby/2.3.0/bin:$PATH
RUN mkdir -p /etc/fluent /var/log/fluentd
RUN chmod 775 /var/log/fluentd
# cutomize following "gem install fluent-plugin-..." line as you wish

USER root
RUN apk --no-cache --update add sudo build-base ruby-dev && \

    sudo -u fluent gem install fluent-plugin-prometheus fluent-plugin-record-reformer fluent-plugin-elasticsearch && \

    rm -rf /home/fluent/.gem/ruby/2.3.0/cache/*.gem && sudo -u fluent gem sources -c && \
    apk del sudo build-base ruby-dev && rm -rf /var/cache/apk/*

ADD fluentd.conf /etc/fluent

EXPOSE 24224

USER fluent
# USER root
CMD exec fluentd -c /etc/fluent/fluentd.conf
#-p /fluentd/plugins $FLUENTD_OPT