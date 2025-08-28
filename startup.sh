#!/bin/sh

# Generate config.toml from environment variables
cat > config.toml << EOF
log_level="${LISTMONK__log_level:-info}"

[server]
address = "${LISTMONK_server__address:-:8082}"
read_timeout = "${LISTMONK_server__read_timeout:-5s}"
write_timeout = "${LISTMONK_server__write_timeout:-5s}"

[messenger.pinpoint]
config = "${LISTMONK_messenger_pinpoint__config}"

# Add other configuration sections as needed for listmonk-messenger
# This should match the structure expected by your listmonk-messenger application

EOF

echo "Generated config.toml:"
cat config.toml

# Start listmonk-messenger
exec ./listmonk-messenger --config config.toml --msgr pinpoint
