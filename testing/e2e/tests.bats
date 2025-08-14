setup() {
  bats_load_library bats-support
  bats_load_library bats-assert

  bats_require_minimum_version 1.5.0
}

wait_for_exporter() {
  kubectl -n default wait --timeout 20m --for=condition=Online --for=condition=Registered \
    exporters.jumpstarter.dev/test-exporter-oidc
  kubectl -n default wait --timeout 20m --for=condition=Online --for=condition=Registered \
    exporters.jumpstarter.dev/test-exporter-sa
  kubectl -n default wait --timeout 20m --for=condition=Online --for=condition=Registered \
    exporters.jumpstarter.dev/test-exporter-legacy
}

@test "can create clients with admin cli" {
  jmp admin create client   test-client-oidc     --unsafe --out /dev/null \
    --oidc-username dex:test-client-oidc
  jmp admin create client   test-client-sa       --unsafe --out /dev/null \
    --oidc-username dex:system:serviceaccount:default:test-client-sa
  jmp admin create client   test-client-legacy   --unsafe --save
}

@test "can create exporters with admin cli" {
  jmp admin create exporter test-exporter-oidc   --out /dev/null \
    --oidc-username dex:test-exporter-oidc \
    --label example.com/board=oidc
  jmp admin create exporter test-exporter-sa     --out /dev/null \
    --oidc-username dex:system:serviceaccount:default:test-exporter-sa \
    --label example.com/board=sa
  jmp admin create exporter test-exporter-legacy --save \
    --label example.com/board=legacy
}

@test "can login with oidc" {
  jmp config client   list
  jmp config exporter list

  jmp login --client test-client-oidc \
    --endpoint "$ENDPOINT" --namespace default --name test-client-oidc \
    --issuer https://dex.dex.svc.cluster.local:5556 \
    --username test-client-oidc@example.com --password password --unsafe

  jmp login --client test-client-oidc-provisioning \
    --endpoint "$ENDPOINT" --namespace default --name "" \
    --issuer https://dex.dex.svc.cluster.local:5556 \
    --username test-client-oidc-provisioning@example.com --password password --unsafe

  jmp login --client test-client-sa \
    --endpoint "$ENDPOINT" --namespace default --name test-client-sa \
    --issuer https://dex.dex.svc.cluster.local:5556 \
    --connector-id kubernetes \
    --token $(kubectl create -n default token test-client-sa) --unsafe

  jmp login --exporter test-exporter-oidc \
    --endpoint "$ENDPOINT" --namespace default --name test-exporter-oidc \
    --issuer https://dex.dex.svc.cluster.local:5556 \
    --username test-exporter-oidc@example.com --password password

  jmp login --exporter test-exporter-sa \
    --endpoint "$ENDPOINT" --namespace default --name test-exporter-sa \
    --issuer https://dex.dex.svc.cluster.local:5556 \
    --connector-id kubernetes \
    --token $(kubectl create -n default token test-exporter-sa)

  go run github.com/mikefarah/yq/v4@latest -i ". * load(\"$GITHUB_ACTION_PATH/exporter.yaml\")" \
    /etc/jumpstarter/exporters/test-exporter-oidc.yaml
  go run github.com/mikefarah/yq/v4@latest -i ". * load(\"$GITHUB_ACTION_PATH/exporter.yaml\")" \
    /etc/jumpstarter/exporters/test-exporter-sa.yaml
  go run github.com/mikefarah/yq/v4@latest -i ". * load(\"$GITHUB_ACTION_PATH/exporter.yaml\")" \
    /etc/jumpstarter/exporters/test-exporter-legacy.yaml
 
  jmp config client   list
  jmp config exporter list
}

@test "can run exporters" {
  cat <<EOF | bash 3>&- &
while true; do
  jmp run --exporter test-exporter-oidc
done
EOF

  cat <<EOF | bash 3>&- &
while true; do
  jmp run --exporter test-exporter-sa
done
EOF

  cat <<EOF | bash 3>&- &
while true; do
  jmp run --exporter test-exporter-legacy
done
EOF


  wait_for_exporter
}

@test "can specify client config only using environment variables" {
  wait_for_exporter

  JMP_NAMESPACE=default \
  JMP_NAME=test-exporter-legacy \
  JMP_ENDPOINT=$(kubectl get clients.jumpstarter.dev -n default test-client-legacy -o 'jsonpath={.status.endpoint}') \
  JMP_TOKEN=$(kubectl get secrets -n default test-client-legacy-client -o 'jsonpath={.data.token}' | base64 -d) \
  jmp shell --selector example.com/board=oidc j power on
}

@test "can operate on leases" {
  wait_for_exporter

  jmp config client use test-client-oidc

  jmp create lease     --selector example.com/board=oidc --duration 1d
  jmp get    leases
  jmp get    exporters
  jmp delete leases    --all
}

@test "can lease and connect to exporters" {
  wait_for_exporter

  jmp shell --client test-client-oidc   --selector example.com/board=oidc   j power on
  jmp shell --client test-client-sa     --selector example.com/board=sa     j power on
  jmp shell --client test-client-legacy --selector example.com/board=legacy j power on

  wait_for_exporter
  jmp shell --client test-client-oidc-provisioning --selector example.com/board=oidc j power on
}

@test "can get crds with admin cli" {
  jmp admin get client
  jmp admin get exporter
  jmp admin get lease
}

@test "can delete clients with admin cli" {
  kubectl -n default get secret test-client-oidc-client
  kubectl -n default get clients.jumpstarter.dev/test-client-oidc
  kubectl -n default get clients.jumpstarter.dev/test-client-sa
  kubectl -n default get clients.jumpstarter.dev/test-client-legacy

  jmp admin delete client   test-client-oidc   --delete
  jmp admin delete client   test-client-sa     --delete
  jmp admin delete client   test-client-legacy --delete

  run ! kubectl -n default get secret test-client-oidc-client
  run ! kubectl -n default get clients.jumpstarter.dev/test-client-oidc
  run ! kubectl -n default get clients.jumpstarter.dev/test-client-sa
  run ! kubectl -n default get clients.jumpstarter.dev/test-client-legacy
}

@test "can delete exporters with admin cli" {
  kubectl -n default get secret test-exporter-oidc-exporter
  kubectl -n default get exporters.jumpstarter.dev/test-exporter-oidc
  kubectl -n default get exporters.jumpstarter.dev/test-exporter-sa
  kubectl -n default get exporters.jumpstarter.dev/test-exporter-legacy

  jmp admin delete exporter test-exporter-oidc   --delete
  jmp admin delete exporter test-exporter-sa     --delete
  jmp admin delete exporter test-exporter-legacy --delete

  run ! kubectl -n default get secret test-exporter-oidc-exporter
  run ! kubectl -n default get exporters.jumpstarter.dev/test-exporter-oidc
  run ! kubectl -n default get exporters.jumpstarter.dev/test-exporter-sa
  run ! kubectl -n default get exporters.jumpstarter.dev/test-exporter-legacy
}
