# gsd

## Running with Docker

You can run the application using Docker. The image is available from GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/2easy/gsd:latest

# Run the container
# Replace /path/to/data with your desired path to store the database
# Replace 8081 with your desired port
docker run -d \
  -p 8081:8081 \
  -v /path/to/data:/data \
  ghcr.io/2easy/gsd:latest \
  --port 8081 \
  --db /data/gsd.db
```

The application will be available at http://localhost:8081 (or your specified port).

### Configuration

The application accepts the following command line flags:

- `--port`: HTTP server port (default: 8081)
- `--db`: Path to the SQLite database file (default: ./gsd.db)

### Building locally

To build the Docker image locally:

```bash
docker build -t gsd .
```