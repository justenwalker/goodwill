concord-agent {
    server {
        apiBaseUrl = "http://concord-server:8001"
        websocketUrl = "ws://concord-server:8001/websocket"
        apiKey = "${concord_api_key}"
    }
    docker {
        host = "tcp://dind:6666"
    }
    runner {
        jvmParams = [
            "-Xmx128m",
            "-XX:+HeapDumpOnOutOfMemoryError",
            "-XX:HeapDumpPath=/tmp"
        ]
    }
}