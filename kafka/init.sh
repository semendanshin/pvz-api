echo "Starting init.sh"

echo "Creating topics"
kafka-topics --create --replication-factor 1 --partitions 1 --topic pvz.events-log --if-not-exists --bootstrap-server kafka0:29092