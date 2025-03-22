# gsd

## Running with Docker

You can run the application using Docker. The image is available from GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/2easy/gsd:latest

# Run the container
# Replace /path/to/data with your desired path to store the database
# Replace 8080 with your desired port
docker run -d \
  -p 8080:8080 \
  -v /path/to/data:/data \
  -e GSD_PORT=8080 \
  -e GSD_DB_PATH=/data/gsd.db \
  ghcr.io/2easy/gsd:latest
```

The application will be available at http://localhost:8080 (or your specified port).

### Configuration

The following environment variables can be used to configure the application:

- `GSD_PORT`: HTTP server port (default: 8080)
- `GSD_DB_PATH`: Path to the SQLite database file (default: /data/gsd.db)

### Building locally

To build the Docker image locally:

```bash
docker build -t gsd .
```