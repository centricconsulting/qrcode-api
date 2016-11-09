# Usage: deploy.sh [dev|qa|demo|stage|prod|master] [app_name] [port]

# Pull the repo down.
echo "Pulling Centric QR Code repository..."
docker pull centric/qrcode-api:$1
# Stop the currently executing container.
echo "Stopping any running containers..."
docker stop $2_$1
docker rm $2_$1
# Remove all of the orphaned containers.
echo "Remove all orphaned and exited containers..."
docker rm $(docker ps -q -f status=exited)
# Start the application.
echo "Starting the application..."
docker run -p $3:3200 --name $2_$1 centric/qrcode-api:$1
