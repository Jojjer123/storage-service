#!/bin/bash

hostname=${HOSTNAME:-localhost}

sed -i -e "s/replace-device-name/"$hostname"/g" $HOME/target_configs/typical_ofsw_config.json && \
sed -i -e "s/replace-motd-banner/Welcome to gNMI service on "$hostname":"$GNMI_PORT"/g" $HOME/target_configs/typical_ofsw_config.json

gnmi_target \
    -bind_address :$GNMI_INSECURE_PORT \
    -alsologtostderr \
    -notls \
    -insecure \
    -config $HOME/target_configs/typical_ofsw_config.json &

gnmi_target \
    -bind_address :$GNMI_PORT \
    -key $HOME/certs/localhost.key \
    -cert $HOME/certs/localhost.crt \
    -ca $HOME/certs/onfca.crt \
    -alsologtostderr \
    -config $HOME/target_configs/typical_ofsw_config.json
