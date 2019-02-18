pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                label 'Oversea'
            }
            environment {
                GOPATH="${HOME}/go"
                PATH="${GOPATH}/bin:$PATH"
                GO111MODULE='on'
                GOFLAGS='-mod=vendor'
            }
            steps {
                echo 'Building..'
                sh 'make install'
                sh 'mkdir -p ${HOME}/.qscamel'
                sh 'mkdir -p ${GOPATH}'
                sh 'ls'
                sh 'tar -czvf vendor.tar.gz vendor '
                stash includes: "vendor.tar.gz", name: 'pkg'
                stash includes: "go.*", name: 'module'

            }
        }
        stage('Test') {
            agent {
                label 'master'
            }
            environment {
                GOPATH="${HOME}/go"
                PATH="${GOPATH}/bin:$PATH"
                GO111MODULE='on'
                GOFLAGS='-mod=vendor'
            }
            steps {
                echo 'Testing..'
                sh 'echo $GOFLAGS'
                unstash 'pkg'
                unstash 'module'
                sh 'tar -xzvf vendor.tar.gz -C .'
                sh 'make install-after-check'
                sh 'make integration-test'
            }
        }
    }
}
