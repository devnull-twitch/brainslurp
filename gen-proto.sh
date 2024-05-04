#!/usr/bin/env bash

protoc -I=. --go_out=. --go_opt=paths=import --go_opt=module=github.com/devnull-twitch/brainslurp \
  lib/proto/issue.proto \
  lib/proto/view.proto \
  lib/proto/project.proto \
  lib/proto/user.proto \
  lib/proto/flow.proto