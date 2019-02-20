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
            }
            steps {
                sh 'make build'

            }
        }
        stage('Test') {
            agent {
                label 'Oversea'
            }
            environment {
                GOPATH="${HOME}/go"
                PATH="${GOPATH}/bin:$PATH"
                GO111MODULE='on'
            }
            steps {
                sh 'make install'
                sh 'make integration-test'
            }
        }
    }
}
