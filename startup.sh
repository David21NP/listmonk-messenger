#!/bin/sh

esc() {
    printf "%s\\n" "$1" | sed -e "s/'/'\"'\"'/g" -e "1s/^/'/" -e "\$s/\$/'/"
}

quoted_config=$(esc "$LISTMONK_messenger_pinpoint__config")

# Generate config.toml from environment variables
cat > config.toml << EOF
log_level="${LISTMONK__log_level:-info}"

[server]
address = "${LISTMONK_server__address:-:8082}"
read_timeout = "${LISTMONK_server__read_timeout:-5s}"
write_timeout = "${LISTMONK_server__write_timeout:-5s}"

[messenger.end_user_messaging]
config = "${quoted_config:-{\}}"

# Add other configuration sections as needed for listmonk-messenger
# This should match the structure expected by your listmonk-messenger application

EOF

echo "Generated config.toml:"
cat config.toml

# Start listmonk-messenger
exec ./listmonk-messenger.bin --config config.toml --msgr end_user_messaging
