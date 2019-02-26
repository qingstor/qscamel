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
                sh 'make build-runner'
                stash includes: '**/bin/qscamel', name: 'qscamel'
                stash includes: '**/bin/qscamel-runner', name: 'runner'
            }
        }
        stage('Intergration-Test') {
            agent {
                label 'Oversea'
            }
            environment {
                GOPATH="${HOME}/go"
                PATH="${GOPATH}/bin:$PATH"
                GO111MODULE='on'
            }
            steps {
                sh 'mkdir -p ./bin'
                unstash 'qscamel'
                unstash 'runner'
                sh 'mv ./bin/qscamel ${GOPATH}/bin/qscamel'
                sh 'mv ./bin/qscamel-runner ${GOPATH}/bin/qscamel-runner'
                sh 'make integration-test'
            }
        }
        stage('Edge-Test') {
            agent {
                label "master"
            }
            environment {
                GOPATH="${HOME}/go"
                PATH="${GOPATH}/bin:$PATH"
                GO111MODULE='on'
            }
            steps {
                sh 'mkdir -p ./bin'
                unstash 'qscamel'
                unstash 'runner'
                sh 'mv ./bin/qscamel ${GOPATH}/bin/qscamel'
                sh 'mv ./bin/qscamel-runner ${GOPATH}/bin/qscamel-runner'
                sh 'make edge-test'
            }
        }
    }
}
