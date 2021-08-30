#!/bin/bash

SIGNIFY_KEY=${SIGNIFY_KEY}
echo "${SIGNIFY_KEY}" > /tmp/signify.sec
export SIGNIFY_KEY=/tmp/signify.sec

mkdir -p ~/.m2
cat <<EOF > ~/.m2/settings.xml
<settings>
  <servers>
    <server>
      <id>ossrh</id>
      <username>${SONATYPE_USERNAME}</username>
      <password>${SONATYPE_PASWORD}</password>
    </server>
  </servers>
</settings>
EOF

exec "$@"
