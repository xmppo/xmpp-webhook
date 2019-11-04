#!/usr/bin/env bash

# applies https://modules.prosody.im/mod_auth_any.html and https://modules.prosody.im/mod_roster_allinall.html for testing

PROSODY_CONF=/etc/prosody/prosody.cfg.lua
AUTH_ANY='authentication = "any"'
MODULES='modules_enabled = { "auth_any", "roster_allinall" }'

function changes_missing {
	if grep -q "$AUTH_ANY" $PROSODY_CONF; then
		return 1
	else
		return 0
	fi
}

function apply_changes {
	echo "$AUTH_ANY" >> $PROSODY_CONF
	echo "$MODULES" >> $PROSODY_CONF
}

if changes_missing; then
	apply_changes
fi

/usr/bin/entrypoint.sh "$@"
