# Badgerodon WWW Site

## Development

Build:

    docker build -t gcr.io/badgerodon-173120/badgerodon-www:VERSION .

Run Local:

    docker run -p 8080:8080 gcr.io/badgerodon-173120/badgerodon-www:VERSION

GCR Push:

    gcloud docker -- push gcr.io/badgerodon-173120/badgerodon-www:VERSION


Update key for travis:

    travis encrypt-file /home/caleb/Dropbox/keys/gcp-doxsey-1-travis-ci.json --add