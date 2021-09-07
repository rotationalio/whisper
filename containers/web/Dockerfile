FROM node:16 AS builder

WORKDIR /usr/src/app

# Dependencies
COPY web/package.json ./
COPY web/yarn.lock ./
RUN yarn

# Set production environment variable for build context
ARG REACT_APP_API_BASE_URL "https://whisper.rotational.dev/v1"
ARG NODE_ENV "production"

# Build app with browserify
COPY web/ ./
RUN yarn build

FROM nginx:stable

LABEL maintainer="Rotational Labs <info@rotational.io>"
LABEL description="Whisper UI, a secret sharing user interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/web/nginx.conf /etc/nginx/conf.d/configfile.template

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
ENV PORT 8080
ENV HOST 0.0.0.0

EXPOSE 8080

CMD sh -c "envsubst '\$PORT' < /etc/nginx/conf.d/configfile.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"