docker build -t builder-image --no-cache .
docker run --name builder builder-image

docker cp builder:/build .

docker rm builder
docker image rm builder-image
