docker build -t auth-service:$VERSION --build-arg DOPPLER_TOKEN=$DOPPLER_TOKEN . && \
docker tag auth-service:$VERSION auth-service:$VERSION && \
docker push <repository_uri>/auth-service:$VERSION