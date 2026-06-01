# Dockerfile
FROM postgis/postgis:16-3.5

# Extensions additionnelles
RUN apt-get update && apt-get install -y \
    postgresql-16-h3 \
    postgresql-contrib \
    && rm -rf /var/lib/apt/lists/*

