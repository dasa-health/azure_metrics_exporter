FROM golang:1.10-alpine

ARG subscription_id
ARG client_id
ARG client_secret
ARG tenant_id
ARG metric_aggregation
ARG resource_query_tag_name
ARG environment
ARG release
ARG elk_host
ARG elk_index
ARG active_log_segregation

ENV subscription_id $subscription_id
ENV client_id $client_id
ENV client_secret $client_secret
ENV tenant_id $tenant_id
ENV metric_aggregation $metric_aggregation
ENV resource_query_tag_name $resource_query_tag_name
ENV environment $environment
ENV release $release
ENV elk_host $elk_host
ENV elk_index $elk_index
ENV active_log_segregation $active_log_segregation

WORKDIR /go/src/app
COPY . .

RUN apk add util-linux && apk add --update git

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9276
CMD ["app"]