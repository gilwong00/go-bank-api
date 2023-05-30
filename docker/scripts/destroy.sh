
# docker kill bank_api_pg
# docker kill bank_api_redis
# docker container rm bank_api_pg
# docker container rm bank_api_redis
containers=$(docker ps -q)
docker kill $containers
docker container rm $containers
