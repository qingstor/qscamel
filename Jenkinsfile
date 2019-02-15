pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'

                sh 'export PATH=${GOPATH}/bin:$PATH'
                sh 'make --version'
                sh 'mkdir -p ${HOME}/.qscamel'

            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sh 'make integration-test'
            }
        }
    }
}
