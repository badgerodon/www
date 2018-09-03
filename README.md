# Badgerodon WWW Site

## Development

Build Local:

    env GO111MODULE=on go build -v -o bin/badgerodon-www
    docker build -t gcr.io/doxsey-1/badgerodon-www .

Run Local:

    docker run -p 8080:8080 gcr.io/doxsey-1/badgerodon-www

Build Remote:

    gcloud builds submit --config cloudbuild.yaml .
