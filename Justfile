install-static:
        go install -ldflags '-w -linkmode external -extldflags "-static"'

package:
        docker build -t mtinside/chain:0.0.1 .
