


docker build -t gdunghi/go_hooks .
docker push gdunghi/go_hooks
docker pull gdunghi/go_hooks
docker run --rm -d --name go_hooks -p 6969:6969 gdunghi/go_hooks

sudo docker run --rm -d --name go_hooks -p 6969:6969 gdunghi/go_hooks